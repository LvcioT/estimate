package repositories

import (
	"LvcioT/estimate/domain/auth/entities"
	"context"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*entities.User, error)
}
