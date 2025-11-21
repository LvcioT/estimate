package internal

import (
	"github.com/google/uuid"
)

type EntityId string

func NewEntityId() EntityId {
	eid, e := uuid.NewV7()
	if e != nil {
		panic("failed to generate ID: " + e.Error())
	}

	return EntityId(eid.String())
}

func NewEntityIdFromString(s string) EntityId {
	return EntityId(s)
}

func (eid EntityId) String() string {
	return string(eid)
}
