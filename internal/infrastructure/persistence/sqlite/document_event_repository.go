package sqlite

import (
	"context"
	"database/sql"

	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type DocumentEventRepository struct {
	db *sql.DB
}

func NewDocumentEventRepository(db *sql.DB) *DocumentEventRepository {
	return &DocumentEventRepository{db: db}
}

func (r *DocumentEventRepository) ReplaceByDocumentID(ctx context.Context, documentID shared.ID, events []document.Event) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM document_events WHERE document_id = ?`, documentID.String()); err != nil {
		return err
	}

	const insertQuery = `
		INSERT INTO document_events (
			id,
			document_id,
			event_type,
			event_date,
			source_text,
			created_at,
			anchor_date,
			date_kind,
			anchor_source,
			relative_days,
			is_business_days,
			trigger_text
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, event := range events {
		isBusinessDays := 0
		if event.IsBusinessDays {
			isBusinessDays = 1
		}

		_, err := tx.ExecContext(
			ctx,
			insertQuery,
			event.ID.String(),
			event.DocumentID.String(),
			event.EventType,
			event.EventDate,
			event.SourceText,
			event.CreatedAt,
			nullIfBlank(event.AnchorDate),
			nullIfBlank(event.DateKind),
			nullIfBlank(event.AnchorSource),
			event.RelativeDays,
			isBusinessDays,
			nullIfBlank(event.TriggerText),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *DocumentEventRepository) ListByDocumentID(ctx context.Context, documentID shared.ID) ([]document.Event, error) {
	const query = `
		SELECT
			id,
			document_id,
			event_type,
			event_date,
			source_text,
			created_at,
			anchor_date,
			date_kind,
			anchor_source,
			relative_days,
			is_business_days,
			trigger_text
		FROM document_events
		WHERE document_id = ?
		ORDER BY event_date ASC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, documentID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]document.Event, 0)
	for rows.Next() {
		var (
			id             string
			rawDocID       string
			eventType      string
			eventDate      string
			sourceText     string
			createdAt      string
			anchorDate     sql.NullString
			dateKind       sql.NullString
			anchorSource   sql.NullString
			relativeDays   int
			isBusinessDays int
			triggerText    sql.NullString
		)

		if err := rows.Scan(
			&id,
			&rawDocID,
			&eventType,
			&eventDate,
			&sourceText,
			&createdAt,
			&anchorDate,
			&dateKind,
			&anchorSource,
			&relativeDays,
			&isBusinessDays,
			&triggerText,
		); err != nil {
			return nil, err
		}

		events = append(events, document.Event{
			ID:             shared.NewID(id),
			DocumentID:     shared.NewID(rawDocID),
			EventType:      eventType,
			EventDate:      eventDate,
			SourceText:     sourceText,
			CreatedAt:      createdAt,
			AnchorDate:     anchorDate.String,
			DateKind:       dateKind.String,
			AnchorSource:   anchorSource.String,
			RelativeDays:   relativeDays,
			IsBusinessDays: isBusinessDays == 1,
			TriggerText:    triggerText.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *DocumentEventRepository) ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]querymodels.CaseFileEventResult, error) {
	const query = `
		SELECT
			de.id,
			de.document_id,
			d.original_name,
			de.event_type,
			de.event_date,
			de.source_text,
			COALESCE(de.anchor_date, ''),
			COALESCE(de.date_kind, ''),
			COALESCE(de.anchor_source, ''),
			COALESCE(de.relative_days, 0),
			COALESCE(de.is_business_days, 0),
			COALESCE(de.trigger_text, '')
		FROM document_events de
		INNER JOIN documents d ON d.id = de.document_id
		WHERE d.case_file_id = ?
		ORDER BY de.event_date ASC, de.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, caseFileID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]querymodels.CaseFileEventResult, 0)
	for rows.Next() {
		var (
			item           querymodels.CaseFileEventResult
			isBusinessDays int
		)

		if err := rows.Scan(
			&item.EventID,
			&item.DocumentID,
			&item.OriginalName,
			&item.EventType,
			&item.EventDate,
			&item.SourceText,
			&item.AnchorDate,
			&item.DateKind,
			&item.AnchorSource,
			&item.RelativeDays,
			&isBusinessDays,
			&item.TriggerText,
		); err != nil {
			return nil, err
		}

		item.IsBusinessDays = isBusinessDays == 1
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *DocumentEventRepository) ListUpcoming(ctx context.Context, fromDate string, caseFileID shared.ID, eventType string) ([]querymodels.CaseFileEventResult, error) {
	query := `
		SELECT
			de.id,
			de.document_id,
			d.original_name,
			de.event_type,
			de.event_date,
			de.source_text,
			COALESCE(de.anchor_date, ''),
			COALESCE(de.date_kind, ''),
			COALESCE(de.anchor_source, ''),
			COALESCE(de.relative_days, 0),
			COALESCE(de.is_business_days, 0),
			COALESCE(de.trigger_text, '')
		FROM document_events de
		INNER JOIN documents d ON d.id = de.document_id
		WHERE 1=1
	`

	args := []any{}

	if caseFileID != "" {
		query += ` AND d.case_file_id = ?`
		args = append(args, caseFileID.String())
	}

	if eventType != "" {
		query += ` AND de.event_type = ?`
		args = append(args, eventType)
	}

	if fromDate != "" {
		query += ` AND de.event_date >= ?`
		args = append(args, fromDate)
	}

	query += ` ORDER BY de.event_date ASC, de.id ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]querymodels.CaseFileEventResult, 0)
	for rows.Next() {
		var (
			item           querymodels.CaseFileEventResult
			isBusinessDays int
		)

		if err := rows.Scan(
			&item.EventID,
			&item.DocumentID,
			&item.OriginalName,
			&item.EventType,
			&item.EventDate,
			&item.SourceText,
			&item.AnchorDate,
			&item.DateKind,
			&item.AnchorSource,
			&item.RelativeDays,
			&isBusinessDays,
			&item.TriggerText,
		); err != nil {
			return nil, err
		}

		item.IsBusinessDays = isBusinessDays == 1
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func nullIfBlank(value string) any {
	if value == "" {
		return nil
	}
	return value
}
