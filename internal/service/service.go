package service

import (
	"context"
	"errors"

	"github.com/ARM-auto-place/auth/internal/client"
	commonErrors "github.com/ARM-auto-place/auth/internal/models/errors"
	"github.com/ARM-auto-place/auth/internal/storage/postgres"
	"github.com/ARM-auto-place/auth/internal/storage/redis"
	interfaces "github.com/ARM-auto-place/auth/pkg/interfaces/externalapi"
	customErrors "github.com/ARM-auto-place/auth/pkg/models/errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type authService struct {
	storage      postgres.AuthStorage
	notification client.Notifications
	redisStorage *redis.CacheAuth
	logger       zerolog.Logger
}

func (a *authService) UserRegister(ctx context.Context, email string, password string) (err error) {
	a.logger.Info().Msg("call user register")

	user, err := a.storage.GetUserByEmail(ctx, email)
	if err != nil {
		return customErrors.InternalServerError().WithWrapped(err, "failed to get user by email")
	}
	if user != nil {
		return customErrors.ErrUserAlreadyExists()
	}
	code, err := generateVerificationCode()
	if err != nil {
		return customErrors.InternalServerError().WithWrapped(err, "failed to generate verification code")
	}
	if err = a.redisStorage.SaveVerificationCode(ctx, email, password, code); err != nil {
		return customErrors.InternalServerError().WithWrapped(err, "failed to save verification code")
	}
	if err = a.notification.SendEmailCode(ctx, email, code); err != nil {
		return customErrors.ErrFailedSendVerificationCode().WithWrapped(err, "failed to send email")
	}

	return nil

}

func (a *authService) CheckCode(ctx context.Context, email string, code string) (userID uuid.UUID, err error) {
	redisCode, password, err := a.redisStorage.GetVerificationCode(ctx, email)
	if err != nil {
		if errors.Is(err, commonErrors.ErrCodeExpiredOrNotFound) {
			return uuid.Nil, customErrors.ErrVerificationCodeNotFound()
		}
		return uuid.Nil, customErrors.InternalServerError().WithWrapped(err, "failed to get verification code")
	}
	if redisCode != code {
		return uuid.Nil, customErrors.ErrInvalidVerificationCode()
	}
	if err = a.redisStorage.DeleteVerificationCode(ctx, email); err != nil {
		return uuid.Nil, customErrors.InternalServerError().WithWrapped(err, "failed to delete verification code")
	}
	hashPass, err := hashPassword(password)
	if err != nil {
		return uuid.Nil, customErrors.InternalServerError().WithWrapped(err, "failed to hash password")
	}
	userID, err = a.storage.UserRegister(ctx, email, hashPass)
	if err != nil {
		return uuid.Nil, customErrors.InternalServerError().WithWrapped(err, "failed to register user")
	}
	return userID, nil

}

func NewAuthService(logger zerolog.Logger, storage postgres.AuthStorage, notification client.Notifications, redisStorage *redis.CacheAuth) interfaces.ExternalAPI {
	return &authService{
		storage:      storage,
		notification: notification,
		redisStorage: redisStorage,
		logger:       logger,
	}
}

func (a *authService) RenewVerifyCode(ctx context.Context, email string) error {
	_, password, err := a.redisStorage.GetVerificationCode(ctx, email)
	if err != nil {
		if errors.Is(err, commonErrors.ErrCodeExpiredOrNotFound) {
			return customErrors.ErrVerificationCodeNotFound()
		}
		return customErrors.InternalServerError().WithWrapped(
			err,
			"failed to get verification code",
		)
	}

	allowed, err := a.redisStorage.TryAcquireRenewCodeLimit(ctx, email)
	if err != nil {
		return customErrors.InternalServerError()
	}

	if !allowed {
		return customErrors.ErrTooManyRequests()
	}

	code, err := generateVerificationCode()
	if err != nil {
		return customErrors.InternalServerError().WithWrapped(
			err,
			"failed to generate verification code",
		)
	}

	if err := a.redisStorage.SaveVerificationCode(ctx, email, password, code); err != nil {
		return customErrors.InternalServerError().WithWrapped(
			err,
			"failed to save verification code",
		)
	}
	if err := a.notification.SendEmailCode(ctx, email, code); err != nil {
		return customErrors.ErrFailedSendVerificationCode().WithWrapped(
			err,
			"failed to send email",
		)
	}

	return nil
}

func (a *authService) Login(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := a.storage.GetUserCredentialsByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, customErrors.InternalServerError().WithWrapped(
			err,
			"failed to get user credentials",
		)
	}

	if user == nil {
		return uuid.Nil, customErrors.ErrInvalidCredentials()
	}

	if !checkPassword(password, user.HashPassword) {
		return uuid.Nil, customErrors.ErrInvalidCredentials()
	}

	return user.UserID, nil
}
