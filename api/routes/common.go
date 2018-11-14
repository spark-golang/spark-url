package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var APP *gin.Engine

func InitAPP(mode string) {
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	APP = gin.New()
	if mode != "production" {
		APP.Use(gin.Logger())
	}

	APP.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

}
