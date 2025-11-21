package sqlite_gorm

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pes/sprint-names/domain/sprint/entities"
	"pes/sprint-names/domain/sprint/repositories/sqlite_gorm/models"
	"pes/sprint-names/internal"
)

type SprintSqliteGormRepository struct {
	db *gorm.DB
}

func NewSprintSqliteGormRepository(db *gorm.DB) (SprintSqliteGormRepository, error) {
	return SprintSqliteGormRepository{
		db: db,
	}, nil
}

func (r SprintSqliteGormRepository) Create(ctx context.Context, sprint entities.Sprint) (entities.Sprint, error) {
	model := entityToModel(sprint)

	result := r.db.Create(&model)
	if result.Error != nil {
		return entities.Sprint{}, errors.New("failed to create sprint: " + result.Error.Error())
	}

	return modelToEntity(model), nil
}

func (r SprintSqliteGormRepository) GetAll(ctx context.Context) ([]*entities.Sprint, error) {
	var modelAll []*models.Sprint

	result := r.db.Find(&modelAll)
	if result.Error != nil {
		return []*entities.Sprint{}, errors.New("failed to retrieve all sprints: " + result.Error.Error())
	}

	sprints := make([]*entities.Sprint, len(modelAll))
	for i, m := range modelAll {
		e := modelToEntity(*m)
		sprints[i] = &e
	}

	return sprints, nil
}

func (r SprintSqliteGormRepository) GetByID(ctx context.Context, id internal.EntityId) (entities.Sprint, error) {
	var model models.Sprint

	if result := r.db.Where("eid='?'", id).First(&model); result.Error != nil {
		return entities.Sprint{}, fmt.Errorf("failed to retrieve 'sprint' ID '%s': %s", id, result.Error.Error())
	}

	return modelToEntity(model), nil
}

func (r SprintSqliteGormRepository) Update(ctx context.Context, sprint entities.Sprint) (entities.Sprint, error) {
	var model models.Sprint
	if result := r.db.Model(&model).Where("eid=?", sprint.ID).Updates(entityToModel(sprint)); result.Error != nil {
		return entities.Sprint{}, fmt.Errorf("failed to update 'sprint' '%s': %s", sprint.ID, result.Error.Error())
	}

	return modelToEntity(model), nil
}

func (r SprintSqliteGormRepository) Delete(ctx context.Context, id internal.EntityId) error {
	var model models.Sprint
	if result := r.db.Where("eid=?", id).Delete(&model); result.Error != nil {
		return fmt.Errorf("failed to delete 'sprint' '%s': %s", id, result.Error.Error())
	}

	return nil
}

func entityToModel(entity entities.Sprint) models.Sprint {
	return models.Sprint{
		EID:     entity.ID,
		N:       entity.N,
		Name:    entity.Name,
		Letter:  entity.Letter,
		StartAt: entity.Period.StartAt,
		EndAt:   entity.Period.EndAt,
	}
}

func modelToEntity(model models.Sprint) entities.Sprint {
	return entities.Sprint{
		ID:     model.EID,
		N:      model.N,
		Name:   model.Name,
		Letter: model.Letter,
		Period: internal.Period{
			StartAt: model.StartAt,
			EndAt:   model.EndAt,
		},
	}
}

func (r SprintSqliteGormRepository) AutoMigrate() error {
	if e := r.db.AutoMigrate(&models.Sprint{}); e != nil {
		return errors.New("failed to migrate 'sprint': " + e.Error())
	}

	return nil
}
