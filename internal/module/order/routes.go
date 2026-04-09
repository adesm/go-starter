package order

import (
	"github.com/gin-gonic/gin"

	"boilerplate/internal/config"
	"boilerplate/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, cfg *config.Config) {
	order := router.Group("/order")
	{
		order.POST("", middleware.Auth(cfg.JWT.Secret), handler.Create)
		order.GET("/:id", handler.GetByID)
		order.GET("", handler.List)
		// Add PUT and DELETE routes
	}
}
