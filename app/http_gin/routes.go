package http_gin

import "github.com/gin-gonic/gin"

func DeclareRoutes(r *gin.Engine) error {
	r.GET("/", ApplicationInfosHandler)

	return nil
}
