package flight

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateFlight(ctx context.Context, username string, tgID int64, desc, origin, destination string, flightDate time.Time) (*Flight, error) {
	now := time.Now()
	num := generateFlightNumber(now, tgID)
	flight := &Flight{
		FlightNumber:      num,
		PublisherUsername: username,
		PublisherTgID:     tgID,
		PublishedAt:       now,
		Description:       desc,
		Origin:            origin,
		Destination:       destination,
		FlightDate:        flightDate,
		Status:            StatusPending,
	}
	if err := s.repo.Create(flight); err != nil {
		return nil, err
	}
	return flight, nil
}

func (s *Service) GetFlightByID(ctx context.Context, id int) (*Flight, error) {
	return s.repo.GetByID(id)
}

func (s *Service) SetMapURL(ctx context.Context, id int, url string) error {
	return s.repo.UpdateMapURL(id, url)
}

func (s *Service) SetStatus(ctx context.Context, id int, status ModerationStatus) error {
	return s.repo.UpdateStatus(id, status)
}

func generateFlightNumber(now time.Time, tgID int64) string {
	return fmt.Sprintf("Рейс #%04d-%04d", time.Now().Unix()%10000, tgID%10000)
}

func (s *Service) GetAllFlights(ctx context.Context) ([]*Flight, error) {
	return s.repo.GetAllFlights(ctx)
}

func (s *Service) GetFlightsByStatus(ctx context.Context, status ModerationStatus) ([]*Flight, error) {
	return s.repo.GetByStatus(ctx, status)
}
