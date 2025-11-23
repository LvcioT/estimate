package entities

import "LvcioT/estimate/shared/providers/entity_id"

type User struct {
	ID    entity_id.EntityId
	Email string
	Name  string
}
