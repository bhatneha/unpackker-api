package api

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

func (r *Router) AllRoutes() {

	v1 := (r.Router).Group("/v1")
	{
		v1.POST("/packker", packker)
		v1.POST("/unpackker", unpackker)
		v1.POST("/upload", upload)
		v1.DELETE("/delete", delete)
		v1.GET("/download", download)
	}
}
