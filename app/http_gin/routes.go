package http_gin

import (
	"LvcioT/estimate/domain/auth/ports/http_gin"

	"github.com/gin-gonic/gin"
)

func RouteGroups(r *gin.Engine) error {
	r.GET("/", ApplicationInfosHandler)

	authGroup := r.Group("/auth")
	http_gin.RouteGroups(authGroup)

	return nil
}
