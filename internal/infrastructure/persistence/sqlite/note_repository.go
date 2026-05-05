package sqlite

import (
	"context"
	"database/sql"
	"time"

	"lexbox/internal/domain/note"
	"lexbox/internal/domain/shared"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Save(ctx context.Context, n note.Note) error {
	const query = `
		INSERT INTO notes (
			id, case_file_id, title, content, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		n.ID.String(),
		n.CaseFileID.String(),
		n.Title,
		n.Content,
		n.CreatedAt.Time().Format(time.RFC3339),
		n.UpdatedAt.Time().Format(time.RFC3339),
	)

	return err
}

func (r *NoteRepository) ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]note.Note, error) {
	const query = `
		SELECT id, case_file_id, title, content, created_at, updated_at
		FROM notes
		WHERE case_file_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, caseFileID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []note.Note

	for rows.Next() {
		var rawID, cfID, title, content, createdAt, updatedAt string

		if err := rows.Scan(&rawID, &cfID, &title, &content, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		ct, _ := time.Parse(time.RFC3339, createdAt)
		ut, _ := time.Parse(time.RFC3339, updatedAt)

		result = append(result, note.Note{
			ID:         shared.NewID(rawID),
			CaseFileID: shared.NewID(cfID),
			Title:      title,
			Content:    content,
			CreatedAt:  shared.Timestamp(ct),
			UpdatedAt:  shared.Timestamp(ut),
		})
	}

	return result, rows.Err()
}

func (r *NoteRepository) Delete(ctx context.Context, id shared.ID) error {
	const query = `DELETE FROM notes WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}
