package clientapp

import (
	"context"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/client"
)

type ListClients struct {
	Clients ports.ClientRepository
}

func (uc ListClients) Execute(ctx context.Context) ([]client.Client, error) {
	return uc.Clients.List(ctx)
}
