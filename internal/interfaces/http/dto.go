package httpapi

import (
	casefileapp "lexbox/internal/application/casefile"
	documentapp "lexbox/internal/application/document"
	"lexbox/internal/application/querymodels"
	domaincasefile "lexbox/internal/domain/casefile"
	domaindocument "lexbox/internal/domain/document"
	domainnote "lexbox/internal/domain/note"
)

type CaseFileResponse struct {
	ID                string `json:"id"`
	ClientID          string `json:"client_id"`
	Reference         string `json:"reference"`
	Title             string `json:"title"`
	Type              string `json:"type"`
	Status            string `json:"status"`
	Description       string `json:"description"`
	CalendarScope     string `json:"calendar_scope"`
	AugustNonBusiness bool   `json:"august_non_business"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type NoteResponse struct {
	ID         string `json:"id"`
	CaseFileID string `json:"case_file_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type DocumentDetailResponse struct {
	Document             DocumentResponse      `json:"document"`
	FileExists           bool                  `json:"file_exists"`
	HasExtractedText     bool                  `json:"has_extracted_text"`
	ExtractedText        string                `json:"extracted_text"`
	ExtractedTextLength  int                   `json:"extracted_text_length"`
	ExtractedTextPreview string                `json:"extracted_text_preview"`
	HasMetadata          bool                  `json:"has_metadata"`
	DocumentType         string                `json:"document_type"`
	LegalArea            string                `json:"legal_area"`
	MetadataAnalyzedAt   string                `json:"metadata_analyzed_at"`
	HasEvents            bool                  `json:"has_events"`
	Events               []DocumentEventDetail `json:"events"`
}

type DocumentEventDetail struct {
	EventID        string `json:"event_id"`
	DocumentID     string `json:"document_id"`
	EventType      string `json:"event_type"`
	EventDate      string `json:"event_date"`
	SourceText     string `json:"source_text"`
	CreatedAt      string `json:"created_at"`
	AnchorDate     string `json:"anchor_date,omitempty"`
	DateKind       string `json:"date_kind,omitempty"`
	AnchorSource   string `json:"anchor_source,omitempty"`
	RelativeDays   int    `json:"relative_days,omitempty"`
	IsBusinessDays bool   `json:"is_business_days,omitempty"`
	AddExtraDay    bool   `json:"add_extra_day,omitempty"`
	CalendarScope  string `json:"calendar_scope,omitempty"`
	TriggerText    string `json:"trigger_text,omitempty"`
	Computation    string `json:"computation,omitempty"`
	ReviewStatus   string `json:"review_status"`
	ReviewedAt     string `json:"reviewed_at,omitempty"`
	ResolvedAt     string `json:"resolved_at,omitempty"`
	ResolutionNote string `json:"resolution_note,omitempty"`
}

