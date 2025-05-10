package utils

import (
	"net/http"

	errs "github.com/0xanonydxck/simple-bookstore/internal/error"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool `json:"success"`
	Error   any  `json:"error,omitempty"`
	Result  any  `json:"result,omitempty"`
}

func ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Result:  data,
	})
}

func ResponseCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Result:  data,
	})
}

func ResponseErrorWithStatus(c *gin.Context, status int, errorMessage string) {
	c.JSON(status, Response{Success: false, Error: errorMessage})
}

func ResponseError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errs.AppError:
		ResponseErrorWithStatus(c, e.Code, e.Message)
	default:
		ResponseErrorWithStatus(c, http.StatusInternalServerError, "internal server error")
	}
}
