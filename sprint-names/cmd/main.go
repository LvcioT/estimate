package main

import (
	"pes/sprint-names/application/config"
	httpGinApp "pes/sprint-names/application/http_gin"
	sqliteGormApp "pes/sprint-names/application/sqlite_gorm"
	httpGinInfra "pes/sprint-names/infrastructure/http_gin"
	sqliteGormInfra "pes/sprint-names/infrastructure/sqlite_gorm"
)

func main() {
	var err error
	c := config.Get()

	if err = initSqliteGorm(c.Application.Environment, c.Database); err != nil {
		panic(err)
	}

	if err = initHttpGin(c.Application.Environment, c.Http); err != nil {
		panic(err)
	}
}

func initSqliteGorm(env config.Environment, config config.Database) error {
	db, err := sqliteGormInfra.InitConnection(config.Filename)
	if err != nil {
		return err
	}

	if err := sqliteGormApp.MigrateRepositories(db); err != nil {
		return err
	}

	return nil
}

func initHttpGin(env config.Environment, config config.Http) error {
	http, err := httpGinInfra.InitEngine(env)
	if err != nil {
		return err
	}

	err = httpGinApp.DeclareRoutes(http)
	if err != nil {
		return err
	}

	err = httpGinInfra.StartEngine(config.Port)
	if err != nil {
		return err
	}

	return nil
}
