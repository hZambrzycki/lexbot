package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"lexbox/internal/domain/client"
	"lexbox/internal/domain/shared"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Save(ctx context.Context, c client.Client) error {
	const query = `
		INSERT INTO clients (
			id, name, email, phone, identifier, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		c.ID.String(),
		c.Name,
		c.Email,
		c.Phone,
		c.Identifier,
		c.CreatedAt.Time().Format(time.RFC3339),
		c.UpdatedAt.Time().Format(time.RFC3339),
	)

	return err
}

func (r *ClientRepository) GetByID(ctx context.Context, id shared.ID) (client.Client, error) {
	const query = `
		SELECT id, name, email, phone, identifier, created_at, updated_at
		FROM clients
		WHERE id = ?
	`

	var (
		rawID      string
		name       string
		email      string
		phone      string
		identifier string
		createdAt  string
		updatedAt  string
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&rawID,
		&name,
		&email,
		&phone,
		&identifier,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return client.Client{}, shared.ErrNotFound
		}
		return client.Client{}, err
	}

	createdTime, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return client.Client{}, err
	}

	updatedTime, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return client.Client{}, err
	}

	return client.Client{
		ID:         shared.NewID(rawID),
		Name:       name,
		Email:      email,
		Phone:      phone,
		Identifier: identifier,
		CreatedAt:  shared.Timestamp(createdTime),
		UpdatedAt:  shared.Timestamp(updatedTime),
	}, nil
}

func (r *ClientRepository) List(ctx context.Context) ([]client.Client, error) {
	const query = `
		SELECT id, name, email, phone, identifier, created_at, updated_at
		FROM clients
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []client.Client

	for rows.Next() {
		var (
			rawID      string
			name       string
			email      string
			phone      string
			identifier string
			createdAt  string
			updatedAt  string
		)

		if err := rows.Scan(
			&rawID,
			&name,
			&email,
			&phone,
			&identifier,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		createdTime, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}

		updatedTime, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, client.Client{
			ID:         shared.NewID(rawID),
			Name:       name,
			Email:      email,
			Phone:      phone,
			Identifier: identifier,
			CreatedAt:  shared.Timestamp(createdTime),
			UpdatedAt:  shared.Timestamp(updatedTime),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
