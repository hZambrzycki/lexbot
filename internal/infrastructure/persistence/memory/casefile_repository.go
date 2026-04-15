package memory

import (
	"context"
	"sync"

	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/shared"
)

type CaseFileRepository struct {
	mu      sync.RWMutex
	storage map[shared.ID]casefile.CaseFile
}

func NewCaseFileRepository() *CaseFileRepository {
	return &CaseFileRepository{
		storage: make(map[shared.ID]casefile.CaseFile),
	}
}

func (r *CaseFileRepository) Save(ctx context.Context, cf casefile.CaseFile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[cf.ID] = cf
	return nil
}

func (r *CaseFileRepository) GetByID(ctx context.Context, id shared.ID) (casefile.CaseFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cf, ok := r.storage[id]
	if !ok {
		return casefile.CaseFile{}, shared.ErrNotFound
	}

	return cf, nil
}

func (r *CaseFileRepository) List(ctx context.Context) ([]casefile.CaseFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []casefile.CaseFile
	for _, cf := range r.storage {
		result = append(result, cf)
	}

	return result, nil
}

func (r *CaseFileRepository) ListByClientID(ctx context.Context, clientID shared.ID) ([]casefile.CaseFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []casefile.CaseFile
	for _, cf := range r.storage {
		if cf.ClientID == clientID {
			result = append(result, cf)
		}
	}

	return result, nil
}
