package entities

import "time"

type RetriveOriginalUrlRes struct {
	Id          uint      `json:"id"`
	OriginalUrl string    `json:"original_url"`
	ShortUrl    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUrlReq struct {
	Url string `json:"url"`
}
type UpdateUrlRes struct {
	Id          string    `json:"id"`
	ShortUrl    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateShortenUrlRes struct {
	Id          string    `json:"id"`
	ShortUrl    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type CreateShortenUrlReq struct {
	OriginalUrl string `json:"original_url"`
}

type UrlStaticRes struct {
	Id          string    `json:"id"`
	Url         string    `json:"url"`
	ShortCode   string    `json:"shortCode"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	AccessCount int       `json:"accessCount"`
}
