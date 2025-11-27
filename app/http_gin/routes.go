package http_gin

import (
	"LvcioT/estimate/domain/auth/ports/http_gin"

	"github.com/gin-gonic/gin"
)

func RouteGroups(r *gin.Engine) error {
	r.LoadHTMLGlob("templates/**/*")
	r.GET("/", infosHandler)

	authGroup := r.Group("/auth")
	http_gin.RouteGroups(authGroup)

	examplePagesGroup := r.Group("/example")
	{
		examplePagesGroup.GET("/page1", page1Handler)
		examplePagesGroup.GET("/page2", page2Handler)
	}

	return nil
}
