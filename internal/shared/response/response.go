package response

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

type MetaData struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func Success(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func SuccessWithData(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// JSONResponse is a helper to send consistent JSON responses
func JSONResponse(c *gin.Context, statusCode int, success bool, data interface{}, message string, err string, details interface{}) {
	c.JSON(statusCode, Response{
		Success: success,
		Data:    data,
		Message: message,
		Error:   err,
		Details: details,
	})
}

// JSONError is a smart error handler that masks 500 errors and logs them
func JSONError(c *gin.Context, statusCode int, err error) {
	if statusCode >= http.StatusInternalServerError {
		slog.Error("INTERNAL SERVER ERROR",
			slog.String("path", c.Request.URL.Path),
			slog.String("method", c.Request.Method),
			slog.Any("error", err.Error()),
		)

		JSONResponse(c, statusCode, false, nil, "", "Internal Server Error", nil)
		return
	}

	JSONResponse(c, statusCode, false, nil, "", err.Error(), nil)
}

// ValidationError handles gin validation errors and returns structured field details
func ValidationError(c *gin.Context, err error) {
	details := make(map[string]string)
	
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, f := range verrs {
			field := f.Field()
			// You can add more complex tag mapping here
			switch f.Tag() {
			case "required":
				details[field] = fmt.Sprintf("%s is required", field)
			case "email":
				details[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				details[field] = fmt.Sprintf("%s must be at least %s characters", field, f.Param())
			default:
				details[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	if len(details) == 0 {
		JSONError(c, http.StatusBadRequest, err)
		return
	}

	JSONResponse(c, http.StatusBadRequest, false, nil, "", "Validation failed", details)
}
