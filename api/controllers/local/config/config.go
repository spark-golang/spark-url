package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spark-golang/spark-url/models/config"
)

func Get(c *gin.Context) {
	name := c.Query("name")

	value := config.Get(name)

	c.JSON(http.StatusOK, gin.H{"status_code": 200, "message": "success", "data": value})
}

func Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status_code": 200, "message": "success", "data": "new"})
}
