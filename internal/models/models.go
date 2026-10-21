package models

import "github.com/google/uuid"

type UserProfile struct {
	UserID uuid.UUID
	Email  string
}
type UserCredentials struct {
	UserID       uuid.UUID
	Email        string
	HashPassword string
}
