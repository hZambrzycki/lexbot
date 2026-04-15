package note

import (
	"strings"

	"lexbox/internal/domain/shared"
)

type Note struct {
	ID         shared.ID
	CaseFileID shared.ID
	Title      string
	Content    string

	CreatedAt shared.Timestamp
	UpdatedAt shared.Timestamp
}

func NewNote(id shared.ID, caseFileID shared.ID, title string, content string) (Note, error) {
	if id == "" {
		return Note{}, shared.ErrInvalidID
	}

	if caseFileID == "" {
		return Note{}, shared.ErrInvalidAssociation
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if content == "" {
		return Note{}, shared.ErrEmptyField
	}

	now := shared.Now()

	return Note{
		ID:         id,
		CaseFileID: caseFileID,
		Title:      title,
		Content:    content,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
