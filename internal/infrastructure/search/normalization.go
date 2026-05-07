package search

import (
	"strings"
)

func NormalizeText(value string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ä", "a", "â", "a",
		"é", "e", "è", "e", "ë", "e", "ê", "e",
		"í", "i", "ì", "i", "ï", "i", "î", "i",
		"ó", "o", "ò", "o", "ö", "o", "ô", "o",
		"ú", "u", "ù", "u", "ü", "u", "û", "u",
		"ñ", "n",

		"Á", "a", "À", "a", "Ä", "a", "Â", "a",
		"É", "e", "È", "e", "Ë", "e", "Ê", "e",
		"Í", "i", "Ì", "i", "Ï", "i", "Î", "i",
		"Ó", "o", "Ò", "o", "Ö", "o", "Ô", "o",
		"Ú", "u", "Ù", "u", "Ü", "u", "Û", "u",
		"Ñ", "n",
	)

	return strings.ToLower(strings.TrimSpace(replacer.Replace(value)))
}

func NormalizeTextWithIndex(value string) (string, []int) {
	var normalized strings.Builder
	indexMap := make([]int, 0, len(value))

	for index, char := range value {
		replaced := NormalizeText(string(char))

		for _, normalizedChar := range replaced {
			normalized.WriteRune(normalizedChar)
			indexMap = append(indexMap, index)
		}
	}

	return normalized.String(), indexMap
}
