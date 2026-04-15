package casefile

import (
	"strings"

	"lexbox/internal/domain/shared"
)

type Status string
type Type string

const (
	StatusOpen     Status = "open"
	StatusPending  Status = "pending"
	StatusClosed   Status = "closed"
	StatusArchived Status = "archived"
)

const (
	TypeCivil          Type = "civil"
	TypeLaboral        Type = "laboral"
	TypeExtranjeria    Type = "extranjeria"
	TypeMercantil      Type = "mercantil"
	TypeAdministrativo Type = "administrativo"
	TypeOtros          Type = "otros"
)

func IsValidStatus(s Status) bool {
	switch s {
	case StatusOpen, StatusPending, StatusClosed, StatusArchived:
		return true
	default:
		return false
	}
}

func IsValidType(t Type) bool {
	switch t {
	case TypeCivil, TypeLaboral, TypeExtranjeria, TypeMercantil, TypeAdministrativo, TypeOtros:
		return true
	default:
		return false
	}
}

type CaseFile struct {
	ID          shared.ID
	ClientID    shared.ID
	Reference   string
	Title       string
	Type        Type
	Status      Status
	Description string

	CreatedAt shared.Timestamp
	UpdatedAt shared.Timestamp
}

func NewCaseFile(id shared.ID, clientID shared.ID, title string) (CaseFile, error) {
	if id == "" {
		return CaseFile{}, shared.ErrInvalidID
	}

	if clientID == "" {
		return CaseFile{}, shared.ErrInvalidAssociation
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return CaseFile{}, shared.ErrEmptyField
	}

	now := shared.Now()

	return CaseFile{
		ID:        id,
		ClientID:  clientID,
		Title:     title,
		Type:      TypeOtros,
		Status:    StatusOpen,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (cf *CaseFile) SetType(t Type) error {
	if !IsValidType(t) {
		return shared.ErrInvalidState
	}

	cf.Type = t
	cf.touch()
	return nil
}

func (cf *CaseFile) SetStatus(s Status) error {
	if !IsValidStatus(s) {
		return shared.ErrInvalidState
	}

	cf.Status = s
	cf.touch()
	return nil
}

func (cf *CaseFile) SetReference(reference string) {
	cf.Reference = strings.TrimSpace(reference)
	cf.touch()
}

func (cf *CaseFile) SetDescription(description string) {
	cf.Description = strings.TrimSpace(description)
	cf.touch()
}

func (cf *CaseFile) Rename(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return shared.ErrEmptyField
	}

	cf.Title = title
	cf.touch()
	return nil
}

func (cf *CaseFile) touch() {
	cf.UpdatedAt = shared.Now()
}
