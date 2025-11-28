package models

import (
	"gorm.io/gorm"
	"pes/sprint-names/internal"
	"time"
)

type Sprint struct {
	gorm.Model
	EID     internal.EntityId `gorm:"index,unique"`
	N       uint              `gorm:"unique,required"`
	Name    string            `gorm:"unique,required"`
	Letter  string            `gorm:"size:1,required"`
	StartAt time.Time         `gorm:"type:date,required"`
	EndAt   time.Time         `gorm:"type:date,required"`
}
