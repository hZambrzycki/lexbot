package casefileapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	domaincalendar "lexbox/internal/domain/calendar"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/shared"
)

type UpdateCaseFileConfigInput struct {
	CaseFileID        string
	CalendarScope     *string
	AugustNonBusiness *bool
}

type UpdateCaseFileConfig struct {
	CaseFiles ports.CaseFileRepository
}

func (uc UpdateCaseFileConfig) Execute(ctx context.Context, in UpdateCaseFileConfigInput) (casefile.CaseFile, error) {
	id := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if id == "" {
		return casefile.CaseFile{}, shared.ErrInvalidID
	}

	cf, err := uc.CaseFiles.GetByID(ctx, id)
	if err != nil {
		return casefile.CaseFile{}, err
	}

	if in.CalendarScope != nil {
		scope := strings.TrimSpace(*in.CalendarScope)
		switch scope {
		case domaincalendar.ScopeMadrid, domaincalendar.ScopeState:
			cf.CalendarScope = scope
		default:
			return casefile.CaseFile{}, shared.ErrInvalidState
		}
	}

	if in.AugustNonBusiness != nil {
		cf.AugustNonBusiness = *in.AugustNonBusiness
	}

	// tocamos updated_at
	cf.UpdatedAt = shared.Now()

	if err := uc.CaseFiles.Update(ctx, cf); err != nil {
		return casefile.CaseFile{}, err
	}

	return cf, nil
}
