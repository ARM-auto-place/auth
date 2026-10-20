package redis

import (
	"context"
	"time"

	commonErrors "github.com/ARM-auto-place/auth/internal/models/errors"
	customErrors "github.com/ARM-auto-place/auth/pkg/models/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	timeout = 5 * time.Second
)

type CacheAuth struct {
	client redis.UniversalClient
	logger zerolog.Logger
}

func New(logger zerolog.Logger, client redis.UniversalClient) (*CacheAuth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, customErrors.Wrap(err, "New: redis ping error")
	}

	return &CacheAuth{
		client: client,
		logger: logger,
	}, nil
}
func (r *CacheAuth) SaveVerificationCode(ctx context.Context, userID uuid.UUID, email string, password string, code string) error {
	err := r.client.Set(ctx, verificationCodeKey(userID), code, timeout).Err()
	if err != nil {
		r.logger.Err(err).Msgf("failed to save verification code:%s", userID.String())
		return err
	}
	return nil
}

func (r *CacheAuth) GetVerificationCode(ctx context.Context, userID uuid.UUID) (string, error) {
	val, err := r.client.Get(ctx, verificationCodeKey(userID)).Result()
	if err == redis.Nil {
		return "", commonErrors.ErrCodeExpiredOrNotFound
	} else if err != nil {
		r.logger.Err(err).Msgf("failed to get verification code:%s", userID.String())
		return "", err
	}
	return val, nil
}

func (r *CacheAuth) DeleteVerificationCode(ctx context.Context, userID uuid.UUID) error {
	err := r.client.Del(ctx, verificationCodeKey(userID)).Err()
	if err != nil {
		r.logger.Err(err).Msgf("failed to delete verification code:%s", userID.String())
		return err
	}
	r.logger.Info().Msgf("verification code deleted:%s", userID.String())
	return nil
}

func verificationCodeKey(userID uuid.UUID) string {
	return "verify_code:" + userID.String()
}
