package memory

import (
	"context"
	"sync"

	"lexbox/internal/domain/note"
	"lexbox/internal/domain/shared"
)

type NoteRepository struct {
	mu      sync.RWMutex
	storage map[shared.ID]note.Note
}

func NewNoteRepository() *NoteRepository {
	return &NoteRepository{
		storage: make(map[shared.ID]note.Note),
	}
}

func (r *NoteRepository) Save(ctx context.Context, n note.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[n.ID] = n
	return nil
}

func (r *NoteRepository) ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]note.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []note.Note
	for _, n := range r.storage {
		if n.CaseFileID == caseFileID {
			result = append(result, n)
		}
	}

	return result, nil
}

func (r *NoteRepository) Delete(ctx context.Context, id shared.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.storage, id)
	return nil
}
