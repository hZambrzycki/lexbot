package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/shared"
)

type ListEventsByCaseFileInput struct {
	CaseFileID string
}

type ListEventsByCaseFile struct {
	Events ports.DocumentEventRepository
}

func (uc ListEventsByCaseFile) Execute(ctx context.Context, in ListEventsByCaseFileInput) ([]querymodels.CaseFileEventResult, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return nil, shared.ErrInvalidID
	}

	return uc.Events.ListByCaseFileID(ctx, caseFileID)
}
