package repository

import (
	"context"
	"errors"
	"log"
	"shorten-url/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

// Example repository interface - modify as needed
type URLRepository interface {
	Create(ctx context.Context, url *model.URL) (*model.URLInterpeter, error)
	GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	UpdateShortUrl(pctx context.Context, shortCode string, updatedUrl string) (*model.URL, error)
	DeleteByShortCode(ctx context.Context, shortCode string) error
	UpdateShortUrlCount(pctx context.Context, shortCode string) error
	IsShortCodeExists(pctx context.Context, shortCode string) bool
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

func (r *urlRepository) UpdateShortUrlCount(pctx context.Context, shortCode string) error {
	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	query := `UPDATE urls 
              SET click_count = click_count + 1, updated_at = CURRENT_TIMESTAMP 
              WHERE short_code = $1`

	result, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		log.Printf("Error updating click count for short code %s: %v", shortCode, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected for short code %s: %v", shortCode, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No URL found with short code: %s", shortCode)
		return errors.New("no URL found with the given short code")
	}

	return nil

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

func (r *urlRepository) UpdateShortUrl(pctx context.Context, shortCode string, updatedUrl string) (*model.URL, error) {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	query := `UPDATE urls 
              SET original_url = $1, updated_at = CURRENT_TIMESTAMP 
              WHERE short_code = $2 
              RETURNING id, short_code, original_url, click_count, created_at, updated_at`

	urlData := new(model.URL)

	if err := r.db.QueryRowContext(ctx, query, updatedUrl, shortCode).Scan(
		&urlData.ID,
		&urlData.ShortCode,
		&urlData.OriginalURL,
		&urlData.ClickCount,
		&urlData.CreatedAt,
		&urlData.UpdatedAt,
	); err != nil {
		log.Printf("Error updating short url: %v", err)
		return nil, err
	}

	return urlData, nil
}

func (r *urlRepository) DeleteByShortCode(ctx context.Context, shortCode string) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := `DELETE FROM urls WHERE short_code = $1 `

	result, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		log.Printf("Error deleting URL by short code: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected for delete by shortcode %s: %v", shortCode, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("no URL found with short code: %s", shortCode)
		return errors.New("no URL found with the given short code")
	}

	log.Printf("Successfully deleted URL with short code: %s", shortCode)
	return nil
}

func (r *urlRepository) Create(ctx context.Context, url *model.URL) (*model.URLInterpeter, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	url.ClickCount = 1

	query := `INSERT INTO urls (short_code, original_url,click_count) VALUES ($1, $2 ,$3) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, url.ShortCode, url.OriginalURL, url.ClickCount).Scan(&url.ID, &url.CreatedAt, &url.UpdatedAt)

	return &model.URLInterpeter{
		ID:        url.ID,
		CreatedAt: url.CreatedAt,
		UpdatedAt: url.UpdatedAt,
	}, err
}

func (r *urlRepository) IsShortCodeExists(pctx context.Context, shortCode string) bool {

	ctx, cancel := context.WithTimeout(pctx, time.Second*5)
	defer cancel()

	query := `SELECT COUNT(1) FROM urls WHERE short_code = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(&count)
	if err != nil {
		log.Printf("Error checking if short code exists: %v", err)
		return false
	}

	return count > 0
}
