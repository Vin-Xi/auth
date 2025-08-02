package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/Vin-Xi/auth/internal/database"
	internal "github.com/Vin-Xi/auth/internal/database"
	"github.com/Vin-Xi/auth/internal/transport"
	util "github.com/Vin-Xi/auth/internal/util"
	"github.com/gin-gonic/gin"

	service "github.com/Vin-Xi/auth/internal/service"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	ctx := context.Background()
	fmt.Print(databaseUrl)
	pool, err := internal.InitDB(ctx, databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "db initialization failed: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Connection is successful!")
	}

	defer pool.Close()

	userRepo := database.NewPostresRepository(pool)
	userService := service.NewService(userRepo)
	jwtEngine := util.NewJWTEngine(jwtSecret, 15*time.Minute)
	httpHandler := transport.UserHandler{UserService: userService, JwtEngine: jwtEngine}

	router := gin.Default()
	httpHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: router,
	}

	go func() {
		log.Printf("Server Listing on 3001")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shut down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shut down")
	}

	log.Println("Server exiting")
}
