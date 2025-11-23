package entity_id

type EntityIdProvider interface {
	Generate() (EntityId, error)
	FromString(string) (EntityId, error)
}
