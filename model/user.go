package model

import "github.com/google/uuid"

// User struct holds user related information
type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Type  string    `json:"type"`
}
