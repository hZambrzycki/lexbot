package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type ListEventsByCaseFileInput struct {
	CaseFileID   string
	ReviewStatus string
}

type ListEventsByCaseFile struct {
	Events ports.DocumentEventRepository
}

func (uc ListEventsByCaseFile) Execute(ctx context.Context, in ListEventsByCaseFileInput) ([]querymodels.CaseFileEventResult, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return nil, shared.ErrInvalidID
	}

	reviewStatus := normalizeReviewStatusFilter(in.ReviewStatus)
	if reviewStatus != "" && !document.IsValidReviewStatus(reviewStatus) {
		return nil, shared.ErrInvalidAssociation
	}

	return uc.Events.ListByCaseFileID(ctx, caseFileID, reviewStatus)
}

func normalizeReviewStatusFilter(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
