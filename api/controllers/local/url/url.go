package url

import (
	"github.com/gin-gonic/gin"
	pb "github.com/spark-golang/spark-url/grpc/pb"
	"github.com/spark-golang/spark-url/models/url"
	"google.golang.org/grpc"
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

func Hello(c *gin.Context) {
	name := c.PostForm("name")
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	req := &pb.HelloRequest{Name: name}
	res, err := client.SayHello(c, req)


	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"status_code": 200, "message": "success", "url": res.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status_code": 200, "message": "success", "url": res.Message})
}
