package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manish-npx/go-auth-api/internal/config"
	handlers "github.com/manish-npx/go-auth-api/internal/http/handlers/auth" // <-- import your custom middleware
	httpmw "github.com/manish-npx/go-auth-api/internal/http/middleware"
	"github.com/manish-npx/go-auth-api/internal/repository"
	"github.com/manish-npx/go-auth-api/internal/service"
	"github.com/rs/zerolog/log"
)

func main() {
	// Config
	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("config load failed")
	}
	logger := config.ConfigureLogger(cfg)

	// DB
	ctx := context.Background()
	pool, err := config.NewDBPool(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("db connection failed")
	}
	defer pool.Close()

	// Dependencies
	userRepo := repository.User(pool)
	authSvc := service.NewAuthService(cfg, userRepo)
	authHandler := handlers.NewAuthHandler(authSvc)

	// Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Logger())

	//routes
	// Simple health check
	e.GET("/ok", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Ok"})
	})

	// API root
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "API root – available endpoints: /auth/register, /auth/login, /me",
		})
	})

	api := e.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected route
	protected := api.Group("/me", httpmw.RequireAuth(cfg))
	protected.GET("", func(c echo.Context) error {
		userID := c.Get("user_id").(int64)
		return c.JSON(200, map[string]interface{}{"message": "protected route", "user_id": userID})
	})

	// Graceful shutdown
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	go func() {
		if err := e.Start(addr); err != nil {
			logger.Info().Msg("server shutdown")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctxShutdown); err != nil {
		logger.Fatal().Err(err).Msg("server forced to shutdown")
	}
}
