package model

import (
	"time"

	"github.com/google/uuid"
)

// Book struct holds book related information
type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Summary     string    `json:"summary"`
	Genre       string    `json:"genre"`
	PublishDate time.Time `json:"publishDate"`
	RegNum      string    `json:"regNum"`
	Publisher   string    `json:"publisher"`
	Language    string    `json:"language"`
	PageCount   int       `json:"pageCount"`
	URL         string    `json:"url"`
	Image       Media     `json:"image"`
}

// Media struct holds multimedia data
type Media struct {
	Data        []byte `json:"data"`
	ContentType string `json:"contentType"`
}
