package entities

import "time"

// User model (domain object)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Hide password in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateShortenUrlReq struct {
	OriginalUrl string `json:"original_url"`
}

type CreateShortenUrlRes struct {
	Id          string    `json:"id"`
	ShortUrl    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