type DocumentResponse struct {
	ID           string `json:"id"`
	CaseFileID   string `json:"case_file_id"`
	OriginalName string `json:"original_name"`
	StoragePath  string `json:"storage_path"`
	MimeType     string `json:"mime_type"`
	FileHash     string `json:"file_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CaseFileDetailResponse struct {
	CaseFile  CaseFileResponse   `json:"case_file"`
	Notes     []NoteResponse     `json:"notes"`
	Documents []DocumentResponse `json:"documents"`
}

type EventResponse struct {
	EventID      string `json:"event_id"`
	DocumentID   string `json:"document_id,omitempty"`
	OriginalName string `json:"original_name,omitempty"`

	CaseFileID        string `json:"case_file_id,omitempty"`
	CaseFileReference string `json:"case_file_reference,omitempty"`
	CaseFileTitle     string `json:"case_file_title,omitempty"`

	EventType      string   `json:"event_type"`
	EventDate      string   `json:"event_date"`
	SourceText     string   `json:"source_text"`
	DaysRemaining  int      `json:"days_remaining,omitempty"`
	Status         string   `json:"status,omitempty"`
	Priority       string   `json:"priority,omitempty"`
	DuplicateCount int      `json:"duplicate_count,omitempty"`
	DocumentNames  []string `json:"document_names,omitempty"`
	DocumentIDs    []string `json:"document_ids,omitempty"`
	AnchorDate     string   `json:"anchor_date,omitempty"`
	DateKind       string   `json:"date_kind,omitempty"`
	AnchorSource   string   `json:"anchor_source,omitempty"`
	RelativeDays   int      `json:"relative_days,omitempty"`
	IsBusinessDays bool     `json:"is_business_days,omitempty"`
	AddExtraDay    bool     `json:"add_extra_day,omitempty"`
	CalendarScope  string   `json:"calendar_scope,omitempty"`
	TriggerText    string   `json:"trigger_text,omitempty"`
	Computation    string   `json:"computation,omitempty"`
	ReviewStatus   string   `json:"review_status"`
	ReviewedAt     string   `json:"reviewed_at,omitempty"`
	ResolvedAt     string   `json:"resolved_at,omitempty"`
	ResolutionNote string   `json:"resolution_note,omitempty"`
}

type EventActionResponse struct {
	EventID        string `json:"event_id"`
	ReviewStatus   string `json:"review_status"`
	ReviewedAt     string `json:"reviewed_at,omitempty"`
	ResolvedAt     string `json:"resolved_at,omitempty"`
	ResolutionNote string `json:"resolution_note,omitempty"`
}

type ImportDocumentResponse struct {
	Document         DocumentResponse `json:"document"`
	TextExtracted    bool             `json:"text_extracted"`
	MetadataAnalyzed bool             `json:"metadata_analyzed"`
	EventsAnalyzed   bool             `json:"events_analyzed"`
	EventsDetected   int              `json:"events_detected"`
}

type DashboardResponse struct {
	CaseFile                         CaseFileResponse `json:"case_file"`
	NoteCount                        int              `json:"note_count"`
	DocumentCount                    int              `json:"document_count"`
	UpcomingEvents                   []EventResponse  `json:"upcoming_events"`
	DocumentsWithoutText             int              `json:"documents_without_text"`
	DocumentsWithoutTextList         []string         `json:"documents_without_text_list"`
	DocumentsWithUnknownMetadata     int              `json:"documents_with_unknown_metadata"`
	DocumentsWithUnknownMetadataList []string         `json:"documents_with_unknown_metadata_list"`
	DocumentsWithoutEvents           int              `json:"documents_without_events"`
	DocumentsWithoutEventsList       []string         `json:"documents_without_events_list"`
	OverdueCount                     int              `json:"overdue_count"`
	TodayCount                       int              `json:"today_count"`
	UpcomingCount                    int              `json:"upcoming_count"`
	CriticalCount                    int              `json:"critical_count"`
	HighCount                        int              `json:"high_count"`
	MediumCount                      int              `json:"medium_count"`
	LowCount                         int              `json:"low_count"`
	PendingReviewCount               int              `json:"pending_review_count"`
	ReviewedCount                    int              `json:"reviewed_count"`
	ResolvedCount                    int              `json:"resolved_count"`
	ActiveEventCount                 int              `json:"active_event_count"`
	ResolvedEventCount               int              `json:"resolved_event_count"`
	NeedsAttention                   bool             `json:"needs_attention"`
	TopAlert                         string           `json:"top_alert"`
	RecommendedNextAction            string           `json:"recommended_next_action"`
	ProceduralHint                   string           `json:"procedural_hint"`
}

func toCaseFileResponse(cf domaincasefile.CaseFile) CaseFileResponse {
	return CaseFileResponse{
		ID:                cf.ID.String(),
		ClientID:          cf.ClientID.String(),
		Reference:         cf.Reference,
		Title:             cf.Title,
		Type:              string(cf.Type),
		Status:            string(cf.Status),
		Description:       cf.Description,
		CalendarScope:     cf.CalendarScope,
		AugustNonBusiness: cf.AugustNonBusiness,
		CreatedAt:         cf.CreatedAt.Time().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         cf.UpdatedAt.Time().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toNoteResponse(n domainnote.Note) NoteResponse {
	return NoteResponse{
		ID:         n.ID.String(),
		CaseFileID: n.CaseFileID.String(),
		Title:      n.Title,
		Content:    n.Content,
		CreatedAt:  n.CreatedAt.Time().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  n.UpdatedAt.Time().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toDocumentResponse(d domaindocument.Document) DocumentResponse {
	return DocumentResponse{
		ID:           d.ID.String(),
		CaseFileID:   d.CaseFileID.String(),
		OriginalName: d.OriginalName,
		StoragePath:  d.StoragePath,
		MimeType:     d.MimeType,
		FileHash:     d.FileHash,
		CreatedAt:    d.CreatedAt.Time().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    d.UpdatedAt.Time().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toCaseFileDetailResponse(in casefileapp.CaseFileDetail) CaseFileDetailResponse {
	notes := make([]NoteResponse, 0, len(in.Notes))
	for _, n := range in.Notes {
		notes = append(notes, toNoteResponse(n))
	}

	documents := make([]DocumentResponse, 0, len(in.Documents))
	for _, d := range in.Documents {
		documents = append(documents, toDocumentResponse(d))
	}

	return CaseFileDetailResponse{
		CaseFile:  toCaseFileResponse(in.CaseFile),
		Notes:     notes,
		Documents: documents,
	}
}

func toCaseFileListResponse(items []domaincasefile.CaseFile) []CaseFileResponse {
	out := make([]CaseFileResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toCaseFileResponse(item))
	}
	return out
}

func toCaseFileEventResponse(e querymodels.CaseFileEventResult) EventResponse {
	return EventResponse{
		EventID:      e.EventID,
		DocumentID:   e.DocumentID,
		OriginalName: e.OriginalName,

		CaseFileID:        e.CaseFileID,
		CaseFileReference: e.CaseFileReference,
		CaseFileTitle:     e.CaseFileTitle,

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
		ReviewStatus:   e.ReviewStatus,
		ReviewedAt:     e.ReviewedAt,
		ResolvedAt:     e.ResolvedAt,
		ResolutionNote: e.ResolutionNote,
	}
}

func toCaseFileEventsResponse(items []querymodels.CaseFileEventResult) []EventResponse {
	out := make([]EventResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toCaseFileEventResponse(item))
	}
	return out
}

func toUpcomingEventResponse(e documentapp.UpcomingEvent) EventResponse {
	return EventResponse{
		EventID:      e.EventID,
		DocumentID:   e.DocumentID,
		OriginalName: e.OriginalName,

		CaseFileID:        e.CaseFileID,
		CaseFileReference: e.CaseFileReference,
		CaseFileTitle:     e.CaseFileTitle,

		EventType:      e.EventType,
		EventDate:      e.EventDate,
		SourceText:     e.SourceText,
		DaysRemaining:  e.DaysRemaining,
		Status:         e.Status,
		Priority:       e.Priority,
		DuplicateCount: e.DuplicateCount,
		DocumentNames:  e.DocumentNames,
		DocumentIDs:    e.DocumentIDs,
		AnchorDate:     e.AnchorDate,
		DateKind:       e.DateKind,
		AnchorSource:   e.AnchorSource,
		RelativeDays:   e.RelativeDays,
		IsBusinessDays: e.IsBusinessDays,
		AddExtraDay:    e.AddExtraDay,
		CalendarScope:  e.CalendarScope,
		TriggerText:    e.TriggerText,
		Computation:    e.Computation,
		ReviewStatus:   e.ReviewStatus,
		ReviewedAt:     e.ReviewedAt,
		ResolvedAt:     e.ResolvedAt,
		ResolutionNote: e.ResolutionNote,
	}
}

func toUpcomingEventsResponse(items []documentapp.UpcomingEvent) []EventResponse {
	out := make([]EventResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toUpcomingEventResponse(item))
	}
	return out
}

func toEventActionResponse(eventID, reviewStatus, reviewedAt, resolvedAt, resolutionNote string) EventActionResponse {
	return EventActionResponse{
		EventID:        eventID,
		ReviewStatus:   reviewStatus,
		ReviewedAt:     reviewedAt,
		ResolvedAt:     resolvedAt,
		ResolutionNote: resolutionNote,
	}
}

func toDashboardResponse(in casefileapp.CaseFileDashboardResult) DashboardResponse {
	return DashboardResponse{
		CaseFile:                         toCaseFileResponse(in.CaseFile),
		NoteCount:                        in.NoteCount,
		DocumentCount:                    in.DocumentCount,
		UpcomingEvents:                   toUpcomingEventsResponse(in.UpcomingEvents),
		DocumentsWithoutText:             in.DocumentsWithoutText,
		DocumentsWithoutTextList:         in.DocumentsWithoutTextList,
		DocumentsWithUnknownMetadata:     in.DocumentsWithUnknownMetadata,
		DocumentsWithUnknownMetadataList: in.DocumentsWithUnknownMetadataList,
		DocumentsWithoutEvents:           in.DocumentsWithoutEvents,
		DocumentsWithoutEventsList:       in.DocumentsWithoutEventsList,
		OverdueCount:                     in.OverdueCount,
		TodayCount:                       in.TodayCount,
		UpcomingCount:                    in.UpcomingCount,
		CriticalCount:                    in.CriticalCount,
		HighCount:                        in.HighCount,
		MediumCount:                      in.MediumCount,
		LowCount:                         in.LowCount,
		PendingReviewCount:               in.PendingReviewCount,
		ReviewedCount:                    in.ReviewedCount,
		ResolvedCount:                    in.ResolvedCount,
		ActiveEventCount:                 in.ActiveEventCount,
		ResolvedEventCount:               in.ResolvedEventCount,
		NeedsAttention:                   in.NeedsAttention,
		TopAlert:                         in.TopAlert,
		RecommendedNextAction:            in.RecommendedNextAction,
		ProceduralHint:                   in.ProceduralHint,
	}
}

func toImportDocumentResponse(result documentapp.ImportDocumentResult) ImportDocumentResponse {
	return ImportDocumentResponse{
		Document:         toDocumentResponse(result.Document),
		TextExtracted:    result.TextExtracted,
		MetadataAnalyzed: result.MetadataAnalyzed,
		EventsAnalyzed:   result.EventsAnalyzed,
		EventsDetected:   result.EventsDetected,
	}
}

func toDocumentEventDetail(e domaindocument.Event) DocumentEventDetail {
	return DocumentEventDetail{
		EventID:        e.ID.String(),
		DocumentID:     e.DocumentID.String(),
		EventType:      e.EventType,
		EventDate:      e.EventDate,
		SourceText:     e.SourceText,
		CreatedAt:      e.CreatedAt,
		AnchorDate:     e.AnchorDate,
		DateKind:       e.DateKind,
		AnchorSource:   e.AnchorSource,
		RelativeDays:   e.RelativeDays,
		IsBusinessDays: e.IsBusinessDays,
		AddExtraDay:    e.AddExtraDay,
		CalendarScope:  e.CalendarScope,
		TriggerText:    e.TriggerText,
		Computation:    e.Computation,
		ReviewStatus:   e.ReviewStatus,
		ReviewedAt:     e.ReviewedAt,
		ResolvedAt:     e.ResolvedAt,
		ResolutionNote: e.ResolutionNote,
	}
}

func toDocumentDetailResponse(in documentapp.GetDocumentDetailResult) DocumentDetailResponse {
	events := make([]DocumentEventDetail, 0, len(in.Events))
	for _, event := range in.Events {
		events = append(events, toDocumentEventDetail(event))
	}

	return DocumentDetailResponse{
		Document:             toDocumentResponse(in.Document),
		FileExists:           in.FileExists,
		HasExtractedText:     in.HasExtractedText,
		ExtractedText:        in.ExtractedText,
		ExtractedTextLength:  in.ExtractedTextLength,
		ExtractedTextPreview: in.ExtractedTextPreview,
		HasMetadata:          in.HasMetadata,
		DocumentType:         in.DocumentType,
		LegalArea:            in.LegalArea,
		MetadataAnalyzedAt:   in.MetadataAnalyzedAt,
		HasEvents:            in.HasEvents,
		Events:               events,
	}
}
