package documentapp

import (
	"context"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/shared"
)

type ReprocessDocumentInput struct {
	DocumentID string
}

type ReprocessDocumentResult struct {
	DocumentID       string `json:"document_id"`
	TextReindexed    bool   `json:"text_reindexed"`
	ExtractedLength  int    `json:"extracted_length"`
	MetadataAnalyzed bool   `json:"metadata_analyzed"`
	EventsAnalyzed   bool   `json:"events_analyzed"`
	EventsDetected   int    `json:"events_detected"`
}

type ReprocessDocument struct {
	Documents               ports.DocumentRepository
	DocumentContents        ports.DocumentContentRepository
	SearchIndex             ports.DocumentSearchIndexRepository
	ReindexDocument         ReindexDocument
	AnalyzeDocumentMetadata AnalyzeDocumentMetadata
	AnalyzeDocumentEvents   AnalyzeDocumentEvents
}

func (uc ReprocessDocument) Execute(ctx context.Context, in ReprocessDocumentInput) (ReprocessDocumentResult, error) {
	documentID := shared.NewID(in.DocumentID)
	if documentID == "" {
		return ReprocessDocumentResult{}, shared.ErrInvalidID
	}

	doc, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return ReprocessDocumentResult{}, err
	}

	reindexResult, err := uc.ReindexDocument.Execute(ctx, ReindexDocumentInput{
		DocumentID: in.DocumentID,
	})
	if err != nil {
		return ReprocessDocumentResult{}, err
	}

	metadataResult, err := uc.AnalyzeDocumentMetadata.Execute(ctx, AnalyzeDocumentMetadataInput{
		DocumentID: in.DocumentID,
	})
	if err != nil {
		return ReprocessDocumentResult{}, err
	}

	if uc.SearchIndex != nil && uc.DocumentContents != nil {
		content, err := uc.DocumentContents.GetByDocumentID(ctx, doc.ID.String())
		if err != nil {
			return ReprocessDocumentResult{}, err
		}

		if err := uc.SearchIndex.UpsertDocument(
			ctx,
			doc.ID.String(),
			doc.CaseFileID.String(),
			doc.OriginalName,
			content,
			metadataResult.DocumentType,
			metadataResult.LegalArea,
		); err != nil {
			return ReprocessDocumentResult{}, err
		}
	}

	eventsResult, err := uc.AnalyzeDocumentEvents.Execute(ctx, AnalyzeDocumentEventsInput{
		DocumentID: in.DocumentID,
	})
	if err != nil {
		return ReprocessDocumentResult{}, err
	}

	return ReprocessDocumentResult{
		DocumentID:       reindexResult.DocumentID,
		TextReindexed:    true,
		ExtractedLength:  reindexResult.ExtractedLength,
		MetadataAnalyzed: true,
		EventsAnalyzed:   true,
		EventsDetected:   eventsResult.Detected,
	}, nil
}
