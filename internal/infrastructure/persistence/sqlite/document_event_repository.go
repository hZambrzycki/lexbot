package sqlite

import (
	"context"
	"database/sql"
	"strings"

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
			add_extra_day,
			calendar_scope,
			trigger_text,
			computation,
			review_status,
			reviewed_at,
			resolved_at,
			resolution_note
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, event := range events {
		isBusinessDays := 0
		if event.IsBusinessDays {
			isBusinessDays = 1
		}

		addExtraDay := 0
		if event.AddExtraDay {
			addExtraDay = 1
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
			addExtraDay,
			emptyDefault(event.CalendarScope, "madrid"),
			nullIfBlank(event.TriggerText),
			emptyDefault(event.Computation, ""),
			emptyDefault(event.ReviewStatus, document.ReviewStatusPending),
			emptyDefault(event.ReviewedAt, ""),
			emptyDefault(event.ResolvedAt, ""),
			emptyDefault(event.ResolutionNote, ""),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *DocumentEventRepository) GetByID(ctx context.Context, eventID shared.ID) (document.Event, error) {
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
			add_extra_day,
			calendar_scope,
			trigger_text,
			computation,
			review_status,
			reviewed_at,
			resolved_at,
			resolution_note
		FROM document_events
		WHERE id = ?
	`

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
		addExtraDay    int
		calendarScope  string
		triggerText    sql.NullString
		computation    string
		reviewStatus   string
		reviewedAt     string
		resolvedAt     string
		resolutionNote string
	)

	err := r.db.QueryRowContext(ctx, query, eventID.String()).Scan(
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
		&addExtraDay,
		&calendarScope,
		&triggerText,
		&computation,
		&reviewStatus,
		&reviewedAt,
		&resolvedAt,
		&resolutionNote,
	)
	if err != nil {
		return document.Event{}, err
	}

	return document.Event{
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
		AddExtraDay:    addExtraDay == 1,
		CalendarScope:  calendarScope,
		TriggerText:    triggerText.String,
		Computation:    computation,
		ReviewStatus:   reviewStatus,
		ReviewedAt:     reviewedAt,
		ResolvedAt:     resolvedAt,
		ResolutionNote: resolutionNote,
	}, nil
}

func (r *DocumentEventRepository) UpdateReviewState(ctx context.Context, eventID shared.ID, reviewStatus, reviewedAt, resolvedAt, resolutionNote string) error {
	const query = `
		UPDATE document_events
		SET
			review_status = ?,
			reviewed_at = ?,
			resolved_at = ?,
			resolution_note = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		reviewStatus,
		emptyDefault(reviewedAt, ""),
		emptyDefault(resolvedAt, ""),
		emptyDefault(resolutionNote, ""),
		eventID.String(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
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
			add_extra_day,
			calendar_scope,
			trigger_text,
			computation,
			review_status,
			reviewed_at,
			resolved_at,
			resolution_note
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
			addExtraDay    int
			calendarScope  string
			triggerText    sql.NullString
			computation    string
			reviewStatus   string
			reviewedAt     string
			resolvedAt     string
			resolutionNote string
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
			&addExtraDay,
			&calendarScope,
			&triggerText,
			&computation,
			&reviewStatus,
			&reviewedAt,
			&resolvedAt,
			&resolutionNote,
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
			AddExtraDay:    addExtraDay == 1,
			CalendarScope:  calendarScope,
			TriggerText:    triggerText.String,
			Computation:    computation,
			ReviewStatus:   reviewStatus,
			ReviewedAt:     reviewedAt,
			ResolvedAt:     resolvedAt,
			ResolutionNote: resolutionNote,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *DocumentEventRepository) ListByCaseFileID(ctx context.Context, caseFileID shared.ID, reviewStatus string) ([]querymodels.CaseFileEventResult, error) {
	query := `
		SELECT
			de.id,
			de.document_id,
			d.original_name,
			cf.id,
			cf.reference,
			cf.title,
			de.event_type,
			de.event_date,
			de.source_text,
			COALESCE(de.anchor_date, ''),
			COALESCE(de.date_kind, ''),
			COALESCE(de.anchor_source, ''),
			COALESCE(de.relative_days, 0),
			COALESCE(de.is_business_days, 0),
			COALESCE(de.add_extra_day, 0),
			COALESCE(de.calendar_scope, 'madrid'),
			COALESCE(de.trigger_text, ''),
			COALESCE(de.computation, ''),
			COALESCE(de.review_status, 'pending'),
			COALESCE(de.reviewed_at, ''),
			COALESCE(de.resolved_at, ''),
			COALESCE(de.resolution_note, '')
		FROM document_events de
		INNER JOIN documents d ON d.id = de.document_id
		INNER JOIN case_files cf ON cf.id = d.case_file_id
		WHERE d.case_file_id = ?
	`

	args := []any{caseFileID.String()}

	if strings.TrimSpace(reviewStatus) != "" {
		query += ` AND de.review_status = ?`
		args = append(args, strings.TrimSpace(reviewStatus))
	}

	query += ` ORDER BY de.event_date ASC, de.id ASC`

	return r.scanCaseFileEventResults(ctx, query, args...)
}

func (r *DocumentEventRepository) ListUpcoming(ctx context.Context, fromDate string, caseFileID shared.ID, eventType string, reviewStatus string) ([]querymodels.CaseFileEventResult, error) {
	query := `
		SELECT
			de.id,
			de.document_id,
			d.original_name,
			cf.id,
			cf.reference,
			cf.title,
			de.event_type,
			de.event_date,
			de.source_text,
			COALESCE(de.anchor_date, ''),
			COALESCE(de.date_kind, ''),
			COALESCE(de.anchor_source, ''),
			COALESCE(de.relative_days, 0),
			COALESCE(de.is_business_days, 0),
			COALESCE(de.add_extra_day, 0),
			COALESCE(de.calendar_scope, 'madrid'),
			COALESCE(de.trigger_text, ''),
			COALESCE(de.computation, ''),
			COALESCE(de.review_status, 'pending'),
			COALESCE(de.reviewed_at, ''),
			COALESCE(de.resolved_at, ''),
			COALESCE(de.resolution_note, '')
		FROM document_events de
		INNER JOIN documents d ON d.id = de.document_id
		INNER JOIN case_files cf ON cf.id = d.case_file_id
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

	if strings.TrimSpace(reviewStatus) != "" {
		query += ` AND de.review_status = ?`
		args = append(args, strings.TrimSpace(reviewStatus))
	}

	if fromDate != "" {
		query += ` AND de.event_date >= ?`
		args = append(args, fromDate)
	}

	query += ` ORDER BY de.event_date ASC, de.id ASC`

	return r.scanCaseFileEventResults(ctx, query, args...)
}

func (r *DocumentEventRepository) scanCaseFileEventResults(ctx context.Context, query string, args ...any) ([]querymodels.CaseFileEventResult, error) {
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
			addExtraDay    int
		)

		if err := rows.Scan(
			&item.EventID,
			&item.DocumentID,
			&item.OriginalName,
			&item.CaseFileID,
			&item.CaseFileReference,
			&item.CaseFileTitle,
			&item.EventType,
			&item.EventDate,
			&item.SourceText,
			&item.AnchorDate,
			&item.DateKind,
			&item.AnchorSource,
			&item.RelativeDays,
			&isBusinessDays,
			&addExtraDay,
			&item.CalendarScope,
			&item.TriggerText,
			&item.Computation,
			&item.ReviewStatus,
			&item.ReviewedAt,
			&item.ResolvedAt,
			&item.ResolutionNote,
		); err != nil {
			return nil, err
		}

		item.IsBusinessDays = isBusinessDays == 1
		item.AddExtraDay = addExtraDay == 1

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

func emptyDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func (r *DocumentEventRepository) GetDetailByID(ctx context.Context, eventID shared.ID) (querymodels.CaseFileEventResult, error) {
	const query = `
		SELECT
			de.id,
			de.document_id,
			d.original_name,
			cf.id,
			cf.reference,
			cf.title,
			de.event_type,
			de.event_date,
			de.source_text,
			COALESCE(de.anchor_date, ''),
			COALESCE(de.date_kind, ''),
			COALESCE(de.anchor_source, ''),
			COALESCE(de.relative_days, 0),
			COALESCE(de.is_business_days, 0),
			COALESCE(de.add_extra_day, 0),
			COALESCE(de.calendar_scope, 'madrid'),
			COALESCE(de.trigger_text, ''),
			COALESCE(de.computation, ''),
			COALESCE(de.review_status, 'pending'),
			COALESCE(de.reviewed_at, ''),
			COALESCE(de.resolved_at, ''),
			COALESCE(de.resolution_note, '')
		FROM document_events de
		INNER JOIN documents d ON d.id = de.document_id
		INNER JOIN case_files cf ON cf.id = d.case_file_id
		WHERE de.id = ?
	`

	results, err := r.scanCaseFileEventResults(ctx, query, eventID.String())
	if err != nil {
		return querymodels.CaseFileEventResult{}, err
	}

	if len(results) == 0 {
		return querymodels.CaseFileEventResult{}, sql.ErrNoRows
	}

	return results[0], nil
}
