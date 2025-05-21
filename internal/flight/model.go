package flight

import "time"

type ModerationStatus string

const (
	StatusPending  ModerationStatus = "pending"
	StatusApproved ModerationStatus = "approved"
	StatusRejected ModerationStatus = "rejected"
	StatusDeleted  ModerationStatus = "deleted"
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
