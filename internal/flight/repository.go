package flight

import (
	"context"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
)

type Repository interface {
	Create(flight *Flight) error
	GetByID(id int) (*Flight, error)
	UpdateMapURL(id int, url string) error
	UpdateStatus(id int, status ModerationStatus) error
	GetAllFlights(ctx context.Context) ([]*Flight, error)
	GetByStatus(ctx context.Context, status ModerationStatus) ([]*Flight, error) // ← вот этого не хватает
}

type PostgresRepository struct {
	DB *db.DB
}

func NewPostgresRepository(db *db.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) Create(f *Flight) error {
	query := `
		INSERT INTO flights (
			flight_number,
			publisher_username,
			publisher_tg_id,
			published_at,
			flight_date,
			description,
			origin,
			destination,
			status,
			map_url
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id`
	return r.DB.Pool.QueryRow(
		context.Background(),
		query,
		f.FlightNumber,
		f.PublisherUsername,
		f.PublisherTgID,
		f.PublishedAt,
		f.FlightDate,
		f.Description,
		f.Origin,
		f.Destination,
		f.Status,
		f.MapURL,
	).Scan(&f.ID)
}

func (r *PostgresRepository) GetByID(id int) (*Flight, error) {
	query := `
		SELECT id, flight_number, publisher_username, publisher_tg_id, 
		       published_at, flight_date, description, origin, destination, status, map_url 
		FROM flights 
		WHERE id = $1`
	row := r.DB.Pool.QueryRow(context.Background(), query, id)

	var f Flight
	err := row.Scan(
		&f.ID,
		&f.FlightNumber,
		&f.PublisherUsername,
		&f.PublisherTgID,
		&f.PublishedAt,
		&f.FlightDate,
		&f.Description,
		&f.Origin,
		&f.Destination,
		&f.Status,
		&f.MapURL,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *PostgresRepository) UpdateMapURL(id int, url string) error {
	query := `UPDATE flights SET map_url = $1 WHERE id = $2`
	_, err := r.DB.Pool.Exec(context.Background(), query, url, id)
	return err
}

func (r *PostgresRepository) UpdateStatus(id int, status ModerationStatus) error {
	query := `UPDATE flights SET status = $1 WHERE id = $2`
	_, err := r.DB.Pool.Exec(context.Background(), query, status, id)
	return err
}

func (r *PostgresRepository) GetAllFlights(ctx context.Context) ([]*Flight, error) {
	query := `
		SELECT id, flight_number, publisher_username, publisher_tg_id,
		       published_at, flight_date, description, origin, destination, status, map_url
		FROM flights
		ORDER BY published_at DESC
	`

	rows, err := r.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []*Flight
	for rows.Next() {
		var f Flight
		if err := rows.Scan(
			&f.ID,
			&f.FlightNumber,
			&f.PublisherUsername,
			&f.PublisherTgID,
			&f.PublishedAt,
			&f.FlightDate,
			&f.Description,
			&f.Origin,
			&f.Destination,
			&f.Status,
			&f.MapURL,
		); err != nil {
			return nil, err
		}
		flights = append(flights, &f)
	}

	return flights, nil
}

func (r *PostgresRepository) GetByStatus(ctx context.Context, status ModerationStatus) ([]*Flight, error) {
	query := `
		SELECT id, flight_number, publisher_username, publisher_tg_id,
		       published_at, flight_date, description, origin, destination,
		       status, map_url
		FROM flights
		WHERE status = $1
		ORDER BY published_at DESC
	`

	rows, err := r.DB.Pool.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []*Flight
	for rows.Next() {
		var f Flight
		if err := rows.Scan(
			&f.ID, &f.FlightNumber, &f.PublisherUsername, &f.PublisherTgID,
			&f.PublishedAt, &f.FlightDate, &f.Description, &f.Origin, &f.Destination,
			&f.Status, &f.MapURL,
		); err != nil {
			return nil, err
		}
		flights = append(flights, &f)
	}
	return flights, nil
}
