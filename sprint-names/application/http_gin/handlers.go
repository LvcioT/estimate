package http_gin

import (
	"github.com/gin-gonic/gin"
	"pes/sprint-names/application/config"
)

func ApplicationInfosHandler(c *gin.Context) {
	conf := config.Get().Application

	c.JSON(200, gin.H{
		"name":        conf.Name,
		"version":     conf.Version,
		"environment": conf.Environment,
	})
}
