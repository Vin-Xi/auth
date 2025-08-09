package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/Vin-Xi/auth/internal/database"
	internal "github.com/Vin-Xi/auth/internal/database"
	"github.com/Vin-Xi/auth/internal/transport"
	util "github.com/Vin-Xi/auth/internal/util"
	"github.com/Vin-Xi/auth/pkg/logger"
	"github.com/gin-gonic/gin"

	service "github.com/Vin-Xi/auth/internal/service"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	ctx := context.Background()

	logger.Init()

	pool, err := internal.InitDB(ctx, databaseUrl)

	if err != nil {
		logger.Log.ErrorWithStack("Operation failed", err)
		os.Exit(1)
	}

	defer pool.Close()

	userRepo := database.NewPostresRepository(pool)
	userService := service.NewService(userRepo)
	jwtEngine := util.NewJWTEngine(jwtSecret, 15*time.Minute)
	userHttpHandler := transport.UserHandler{UserService: userService, JwtEngine: jwtEngine}
	webHttpHandler := transport.WebHandler{}

	router := gin.Default()

	transport.SetTemplateFS(router)
	webHttpHandler.RegisterRoutes(router)
	userHttpHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: router,
	}

	go func() {
		logger.Log.Info("Server started listening on 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.ErrorWithStack("Fatal error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shut down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Info("Server forced to shut down")
	}

	logger.Log.Info("Server exiting")
}
