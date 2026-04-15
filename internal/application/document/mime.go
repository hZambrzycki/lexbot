package documentapp

import "strings"

func normalizeMimeType(value string) string {
	mimeType := strings.TrimSpace(strings.ToLower(value))

	for strings.HasSuffix(mimeType, "~") {
		mimeType = strings.TrimSuffix(mimeType, "~")
		mimeType = strings.TrimSpace(mimeType)
	}

	return mimeType
}
