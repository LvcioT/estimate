package dtos

import (
	"pes/sprint-names/domain/sprint/entities"
	"pes/sprint-names/internal"
)

type Sprint struct {
	ID      internal.EntityId `json:"id"`
	N       uint              `json:"n"`
	Name    string            `json:"name"`
	Letter  string            `json:"letter"`
	StartAt string            `json:"start_at"`
	EndAt   string            `json:"end_at"`
}

func NewSprintFromEntity(se entities.Sprint) Sprint {
	return Sprint{
		ID:      se.ID,
		N:       se.N,
		Name:    se.Name,
		Letter:  se.Letter,
		StartAt: se.Period.StartAt.String(),
		EndAt:   se.Period.EndAt.String(),
	}
}
