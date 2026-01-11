package models

import "time"

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	HasRead   bool      `json:"has_read"`
	CreatedAt time.Time `json:"created_at"`
}
