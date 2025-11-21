package config

import "time"

type Environment string

const (
	Local Environment = "local"
	Test  Environment = "test"
	Prod  Environment = "prod"
)

type Config struct {
	Application Application
	Http        Http
	Database    Database
}

type Application struct {
	Name        string
	Version     string
	Environment Environment
}

type Http struct {
	Port    string
	Timeout time.Time
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	Filename string
}
