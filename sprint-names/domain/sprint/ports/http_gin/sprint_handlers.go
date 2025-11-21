package http_gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pes/sprint-names/domain/sprint/ports/http_gin/dtos"
	"pes/sprint-names/domain/sprint/repositories"
)

type SprintHandler struct {
	r repositories.SprintRepository
}

func NewSprintHandler(repository repositories.SprintRepository) SprintHandler {
	return SprintHandler{r: repository}
}

func (sh SprintHandler) Index(c *gin.Context) {
	sprints, err := sh.r.GetAll(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	sprintDTOs := make([]dtos.Sprint, len(sprints))

	for i, sprint := range sprints {
		sprintDTOs[i] = dtos.NewSprintFromEntity(*sprint)
	}

	c.JSON(http.StatusOK, sprintDTOs)
}
