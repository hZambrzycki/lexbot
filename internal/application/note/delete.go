package noteapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/shared"
)

type DeleteNoteInput struct {
	ID string
}

type DeleteNoteOutput struct {
	ID      string
	Deleted bool
}

type DeleteNote struct {
	Notes ports.NoteRepository
}

func (uc DeleteNote) Execute(ctx context.Context, in DeleteNoteInput) (DeleteNoteOutput, error) {
	id := shared.NewID(strings.TrimSpace(in.ID))
	if id.String() == "" {
		return DeleteNoteOutput{}, shared.ErrInvalidID
	}

	if err := uc.Notes.Delete(ctx, id); err != nil {
		return DeleteNoteOutput{}, err
	}

	return DeleteNoteOutput{
		ID:      id.String(),
		Deleted: true,
	}, nil
}
