package handler

import (
	"book-store/internal/api"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is an alias for api.ErrorResponse kept here for Swagger annotations.
type ErrorResponse = api.ErrorResponse

func errResp(c *gin.Context, code int, msg string) {
	api.ErrResp(c, code, msg)
}
