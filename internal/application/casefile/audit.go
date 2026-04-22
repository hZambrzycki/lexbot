package casefileapp

import (
	"context"
	"strings"

	documentapp "lexbox/internal/application/document"
	"lexbox/internal/domain/shared"
)

type AuditCaseFileInput struct {
	CaseFileID string
}

type AuditCaseFileResult struct {
	CaseFileID string

	VerifyResult documentapp.VerifyAllDocumentsResult

	DocumentsWithUnknownMetadata     int
	DocumentsWithUnknownMetadataList []string

	DocumentsWithoutEvents     int
	DocumentsWithoutEventsList []string

	IsHealthy          bool
	TopIssue           string
	RecommendedActions []string
}

type AuditCaseFile struct {
	VerifyAllDocuments documentapp.VerifyAllDocuments
	GetDashboard       GetCaseFileDashboard
}

func (uc AuditCaseFile) Execute(ctx context.Context, in AuditCaseFileInput) (AuditCaseFileResult, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return AuditCaseFileResult{}, shared.ErrInvalidID
	}

	verifyResult, err := uc.VerifyAllDocuments.Execute(ctx, documentapp.VerifyAllDocumentsInput{
		CaseFileID: caseFileID.String(),
	})
	if err != nil {
		return AuditCaseFileResult{}, err
	}

	dashboard, err := uc.GetDashboard.Execute(ctx, GetCaseFileDashboardInput{
		CaseFileID: caseFileID.String(),
	})
	if err != nil {
		return AuditCaseFileResult{}, err
	}

	result := AuditCaseFileResult{
		CaseFileID: caseFileID.String(),

		VerifyResult: verifyResult,

		DocumentsWithUnknownMetadata:     dashboard.DocumentsWithUnknownMetadata,
		DocumentsWithUnknownMetadataList: dashboard.DocumentsWithUnknownMetadataList,

		DocumentsWithoutEvents:     dashboard.DocumentsWithoutEvents,
		DocumentsWithoutEventsList: dashboard.DocumentsWithoutEventsList,
	}

	result.IsHealthy, result.TopIssue = buildAuditStatus(result)
	result.RecommendedActions = buildAuditRecommendedActions(result)

	return result, nil
}

func buildAuditStatus(result AuditCaseFileResult) (bool, string) {
	switch {
	case result.VerifyResult.MissingFile > 0:
		return false, "some documents are missing from storage"
	case result.VerifyResult.MissingText > 0:
		return false, "some documents have no extracted text"
	case result.VerifyResult.InvalidMime > 0:
		return false, "some documents have invalid mime types"
	case result.DocumentsWithUnknownMetadata > 0:
		return false, "some documents still have unknown metadata"
	case result.DocumentsWithoutEvents > 0:
		return false, "some documents have no detected events"
	default:
		return true, "no technical issues detected"
	}
}

func buildAuditRecommendedActions(result AuditCaseFileResult) []string {
	actions := make([]string, 0)

	switch {
	case result.VerifyResult.MissingFile > 0:
		actions = append(actions, "review missing files in storage")
	case result.VerifyResult.MissingText > 0:
		actions = append(actions, "run case-fix to reindex and extract text")
	case result.VerifyResult.InvalidMime > 0:
		actions = append(actions, "run mime-normalize and review invalid mime types")
	}

	if result.DocumentsWithUnknownMetadata > 0 {
		actions = append(actions, "review unknown metadata classifications")
	}

	if result.DocumentsWithoutEvents > 0 {
		actions = append(actions, "review documents without detected events")
	}

	if len(actions) == 0 {
		actions = append(actions, "no action needed")
	}

	return actions
}
