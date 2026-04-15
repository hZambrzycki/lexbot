package casefileapp

import (
	"context"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/note"
	"lexbox/internal/domain/shared"
)

type GetCaseFileDetailInput struct {
	ID string
}

type CaseFileDetail struct {
	CaseFile  casefile.CaseFile
	Notes     []note.Note
	Documents []document.Document
}

type GetCaseFileDetail struct {
	CaseFiles ports.CaseFileRepository
	Notes     ports.NoteRepository
	Documents ports.DocumentRepository
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

	return CaseFileDetail{
		CaseFile:  cf,
		Notes:     notes,
		Documents: documents,
	}, nil
}
