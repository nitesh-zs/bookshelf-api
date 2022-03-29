package model

import "github.com/google/uuid"

type BookRes struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Summary     string    `json:"summary,omitempty"`
	Genre       string    `json:"genre"`
	PublishYear int       `json:"year,omitempty"`
	Publisher   string    `json:"publisher,omitempty"`
	ImageURI    string    `json:"imageUri"`
}
