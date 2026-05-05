package casefileapp

import (
	"context"
	"errors"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/note"
	"lexbox/internal/domain/shared"
)

type GetCaseFileDetailInput struct {
	ID string
}

type DocumentSummary struct {
	Document document.Document

	HasExtractedText bool
	HasMetadata      bool
	HasEvents        bool

	DocumentType string
	LegalArea    string

	EventCount int
}

type CaseFileDetail struct {
	CaseFile  casefile.CaseFile
	Notes     []note.Note
	Documents []DocumentSummary
}

type GetCaseFileDetail struct {
	CaseFiles        ports.CaseFileRepository
	Notes            ports.NoteRepository
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	Metadata         ports.DocumentMetadataRepository
	Events           ports.DocumentEventRepository
}

func (uc GetCaseFileDetail) Execute(ctx context.Context, in GetCaseFileDetailInput) (CaseFileDetail, error) {
	id := shared.NewID(in.ID)
	if id == "" {
		return CaseFileDetail{}, shared.ErrInvalidID
	}

	cf, err := uc.CaseFiles.GetByID(ctx, id)
	if err != nil {
		return CaseFileDetail{}, err
	}

	notes, err := uc.Notes.ListByCaseFileID(ctx, id)
	if err != nil {
		return CaseFileDetail{}, err
	}

	documents, err := uc.Documents.ListByCaseFileID(ctx, id)
	if err != nil {
		return CaseFileDetail{}, err
	}

	summaries := make([]DocumentSummary, 0, len(documents))

	for _, doc := range documents {
		summary := DocumentSummary{
			Document:     doc,
			DocumentType: "unknown",
			LegalArea:    "unknown",
		}

		if uc.DocumentContents != nil {
			content, err := uc.DocumentContents.GetByDocumentID(ctx, doc.ID.String())
			if err == nil && strings.TrimSpace(content) != "" {
				summary.HasExtractedText = true
			} else if err != nil && !errors.Is(err, shared.ErrNotFound) {
				return CaseFileDetail{}, err
			}
		}

		if uc.Metadata != nil {
			meta, err := uc.Metadata.GetByDocumentID(ctx, doc.ID)
			if err == nil {
				summary.HasMetadata = true
				summary.DocumentType = meta.DocumentType
				summary.LegalArea = meta.LegalArea
			} else if err != nil && !errors.Is(err, shared.ErrNotFound) {
				return CaseFileDetail{}, err
			}
		}

		if uc.Events != nil {
			events, err := uc.Events.ListByDocumentID(ctx, doc.ID)
			if err != nil {
				return CaseFileDetail{}, err
			}

			summary.EventCount = len(events)
			summary.HasEvents = len(events) > 0
		}

		summaries = append(summaries, summary)
	}

	return CaseFileDetail{
		CaseFile:  cf,
		Notes:     notes,
		Documents: summaries,
	}, nil
}
