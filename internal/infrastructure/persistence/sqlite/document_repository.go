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
			id, case_file_id, original_name, storage_path, mime_type, file_hash, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
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
			updated_at = ?
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
		SELECT id, case_file_id, original_name, storage_path, mime_type, file_hash, created_at, updated_at
		FROM documents
		WHERE id = ?
	`

	var rawID, caseFileID, originalName, storagePath, mimeType, fileHash, createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&rawID,
		&caseFileID,
		&originalName,
		&storagePath,
		&mimeType,
		&fileHash,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return document.Document{}, shared.ErrNotFound
		}
		return document.Document{}, err
	}

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
	}, nil
}

func (r *DocumentRepository) ListAll(ctx context.Context) ([]document.Document, error) {
	const query = `
		SELECT id, case_file_id, original_name, storage_path, mime_type, file_hash, created_at, updated_at
		FROM documents
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []document.Document

	for rows.Next() {
		var rawID, cfID, name, path, mime, hash, createdAt, updatedAt string

		if err := rows.Scan(&rawID, &cfID, &name, &path, &mime, &hash, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		ct, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}

		ut, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, document.Document{
			ID:           shared.NewID(rawID),
			CaseFileID:   shared.NewID(cfID),
			OriginalName: name,
			StoragePath:  path,
			MimeType:     mime,
			FileHash:     hash,
			CreatedAt:    shared.Timestamp(ct),
			UpdatedAt:    shared.Timestamp(ut),
		})
	}

	return result, rows.Err()
}

func (r *DocumentRepository) ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]document.Document, error) {
	const query = `
		SELECT id, case_file_id, original_name, storage_path, mime_type, file_hash, created_at, updated_at
		FROM documents
		WHERE case_file_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, caseFileID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []document.Document

	for rows.Next() {
		var rawID, cfID, name, path, mime, hash, createdAt, updatedAt string

		if err := rows.Scan(&rawID, &cfID, &name, &path, &mime, &hash, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		ct, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}

		ut, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, document.Document{
			ID:           shared.NewID(rawID),
			CaseFileID:   shared.NewID(cfID),
			OriginalName: name,
			StoragePath:  path,
			MimeType:     mime,
			FileHash:     hash,
			CreatedAt:    shared.Timestamp(ct),
			UpdatedAt:    shared.Timestamp(ut),
		})
	}

	return result, rows.Err()
}

func (r *DocumentRepository) GetByCaseFileIDAndFileHash(ctx context.Context, caseFileID shared.ID, fileHash string) (document.Document, error) {
	const query = `
		SELECT id, case_file_id, original_name, storage_path, mime_type, file_hash, created_at, updated_at
		FROM documents
		WHERE case_file_id = ?
		  AND file_hash = ?
		LIMIT 1
	`

	var rawID, cfID, name, path, mime, hash, createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, caseFileID.String(), fileHash).Scan(
		&rawID,
		&cfID,
		&name,
		&path,
		&mime,
		&hash,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return document.Document{}, shared.ErrNotFound
		}
		return document.Document{}, err
	}

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
		CaseFileID:   shared.NewID(cfID),
		OriginalName: name,
		StoragePath:  path,
		MimeType:     mime,
		FileHash:     hash,
		CreatedAt:    shared.Timestamp(ct),
		UpdatedAt:    shared.Timestamp(ut),
	}, nil
}
