package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go-kv/internal/config"
	"go-kv/internal/middleware"
	"go-kv/internal/repository"
	"go-kv/internal/routes"
	"go-kv/internal/services"
)

func main() {
	cfg := config.LoadConfig() // Load configuration
	kvRepo := repository.NewKVRepository()
	kvService := services.NewKVService(kvRepo)

	router := gin.Default()
	//Add middleware - logging and security plugin
	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.SecureMiddleware())

	routes.SetupRoutes(router, kvService)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		//Boot the webserver
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	//Listen for termination signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println("Shutting down server...")
	if err := srv.Close(); err != nil {
		log.Fatalf("Server close failed: %v", err)
	}
}
