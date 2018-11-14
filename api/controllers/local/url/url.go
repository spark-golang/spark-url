package url

import (
	"github.com/gin-gonic/gin"
	"github.com/spark-golang/spark-url/models/url"
	"log"
	"net/http"
)

func Create(c *gin.Context) {

	purl := c.PostForm("url")
	value := url.Create(purl)

	c.JSON(http.StatusOK, gin.H{"status_code": 200, "message": "success", "url": value})
}

func Get(c *gin.Context) {
	str := c.Param("str")
	log.Println(str)
	value := url.Get(str)
	c.Redirect(http.StatusMovedPermanently, value)

}
