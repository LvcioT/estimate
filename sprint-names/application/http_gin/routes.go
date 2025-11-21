package http_gin

import (
	sprint "pes/sprint-names/domain/sprint/ports/http_gin"

	"github.com/gin-gonic/gin"
)

func DeclareRoutes(router *gin.Engine) error {
	var err error

	router.GET("/", ApplicationInfosHandler)

	sprintsGroup := router.Group("/sprints")
	if err = sprint.GroupSprint(sprintsGroup); err != nil {
		return err
	}

	return nil
}
