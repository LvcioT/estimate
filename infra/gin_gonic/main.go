package gin_gonic

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() error {
	router = gin.Default()
	return nil
}

func GetRouter() *gin.Engine {
	return router
}

func StartServer(o Options) error {
	if o.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	addr := fmt.Sprintf(":%d", o.Port)

	err := router.Run(addr)
	return fmt.Errorf("cannot start server: %w", err)
}
