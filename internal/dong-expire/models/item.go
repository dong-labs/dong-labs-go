package models

import "time"

type Item struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	ExpireDate   string    `json:"expire_date"`
	ReminderDays int       `json:"reminder_days"`
	Tags         string    `json:"tags"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RenewHistory struct {
	ID             int       `json:"id"`
	ItemID         int       `json:"item_id"`
	OldExpireDate  string    `json:"old_expire_date"`
	NewExpireDate  string    `json:"new_expire_date"`
	RenewedAt      time.Time `json:"renewed_at"`
}
