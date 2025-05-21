package http

import (
	"github.com/Vovarama1992/go-bagdoor-bot/internal/flight"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/storage"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
)

type FlightDeps struct {
	FlightService *flight.Service
	UserService   *user.Service
}

type OrderDeps struct {
	OrderService *order.Service
	UserService  *user.Service
	Uploader     *storage.S3Uploader
}
