package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/shared"
	searchinfra "lexbox/internal/infrastructure/search"
	"strings"
	"time"
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
			id,
			client_id,
			reference,
			title,
			type,
			status,
			description,
			calendar_scope,
			august_non_business,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	augustNonBusiness := 0
	if cf.AugustNonBusiness {
		augustNonBusiness = 1
	}

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
		cf.CalendarScope,
		augustNonBusiness,
		cf.CreatedAt.Time().Format(time.RFC3339),
		cf.UpdatedAt.Time().Format(time.RFC3339),
	)

	return err
}

func (r *CaseFileRepository) GetByID(ctx context.Context, id shared.ID) (casefile.CaseFile, error) {
	const query = `
		SELECT
			id,
			client_id,
			reference,
			title,
			type,
			status,
			description,
			calendar_scope,
			august_non_business,
			created_at,
			updated_at
		FROM case_files
		WHERE id = ?
	`

	var (
		rawID             string
		clientID          string
		reference         string
		title             string
		t                 string
		status            string
		desc              string
		calendarScope     string
		augustNonBusiness int
		createdAt         string
		updatedAt         string
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&rawID,
		&clientID,
		&reference,
		&title,
		&t,
		&status,
		&desc,
		&calendarScope,
		&augustNonBusiness,
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
		ID:                shared.NewID(rawID),
		ClientID:          shared.NewID(clientID),
		Reference:         reference,
		Title:             title,
		Type:              casefile.Type(t),
		Status:            casefile.Status(status),
		Description:       desc,
		CalendarScope:     calendarScope,
		AugustNonBusiness: augustNonBusiness == 1,
		CreatedAt:         shared.Timestamp(createdTime),
		UpdatedAt:         shared.Timestamp(updatedTime),
	}, nil
}

func (r *CaseFileRepository) List(ctx context.Context) ([]casefile.CaseFile, error) {
	const query = `
		SELECT
			id,
			client_id,
			reference,
			title,
			type,
			status,
			description,
			calendar_scope,
			august_non_business,
			created_at,
			updated_at
		FROM case_files
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []casefile.CaseFile

	for rows.Next() {
		var (
			rawID             string
			clientID          string
			reference         string
			title             string
			t                 string
			status            string
			desc              string
			calendarScope     string
			augustNonBusiness int
			createdAt         string
			updatedAt         string
		)

		if err := rows.Scan(
			&rawID,
			&clientID,
			&reference,
			&title,
			&t,
			&status,
			&desc,
			&calendarScope,
			&augustNonBusiness,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		createdTime, _ := time.Parse(time.RFC3339, createdAt)
		updatedTime, _ := time.Parse(time.RFC3339, updatedAt)

		result = append(result, casefile.CaseFile{
			ID:                shared.NewID(rawID),
			ClientID:          shared.NewID(clientID),
			Reference:         reference,
			Title:             title,
			Type:              casefile.Type(t),
			Status:            casefile.Status(status),
			Description:       desc,
			CalendarScope:     calendarScope,
			AugustNonBusiness: augustNonBusiness == 1,
			CreatedAt:         shared.Timestamp(createdTime),
			UpdatedAt:         shared.Timestamp(updatedTime),
		})
	}

	return result, rows.Err()
}

func (r *CaseFileRepository) ListByClientID(ctx context.Context, clientID shared.ID) ([]casefile.CaseFile, error) {
	const query = `
		SELECT
			id,
			client_id,
			reference,
			title,
			type,
			status,
			description,
			calendar_scope,
			august_non_business,
			created_at,
			updated_at
		FROM case_files
		WHERE client_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, clientID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []casefile.CaseFile

	for rows.Next() {
		var (
			rawID             string
			cID               string
			reference         string
			title             string
			t                 string
			status            string
			desc              string
			calendarScope     string
			augustNonBusiness int
			createdAt         string
			updatedAt         string
		)

		if err := rows.Scan(
			&rawID,
			&cID,
			&reference,
			&title,
			&t,
			&status,
			&desc,
			&calendarScope,
			&augustNonBusiness,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		createdTime, _ := time.Parse(time.RFC3339, createdAt)
		updatedTime, _ := time.Parse(time.RFC3339, updatedAt)

		result = append(result, casefile.CaseFile{
			ID:                shared.NewID(rawID),
			ClientID:          shared.NewID(cID),
			Reference:         reference,
			Title:             title,
			Type:              casefile.Type(t),
			Status:            casefile.Status(status),
			Description:       desc,
			CalendarScope:     calendarScope,
			AugustNonBusiness: augustNonBusiness == 1,
			CreatedAt:         shared.Timestamp(createdTime),
			UpdatedAt:         shared.Timestamp(updatedTime),
		})
	}

	return result, rows.Err()
}

func (r *CaseFileRepository) Update(ctx context.Context, cf casefile.CaseFile) error {
	const query = `
		UPDATE case_files
		SET
			reference = ?,
			title = ?,
			type = ?,
			status = ?,
			description = ?,
			calendar_scope = ?,
			august_non_business = ?,
			updated_at = ?
		WHERE id = ?
	`

	augustNonBusiness := 0
	if cf.AugustNonBusiness {
		augustNonBusiness = 1
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		cf.Reference,
		cf.Title,
		string(cf.Type),
		string(cf.Status),
		cf.Description,
		cf.CalendarScope,
		augustNonBusiness,
		cf.UpdatedAt.Time().Format(time.RFC3339),
		cf.ID.String(),
	)

	return err
}

func (r *CaseFileRepository) SearchCaseFiles(ctx context.Context, rawQuery string, limit int) ([]querymodels.GlobalSearchResult, error) {
	query := "%" + searchinfra.NormalizeText(rawQuery) + "%"

	const sqlQuery = `
		SELECT id, reference, title, type, status, description
		FROM case_files
		WHERE
			lower(reference) LIKE ?
			OR lower(title) LIKE ?
			OR lower(type) LIKE ?
			OR lower(status) LIKE ?
			OR lower(description) LIKE ?
		ORDER BY updated_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, query, query, query, query, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []querymodels.GlobalSearchResult{}

	for rows.Next() {
		var id, reference, title, typ, status, description string

		if err := rows.Scan(&id, &reference, &title, &typ, &status, &description); err != nil {
			return nil, err
		}

		subtitle := strings.TrimSpace(reference + " · " + typ + " · " + status)

		results = append(results, querymodels.GlobalSearchResult{
			Type:     "case_file",
			ID:       id,
			Title:    title,
			Subtitle: subtitle,
			Href:     "/case-files/" + id,
			Snippet:  description,
			Score:    80,
		})
	}

	return results, rows.Err()
}
