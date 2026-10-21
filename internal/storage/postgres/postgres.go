package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ARM-auto-place/auth/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

var defaultTimeout = 5 * time.Second

type storage struct {
	connectManager ConnectManager
	logger         zerolog.Logger
}

func NewAuthStorage(logger zerolog.Logger, connectManager ConnectManager) AuthStorage {
	return &storage{
		connectManager: connectManager,
		logger:         logger,
	}
}
func (s *storage) UserRegister(ctx context.Context, email, password string) (userID uuid.UUID, err error) {
	masterConn := s.connectManager.GetMasterConn()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	err = masterConn.QueryRow(ctx, sqlInsertUser, email, password).Scan(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("UserRegister: postgresql: %w", err)
	}

	return userID, nil
}
func (s *storage) GetUserByEmail(ctx context.Context, email string) (*models.UserProfile, error) {
	replicaConn := s.connectManager.GetReplicaConn()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	var user models.UserProfile
	err := replicaConn.QueryRow(ctx, sqlSelectUserByEmail, email).Scan(
		&user.UserID,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByEmail: postgresql: %w", err)
	}
	return &user, nil
}

func (s *storage) GetUserCredentialsByEmail(ctx context.Context, email string) (*models.UserCredentials, error) {
	replicaConn := s.connectManager.GetReplicaConn()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	var user models.UserCredentials
	err := replicaConn.QueryRow(ctx, sqlSelectUserCredentialsByEmail, email).Scan(
		&user.UserID,
		&user.Email,
		&user.HashPassword,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("GetUserCredentialsByEmail: postgresql: %w", err)
	}
	return &user, nil
}
