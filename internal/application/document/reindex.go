package documentapp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"lexbox/internal/domain/shared"
)

var (
	ErrReindexFileNotFound = errors.New("document file not found")
	ErrReindexUnsupported  = errors.New("document file type is not supported for extraction")
	ErrReindexEmptyContent = errors.New("document extraction returned empty content")
	ErrReindexReadFailed   = errors.New("document file could not be read")
)

type ReindexDocumentInput struct {
	DocumentID string
}

type ReindexDocumentResult struct {
	DocumentID      string
	ExtractedLength int
}

type ReindexDocument struct {
	ExtractDocumentText ExtractDocumentText
}

func (uc ReindexDocument) Execute(ctx context.Context, in ReindexDocumentInput) (ReindexDocumentResult, error) {
	documentID := strings.TrimSpace(in.DocumentID)
	if documentID == "" {
		return ReindexDocumentResult{}, shared.ErrInvalidID
	}

	content, err := uc.ExtractDocumentText.Execute(ctx, ExtractDocumentTextInput{
		DocumentID: documentID,
	})
	if err != nil {
		return ReindexDocumentResult{}, classifyReindexError(err)
	}

	return ReindexDocumentResult{
		DocumentID:      documentID,
		ExtractedLength: len(content),
	}, nil
}

func classifyReindexError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %v", ErrReindexFileNotFound, err)
	}

	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "no such file"),
		strings.Contains(msg, "file not found"),
		strings.Contains(msg, "cannot read file: open"):
		return fmt.Errorf("%w: %v", ErrReindexFileNotFound, err)

	case strings.Contains(msg, "unsupported file type"),
		strings.Contains(msg, "not supported for extraction"):
		return fmt.Errorf("%w: %v", ErrReindexUnsupported, err)

	case strings.Contains(msg, "extracted content is empty"):
		return fmt.Errorf("%w: %v", ErrReindexEmptyContent, err)

	case strings.Contains(msg, "cannot read file"),
		strings.Contains(msg, "permission denied"),
		strings.Contains(msg, "input/output error"),
		strings.Contains(msg, "pdf extraction failed"),
		strings.Contains(msg, "cannot open docx file"),
		strings.Contains(msg, "cannot open docx document.xml"),
		strings.Contains(msg, "cannot read docx document.xml"),
		strings.Contains(msg, "docx extraction failed"):
		return fmt.Errorf("%w: %v", ErrReindexReadFailed, err)

	default:
		return fmt.Errorf("%w: %v", ErrReindexReadFailed, err)
	}
}
