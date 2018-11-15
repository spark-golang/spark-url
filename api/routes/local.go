package routes

import (
	"github.com/spark-golang/spark-url/api/controllers/local/url"
	"github.com/spark-golang/spark-url/api/middlewares"
)

func DispatchForLocal() {
	local()
}

func local() {
	local := APP.Group("/local")
	local.Use(middlewares.InternalAvailable())
	{
		local.POST("/url_create", url.Create)
		local.POST("/hello", url.Hello)

	}

	APP.GET("/:str", url.Get)

}
