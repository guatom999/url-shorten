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
	RetrieveOriginalURL(pctx context.Context, shortCode string) (*entities.RetriveOriginalUrlRes, error)
	UpdateShortUrl(pctx context.Context, shortCode string, updatedUrl string) (*model.URL, error)
	DeleteShortUrl(pctx context.Context, shortCode string) error
	GetUrlStatic(pctx context.Context, shortCode string) (*entities.UrlStaticRes, error)
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
		return "", err
	}

	if err := s.repo.UpdateShortUrlCount(pctx, shortCode); err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

func (s *urlService) DeleteShortUrl(pctx context.Context, shortCode string) error {

	if err := s.repo.DeleteByShortCode(pctx, shortCode); err != nil {
		return errors.New("failed to delete short url ")
	}

	return nil
}

func (s *urlService) RetrieveOriginalURL(pctx context.Context, shortCode string) (*entities.RetriveOriginalUrlRes, error) {

	url, err := s.repo.GetByShortCode(pctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to retrieve original url %s", err.Error())
		return nil, errors.New("failed to retrieve original url ")
	}

	return &entities.RetriveOriginalUrlRes{
		Id:          url.ID,
		OriginalUrl: url.OriginalURL,
		ShortUrl:    url.ShortCode,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
	}, nil
}

func (s *urlService) UpdateShortUrl(pctx context.Context, shortCode string, updatedUrl string) (*model.URL, error) {

	url, err := s.repo.UpdateShortUrl(pctx, shortCode, updatedUrl)
	if err != nil {
		log.Printf("Error: failed to update short url %s", err.Error())
		return nil, errors.New("failed to update short url ")
	}

	return url, nil

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

func (s *urlService) GetUrlStatic(pctx context.Context, shortCode string) (*entities.UrlStaticRes, error) {

	url, err := s.repo.GetByShortCode(pctx, shortCode)
	if err != nil {
		return nil, err
	}

	return &entities.UrlStaticRes{
		Id:          strconv.Itoa(int(url.ID)),
		Url:         url.OriginalURL,
		ShortCode:   url.ShortCode,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
		AccessCount: url.ClickCount,
	}, nil
}
