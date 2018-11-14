package middlewares

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spark-golang/spark-url/conf/constant"
	"github.com/spark-golang/spark-url/utils"
)

// CSRF 安全中间件
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(constant.CSRF)
		var token string
		if err != nil {
			token = utils.RandomString(32)
		} else {
			token = cookie
		}
		switch c.Request.Method {
		case "GET", "OPTIONS", "HEAD":
		default:
			uid := c.MustGet("uid").(uint32)
			if uid <= 0 {
				c.JSON(http.StatusForbidden, gin.H{"status_code": 4206, "message": "未登录"})
				c.Abort()
				return
			}
			clientToken := c.Request.Header.Get("X-CSRF-TOKEN")
			if clientToken == "" {
				c.JSON(http.StatusForbidden, gin.H{"status_code": 1001, "message": "forbiden"})
				c.Abort()
				return
			}
			if !validateCSRFToken(token, clientToken) {
				c.JSON(http.StatusForbidden, gin.H{"status_code": 1001, "message": "forbiden"})
				c.Abort()
				return
			}
		}
		domain := c.Request.Header.Get("X-Server-Domain")
		c.SetCookie(constant.CSRF, token, 86400, "/", "."+domain, false, false)
	}
}

func validateCSRFToken(token, clientToken string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(clientToken)) == 1
}
