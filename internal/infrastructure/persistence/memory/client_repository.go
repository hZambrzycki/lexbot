package memory

import (
	"context"
	"sync"

	"lexbox/internal/domain/client"
	"lexbox/internal/domain/shared"
)

type ClientRepository struct {
	mu      sync.RWMutex
	storage map[shared.ID]client.Client
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{
		storage: make(map[shared.ID]client.Client),
	}
}

func (r *ClientRepository) Save(ctx context.Context, c client.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[c.ID] = c
	return nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id shared.ID) (client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.storage[id]
	if !ok {
		return client.Client{}, shared.ErrNotFound
	}

	return c, nil
}

func (r *ClientRepository) List(ctx context.Context) ([]client.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []client.Client
	for _, c := range r.storage {
		result = append(result, c)
	}

	return result, nil
}
