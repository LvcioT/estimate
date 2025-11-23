package http_gin

import (
	"LvcioT/estimate/domain/auth/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		Name string
	}
)

type UserHandler struct {
	repository repositories.UserRepository
}

func NewUserHandler(repository repositories.UserRepository) UserHandler {
	return UserHandler{repository: repository}
}

func (uh *UserHandler) Index(ctx *gin.Context) {
	userEntities, err := uh.repository.GetAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	users := make([]User, len(userEntities))
	for i, userEntity := range userEntities {
		users[i] = User{Name: userEntity.Name}
	}

	ctx.JSON(http.StatusOK, users)
}
