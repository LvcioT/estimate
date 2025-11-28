package config

import (
	"github.com/joho/godotenv"
	"os"
)

var config Config

func init() {
	if e := godotenv.Load(); e != nil {
		panic("fail to load .env file: " + e.Error())
	}

	config = Config{
		Application: Application{
			Name:        AppName,
			Version:     AppVersion,
			Environment: Local,
		},
		Http: Http{
			Port: os.Getenv("HTTP_PORT"),
		},
		Database: Database{
			Filename: os.Getenv("DB_FILENAME"),
		},
	}
}

func Get() Config {
	return config
}
