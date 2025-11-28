package entity

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Entity struct {
	ID uuid.UUID `json:"id"`
}

type IEntity interface {
	GetID() uuid.UUID
	ToJSON() ([]byte, error)
	ToJSONString() (string, error)
}

func NewEntity() *Entity {

	id, err := uuid.NewV6()

	if err != nil {
		panic(err)
	}

	return &Entity{
		ID: id,
	}
}

// GetID returns the entity's ID
func (e *Entity) GetID() uuid.UUID {
	return e.ID
}

// Método para retornar JSON como bytes
func (e *Entity) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// Opcional: método para string JSON
func (e *Entity) ToJSONString() (string, error) {
	jsonBytes, err := e.ToJSON()
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
