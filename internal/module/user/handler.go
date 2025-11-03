package user

import (
	"log"
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
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrDuplicateEmail {
			status = http.StatusConflict
		}
		c.JSON(status, response.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.Success(user, "User registered successfully"))
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	loginResp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusUnauthorized
		if err == ErrInvalidCredentials {
			status = http.StatusUnauthorized
		}
		c.JSON(status, response.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(loginResp))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	authUserID, _ := c.Get("user_id")
	log.Printf("User %v is accessing user %d profile", authUserID, id)

	user, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, response.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(user))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(ErrNotFound))
		return
	}

	// Optional: Check if user is updating their own profile
	// Uncomment if you want users to only update their own profile
	// if uint(id) != authUserID.(uint) {
	// 	c.JSON(http.StatusForbidden, gin.H{
	// 		"success": false,
	// 		"error":   "You can only update your own profile",
	// 	})
	// 	return
	// }

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
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
		c.JSON(status, response.Error(err))
		return
	}

	log.Printf("User %v updated user %d", authUserID, id)

	c.JSON(http.StatusOK, response.Success(user, "User updated successfully"))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}

	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(ErrNotFound))
		return
	}

	// Optional: Check if user is deleting their own account
	// Uncomment if you want users to only delete their own account
	// if uint(id) != authUserID.(uint) {
	// 	c.JSON(http.StatusForbidden, gin.H{
	// 		"success": false,
	// 		"error":   "You can only delete your own account",
	// 	})
	// 	return
	// }

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, response.Error(err))
		return
	}

	log.Printf("User %v deleted user %d", authUserID, id)

	c.JSON(http.StatusOK, response.Success(nil, "User deleted successfully"))
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Get authenticated user from context
	authUserID, _ := c.Get("user_id")
	log.Printf("User %v is listing users (page: %d, size: %d)", authUserID, page, pageSize)

	users, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithData(users))
}
