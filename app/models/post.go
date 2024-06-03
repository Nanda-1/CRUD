package models

import "time"

type Post struct {
	ID          int       `json:"PostId"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      string    `json:"status"`
	PublishDate time.Time `json:"publish_date"`
	UpdatedAt   time.Time `json:"updated_at"`
}
