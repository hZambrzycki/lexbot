package document

import "lexbox/internal/domain/shared"

const (
	EventTypeUnknown      = "unknown"
	EventTypeHearing      = "hearing"
	EventTypeAppearance   = "appearance"
	EventTypeDeadline     = "deadline"
	EventTypeRequirement  = "requirement"
	EventTypeNotification = "notification"
	EventTypeFiling       = "filing"
)

const (
	ReviewStatusPending  = "pending"
	ReviewStatusReviewed = "reviewed"
	ReviewStatusResolved = "resolved"
)

func IsValidReviewStatus(status string) bool {
	switch status {
	case ReviewStatusPending, ReviewStatusReviewed, ReviewStatusResolved:
		return true
	default:
		return false
	}
}

type Event struct {
	ID             shared.ID
	DocumentID     shared.ID
	EventType      string
	EventDate      string
	SourceText     string
	CreatedAt      string
	AnchorDate     string
	DateKind       string
	AnchorSource   string
	RelativeDays   int
	IsBusinessDays bool
	AddExtraDay    bool
	CalendarScope  string
	TriggerText    string
	Computation    string

	ReviewStatus   string
	ReviewedAt     string
	ResolvedAt     string
	ResolutionNote string
}

func NewEvent(
	id shared.ID,
	documentID shared.ID,
	eventType string,
	eventDate string,
	sourceText string,
	createdAt string,
	anchorDate string,
	dateKind string,
	anchorSource string,
	relativeDays int,
	isBusinessDays bool,
	addExtraDay bool,
	calendarScope string,
	triggerText string,
	computation string,
) (Event, error) {
	if id == "" || documentID == "" {
		return Event{}, shared.ErrInvalidID
	}

	if eventType == "" || eventDate == "" || sourceText == "" || createdAt == "" {
		return Event{}, shared.ErrEmptyField
	}

	if relativeDays < 0 {
		return Event{}, shared.ErrInvalidAssociation
	}

	return Event{
		ID:             id,
		DocumentID:     documentID,
		EventType:      eventType,
		EventDate:      eventDate,
		SourceText:     sourceText,
		CreatedAt:      createdAt,
		AnchorDate:     anchorDate,
		DateKind:       dateKind,
		AnchorSource:   anchorSource,
		RelativeDays:   relativeDays,
		IsBusinessDays: isBusinessDays,
		AddExtraDay:    addExtraDay,
		CalendarScope:  calendarScope,
		TriggerText:    triggerText,
		Computation:    computation,
		ReviewStatus:   ReviewStatusPending,
		ReviewedAt:     "",
		ResolvedAt:     "",
		ResolutionNote: "",
	}, nil
}
