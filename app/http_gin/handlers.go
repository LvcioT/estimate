package http_gin

import (
	"LvcioT/estimate/app/config"

	"github.com/gin-gonic/gin"
)

func ApplicationInfosHandler(c *gin.Context) {
	conf := config.GetConfig().App

	c.JSON(200, gin.H{
		"name":      conf.Name,
		"full_name": conf.FullName,
		"version":   conf.Version,
	})
}
