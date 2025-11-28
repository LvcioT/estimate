package http_gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pes/sprint-names/domain/sprint/repositories"
)

func GroupSprint(g *gin.RouterGroup) error {
	sr, e := repositories.NewSprintRepository()
	if e != nil {
		return fmt.Errorf("filed to instantiate sprint repository: %w", e)
	}

	sh := NewSprintHandler(sr)

	g.GET("", sh.Index)

	return nil
}
