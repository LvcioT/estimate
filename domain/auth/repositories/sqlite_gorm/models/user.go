package models

import (
	"LvcioT/estimate/shared/providers/entity_id"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	EID   entity_id.EntityId `gorm:"index,unique,required"`
	Email string             `gorm:"index,unique,required"`
	Name  string             `gorm:"size:255,required"`
}
