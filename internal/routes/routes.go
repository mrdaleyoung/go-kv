package routes

import (
	"github.com/gin-gonic/gin"
	"go-kv/internal/config"
	"go-kv/internal/handlers"
	"go-kv/internal/services"
)

func SetupRoutes(router *gin.Engine, kvService *services.KVService) {
	cfg := config.LoadConfig()
	api := router.Group(cfg.APIPath)
	{
		api.GET("/:key", handlers.HandleGet(kvService))
		api.PUT("/:key", handlers.HandlePut(kvService))
		api.DELETE("/:key", handlers.HandleDelete(kvService))
		api.GET("/", handlers.HandleListKeys(kvService))
	}
}
