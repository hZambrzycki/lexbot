package documentapp

import (
	"strings"

	"lexbox/internal/application/querymodels"
)

type ParsedSearchDocumentsQuery struct {
	Terms   string
	Filters querymodels.SearchDocumentFilters
}

func ParseSearchDocumentsQuery(raw string) ParsedSearchDocumentsQuery {
	parts := strings.Fields(strings.TrimSpace(raw))

	freeTerms := make([]string, 0, len(parts))
	filters := querymodels.SearchDocumentFilters{}

	for _, part := range parts {
		key, value, ok := strings.Cut(part, ":")
		if !ok || strings.TrimSpace(value) == "" {
			freeTerms = append(freeTerms, part)
			continue
		}

		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.ToLower(strings.TrimSpace(value))

		switch key {
		case "type":
			filters.DocumentType = value
		case "area":
			filters.LegalArea = value
		case "review":
			filters.ReviewStatus = value
		case "doc":
			filters.DocType = value
		case "has":
			filters.Has = value
		default:
			freeTerms = append(freeTerms, part)
		}
	}

	return ParsedSearchDocumentsQuery{
		Terms:   strings.Join(freeTerms, " "),
		Filters: filters,
	}
}
