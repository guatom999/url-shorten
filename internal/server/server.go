package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"shorten-url/configs"
	"shorten-url/internal/handler"
	"shorten-url/internal/repository"
	"shorten-url/internal/service"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Server interface {
		Start()
	}

	server struct {
		app *echo.Echo
		cfg *configs.Config
		db  *sqlx.DB
	}
)

func NewEchoServer(cfg *configs.Config, db *sqlx.DB) Server {
	return &server{
		app: echo.New(),
		cfg: cfg,
		db:  db,
	}
}

func (s *server) Start() {

	ctx := context.Background()

	// Request Timeout
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      time.Second * 10,
	}))

	//Cors
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	s.app.Use(middleware.Logger())

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	go s.gracefulShutdown(ctx, close)

	s.ShortenModules()

	if err := s.app.Start(fmt.Sprintf(":%s", s.cfg.Server.Port)); err != nil {
		log.Printf("Server stopped: %v", err)
	}

}

func (s *server) ShortenModules() {

	shortenRepo := repository.NewURLRepository(s.db)
	shortenService := service.NewURLService(shortenRepo)
	shortenHandler := handler.NewHandler(shortenService)

	route := s.app.Group("/v1/shorten")

	route.GET("/:short_code", shortenHandler.GetShortenURL)

	route.POST("/", shortenHandler.CreateShortenURL)

	// implement leayer of shorten

}

func (s *server) gracefulShutdown(pctx context.Context, close <-chan os.Signal) {

	<-close

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		s.app.Logger.Fatal(err)
	}

	log.Println("Shuttung Down Server....")

}
