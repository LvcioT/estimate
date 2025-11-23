package main

import (
	"LvcioT/estimate/app/config"
	httpGinApp "LvcioT/estimate/app/http_gin"
	"LvcioT/estimate/app/sqlite_gorm"
	httpGinInfra "LvcioT/estimate/infra/http_gin"
	sqliteGormInfra "LvcioT/estimate/infra/sqlite_gorm"
	"fmt"
)

func main() {
	cfg := config.GetConfig()

	fmt.Println(cfg)

	initGormSqlite(cfg)
	initGin(cfg)
}

func initGormSqlite(cfg config.Config) {
	if err := sqliteGormInfra.Init(sqliteGormInfra.Options{
		File: cfg.Sqlite.File,
	}); err != nil {
		panic(fmt.Errorf("failed to init gorm with sqlite: %w", err))
	}

	db := sqliteGormInfra.GetConnection()

	if err := sqlite_gorm.AutoMigrate(db); err != nil {
		panic(fmt.Errorf("failed to migrate database: %w", err))
	}
}

func initGin(cfg config.Config) {
	if err := httpGinInfra.Init(); err != nil {
		panic(fmt.Errorf("failed to init gin: %w", err))
	}

	r := httpGinInfra.GetRouter()

	if err := httpGinApp.RouteGroups(r); err != nil {
		panic(fmt.Errorf("failed to declare routes: %w", err))
	}

	fmt.Println(r.Routes())

	if err := httpGinInfra.StartServer(httpGinInfra.Options{
		Port:  cfg.Gin.Port,
		Debug: cfg.Gin.Debug,
	}); err != nil {
		panic(fmt.Errorf("failed to start gin: %w", err))
	}
}
