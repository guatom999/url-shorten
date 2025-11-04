package handler

import (
	"context"
	"log"
	"net/http"
	"shorten-url/internal/entities"
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

func (h *shortenHandler) GetShortenURL(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	originalUrl, err := h.shortenService.GetOriginalURL(ctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to get original url %s", err.Error())
		return c.JSON(404, "shortened URL not found")
	}

	return c.Redirect(http.StatusMovedPermanently, originalUrl)
}

func (h *shortenHandler) UpdateShortenURL(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	updateUrlReq := new(entities.UpdateUrlReq)

	if err := c.Bind(updateUrlReq); err != nil {
		return c.JSON(400, "invalid request")
	}

	res, err := h.shortenService.UpdateShortUrl(ctx, shortCode, updateUrlReq.Url)
	if err != nil {
		log.Printf("Error: failed to update short url %s", err.Error())
		return c.JSON(500, "failed to update short url")
	}

	return c.JSON(http.StatusOK, res)
}

func (h *shortenHandler) DeleteUrl(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	if err := h.shortenService.DeleteShortUrl(ctx, shortCode); err != nil {
		return c.JSON(500, "failed to delete short url")
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *shortenHandler) RetrieveOriginalURL(c echo.Context) error {
	ctx := context.Background()

	shortCode := c.Param("short_code")

	retriveUrlRes, err := h.shortenService.RetrieveOriginalURL(ctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to retrieve original url %s", err.Error())
		return c.JSON(404, "shortened URL not found")
	}

	return c.JSON(200, retriveUrlRes)
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

func (h *shortenHandler) GetUrlStatic(c echo.Context) error {

	ctx := context.Background()

	shortCode := c.Param("short_code")

	urlStat, err := h.shortenService.GetUrlStatic(ctx, shortCode)
	if err != nil {
		log.Printf("Error: failed to get url stat %s", err.Error())
		return c.JSON(404, "shortened URL not found")
	}

	return c.JSON(200, urlStat)

}
