package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeICS(w http.ResponseWriter, status int, filename string, content string) {
	filename = sanitizeDownloadFilename(filename)

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(status)

	_, _ = w.Write([]byte(content))
}

func sanitizeDownloadFilename(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return "lexbox-calendar.ics"
	}

	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		`"`, "-",
		"<", "-",
		">", "-",
		"|", "-",
	)

	filename = replacer.Replace(filename)

	if !strings.HasSuffix(strings.ToLower(filename), ".ics") {
		filename += ".ics"
	}

	return filename
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
}

func writeMethodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": "method not allowed",
	})
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
