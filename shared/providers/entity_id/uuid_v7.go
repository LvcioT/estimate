package entity_id

import (
	"fmt"

	"github.com/google/uuid"
)

type EntityIdV7Provider struct{}

func (ep EntityIdV7Provider) Generate() (EntityId, error) {
	eid, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate EID: %w", err)
	}

	return EntityId(eid.String()), nil
}

func (ep EntityIdV7Provider) FromString(s string) (EntityId, error) {
	return EntityId(s), nil
}
