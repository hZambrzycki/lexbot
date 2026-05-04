package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Save(ctx context.Context, d document.Document) error {
	const query = `
		INSERT INTO documents (
			id, case_file_id, original_name, storage_path, mime_type, file_hash,
			created_at, updated_at, review_status, reviewed_at, review_note
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		d.ID.String(),
		d.CaseFileID.String(),
		d.OriginalName,
		d.StoragePath,
		d.MimeType,
		d.FileHash,
		d.CreatedAt.Time().Format(time.RFC3339),
		d.UpdatedAt.Time().Format(time.RFC3339),
		emptyDefault(d.ReviewStatus, document.DocumentReviewStatusPending),
		d.ReviewedAt,
		d.ReviewNote,
	)

	return err
}

func (r *DocumentRepository) Update(ctx context.Context, d document.Document) error {
	const query = `
		UPDATE documents
		SET
			case_file_id = ?,
			original_name = ?,
			storage_path = ?,
			mime_type = ?,
			file_hash = ?,
			updated_at = ?,
			review_status = ?,
			reviewed_at = ?,
			review_note = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		d.CaseFileID.String(),
		d.OriginalName,
		d.StoragePath,
		d.MimeType,
		d.FileHash,
		d.UpdatedAt.Time().Format(time.RFC3339),
		emptyDefault(d.ReviewStatus, document.DocumentReviewStatusPending),
		d.ReviewedAt,
		d.ReviewNote,
		d.ID.String(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return shared.ErrNotFound
	}

	return nil
}

func (r *DocumentRepository) GetByID(ctx context.Context, id shared.ID) (document.Document, error) {
	const query = `
		SELECT
			id, case_file_id, original_name, storage_path, mime_type, file_hash,
			created_at, updated_at,
			COALESCE(review_status, 'pending_review'),
			COALESCE(reviewed_at, ''),
			COALESCE(review_note, '')
		FROM documents
		WHERE id = ?
	`

	var rawID, caseFileID, originalName, storagePath, mimeType, fileHash string
	var createdAt, updatedAt, reviewStatus, reviewedAt, reviewNote string

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&rawID,
		&caseFileID,
		&originalName,
		&storagePath,
		&mimeType,
		&fileHash,
		&createdAt,
		&updatedAt,
		&reviewStatus,
		&reviewedAt,
		&reviewNote,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return document.Document{}, shared.ErrNotFound
		}
		return document.Document{}, err
	}

	return buildDocument(rawID, caseFileID, originalName, storagePath, mimeType, fileHash, createdAt, updatedAt, reviewStatus, reviewedAt, reviewNote)
}

func (r *DocumentRepository) ListAll(ctx context.Context) ([]document.Document, error) {
	const query = `
		SELECT
			id, case_file_id, original_name, storage_path, mime_type, file_hash,
			created_at, updated_at,
			COALESCE(review_status, 'pending_review'),
			COALESCE(reviewed_at, ''),
			COALESCE(review_note, '')
		FROM documents
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDocuments(rows)
}

func (r *DocumentRepository) ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]document.Document, error) {
	const query = `
		SELECT
			id, case_file_id, original_name, storage_path, mime_type, file_hash,
			created_at, updated_at,
			COALESCE(review_status, 'pending_review'),
			COALESCE(reviewed_at, ''),
			COALESCE(review_note, '')
		FROM documents
		WHERE case_file_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, caseFileID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDocuments(rows)
}

func (r *DocumentRepository) GetByCaseFileIDAndFileHash(ctx context.Context, caseFileID shared.ID, fileHash string) (document.Document, error) {
	const query = `
		SELECT
			id, case_file_id, original_name, storage_path, mime_type, file_hash,
			created_at, updated_at,
			COALESCE(review_status, 'pending_review'),
			COALESCE(reviewed_at, ''),
			COALESCE(review_note, '')
		FROM documents
		WHERE case_file_id = ?
		  AND file_hash = ?
		LIMIT 1
	`

	var rawID, cfID, name, path, mime, hash string
	var createdAt, updatedAt, reviewStatus, reviewedAt, reviewNote string

	err := r.db.QueryRowContext(ctx, query, caseFileID.String(), fileHash).Scan(
		&rawID,
		&cfID,
		&name,
		&path,
		&mime,
		&hash,
		&createdAt,
		&updatedAt,
		&reviewStatus,
		&reviewedAt,
		&reviewNote,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return document.Document{}, shared.ErrNotFound
		}
		return document.Document{}, err
	}

	return buildDocument(rawID, cfID, name, path, mime, hash, createdAt, updatedAt, reviewStatus, reviewedAt, reviewNote)
}

func (r *DocumentRepository) UpdateReviewState(ctx context.Context, id shared.ID, reviewStatus, reviewedAt, reviewNote string) error {
	const query = `
		UPDATE documents
		SET
			review_status = ?,
			reviewed_at = ?,
			review_note = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		reviewStatus,
		reviewedAt,
		reviewNote,
		time.Now().Format(time.RFC3339),
		id.String(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return shared.ErrNotFound
	}

	return nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id shared.ID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	documentID := id.String()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM document_events
		WHERE document_id = ?
	`, documentID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM document_metadata
		WHERE document_id = ?
	`, documentID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM document_contents
		WHERE document_id = ?
	`, documentID); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		DELETE FROM documents
		WHERE id = ?
	`, documentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return shared.ErrNotFound
	}

	return tx.Commit()
}

func scanDocuments(rows *sql.Rows) ([]document.Document, error) {
	var result []document.Document

	for rows.Next() {
		var rawID, cfID, name, path, mime, hash string
		var createdAt, updatedAt, reviewStatus, reviewedAt, reviewNote string

		if err := rows.Scan(
			&rawID,
			&cfID,
			&name,
			&path,
			&mime,
			&hash,
			&createdAt,
			&updatedAt,
			&reviewStatus,
			&reviewedAt,
			&reviewNote,
		); err != nil {
			return nil, err
		}

		doc, err := buildDocument(rawID, cfID, name, path, mime, hash, createdAt, updatedAt, reviewStatus, reviewedAt, reviewNote)
		if err != nil {
			return nil, err
		}

		result = append(result, doc)
	}

	return result, rows.Err()
}

func buildDocument(
	rawID string,
	caseFileID string,
	originalName string,
	storagePath string,
	mimeType string,
	fileHash string,
	createdAt string,
	updatedAt string,
	reviewStatus string,
	reviewedAt string,
	reviewNote string,
) (document.Document, error) {
	ct, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return document.Document{}, err
	}

	ut, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return document.Document{}, err
	}

	return document.Document{
		ID:           shared.NewID(rawID),
		CaseFileID:   shared.NewID(caseFileID),
		OriginalName: originalName,
		StoragePath:  storagePath,
		MimeType:     mimeType,
		FileHash:     fileHash,
		CreatedAt:    shared.Timestamp(ct),
		UpdatedAt:    shared.Timestamp(ut),

		ReviewStatus: emptyDefault(reviewStatus, document.DocumentReviewStatusPending),
		ReviewedAt:   reviewedAt,
		ReviewNote:   reviewNote,
	}, nil
}
