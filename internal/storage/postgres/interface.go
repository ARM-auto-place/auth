package postgres

import "github.com/google/uuid"

type AuthStorage interface {
	UserRegister(email, password string) (userID uuid.UUID, code string, err error)
}
