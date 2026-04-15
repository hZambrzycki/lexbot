package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/shared"
)

type CaseFileRepository struct {
	db *sql.DB
}

func NewCaseFileRepository(db *sql.DB) *CaseFileRepository {
	return &CaseFileRepository{db: db}
}

func (r *CaseFileRepository) Save(ctx context.Context, cf casefile.CaseFile) error {
	const query = `
		INSERT INTO case_files (
			id, client_id, reference, title, type, status, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		cf.ID.String(),
		cf.ClientID.String(),
		cf.Reference,
		cf.Title,
		string(cf.Type),
		string(cf.Status),
		cf.Description,
		cf.CreatedAt.Time().Format(time.RFC3339),
		cf.UpdatedAt.Time().Format(time.RFC3339),
	)

	return err
}

func (r *CaseFileRepository) GetByID(ctx context.Context, id shared.ID) (casefile.CaseFile, error) {
	const query = `
		SELECT id, client_id, reference, title, type, status, description, created_at, updated_at
		FROM case_files
		WHERE id = ?
	`

	var (
		rawID     string
		clientID  string
		reference string
		title     string
		t         string
		status    string
		desc      string
		createdAt string
		updatedAt string
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&rawID,
		&clientID,
		&reference,
		&title,
		&t,
		&status,
		&desc,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return casefile.CaseFile{}, shared.ErrNotFound
		}
		return casefile.CaseFile{}, err
	}

	createdTime, _ := time.Parse(time.RFC3339, createdAt)
	updatedTime, _ := time.Parse(time.RFC3339, updatedAt)

	return casefile.CaseFile{
		ID:          shared.NewID(rawID),
		ClientID:    shared.NewID(clientID),
		Reference:   reference,
		Title:       title,
		Type:        casefile.Type(t),
		Status:      casefile.Status(status),
		Description: desc,
		CreatedAt:   shared.Timestamp(createdTime),
		UpdatedAt:   shared.Timestamp(updatedTime),
	}, nil
}

func (r *CaseFileRepository) List(ctx context.Context) ([]casefile.CaseFile, error) {
	const query = `SELECT id, client_id, reference, title, type, status, description, created_at, updated_at FROM case_files`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []casefile.CaseFile

	for rows.Next() {
		var (
			rawID, clientID, reference, title, t, status, desc, createdAt, updatedAt string
		)

		if err := rows.Scan(&rawID, &clientID, &reference, &title, &t, &status, &desc, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		createdTime, _ := time.Parse(time.RFC3339, createdAt)
		updatedTime, _ := time.Parse(time.RFC3339, updatedAt)

		result = append(result, casefile.CaseFile{
			ID:          shared.NewID(rawID),
			ClientID:    shared.NewID(clientID),
			Reference:   reference,
			Title:       title,
			Type:        casefile.Type(t),
			Status:      casefile.Status(status),
			Description: desc,
			CreatedAt:   shared.Timestamp(createdTime),
			UpdatedAt:   shared.Timestamp(updatedTime),
		})
	}

	return result, rows.Err()
}

func (r *CaseFileRepository) ListByClientID(ctx context.Context, clientID shared.ID) ([]casefile.CaseFile, error) {
	const query = `SELECT id, client_id, reference, title, type, status, description, created_at, updated_at FROM case_files WHERE client_id = ?`

	rows, err := r.db.QueryContext(ctx, query, clientID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []casefile.CaseFile

	for rows.Next() {
		var (
			rawID, cID, reference, title, t, status, desc, createdAt, updatedAt string
		)

		if err := rows.Scan(&rawID, &cID, &reference, &title, &t, &status, &desc, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		createdTime, _ := time.Parse(time.RFC3339, createdAt)
		updatedTime, _ := time.Parse(time.RFC3339, updatedAt)

		result = append(result, casefile.CaseFile{
			ID:          shared.NewID(rawID),
			ClientID:    shared.NewID(cID),
			Reference:   reference,
			Title:       title,
			Type:        casefile.Type(t),
			Status:      casefile.Status(status),
			Description: desc,
			CreatedAt:   shared.Timestamp(createdTime),
			UpdatedAt:   shared.Timestamp(updatedTime),
		})
	}

	return result, rows.Err()
}
