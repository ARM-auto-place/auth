package postgres

import (
	"context"

	"github.com/ARM-auto-place/auth/internal/models"
	"github.com/google/uuid"
)

type AuthStorage interface {
	UserRegister(ctx context.Context, email, password string) (userID uuid.UUID, err error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserProfile, error)
	GetUserCredentialsByEmail(ctx context.Context, email string) (*models.UserCredentials, error)
}
