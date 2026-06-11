package api

import "github.com/gin-gonic/gin"

// ErrorResponse is the standard error body returned on all failures.
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func ErrResp(c *gin.Context, code int, msg string) {
	c.JSON(code, ErrorResponse{Code: code, Error: msg})
}
