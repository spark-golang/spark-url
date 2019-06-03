package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spark-golang/spark-url/utils/env"
)

// CORS 支持跨域
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		environment := strings.ToLower(env.Getenv("ENV"))
		if environment != "production" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", "x-csrf-token")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	}
}
