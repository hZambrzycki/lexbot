package search

import "strings"

func ComputeScore(content string, terms []string) int {
	normalized := NormalizeText(normalizeWhitespace(content))
	if normalized == "" {
		return 0
	}

	score := 0

	for _, term := range terms {
		normalizedTerm := NormalizeText(term)
		if normalizedTerm == "" {
			continue
		}

		count := strings.Count(normalized, normalizedTerm)
		if count > 0 {
			score += count * len(normalizedTerm)
		}
	}

	return score
}

func SplitTerms(query string) []string {
	fields := strings.Fields(strings.TrimSpace(query))
	terms := make([]string, 0, len(fields))

	seen := make(map[string]struct{}, len(fields))

	for _, field := range fields {
		normalized := NormalizeText(field)
		if normalized == "" {
			continue
		}

		if _, exists := seen[normalized]; exists {
			continue
		}

		seen[normalized] = struct{}{}
		terms = append(terms, field)
	}

	return terms
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
