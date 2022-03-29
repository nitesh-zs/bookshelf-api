package model

import (
	"github.com/google/uuid"
)

// Book struct holds book related information
type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Summary     string    `json:"summary"`
	Genre       string    `json:"genre"`
	PublishYear int       `json:"year"`
	RegNum      string    `json:"regNum"`
	Publisher   string    `json:"publisher"`
	Language    string    `json:"language"`
	PageCount   int       `json:"pageCount"`
	ImageURI    string    `json:"imageUri"`
}

// Media struct holds multimedia data
type Media struct {
	Data        []byte `json:"data"`
	ContentType string `json:"contentType"`
	Path        string `json:"path"`
}
