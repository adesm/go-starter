#!/bin/bash

# Configuration
MODULE_BASE_DIR="internal/module"
NAME=$1

if [ -z "$NAME" ]; then
    echo "Error: Module name is required."
    echo "Usage: ./scripts/genmodule.sh <module_name>"
    exit 1
fi

# Convert module name to lowercase (package name) and capitalized (struct/model name)
PKG_NAME=$(echo "$NAME" | tr '[:upper:]' '[:lower:]')
STRUCT_NAME=$(echo "$PKG_NAME" | awk '{print toupper(substr($0,1,1)) substr($0,2)}')
MODULE_DIR="$MODULE_BASE_DIR/$PKG_NAME"

if [ -d "$MODULE_DIR" ]; then
    echo "Error: Module '$PKG_NAME' already exists at $MODULE_DIR."
    exit 1
fi

echo "Generating module '$PKG_NAME' (Struct: $STRUCT_NAME) in $MODULE_DIR..."

mkdir -p "$MODULE_DIR"

# 1. dto.go
cat > "$MODULE_DIR/dto.go" <<EOF
package $PKG_NAME

type Create${STRUCT_NAME}Request struct {
	// Add fields with binding:"required" for validation
    // Example: Name string \`json:"name" binding:"required"\`
}

type Update${STRUCT_NAME}Request struct {
	// Add fields with binding:"required" for validation
}

type ${STRUCT_NAME}Response struct {
	ID uint \`json:"id"\`
	// Add fields
}
EOF

# 2. error.go
cat > "$MODULE_DIR/error.go" <<EOF
package $PKG_NAME

import "errors"

var (
	Err${STRUCT_NAME}NotFound = errors.New("${PKG_NAME} not found")
)
EOF

# 3. model.go
cat > "$MODULE_DIR/model.go" <<EOF
package $PKG_NAME

import (
	"time"

	"gorm.io/gorm"
)

type ${STRUCT_NAME} struct {
	ID        uint           \`gorm:"primaryKey" json:"id"\`
	CreatedAt time.Time      \`json:"created_at"\`
	UpdatedAt time.Time      \`json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
    // Add fields
}
EOF

# 4. interface.go
cat > "$MODULE_DIR/interface.go" <<EOF
package $PKG_NAME

import "context"

type Repository interface {
	Create(ctx context.Context, item *${STRUCT_NAME}) error
	FindByID(ctx context.Context, id uint) (*${STRUCT_NAME}, error)
	Update(ctx context.Context, item *${STRUCT_NAME}) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*${STRUCT_NAME}, error)
}

type Service interface {
	Create(ctx context.Context, req *Create${STRUCT_NAME}Request) (*${STRUCT_NAME}, error)
	GetByID(ctx context.Context, id uint) (*${STRUCT_NAME}, error)
	Update(ctx context.Context, id uint, req *Update${STRUCT_NAME}Request) (*${STRUCT_NAME}, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*${STRUCT_NAME}, error)
}
EOF

# 5. repository.go
cat > "$MODULE_DIR/repository.go" <<EOF
package $PKG_NAME

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, item *${STRUCT_NAME}) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*${STRUCT_NAME}, error) {
	var item ${STRUCT_NAME}
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, Err${STRUCT_NAME}NotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) Update(ctx context.Context, item *${STRUCT_NAME}) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&${STRUCT_NAME}{}, id).Error
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]*${STRUCT_NAME}, error) {
	var items []*${STRUCT_NAME}
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&items).Error
	return items, err
}
EOF

# 6. service.go
cat > "$MODULE_DIR/service.go" <<EOF
package $PKG_NAME

import (
	"context"

	"boilerplate/internal/config"
)

type service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{repo: repo, cfg: cfg}
}

func (s *service) Create(ctx context.Context, req *Create${STRUCT_NAME}Request) (*${STRUCT_NAME}, error) {
	item := &${STRUCT_NAME}{
		// Map fields from req
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*${STRUCT_NAME}, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uint, req *Update${STRUCT_NAME}Request) (*${STRUCT_NAME}, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, page, pageSize int) ([]*${STRUCT_NAME}, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, pageSize, offset)
}
EOF

# 7. handler.go
cat > "$MODULE_DIR/handler.go" <<EOF
package $PKG_NAME

import (
	"net/http"
	"strconv"

	"boilerplate/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req Create${STRUCT_NAME}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	item, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessWithData(item))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.JSONError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(item))
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	items, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(items))
}

// Add Update and Delete methods...
EOF

# 8. routes.go
cat > "$MODULE_DIR/routes.go" <<EOF
package $PKG_NAME

import (
	"github.com/gin-gonic/gin"

	"boilerplate/internal/config"
	"boilerplate/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, cfg *config.Config) {
	${PKG_NAME} := router.Group("/${PKG_NAME}s")
	{
		${PKG_NAME}.POST("", middleware.AuthMiddleware(cfg), handler.Create)
		${PKG_NAME}.GET("/:id", handler.GetByID)
		${PKG_NAME}.GET("", handler.List)
		// Add PUT and DELETE routes
	}
}
EOF

# 9. module.go
cat > "$MODULE_DIR/module.go" <<EOF
package $PKG_NAME

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
EOF

echo "Module '$PKG_NAME' successfully generated!"
echo "Don't forget to:"
echo "1. Add 'err = db.AutoMigrate(&$PKG_NAME.${STRUCT_NAME}{})' to cmd/migrate/main.go"
echo "2. Add '$PKG_NAME.InitModule(api, db, cfg)' to cmd/api/main.go setupRouter"
