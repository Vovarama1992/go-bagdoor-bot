package order

import "time"

type OrderType string

const (
	OrderTypePersonal OrderType = "personal"
	OrderTypeStore    OrderType = "store"
)

type ModerationStatus string

const (
	StatusPending  ModerationStatus = "PENDING"
	StatusApproved ModerationStatus = "APPROVED"
	StatusRejected ModerationStatus = "REJECTED"
)

type Order struct {
	ID                int
	OrderNumber       string
	UserID            int
	PublisherUsername string
	PublisherTgID     int64
	PublishedAt       time.Time

	OriginCity      string
	DestinationCity string
	StartDate       time.Time
	EndDate         time.Time
	Title           string
	Description     string
	StoreLink       *string // опциональное поле

	Reward  float64
	Deposit *float64
	Cost    *float64

	MediaURLs []string
	Type      OrderType
	Status    ModerationStatus
}
