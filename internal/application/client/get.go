package clientapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/client"
	"lexbox/internal/domain/shared"
)

type GetClientInput struct {
	ID string
}

type GetClient struct {
	Clients ports.ClientRepository
}

func (uc GetClient) Execute(ctx context.Context, in GetClientInput) (client.Client, error) {
	id := shared.NewID(strings.TrimSpace(in.ID))
	if id == "" {
		return client.Client{}, shared.ErrInvalidID
	}

	return uc.Clients.GetByID(ctx, id)
}
