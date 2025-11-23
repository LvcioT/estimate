package sqlite_gorm

import (
	auth "LvcioT/estimate/domain/auth/repositories/sqlite_gorm"
	"fmt"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := auth.AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
