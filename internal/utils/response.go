package utils

import (
	"net/http"

	errs "github.com/chai-rs/simple-bookstore/internal/error"
	"github.com/gin-gonic/gin"
)

// Response represents the response structure.
type Response struct {
	Success bool `json:"success"`
	Error   any  `json:"error,omitempty"`
	Result  any  `json:"result,omitempty"`
}

// ResponseOk sends a success response.
func ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Result:  data,
	})
}

// ResponseCreated sends a created response.
func ResponseCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Result:  data,
	})
}

// ResponseErrorWithStatus sends an error response with a specific status.
func ResponseErrorWithStatus(c *gin.Context, status int, errorMessage string) {
	c.JSON(status, Response{Success: false, Error: errorMessage})
}

// ResponseError sends an error response.
func ResponseError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errs.AppError:
		ResponseErrorWithStatus(c, e.Code, e.Message)
	default:
		ResponseErrorWithStatus(c, http.StatusInternalServerError, "internal server error")
	}
}
