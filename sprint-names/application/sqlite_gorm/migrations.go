package sqlite_gorm

import (
	"gorm.io/gorm"
	sprint "pes/sprint-names/domain/sprint/repositories/sqlite_gorm"
)

func MigrateRepositories(db *gorm.DB) error {
	sprintRepository, err := sprint.NewSprintSqliteGormRepository(db)
	if err != nil {
		return err
	}

	if err := sprintRepository.AutoMigrate(); err != nil {
		return err
	}

	return nil
}
