/* {
  "id": "regra-1",
  "nome": "Ligar luz quando celular conectar",
  "condicao": {
    "tipo": "dispositivo",
    "evento": "conectado",
    "mac_address": "AA:BB:CC:DD:EE:FF"
  },
  "acao": {
    "tipo": "ligar_dispositivo",
    "dispositivo_slug": "luz-sala"
  },
  "ativa": true,
  "created_at": "2025-07-26T20:30:00Z"
} */

// generate entity for rule following json structure

package entity

import (
	"time"

	"github.com/gsousadev/doolar2/internal/domain/valueObject"
)

type Rule struct {
	ID        string
	Name      string
	Condition valueObject.Condition
	Action    valueObject.Action
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
