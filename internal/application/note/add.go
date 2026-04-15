package noteapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/note"
	"lexbox/internal/domain/shared"
)

type AddNoteInput struct {
	CaseFileID string
	Title      string
	Content    string
}

type AddNote struct {
	Notes     ports.NoteRepository
	CaseFiles ports.CaseFileRepository
	IDs       ports.IDGenerator
}

func (uc AddNote) Execute(ctx context.Context, in AddNoteInput) (note.Note, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return note.Note{}, shared.ErrInvalidAssociation
	}

	_, err := uc.CaseFiles.GetByID(ctx, caseFileID)
	if err != nil {
		return note.Note{}, err
	}

	n, err := note.NewNote(
		uc.IDs.NewID(),
		caseFileID,
		strings.TrimSpace(in.Title),
		strings.TrimSpace(in.Content),
	)
	if err != nil {
		return note.Note{}, err
	}

	if err := uc.Notes.Save(ctx, n); err != nil {
		return note.Note{}, err
	}

	return n, nil
}
