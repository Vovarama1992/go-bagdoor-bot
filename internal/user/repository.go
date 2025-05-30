package user

import (
	"context" // Используем стандартный контекст
	"errors"
	"log"

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
		log.Printf("Пользователь уже существует: tg_id=%d", u.TgID)
		return errors.New("пользователь с таким Telegram ID уже существует")
	} else if err != nil && err.Error() != "no rows in result set" {
		log.Printf("Ошибка при проверке существующего пользователя: %v", err)
	}

	// Запрос на добавление нового пользователя
	query := `
		INSERT INTO users (tg_username, tg_id, first_name, last_name, phone_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, registered_at;
	`

	err = r.DB.Pool.QueryRow(ctx, query,
		u.TgUsername,
		u.TgID,
		u.FirstName,
		u.LastName,
		u.PhoneNumber,
	).Scan(&u.ID, &u.RegisteredAt)

	if err != nil {
		log.Printf("Ошибка при вставке пользователя: %v", err)
		return err
	}

	log.Printf("Пользователь успешно создан: id=%d, tg_id=%d", u.ID, u.TgID)
	return nil
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

func (r *PostgresRepository) CreateFromTgToWebApp(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (tg_username, tg_id, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, registered_at;
	`

	err := r.DB.Pool.QueryRow(ctx, query,
		u.TgUsername,
		u.TgID,
		u.FirstName,
		u.LastName,
	).Scan(&u.ID, &u.RegisteredAt)

	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		return err
	}

	return nil
}
