package order

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"boilerplate/internal/config"
)

func InitModule(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	service := NewService(repo, cfg)
	handler := NewHandler(service)

	RegisterRoutes(router, handler, cfg)
}
