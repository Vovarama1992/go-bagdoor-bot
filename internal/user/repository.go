package user

import (
	"context" // Используем стандартный контекст
	"errors"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
)

type PostgresRepository struct {
	DB *db.DB
}

func NewPostgresRepository(db *db.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

// Создание нового пользователя
func (r *PostgresRepository) CreateUser(ctx context.Context, u *User) error {
	// Проверяем, существует ли пользователь с таким tgID
	existingUser, err := r.GetByTgID(ctx, u.TgID)
	if err == nil && existingUser != nil {
		return errors.New("пользователь с таким Telegram ID уже существует")
	}

	// Запрос на добавление нового пользователя
	query := `
		INSERT INTO users (tg_username, tg_id, first_name, last_name, phone_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, registered_at;
	`

	// Выполнение запроса с использованием стандартного context.Context
	return r.DB.Pool.QueryRow(ctx, query,
		u.TgUsername,
		u.TgID,
		u.FirstName,
		u.LastName,
		u.PhoneNumber,
	).Scan(&u.ID, &u.RegisteredAt)
}

// Получение пользователя по tgID
func (r *PostgresRepository) GetByTgID(ctx context.Context, tgID int64) (*User, error) {
	query := `
		SELECT id, tg_username, tg_id, first_name, last_name, phone_number, registered_at
		FROM users WHERE tg_id = $1
	`

	var u User
	// Используем стандартный context.Context
	err := r.DB.Pool.QueryRow(ctx, query, tgID).Scan(
		&u.ID,
		&u.TgUsername,
		&u.TgID,
		&u.FirstName,
		&u.LastName,
		&u.PhoneNumber,
		&u.RegisteredAt,
	)

	// Если пользователя нет в базе
	if err != nil {
		return nil, err
	}
	return &u, nil
}
