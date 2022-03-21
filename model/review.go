package model

import (
	"time"

	"github.com/google/uuid"
)

// Review struct holds review related information
type Review struct {
	ID     uuid.UUID `json:"id"`
	Text   string    `json:"text"`
	Rating int       `json:"rating"`
	CTime  time.Time `json:"ctime"`
	UserID uuid.UUID `json:"userID"`
	BookID uuid.UUID `json:"bookID"`
}
