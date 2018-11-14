package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResMsg(c *gin.Context, code int, msg string) error {
	c.JSON(http.StatusOK, gin.H{
		"status_code": code,
		"message":     msg,
		"data":        "",
	})

	return nil
}
