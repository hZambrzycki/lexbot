package casefileapp

import (
	"context"
	"lexbox/internal/application/ports"
	domaincalendar "lexbox/internal/domain/calendar"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/shared"
	"strings"
)

type CreateCaseFileInput struct {
	ClientID          string
	Reference         string
	Title             string
	Type              casefile.Type
	Description       string
	CalendarScope     string
	AugustNonBusiness *bool
}
type CreateCaseFile struct {
	CaseFiles ports.CaseFileRepository
	Clients   ports.ClientRepository
	IDs       ports.IDGenerator
}

func (uc CreateCaseFile) Execute(ctx context.Context, in CreateCaseFileInput) (casefile.CaseFile, error) {
	clientID := shared.NewID(strings.TrimSpace(in.ClientID))
	if clientID == "" {
		return casefile.CaseFile{}, shared.ErrInvalidAssociation
	}

	// Verificamos que el cliente existe
	_, err := uc.Clients.GetByID(ctx, clientID)
	if err != nil {
		return casefile.CaseFile{}, err
	}

	title := strings.TrimSpace(in.Title)
	reference := strings.TrimSpace(in.Reference)
	description := strings.TrimSpace(in.Description)

	// Creamos entidad de dominio
	cf, err := casefile.NewCaseFile(
		uc.IDs.NewID(),
		clientID,
		title,
	)
	if err != nil {
		return casefile.CaseFile{}, err
	}

	// Campos adicionales
	cf.Reference = reference
	cf.Description = description

	if strings.TrimSpace(in.CalendarScope) != "" {
		scope := strings.TrimSpace(in.CalendarScope)
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
	// Validación de tipo
	if in.Type != "" {
		if !casefile.IsValidType(in.Type) {
			return casefile.CaseFile{}, shared.ErrInvalidState
		}
		cf.Type = in.Type
	} else {
		cf.Type = casefile.TypeOtros
	}

	// Persistencia
	if err := uc.CaseFiles.Save(ctx, cf); err != nil {
		return casefile.CaseFile{}, err
	}

	return cf, nil
}
