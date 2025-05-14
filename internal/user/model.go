package user

import "time"

type User struct {
	ID           int
	TgID         int64
	TgUsername   string
	FirstName    string
	LastName     string
	PhoneNumber  string
	RegisteredAt time.Time
}
