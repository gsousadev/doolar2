package entity

/* {
  "id": "evt-abc123",
  "tipo": "device_connected",           // Exemplo: device_connected, task_completed, reminder_triggered
  "timestamp": "2025-07-26T20:50:00Z",
  "dados": {
    "dispositivo_mac": "00:1A:2B:3C:4D:5E",
    "ip": "192.168.0.12"
  },
  "relacionado_slug": "smartphone-guilherme"
} */

// generate entity for event following json structure

import (
	"time"
)

type Event struct {
	ID          string
	Type        string
	Timestamp   time.Time
	Data        map[string]interface{} // Can hold various data types, e.g., device MAC,
	RelatedSlug string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
