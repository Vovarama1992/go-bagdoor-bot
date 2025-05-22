package order

import (
	"context"
	"errors"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
)

type PostgresRepository struct {
	DB *db.DB
}

func NewPostgresRepository(db *db.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) CreateOrder(ctx context.Context, o *Order) error {
	query := `
		INSERT INTO orders (
			order_number, publisher_username, publisher_tg_id, published_at,
			origin_city, destination_city, start_date, end_date,
			title, description, reward, deposit, cost,
			media_urls, type, status
		)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11, $12, $13,
			$14, $15, $16
		)
		RETURNING id
	`

	return r.DB.Pool.QueryRow(ctx, query,
		o.OrderNumber,
		o.PublisherUsername,
		o.PublisherTgID,
		o.PublishedAt,
		o.OriginCity,
		o.DestinationCity,
		o.StartDate,
		o.EndDate,
		o.Title,
		o.Description,
		o.Reward,
		o.Deposit,
		o.Cost,
		o.MediaURLs,
		o.Type,
		o.Status,
	).Scan(&o.ID)
}

func (r *PostgresRepository) AddMediaURLs(ctx context.Context, orderID int, urls []string) error {
	if len(urls) == 0 {
		return errors.New("empty media URL list")
	}
	_, err := r.DB.Pool.Exec(ctx, `
		UPDATE orders
		SET media_urls = array_cat(media_urls, $1)
		WHERE id = $2
	`, urls, orderID)
	return err
}

func (r *PostgresRepository) UpdateModerationStatus(ctx context.Context, orderID int, status ModerationStatus) error {
	_, err := r.DB.Pool.Exec(ctx, `
		UPDATE orders
		SET moderation_status = $1
		WHERE id = $2
	`, status, orderID)
	return err
}

func (r *PostgresRepository) GetOrderByID(ctx context.Context, orderID int) (*Order, error) {
	query := `
		SELECT id, order_number, publisher_username, publisher_tg_id, published_at,
		       origin_city, destination_city, start_date, end_date,
		       title, description, reward, deposit, cost,
		       media_urls, type, status
		FROM orders
		WHERE id = $1
	`

	var o Order
	err := r.DB.Pool.QueryRow(ctx, query, orderID).Scan(
		&o.ID,
		&o.OrderNumber,
		&o.PublisherUsername,
		&o.PublisherTgID,
		&o.PublishedAt,
		&o.OriginCity,
		&o.DestinationCity,
		&o.StartDate,
		&o.EndDate,
		&o.Title,
		&o.Description,
		&o.Reward,
		&o.Deposit,
		&o.Cost,
		&o.MediaURLs,
		&o.Type,
		&o.Status,
	)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *PostgresRepository) GetAllOrders(ctx context.Context) ([]*Order, error) {
	query := `
		SELECT id, order_number, publisher_username, publisher_tg_id, published_at,
		       origin_city, destination_city, start_date, end_date,
		       title, description, reward, deposit, cost,
		       media_urls, type, status
		FROM orders
		ORDER BY published_at DESC
	`

	rows, err := r.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID,
			&o.OrderNumber,
			&o.PublisherUsername,
			&o.PublisherTgID,
			&o.PublishedAt,
			&o.OriginCity,
			&o.DestinationCity,
			&o.StartDate,
			&o.EndDate,
			&o.Title,
			&o.Description,
			&o.Reward,
			&o.Deposit,
			&o.Cost,
			&o.MediaURLs,
			&o.Type,
			&o.Status,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}
