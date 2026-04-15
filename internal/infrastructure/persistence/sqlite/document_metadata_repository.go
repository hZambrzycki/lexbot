package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type DocumentMetadataRepository struct {
	db *sql.DB
}

func NewDocumentMetadataRepository(db *sql.DB) *DocumentMetadataRepository {
	return &DocumentMetadataRepository{db: db}
}

func (r *DocumentMetadataRepository) Save(ctx context.Context, metadata document.Metadata) error {
	const query = `
		INSERT INTO document_metadata (
			document_id,
			document_type,
			legal_area,
			analyzed_at
		)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(document_id) DO UPDATE SET
			document_type = excluded.document_type,
			legal_area = excluded.legal_area,
			analyzed_at = excluded.analyzed_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		metadata.DocumentID.String(),
		metadata.DocumentType,
		metadata.LegalArea,
		metadata.AnalyzedAt,
	)

	return err
}

func (r *DocumentMetadataRepository) GetByDocumentID(ctx context.Context, documentID shared.ID) (document.Metadata, error) {
	const query = `
		SELECT document_id, document_type, legal_area, analyzed_at
		FROM document_metadata
		WHERE document_id = ?
	`

	var (
		rawDocumentID string
		documentType  string
		legalArea     string
		analyzedAt    string
	)

	err := r.db.QueryRowContext(ctx, query, documentID.String()).Scan(
		&rawDocumentID,
		&documentType,
		&legalArea,
		&analyzedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return document.Metadata{}, shared.ErrNotFound
		}
		return document.Metadata{}, err
	}

	return document.Metadata{
		DocumentID:   shared.NewID(rawDocumentID),
		DocumentType: documentType,
		LegalArea:    legalArea,
		AnalyzedAt:   analyzedAt,
	}, nil
}
