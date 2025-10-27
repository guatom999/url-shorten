package repository

import (
	"context"
	"log"
	"shorten-url/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

// Example repository interface - modify as needed
type URLRepository interface {
	Create(ctx context.Context, url *model.URL) (*model.URLInterpeter, error)
	GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	// GetByID(id uint) (*model.URL, error)
	// Update(url *model.URL) error
	// Delete(id uint) error
}

// Example repository struct - modify as needed
type urlRepository struct {
	db *sqlx.DB
}

// NewURLRepository creates a new URL repository
func NewURLRepository(db *sqlx.DB) URLRepository {
	return &urlRepository{
		db: db,
	}
}

func (r *urlRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `SELECT id , short_code , original_url , click_count , created_at , updated_at FROM urls WHERE short_code = $1`

	url := new(model.URL)
	err := r.db.GetContext(ctx, url, query, shortCode)
	if err != nil {
		log.Printf("Error getting URL by short code: %v", err)
		return nil, err
	}

	return url, nil
}

func (r *urlRepository) Create(ctx context.Context, url *model.URL) (*model.URLInterpeter, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	url.ClickCount = 1

	query := `INSERT INTO urls (short_code, original_url,click_count) VALUES ($1, $2 ,$3) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, url.ShortCode, url.OriginalURL, url.ClickCount).Scan(&url.ID, &url.CreatedAt, &url.UpdatedAt)

	return &model.URLInterpeter{
		ID:        url.ID,
		CreatedAt: url.CreatedAt,
		UpdatedAt: url.UpdatedAt,
	}, err
}
