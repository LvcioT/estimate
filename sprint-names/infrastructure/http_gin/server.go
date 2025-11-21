package http_gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pes/sprint-names/application/config"
)

var engine *gin.Engine

func InitEngine(environment config.Environment) (*gin.Engine, error) {
	var mode string
	switch environment {
	case config.Local:
		mode = gin.DebugMode
	case config.Test:
		mode = gin.TestMode
	case config.Prod:
		mode = gin.ReleaseMode
	default:
		return nil, fmt.Errorf("unknown environment: %s", environment)
	}

	gin.SetMode(mode)
	engine = gin.New()

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return engine, nil
}

func StartEngine(port string) error {
	if e := engine.Run(":" + port); e != nil {
		return fmt.Errorf("error starting server: %w", e)
	}

	return nil
}

func GetEngine() (*gin.Engine, error) {
	if engine == nil {
		return nil, fmt.Errorf("HTTP engine not initialized")
	}

	return engine, nil
}
