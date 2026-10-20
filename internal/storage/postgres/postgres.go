package postgres

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

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
func (s *storage) UserRegister(email, password string) (userID uuid.UUID, code string, err error) {

	return uuid.Nil, "", nil
}
