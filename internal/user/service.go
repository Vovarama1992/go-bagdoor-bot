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

func (s *Service) FindOrCreateFromTelegram(ctx context.Context, tgID int64, username, first, last string) (*User, error) {
	user, err := s.repo.GetByTgID(ctx, tgID)
	if err == nil {
		return user, nil
	}

	// если не найден — создать
	newUser := &User{
		TgID:       tgID,
		TgUsername: username,
		FirstName:  first,
		LastName:   last,
	}
	if err := s.repo.CreateFromTgToWebApp(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}
