package handler

import (
	"context"
	"log"
	"shorten-url/internal/entities"
	"shorten-url/internal/service"

	"github.com/labstack/echo/v4"
)

type (
	ShortenHandler interface {
		CreateShortenURL(c echo.Context) error
		GetShortenURL(c echo.Context) error
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

func (h *shortenHandler) GetShortenURL(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	originalUrl, err := h.shortenService.GetOriginalURL(ctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to get original url %s", err.Error())
		return c.JSON(404, "shortened URL not found")
	}

	return c.Redirect(301, originalUrl)
}

func (h *shortenHandler) CreateShortenURL(c echo.Context) error {

	ctx := context.Background()

	req := new(entities.CreateShortenUrlReq)

	if err := c.Bind(req); err != nil {
		return c.JSON(400, "invalid request")
	}

	shorten, err := h.shortenService.ShortenURL(ctx, req.OriginalUrl)
	if err != nil {
		return c.JSON(500, "failed to shorten url")
	}

	return c.JSON(200, shorten)

}
