package user

import (
	"github.com/gin-gonic/gin"

	"boilerplate/internal/config"
	"boilerplate/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, cfg *config.Config) {
	users := router.Group("/users")
	{
		users.POST("/register", handler.Register)
		users.POST("/login", handler.Login)

		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.GET("/:id", handler.GetByID)
			protected.PUT("/:id", handler.Update)
			protected.DELETE("/:id", handler.Delete)
			protected.GET("", handler.List)
		}
	}
}
