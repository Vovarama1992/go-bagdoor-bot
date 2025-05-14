package user

import (
	"context"
)

type Service struct {
	repo *PostgresRepository
}

func NewService(repo *PostgresRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(ctx context.Context, u *User) error {
	return s.repo.CreateUser(ctx, u)
}

func (s *Service) GetByTgID(ctx context.Context, tgID int64) (*User, error) {
	return s.repo.GetByTgID(ctx, tgID)
}
