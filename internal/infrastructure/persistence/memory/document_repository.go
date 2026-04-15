package memory

import (
	"context"
	"sync"

	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type DocumentRepository struct {
	mu      sync.RWMutex
	storage map[shared.ID]document.Document
}

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		storage: make(map[shared.ID]document.Document),
	}
}

func (r *DocumentRepository) Save(ctx context.Context, d document.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[d.ID] = d
	return nil
}

func (r *DocumentRepository) ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]document.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []document.Document
	for _, d := range r.storage {
		if d.CaseFileID == caseFileID {
			result = append(result, d)
		}
	}

	return result, nil
}
