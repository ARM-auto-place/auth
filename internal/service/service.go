package service

import (
	"context"
	"fmt"

	"github.com/ARM-auto-place/auth/internal/client"
	"github.com/ARM-auto-place/auth/internal/storage/postgres"
	"github.com/ARM-auto-place/auth/internal/storage/redis"
	interfaces "github.com/ARM-auto-place/auth/pkg/interfaces/externalapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type authService struct {
	storage      postgres.AuthStorage
	notification client.Notifications
	redisStorage *redis.CacheAuth
	logger       zerolog.Logger
}

func (a *authService) UserRegister(ctx context.Context, email string, password string) (userID uuid.UUID, err error) {
	a.logger.Info().Msg("call user register")

	userID, err = uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func (a *authService) CheckCode(ctx context.Context, userID uuid.UUID, code string) (err error) {
	redisCode, err := a.redisStorage.GetVerificationCode(ctx, userID)
	if err != nil {
		return err
	}
	if redisCode != code {
		// TODO: ДИМА СДЕЛАЙ НОРМАЛЬНУЮ ОБРАБОТКУ ОШИБКИ ЧЕРЕЗ НАШ ПАКЕТ КОТОРЫЙ В pkg/models/errors
		return fmt.Errorf("some error")
	}
	return nil

}

func NewAuthService(logger zerolog.Logger, storage postgres.AuthStorage, notification client.Notifications, redisStorage *redis.CacheAuth) interfaces.ExternalAPI {
	return &authService{
		storage:      storage,
		notification: notification,
		redisStorage: redisStorage,
		logger:       logger,
	}
}
