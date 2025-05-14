package user

import (
	tele "gopkg.in/telebot.v3"
)

type Service struct {
	repo Repository
}

type Repository interface {
	CreateUser(ctx tele.Context, u *User) error
	GetByTgID(ctx tele.Context, tgID int64) (*User, error)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Метод для регистрации пользователя
func (s *Service) RegisterUser(ctx tele.Context, u *User) error {
	// Проверяем, существует ли уже пользователь
	existing, err := s.repo.GetByTgID(ctx, u.TgID)
	if err == nil && existing != nil {
		// Пользователь уже существует, ничего не делаем
		return nil
	}

	// Создаем нового пользователя, включая его номер телефона
	return s.repo.CreateUser(ctx, u)
}

// Метод для получения пользователя по tgID
func (s *Service) GetByTgID(ctx tele.Context, tgID int64) (*User, error) {
	return s.repo.GetByTgID(ctx, tgID)
}
