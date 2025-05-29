package flight

import "time"

type ModerationStatus string

const (
	StatusPending  ModerationStatus = "ожидает модерации"
	StatusApproved ModerationStatus = "отмодерирован и опубликован"
	StatusRejected ModerationStatus = "отклонён"
	StatusDeleted  ModerationStatus = "удалён"
)

type Flight struct {
	ID                int
	FlightNumber      string
	PublisherUsername string
	PublisherTgID     int64
	FlightDate        time.Time
	PublishedAt       time.Time
	Description       string
	Origin            string
	Destination       string
	Status            ModerationStatus
	MapURL            *string // PDF карта, может быть nil
}
