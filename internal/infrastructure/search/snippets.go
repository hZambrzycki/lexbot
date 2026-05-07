package search

import (
	"sort"
	"strings"
)

type snippetMatchRange struct {
	start int
	end   int
}

func BuildSnippet(content string, rawQuery string, terms []string, maxLen int) string {
	normalized := normalizeWhitespace(content)

	if normalized == "" {
		return ""
	}

	bestTerm := bestMatchingTerm(normalized, terms)

	if bestTerm == "" {
		bestTerm = strings.TrimSpace(rawQuery)
	}

	if bestTerm == "" {
		return truncateSnippet(normalized, maxLen)
	}

	matchStart, _, ok := findAccentInsensitiveMatch(normalized, bestTerm)

	if !ok {
		return truncateSnippet(normalized, maxLen)
	}

	start := matchStart - (maxLen / 3)

	if start < 0 {
		start = 0
	}

	end := start + maxLen

	if end > len(normalized) {
		end = len(normalized)
		start = end - maxLen

		if start < 0 {
			start = 0
		}
	}

	snippet := normalized[start:end]
	highlighted := highlightSnippetTerms(snippet, terms)

	return decorateSnippetBounds(highlighted, start, end, len(normalized))
}

func highlightSnippetTerms(snippet string, terms []string) string {
	ranges := findAllAccentInsensitiveMatches(snippet, terms)

	if len(ranges) == 0 {
		return snippet
	}

	var sb strings.Builder
	cursor := 0

	for _, match := range ranges {
		if match.start < cursor {
			continue
		}

		sb.WriteString(snippet[cursor:match.start])
		sb.WriteString("[")
		sb.WriteString(snippet[match.start:match.end])
		sb.WriteString("]")

		cursor = match.end
	}

	if cursor < len(snippet) {
		sb.WriteString(snippet[cursor:])
	}

	return sb.String()
}

func findAllAccentInsensitiveMatches(content string, terms []string) []snippetMatchRange {
	ranges := make([]snippetMatchRange, 0)

	for _, term := range terms {
		normalizedTerm := NormalizeText(term)

		if normalizedTerm == "" {
			continue
		}

		searchContent := content

		for {
			start, end, ok := findAccentInsensitiveMatch(searchContent, normalizedTerm)

			if !ok {
				break
			}

			absoluteStart := len(content) - len(searchContent) + start
			absoluteEnd := len(content) - len(searchContent) + end

			ranges = append(ranges, snippetMatchRange{
				start: absoluteStart,
				end:   absoluteEnd,
			})

			if end >= len(searchContent) {
				break
			}

			searchContent = searchContent[end:]
		}
	}

	sort.SliceStable(ranges, func(i, j int) bool {
		if ranges[i].start == ranges[j].start {
			return ranges[i].end > ranges[j].end
		}

		return ranges[i].start < ranges[j].start
	})

	filtered := make([]snippetMatchRange, 0, len(ranges))

	for _, current := range ranges {
		if len(filtered) == 0 {
			filtered = append(filtered, current)
			continue
		}

		last := filtered[len(filtered)-1]

		if current.start < last.end {
			continue
		}

		filtered = append(filtered, current)
	}

	return filtered
}

func bestMatchingTerm(content string, terms []string) string {
	normalizedContent := NormalizeText(content)

	bestTerm := ""
	bestCount := 0

	for _, term := range terms {
		normalizedTerm := NormalizeText(term)

		if normalizedTerm == "" {
			continue
		}

		count := strings.Count(normalizedContent, normalizedTerm)

		if count > bestCount {
			bestCount = count
			bestTerm = term
		}
	}

	return bestTerm
}

func truncateSnippet(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}

	return content[:maxLen] + "..."
}

func decorateSnippetBounds(snippet string, start, end, totalLen int) string {
	if start > 0 {
		snippet = "..." + snippet
	}

	if end < totalLen {
		snippet = snippet + "..."
	}

	return snippet
}

func findAccentInsensitiveMatch(content string, query string) (int, int, bool) {
	normalizedContent, indexMap := NormalizeTextWithIndex(content)
	normalizedQuery := NormalizeText(query)

	start := strings.Index(normalizedContent, normalizedQuery)

	if start == -1 {
		return 0, 0, false
	}

	originalStart := indexMap[start]
	originalEnd := indexMap[start+len(normalizedQuery)-1] + 1

	return originalStart, originalEnd, true
}
