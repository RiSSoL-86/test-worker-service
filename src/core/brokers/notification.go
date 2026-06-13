package brokers

import (
	"encoding/json"
	"time"
)

type Notification struct {
	NotificationID string          `json:"notification_id"`
	Type           string          `json:"type"`
	EntityID       string          `json:"entity_id"`
	Payload        json.RawMessage `json:"payload"`
	CreatedAt      time.Time       `json:"created_at"`
}
