package sqlite_gorm

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitConnection(filename string) (*gorm.DB, error) {
	var err error

	db, err = gorm.Open(sqlite.Open(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

func GetConnection() (*gorm.DB, error) {
	var err error

	if db == nil {
		err = errors.New("database connection not initialized")
	}

	return db, err
}
