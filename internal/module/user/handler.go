package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"boilerplate/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrDuplicateEmail {
			status = http.StatusConflict
		}
		response.JSONError(c, status, err)
		return
	}

	c.JSON(http.StatusCreated, response.Success(user, "User registered successfully"))
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	loginResp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusUnauthorized
		response.JSONError(c, status, err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(loginResp))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		response.JSONError(c, status, err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(user))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, err)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		} else if err == ErrDuplicateEmail {
			status = http.StatusConflict
		}
		response.JSONError(c, status, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(user, "User updated successfully"))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		response.JSONError(c, status, err)
		return
	}

	c.JSON(http.StatusOK, response.Success(nil, "User deleted successfully"))
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(users))
}
