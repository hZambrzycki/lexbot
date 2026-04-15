package client

import (
	"lexbox/internal/domain/shared"
)

type Client struct {
	ID         shared.ID
	Name       string
	Email      string
	Phone      string
	Identifier string

	CreatedAt shared.Timestamp
	UpdatedAt shared.Timestamp
}

func NewClient(id shared.ID, name string) (Client, error) {
	if name == "" {
		return Client{}, shared.ErrEmptyField
	}

	now := shared.Now()

	return Client{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
