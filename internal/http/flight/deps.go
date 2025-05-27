package http_flight

import (
	"github.com/Vovarama1992/go-bagdoor-bot/internal/flight"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
)

type FlightDeps struct {
	FlightService *flight.Service
	UserService   *user.Service
}
