package casefileapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/shared"
)

type ListCaseFiles struct {
	CaseFiles ports.CaseFileRepository
}

func (uc ListCaseFiles) Execute(ctx context.Context) ([]casefile.CaseFile, error) {
	return uc.CaseFiles.List(ctx)
}

type ListCaseFilesByClientInput struct {
	ClientID string
}

type ListCaseFilesByClient struct {
	CaseFiles ports.CaseFileRepository
}

func (uc ListCaseFilesByClient) Execute(ctx context.Context, in ListCaseFilesByClientInput) ([]casefile.CaseFile, error) {
	clientID := shared.NewID(strings.TrimSpace(in.ClientID))
	if clientID == "" {
		return nil, shared.ErrInvalidID
	}

	return uc.CaseFiles.ListByClientID(ctx, clientID)
}
