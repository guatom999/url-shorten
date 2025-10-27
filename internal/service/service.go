package service

import (
	"context"
	"errors"
	"log"
	"shorten-url/internal/entities"
	"shorten-url/internal/model"
	"shorten-url/internal/repository"
	"shorten-url/utils"
	"strconv"
)

type URLService interface {
	ShortenURL(pctx context.Context, originalURL string) (*entities.CreateShortenUrlRes, error)
	GetOriginalURL(pctx context.Context, shortCode string) (string, error)
	// GetURLStats(shortCode string) (*model.URLStats, error)
}

// Example service struct - modify as needed
type urlService struct {
	repo repository.URLRepository
}

// NewURLService creates a new URL service
func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) GetOriginalURL(pctx context.Context, shortCode string) (string, error) {

	url, err := s.repo.GetByShortCode(pctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to get original url %s", err.Error())
		return "", errors.New("failed to get original url ")
	}

	return url.OriginalURL, nil
}

func (s *urlService) ShortenURL(pctx context.Context, originalURL string) (*entities.CreateShortenUrlRes, error) {

	newUrl := utils.RandString(6)

	shortenInterpreter, err := s.repo.Create(pctx, &model.URL{
		ShortCode:   newUrl,
		OriginalURL: originalURL,
	})
	if err != nil {
		log.Printf("Error: failed to creat shorten url %s", err.Error())
		return nil, errors.New("failed to created shorten url ")
	}

	return &entities.CreateShortenUrlRes{
		Id:          strconv.Itoa(int(shortenInterpreter.ID)),
		ShortUrl:    newUrl,
		OriginalURL: originalURL,
		CreatedAt:   shortenInterpreter.CreatedAt,
		UpdatedAt:   shortenInterpreter.UpdatedAt,
	}, nil
}
