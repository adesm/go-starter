package order

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
	var req CreateOrderRequest
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
