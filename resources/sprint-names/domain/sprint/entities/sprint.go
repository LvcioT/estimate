package entities

import (
	"pes/sprint-names/internal"
)

type Sprint struct {
	ID     internal.EntityId
	N      uint
	Name   string
	Letter string
	Period internal.Period
}
