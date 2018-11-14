package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spark-golang/spark-url/utils/ip_tools"
)

func InternalAvailable() gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteAddr := ip_tools.InternalIP(c)
		available := ip_tools.IsInternalIP(remoteAddr)

		if !available {
			c.JSON(403, gin.H{"error": "403", "data": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
