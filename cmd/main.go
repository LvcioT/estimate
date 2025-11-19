package main

import (
	"LvcioT/estimate/app/config"
	"LvcioT/estimate/infra/gin_gonic"
	"LvcioT/estimate/infra/gorm_sqlite"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	fmt.Println(cfg)

	initGormSqlite(cfg)
	initGin(cfg)
}

func initGormSqlite(cfg config.Config) {
	err := gorm_sqlite.Init(gorm_sqlite.Options{
		File: cfg.Sqlite.File,
	})
	if err != nil {
		panic(fmt.Errorf("failed to init gorm with sqlite: %w", err))
	}
}

func initGin(cfg config.Config) {
	err := gin_gonic.Init()
	if err != nil {
		panic(fmt.Errorf("failed to init gin: %w", err))
	}

	r := gin_gonic.GetRouter()
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	err = gin_gonic.StartServer(gin_gonic.Options{
		Port:  cfg.Gin.Port,
		Debug: cfg.Gin.Debug,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start gin: %w", err))
	}
}
