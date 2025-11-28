// Criar uma estrutura de entidade para representar um comodo da casa em ingles

/*
{
  "nome": "Sala de Estar",
  "slug": "sala-de-estar",
  "pontos": {
    "a": { "lat": -23.1234, "lng": -46.5678 },
    "b": { "lat": -23.1235, "lng": -46.5679 },
    "c": { "lat": -23.1236, "lng": -46.5680 },
    "d": { "lat": -23.1237, "lng": -46.5681 }
  },
  "dispositivos_mac": [ "00:1A:2B:3C:4D:5E", "00:1A:2B:3C:4D:5F" ],  // Índice secundário para busca rápida
  "created_at": "2025-07-26T20:05:00Z",
  "updated_at": "2025-07-26T21:00:00Z"
}


*/

package entity

import (
	"github.com/gsousadev/doolar2/internal/shared/domain/entity"
	"github.com/gsousadev/doolar2/internal/shared/domain/value_object"
)

type Room struct {
	*entity.Entity
	Name string
	Slug string
}

func NewRoom(name string) *Room {

	slug, err := value_object.NewSlugFromString(name)

	if err != nil {
		panic(err)
	}

	return &Room{
		Entity: entity.NewEntity(),
		Name:   name,
		Slug:   slug.Value(),
	}
}
