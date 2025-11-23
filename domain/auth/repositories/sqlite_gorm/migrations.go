package sqlite_gorm

import (
	"fmt"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	user := NewUserSqliteGormRepository(db)

	err := user.AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate user repository: %w", err)
	}

	return nil
}
