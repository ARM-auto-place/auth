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
	timeout             = 5 * time.Second
	verificationCodeTTL = 10 * time.Minute
	renewCodeLimitTTL   = 60 * time.Second
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
	err := r.client.Set(ctx, verificationCodeKey(email), value, verificationCodeTTL).Err()
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
	res := strings.SplitN(val, ",", 2)
	if len(res) != 2 {
		return "", "", fmt.Errorf("invalid verification code data format")
	}
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

func (r *CacheAuth) TryAcquireRenewCodeLimit(ctx context.Context, email string) (bool, error) {
	return r.client.SetNX(
		ctx,
		renewCodeLimitKey(email),
		"1",
		renewCodeLimitTTL,
	).Result()
}

func verificationCodeKey(email string) string {
	return "verify_code:" + email
}

func renewCodeLimitKey(email string) string {
	return "renew_code_limit" + email
}
