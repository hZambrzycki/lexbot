package document

import "lexbox/internal/domain/shared"

type Metadata struct {
	DocumentID   shared.ID
	DocumentType string
	LegalArea    string
	AnalyzedAt   string
}

func NewMetadata(
	documentID shared.ID,
	documentType string,
	legalArea string,
	analyzedAt string,
) (Metadata, error) {
	if documentID == "" {
		return Metadata{}, shared.ErrInvalidID
	}

	if documentType == "" || legalArea == "" || analyzedAt == "" {
		return Metadata{}, shared.ErrEmptyField
	}

	return Metadata{
		DocumentID:   documentID,
		DocumentType: documentType,
		LegalArea:    legalArea,
		AnalyzedAt:   analyzedAt,
	}, nil
}
