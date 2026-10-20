package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	commonErrors "github.com/ARM-auto-place/auth/internal/models/errors"
	customErrors "github.com/ARM-auto-place/auth/pkg/models/errors"
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
func (r *CacheAuth) SaveVerificationCode(ctx context.Context, email string, password string, code string) error {
	value := fmt.Sprintf("%s,%s", code, password)
	err := r.client.Set(ctx, verificationCodeKey(email), value, timeout).Err()
	if err != nil {
		r.logger.Err(err).Msgf("failed to save verification code:%s", email)
		return err
	}
	return nil
}

func (r *CacheAuth) GetVerificationCode(ctx context.Context, email string) (string, string, error) {

	val, err := r.client.Get(ctx, verificationCodeKey(email)).Result()
	if err == redis.Nil {
		return "", "", commonErrors.ErrCodeExpiredOrNotFound
	} else if err != nil {
		r.logger.Err(err).Msgf("failed to get verification code:%s", email)
		return "", "", err
	}
	res := strings.Split(val, ",")
	return res[0], res[1], nil
}

func (r *CacheAuth) DeleteVerificationCode(ctx context.Context, email string) error {
	err := r.client.Del(ctx, verificationCodeKey(email)).Err()
	if err != nil {
		r.logger.Err(err).Msgf("failed to delete verification code:%s", email)
		return err
	}
	r.logger.Info().Msgf("verification code deleted:%s", email)
	return nil
}

func verificationCodeKey(email string) string {
	return "verify_code:" + email
}
