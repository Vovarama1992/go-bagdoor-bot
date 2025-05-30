package order

import (
	"context"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
)

type Service struct {
	Repo     *PostgresRepository
	UserRepo *user.PostgresRepository
}

// вот так — корректно:
func NewService(repo *PostgresRepository, userRepo *user.PostgresRepository) *Service {
	return &Service{
		Repo:     repo,
		UserRepo: userRepo,
	}
}

func (s *Service) CreateOrder(ctx context.Context, o *Order) error {
	return s.Repo.CreateOrder(ctx, o)
}

func (s *Service) AddMediaURLs(ctx context.Context, orderID int, urls []string) error {
	return s.Repo.AddMediaURLs(ctx, orderID, urls)
}

func (s *Service) UpdateModerationStatus(ctx context.Context, orderID int, status ModerationStatus) error {
	return s.Repo.UpdateModerationStatus(ctx, orderID, status)
}

func (s *Service) GetAllOrders(ctx context.Context) ([]*Order, error) {
	return s.Repo.GetAllOrders(ctx)
}

func (s *Service) GetByStatus(ctx context.Context, status ModerationStatus) ([]*Order, error) {
	return s.Repo.GetByStatus(ctx, status)
}
