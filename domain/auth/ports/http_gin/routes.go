package http_gin

import (
	userGinGorm "LvcioT/estimate/domain/auth/repositories/sqlite_gorm"
	"LvcioT/estimate/infra/sqlite_gorm"

	"github.com/gin-gonic/gin"
)

func RouteGroups(rg *gin.RouterGroup) {
	db := sqlite_gorm.GetConnection()

	// users
	ur := userGinGorm.NewUserSqliteGormRepository(db)
	uh := NewUserHandler(ur)
	{
		ug := rg.Group("users")
		ug.GET("/", uh.Index)
	}
}
