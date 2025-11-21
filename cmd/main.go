package main

import (
	"LvcioT/estimate/app/config"
	httpGinApp "LvcioT/estimate/app/http_gin"
	"LvcioT/estimate/infra/gorm_sqlite"
	httpGinInfra "LvcioT/estimate/infra/http_gin"
	"fmt"
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
	err := httpGinInfra.Init()
	if err != nil {
		panic(fmt.Errorf("failed to init gin: %w", err))
	}

	r := httpGinInfra.GetRouter()

	err = httpGinApp.DeclareRoutes(r)
	if err != nil {
		panic(fmt.Errorf("failed to declare routes: %w", err))
	}

	err = httpGinInfra.StartServer(httpGinInfra.Options{
		Port:  cfg.Gin.Port,
		Debug: cfg.Gin.Debug,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start gin: %w", err))
	}
}
