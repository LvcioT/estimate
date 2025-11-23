package sqlite_gorm

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var connection *gorm.DB

func Init(o Options) error {
	db, err := gorm.Open(sqlite.Open(o.File), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("cannot open sqlite DB: %w", err)
	}

	connection = db
	return nil
}

func GetConnection() *gorm.DB {
	return connection
}
