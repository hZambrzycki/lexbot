package clientapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/client"
)

type CreateClientInput struct {
	Name       string
	Email      string
	Phone      string
	Identifier string
}

type CreateClient struct {
	Clients ports.ClientRepository
	IDs     ports.IDGenerator
}

func (uc CreateClient) Execute(ctx context.Context, in CreateClientInput) (client.Client, error) {
	name := strings.TrimSpace(in.Name)

	c, err := client.NewClient(uc.IDs.NewID(), name)
	if err != nil {
		return client.Client{}, err
	}

	c.Email = strings.TrimSpace(in.Email)
	c.Phone = strings.TrimSpace(in.Phone)
	c.Identifier = strings.TrimSpace(in.Identifier)

	if err := uc.Clients.Save(ctx, c); err != nil {
		return client.Client{}, err
	}

	return c, nil
}
