package casefileapp

import (
	"context"
	"strings"

	documentapp "lexbox/internal/application/document"
	"lexbox/internal/domain/shared"
)

type FixCaseFileInput struct {
	CaseFileID string
}

type FixCaseFileResult struct {
	CaseFileID string

	ReindexResult  documentapp.ReindexAllDocumentsResult
	MetadataResult documentapp.AnalyzeAllDocumentMetadataResult
	EventsResult   documentapp.AnalyzeAllDocumentEventsResult

	NeedsAttention bool
	TopAlert       string
}

type FixCaseFile struct {
	ReindexAllDocuments documentapp.ReindexAllDocuments
	AnalyzeAllMetadata  documentapp.AnalyzeAllDocumentMetadata
	AnalyzeAllEvents    documentapp.AnalyzeAllDocumentEvents
	GetDashboard        GetCaseFileDashboard
}

func (uc FixCaseFile) Execute(ctx context.Context, in FixCaseFileInput) (FixCaseFileResult, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return FixCaseFileResult{}, shared.ErrInvalidID
	}

	reindexResult, err := uc.ReindexAllDocuments.Execute(ctx, documentapp.ReindexAllDocumentsInput{
		CaseFileID: caseFileID.String(),
	})
	if err != nil {
		return FixCaseFileResult{}, err
	}

	metadataResult, err := uc.AnalyzeAllMetadata.Execute(ctx, documentapp.AnalyzeAllDocumentMetadataInput{
		CaseFileID: caseFileID.String(),
	})
	if err != nil {
		return FixCaseFileResult{}, err
	}

	eventsResult, err := uc.AnalyzeAllEvents.Execute(ctx, documentapp.AnalyzeAllDocumentEventsInput{
		CaseFileID: caseFileID.String(),
	})
	if err != nil {
		return FixCaseFileResult{}, err
	}

	result := FixCaseFileResult{
		CaseFileID:     caseFileID.String(),
		ReindexResult:  reindexResult,
		MetadataResult: metadataResult,
		EventsResult:   eventsResult,
	}

	dashboard, err := uc.GetDashboard.Execute(ctx, GetCaseFileDashboardInput{
		CaseFileID: caseFileID.String(),
	})
	if err != nil {
		return FixCaseFileResult{}, err
	}

	result.NeedsAttention = dashboard.NeedsAttention
	result.TopAlert = dashboard.TopAlert

	return result, nil
}
