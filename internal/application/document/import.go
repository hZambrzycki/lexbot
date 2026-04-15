package documentapp

import (
	"context"
	"errors"
	"strings"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type DuplicateDocumentError struct {
	Existing document.Document
}

func (e DuplicateDocumentError) Error() string {
	return "document with same file hash already exists in case file"
}

type ImportDocumentInput struct {
	CaseFileID string
	SourcePath string
	MimeType   string
	FileHash   string
}

type ImportDocumentResult struct {
	Document         document.Document
	TextExtracted    bool
	MetadataAnalyzed bool
	EventsAnalyzed   bool
	EventsDetected   int
}

type ImportDocument struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	Metadata         ports.DocumentMetadataRepository
	CaseFiles        ports.CaseFileRepository
	Storage          ports.FileStorage
	Extractor        ports.TextExtractor
	Hasher           ports.FileHasher
	IDs              ports.IDGenerator
	AnalyzeEvents    AnalyzeDocumentEvents
}

func (uc ImportDocument) Execute(ctx context.Context, in ImportDocumentInput) (ImportDocumentResult, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return ImportDocumentResult{}, shared.ErrInvalidAssociation
	}

	_, err := uc.CaseFiles.GetByID(ctx, caseFileID)
	if err != nil {
		return ImportDocumentResult{}, err
	}

	documentID := uc.IDs.NewID()

	stored, err := uc.Storage.StoreDocument(ctx, documentID.String(), strings.TrimSpace(in.SourcePath))
	if err != nil {
		return ImportDocumentResult{}, err
	}

	fileHash := strings.TrimSpace(in.FileHash)
	if fileHash == "" && uc.Hasher != nil {
		fileHash, err = uc.Hasher.HashFile(ctx, stored.StoragePath)
		if err != nil {
			_ = uc.Storage.DeleteDocument(ctx, documentID.String())
			return ImportDocumentResult{}, err
		}
	}

	if fileHash != "" {
		existing, err := uc.Documents.GetByCaseFileIDAndFileHash(ctx, caseFileID, fileHash)
		if err == nil {
			_ = uc.Storage.DeleteDocument(ctx, documentID.String())
			return ImportDocumentResult{}, DuplicateDocumentError{
				Existing: existing,
			}
		}
		if !errors.Is(err, shared.ErrNotFound) {
			_ = uc.Storage.DeleteDocument(ctx, documentID.String())
			return ImportDocumentResult{}, err
		}
	}

	doc, err := document.NewDocument(
		documentID,
		caseFileID,
		stored.OriginalName,
		stored.StoragePath,
	)
	if err != nil {
		_ = uc.Storage.DeleteDocument(ctx, documentID.String())
		return ImportDocumentResult{}, err
	}

	doc.MimeType = normalizeMimeType(in.MimeType)
	doc.FileHash = fileHash

	if err := uc.Documents.Save(ctx, doc); err != nil {
		_ = uc.Storage.DeleteDocument(ctx, documentID.String())
		return ImportDocumentResult{}, err
	}

	textExtracted := false
	metadataAnalyzed := false
	eventsAnalyzed := false
	eventsDetected := 0

	if uc.Extractor != nil && uc.DocumentContents != nil {
		content, err := uc.Extractor.ExtractText(ctx, doc.StoragePath)
		if err != nil {
			if isUnsupportedExtractionError(err) || isEmptyExtractionError(err) {
				// No indexamos, pero tampoco rompemos la importación.
			} else {
				_ = uc.Storage.DeleteDocument(ctx, documentID.String())
				return ImportDocumentResult{}, err
			}
		} else {
			if err := uc.DocumentContents.Save(ctx, doc.ID.String(), content); err != nil {
				_ = uc.Storage.DeleteDocument(ctx, documentID.String())
				return ImportDocumentResult{}, err
			}

			textExtracted = true

			if strings.TrimSpace(content) != "" {
				if uc.Metadata != nil {
					classification := classifyDocumentMetadata(content)

					metadata, err := document.NewMetadata(
						doc.ID,
						classification.DocumentType,
						classification.LegalArea,
						shared.Now().Time().Format(time.RFC3339),
					)
					if err == nil {
						if err := uc.Metadata.Save(ctx, metadata); err == nil {
							metadataAnalyzed = true
						}
					}
				}

				if uc.AnalyzeEvents.Documents != nil &&
					uc.AnalyzeEvents.DocumentContents != nil &&
					uc.AnalyzeEvents.Events != nil &&
					uc.AnalyzeEvents.IDs != nil {

					analyzeResult, err := uc.AnalyzeEvents.Execute(ctx, AnalyzeDocumentEventsInput{
						DocumentID: doc.ID.String(),
					})
					if err != nil {
						_ = uc.Storage.DeleteDocument(ctx, documentID.String())
						return ImportDocumentResult{}, err
					}

					eventsAnalyzed = true
					eventsDetected = analyzeResult.Detected
				}
			}
		}
	}

	return ImportDocumentResult{
		Document:         doc,
		TextExtracted:    textExtracted,
		MetadataAnalyzed: metadataAnalyzed,
		EventsAnalyzed:   eventsAnalyzed,
		EventsDetected:   eventsDetected,
	}, nil
}

func isUnsupportedExtractionError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(classifyReindexError(err), ErrReindexUnsupported)
}

func isEmptyExtractionError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(strings.ToLower(err.Error()), "extracted content is empty")
}
