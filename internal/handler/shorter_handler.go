package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"shorten-url/internal/entities"
	appErrors "shorten-url/internal/errors"
	"shorten-url/internal/service"

	"github.com/labstack/echo/v4"
)

type (
	ShortenHandler interface {
		CreateShortenURL(c echo.Context) error
		GetShortenURL(c echo.Context) error
		RetrieveOriginalURL(c echo.Context) error
		UpdateShortenURL(c echo.Context) error
		DeleteUrl(c echo.Context) error
		GetUrlStatic(c echo.Context) error
	}

	shortenHandler struct {
		shortenService service.URLService
	}
)

func NewHandler(shortenService service.URLService) ShortenHandler {
	return &shortenHandler{
		shortenService: shortenService,
	}
}

// Helper function to handle errors and return appropriate HTTP status
func (h *shortenHandler) handleError(c echo.Context, err error) error {
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case appErrors.NotFound:
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": appErr.Message,
			})
		case appErrors.InvalidInput:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": appErr.Message,
			})
		case appErrors.Conflict:
			return c.JSON(http.StatusConflict, map[string]string{
				"error": appErr.Message,
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": appErr.Message,
			})
		}
	}

	// Default error handling
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"error": "Internal server error",
	})
}

func (h *shortenHandler) GetShortenURL(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	originalUrl, err := h.shortenService.GetOriginalURL(ctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to get original url %s", err.Error())
		return h.handleError(c, err)
	}

	return c.Redirect(http.StatusMovedPermanently, originalUrl)
}

func (h *shortenHandler) UpdateShortenURL(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	updateUrlReq := new(entities.UpdateUrlReq)

	if err := c.Bind(updateUrlReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	res, err := h.shortenService.UpdateShortUrl(ctx, shortCode, updateUrlReq.Url)
	if err != nil {
		log.Printf("Error: failed to update short url %s", err.Error())
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *shortenHandler) DeleteUrl(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	if err := h.shortenService.DeleteShortUrl(ctx, shortCode); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *shortenHandler) RetrieveOriginalURL(c echo.Context) error {
	ctx := context.Background()

	shortCode := c.Param("short_code")

	retriveUrlRes, err := h.shortenService.RetrieveOriginalURL(ctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to retrieve original url %s", err.Error())
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, retriveUrlRes)
}

func (h *shortenHandler) CreateShortenURL(c echo.Context) error {

	ctx := context.Background()

	req := new(entities.CreateShortenUrlReq)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	shorten, err := h.shortenService.ShortenURL(ctx, req.OriginalUrl)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, shorten)

}

func (h *shortenHandler) GetUrlStatic(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	urlStat, err := h.shortenService.GetUrlStatic(ctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to get url stat %s", err.Error())
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, urlStat)

}
