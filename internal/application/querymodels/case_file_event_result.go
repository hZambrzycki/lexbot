package querymodels

type CaseFileEventResult struct {
	EventID      string
	DocumentID   string
	OriginalName string

	CaseFileID        string
	CaseFileReference string
	CaseFileTitle     string

	EventType  string
	EventDate  string
	SourceText string

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
