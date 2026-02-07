package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uzears/golangcodes/research-api/internal/platform/config"
	"github.com/uzears/golangcodes/research-api/internal/platform/database"
	"github.com/uzears/golangcodes/research-api/internal/platform/http"
	"github.com/uzears/golangcodes/research-api/internal/platform/logger"
	"github.com/uzears/golangcodes/research-api/internal/platform/middleware"
	"github.com/uzears/golangcodes/research-api/internal/auth"
)

func main() {

	// 1. Load environment variables (local dev only)
	_ = godotenv.Load()

	// 2. Load application configuration
	cfg := config.Load()

	// 3. Initialize logger (Zerolog via interface)
	logr := logger.New()
	logr.Info(
		"application starting",
		"app", cfg.AppName,
		"env", cfg.Env,
	)

	// 4. Initialize database connection
	db := database.Connect(cfg.DB.URL)

	// 5. Initialize domain layers
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.JWT.Secret, logr)
	authHandler := auth.NewHandler(authService, logr)

	// 6. Initialize HTTP server (Gin)
	addr := ":" + strconv.Itoa(cfg.Port)
	server := http.NewServer(addr, cfg.Env)
	router := server.Router()
	router.Use(middleware.CORS("http://localhost:3000"))
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logr))

	// 7. Register routes (NO business logic here)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	// Protected routes example
	protected := router.Group("/")
	protected.Use(middleware.JWT(cfg.JWT.Secret))
	{
		protected.GET("/me", authHandler.Me)
	}

	// 8. Start HTTP server (non-blocking)
	go func() {
		logr.Info("http server started", "addr", addr)
		if err := server.Start(); err != nil {
			logr.Error("http server error", "err", err)
		}
	}()
	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logr.Warn("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logr.Error("server shutdown failed", "err", err)
	}

	logr.Info("application stopped gracefully")
}
