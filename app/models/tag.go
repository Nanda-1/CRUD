package models

import "time"

type Tag struct {
	ID        int       `json:"tagId"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
