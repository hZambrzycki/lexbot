package noteapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/note"
	"lexbox/internal/domain/shared"
)

type ListNotesByCaseFileInput struct {
	CaseFileID string
}

type ListNotesByCaseFile struct {
	Notes ports.NoteRepository
}

func (uc ListNotesByCaseFile) Execute(ctx context.Context, in ListNotesByCaseFileInput) ([]note.Note, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return nil, shared.ErrInvalidID
	}

	return uc.Notes.ListByCaseFileID(ctx, caseFileID)
}
