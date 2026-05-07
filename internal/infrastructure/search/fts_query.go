package search

import (
	"strings"
	"unicode"
)

func BuildFTS5Query(rawQuery string) string {
	terms := SplitTermsForFTS5(rawQuery)
	if len(terms) == 0 {
		return ""
	}

	return strings.Join(terms, " ")
}

func SplitTermsForFTS5(rawQuery string) []string {
	fields := strings.FieldsFunc(rawQuery, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	terms := make([]string, 0, len(fields))
	seen := make(map[string]struct{}, len(fields))

	for _, field := range fields {
		term := NormalizeText(field)
		if term == "" {
			continue
		}

		if _, exists := seen[term]; exists {
			continue
		}

		seen[term] = struct{}{}

		if len([]rune(term)) >= 4 {
			term += "*"
		}

		terms = append(terms, term)
	}

	return terms
}
