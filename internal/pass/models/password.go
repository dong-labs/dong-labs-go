package models

import "time"

type Password struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	URL       string    `json:"url,omitempty"`
	Category  string    `json:"category,omitempty"`
	Tags      string    `json:"tags,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PasswordList struct {
	Total int        `json:"total"`
	Items []Password `json:"items"`
}

type StatsResponse struct {
	Total      int          `json:"total"`
	ByCategory []CategoryStat `json:"by_category,omitempty"`
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}
