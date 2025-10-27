package model

import "time"

// Example model - remove or modify as needed
type URL struct {
	ID          uint      `db:"id" json:"id"`
	ShortCode   string    `db:"short_code" json:"short_code"`
	OriginalURL string    `db:"original_url" json:"original_url"`
	ClickCount  int       `db:"click_count" json:"click_count"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type URLInterpeter struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

// // APIResponse represents a standard API response
// type APIResponse struct {
// 	Success bool        `json:"success"`
// 	Message string      `json:"message,omitempty"`
// 	Data    interface{} `json:"data,omitempty"`
// 	Error   string      `json:"error,omitempty"`
// }

// // ErrorResponse represents an error response
// type ErrorResponse struct {
// 	Success bool   `json:"success"`
// 	Error   string `json:"error"`
// 	Code    int    `json:"code,omitempty"`
// }
