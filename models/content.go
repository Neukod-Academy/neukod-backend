package models

import "time"

type Content struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsPublished bool      `json:"is_published"`
	Tags        []string  `json:"tags"`
}
