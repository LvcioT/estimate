package main

import (
	"LvcioT/estimate/app/config"
	"LvcioT/estimate/infra/gin_gonic"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	fmt.Println(cfg)

	r := gin_gonic.GetRouter()
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
	err := gin_gonic.StartServer(gin_gonic.Options{
		Port:  cfg.Gin.Port,
		Debug: cfg.Gin.Debug,
	})
	if err != nil {
		panic(err)
	}
}
