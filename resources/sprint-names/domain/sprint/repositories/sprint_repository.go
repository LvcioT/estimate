package repositories

import (
	"context"
	"fmt"
	"pes/sprint-names/domain/sprint/entities"
	"pes/sprint-names/domain/sprint/repositories/sqlite_gorm"
	sqliteGormInfra "pes/sprint-names/infrastructure/sqlite_gorm"
	"pes/sprint-names/internal"
)

type SprintRepository interface {
	GetAll(ctx context.Context) ([]*entities.Sprint, error)
	GetByID(ctx context.Context, id internal.EntityId) (entities.Sprint, error)
	Create(ctx context.Context, sprint entities.Sprint) (entities.Sprint, error)
	Update(ctx context.Context, sprint entities.Sprint) (entities.Sprint, error)
	Delete(ctx context.Context, id internal.EntityId) error
}

func NewSprintRepository() (SprintRepository, error) {
	db, err := sqliteGormInfra.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	r, err := sqlite_gorm.NewSprintSqliteGormRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate sprint repository: %v", err)
	}

	return r, nil
}
