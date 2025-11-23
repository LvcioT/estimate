package sqlite_gorm

import (
	"LvcioT/estimate/domain/auth/entities"
	"LvcioT/estimate/domain/auth/repositories/sqlite_gorm/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserSqliteGormRepository struct {
	db *gorm.DB
}

func NewUserSqliteGormRepository(db *gorm.DB) UserSqliteGormRepository {
	return UserSqliteGormRepository{db: db}
}

func (ur UserSqliteGormRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	var modelAll []*models.User

	result := ur.db.Find(&modelAll)
	if result.Error != nil {
		return []*entities.User{}, fmt.Errorf("failed to retrieve all users: %s", result.Error.Error())
	}

	entityAll := make([]*entities.User, len(modelAll))
	for i, m := range modelAll {
		e := modelToEntity(*m)
		entityAll[i] = &e
	}

	return entityAll, nil
}

func (ur UserSqliteGormRepository) AutoMigrate() error {
	if e := ur.db.AutoMigrate(&models.User{}); e != nil {
		return fmt.Errorf("failed to migrate 'user': %s", e.Error())
	}

	return nil
}

func modelToEntity(model models.User) entities.User {
	return entities.User{
		ID:    model.EID,
		Email: model.Email,
		Name:  model.Name,
	}
}

func entityToModel(entity entities.User) models.User {
	return models.User{
		EID:   entity.ID,
		Email: entity.Email,
		Name:  entity.Name,
	}
}
