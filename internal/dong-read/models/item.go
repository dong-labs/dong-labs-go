package models

import "time"

type Item struct {
	ID        int       `json:"id"`
	Content   string    `json:"content,omitempty"`
	URL       string    `json:"url,omitempty"`
	Title     string    `json:"title,omitempty"`
	Note      string    `json:"note,omitempty"`
	Source    string    `json:"source,omitempty"`
	Type      string    `json:"type"`
	Tags      string    `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemList struct {
	Total int    `json:"total"`
	Items []Item `json:"items"`
}
