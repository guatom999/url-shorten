package service

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"shorten-url/internal/entities"
	appErrors "shorten-url/internal/errors"
	"shorten-url/internal/model"
	"shorten-url/internal/repository"
	"shorten-url/utils"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"
)

type URLService interface {
	ShortenURL(pctx context.Context, originalURL string) (*entities.CreateShortenUrlRes, error)
	CreateQrCode(pctx context.Context, shortCode string) (*entities.CreateQrCodeRes, error)
	GetOriginalURL(pctx context.Context, shortCode string) (string, error)
	RetrieveOriginalURL(pctx context.Context, shortCode string) (*entities.RetriveOriginalUrlRes, error)
	UpdateShortUrl(pctx context.Context, shortCode string, updatedUrl string) (*model.URL, error)
	DeleteShortUrl(pctx context.Context, shortCode string) error
	GetUrlStatic(pctx context.Context, shortCode string) (*entities.UrlStaticRes, error)
}

type urlService struct {
	repo    repository.URLRepository
	baseURL string
}

func NewURLService(repo repository.URLRepository, baseURL ...string) URLService {
	resolvedBaseURL := "http://localhost:8080"
	if len(baseURL) > 0 && strings.TrimSpace(baseURL[0]) != "" {
		resolvedBaseURL = strings.TrimRight(strings.TrimSpace(baseURL[0]), "/")
	}

	return &urlService{
		repo:    repo,
		baseURL: resolvedBaseURL,
	}
}

func (s *urlService) GetOriginalURL(pctx context.Context, shortCode string) (string, error) {

	url, err := s.repo.GetByShortCode(pctx, shortCode)
	if err != nil {
		return "", appErrors.NewNotFoundError("short url was not found")
	}

	if err := s.repo.UpdateShortUrlCount(pctx, shortCode); err != nil {
		return "", appErrors.NewInternalError("failed to update click count", err)
	}

	return url.OriginalURL, nil
}

func (s *urlService) DeleteShortUrl(pctx context.Context, shortCode string) error {

	if !s.repo.IsShortCodeExists(pctx, shortCode) {
		return appErrors.NewNotFoundError("short url was not found")
	}

	if err := s.repo.DeleteByShortCode(pctx, shortCode); err != nil {
		return appErrors.NewInternalError("failed to delete short url", err)
	}

	return nil
}

func (s *urlService) RetrieveOriginalURL(pctx context.Context, shortCode string) (*entities.RetriveOriginalUrlRes, error) {

	url, err := s.repo.GetByShortCode(pctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to retrieve original url %s", err.Error())
		return nil, appErrors.NewNotFoundError("short url was not found")
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
		return nil, appErrors.NewInternalError("failed to update short url", err)
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
		return nil, appErrors.NewInternalError("failed to created shorten url", err)
	}

	return &entities.CreateShortenUrlRes{
		Id:          strconv.Itoa(int(shortenInterpreter.ID)),
		ShortUrl:    newUrl,
		OriginalURL: originalURL,
		CreatedAt:   shortenInterpreter.CreatedAt,
		UpdatedAt:   shortenInterpreter.UpdatedAt,
	}, nil
}

func (s *urlService) CreateQrCode(pctx context.Context, shortCode string) (*entities.CreateQrCodeRes, error) {

	if strings.TrimSpace(shortCode) == "" {
		return nil, appErrors.NewInvalidInputError("short code is required")
	}

	existingURL, err := s.repo.GetByShortCode(pctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to find short url for qr code %s", err.Error())
		return nil, appErrors.NewNotFoundError("short url was not found")
	}

	if err := os.MkdirAll("temp", 0o755); err != nil {
		log.Printf("Error: failed to create temp directory %s", err.Error())
		return nil, appErrors.NewInternalError("failed to prepare qrcode directory", err)
	}

	targetURL, err := url.JoinPath(s.baseURL, shortCode)
	if err != nil {
		log.Printf("Error: failed to build qrcode target url %s", err.Error())
		return nil, appErrors.NewInternalError("failed to build qrcode url", err)
	}

	png, err := qrcode.Encode(targetURL, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error: failed to encode qr code %s", err.Error())
		return nil, appErrors.NewInternalError("failed to create qrcode", err)
	}

	fileName := fmt.Sprintf("qrcode_%s.png", utils.RandString(6))
	filePath := filepath.Join("temp", fileName)

	if err := os.WriteFile(filePath, png, 0o644); err != nil {
		log.Printf("Error: failed to save qr code file %s", err.Error())
		return nil, appErrors.NewInternalError("failed to save qrcode image", err)
	}

	publicImageURL, err := url.JoinPath(s.baseURL, "temp", fileName)
	if err != nil {
		log.Printf("Error: failed to build public qrcode url %s", err.Error())
		return nil, appErrors.NewInternalError("failed to build public qrcode url", err)
	}

	existingURL.QrCodeUrl = publicImageURL
	updatedURL, err := s.repo.CreateQrCode(pctx, existingURL)
	if err != nil {
		log.Printf("Error: failed to save qrcode url %s", err.Error())
		return nil, appErrors.NewInternalError("failed to save qrcode url", err)
	}

	return &entities.CreateQrCodeRes{
		Id:          strconv.Itoa(int(updatedURL.ID)),
		ShortUrl:    existingURL.ShortCode,
		OriginalURL: existingURL.OriginalURL,
		QrCodeURL:   publicImageURL,
		CreatedAt:   updatedURL.CreatedAt,
		UpdatedAt:   updatedURL.UpdatedAt,
	}, nil
}

func (s *urlService) GetUrlStatic(pctx context.Context, shortCode string) (*entities.UrlStaticRes, error) {

	url, err := s.repo.GetByShortCode(pctx, shortCode)
	if err != nil {
		return nil, appErrors.NewNotFoundError("short url was not found")
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
