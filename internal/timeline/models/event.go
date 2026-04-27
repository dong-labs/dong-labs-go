package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Date        string    `json:"date"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category,omitempty"`
	Tags        string    `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
