package documentapp

import (
	"context"
	"errors"
	"os"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type GetDocumentDetailInput struct {
	DocumentID string
}

type EventTrace struct {
	EventType      string
	EventDate      string
	SourceText     string
	AnchorDate     string
	DateKind       string
	AnchorSource   string
	RelativeDays   int
	IsBusinessDays bool
	AddExtraDay    bool
	CalendarScope  string
	TriggerText    string
	Computation    string
}

type GetDocumentDetailResult struct {
	Document             document.Document
	FileExists           bool
	HasExtractedText     bool
	ExtractedText        string
	ExtractedTextLength  int
	ExtractedTextPreview string

	HasMetadata        bool
	DocumentType       string
	LegalArea          string
	MetadataAnalyzedAt string

	HasEvents bool
	Events    []document.Event

	HasEventTrace bool
	EventTrace    []EventTrace
}

type GetDocumentDetail struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	Metadata         ports.DocumentMetadataRepository
	Events           ports.DocumentEventRepository
}

func (uc GetDocumentDetail) Execute(ctx context.Context, in GetDocumentDetailInput) (GetDocumentDetailResult, error) {
	documentID := shared.NewID(strings.TrimSpace(in.DocumentID))
	if documentID == "" {
		return GetDocumentDetailResult{}, shared.ErrInvalidID
	}

	doc, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return GetDocumentDetailResult{}, err
	}

	fileExists := false
	if _, err := os.Stat(doc.StoragePath); err == nil {
		fileExists = true
	}

	result := GetDocumentDetailResult{
		Document:             doc,
		FileExists:           fileExists,
		HasExtractedText:     false,
		ExtractedText:        "",
		ExtractedTextLength:  0,
		ExtractedTextPreview: "",
		HasMetadata:          false,
		DocumentType:         "",
		LegalArea:            "",
		MetadataAnalyzedAt:   "",
		HasEvents:            false,
		Events:               []document.Event{},
		HasEventTrace:        false,
		EventTrace:           []EventTrace{},
	}

	if uc.DocumentContents != nil {
		content, err := uc.DocumentContents.GetByDocumentID(ctx, documentID.String())
		if err != nil {
			if !errors.Is(err, shared.ErrNotFound) {
				return GetDocumentDetailResult{}, err
			}
		} else {
			result.HasExtractedText = true
			result.ExtractedText = content
			result.ExtractedTextLength = len(content)
			result.ExtractedTextPreview = buildPreview(content, 160)
		}
	}

	if uc.Metadata != nil {
		metadata, err := uc.Metadata.GetByDocumentID(ctx, documentID)
		if err != nil {
			if !errors.Is(err, shared.ErrNotFound) {
				return GetDocumentDetailResult{}, err
			}
		} else {
			result.HasMetadata = true
			result.DocumentType = metadata.DocumentType
			result.LegalArea = metadata.LegalArea
			result.MetadataAnalyzedAt = metadata.AnalyzedAt
		}
	}

	if uc.Events != nil {
		events, err := uc.Events.ListByDocumentID(ctx, documentID)
		if err != nil {
			return GetDocumentDetailResult{}, err
		}
		if len(events) > 0 {
			result.HasEvents = true
			result.Events = events
			result.HasEventTrace = true
			result.EventTrace = buildEventTraceFromPersistedEvents(events)
		}
	}

	return result, nil
}

func buildEventTraceFromPersistedEvents(events []document.Event) []EventTrace {
	trace := make([]EventTrace, 0, len(events))
	for _, e := range events {
		trace = append(trace, EventTrace{
			EventType:      e.EventType,
			EventDate:      e.EventDate,
			SourceText:     e.SourceText,
			AnchorDate:     e.AnchorDate,
			DateKind:       e.DateKind,
			AnchorSource:   e.AnchorSource,
			RelativeDays:   e.RelativeDays,
			IsBusinessDays: e.IsBusinessDays,
			AddExtraDay:    e.AddExtraDay,
			CalendarScope:  e.CalendarScope,
			TriggerText:    e.TriggerText,
			Computation:    e.Computation,
		})
	}
	return trace
}

func buildPreview(content string, maxLen int) string {
	normalized := strings.Join(strings.Fields(content), " ")
	if normalized == "" {
		return ""
	}

	if len(normalized) <= maxLen {
		return normalized
	}

	return normalized[:maxLen] + "..."
}
