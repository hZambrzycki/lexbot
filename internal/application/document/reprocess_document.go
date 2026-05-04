package documentapp

import "context"

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
	ReindexDocument         ReindexDocument
	AnalyzeDocumentMetadata AnalyzeDocumentMetadata
	AnalyzeDocumentEvents   AnalyzeDocumentEvents
}

func (uc ReprocessDocument) Execute(ctx context.Context, in ReprocessDocumentInput) (ReprocessDocumentResult, error) {
	reindexResult, err := uc.ReindexDocument.Execute(ctx, ReindexDocumentInput{
		DocumentID: in.DocumentID,
	})
	if err != nil {
		return ReprocessDocumentResult{}, err
	}

	_, err = uc.AnalyzeDocumentMetadata.Execute(ctx, AnalyzeDocumentMetadataInput{
		DocumentID: in.DocumentID,
	})
	if err != nil {
		return ReprocessDocumentResult{}, err
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
