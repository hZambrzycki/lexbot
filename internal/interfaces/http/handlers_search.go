package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	documentapp "lexbox/internal/application/document"
)

type SearchHandler struct {
	GlobalSearch documentapp.GlobalSearch
}

func (h SearchHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/search/global", h.handleGlobalSearch)
}

func (h SearchHandler) handleGlobalSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	limit := parseGlobalSearchLimit(r.URL.Query().Get("limit"))

	result, err := h.GlobalSearch.Execute(r.Context(), documentapp.GlobalSearchInput{
		Query: query,
		Limit: limit,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toGlobalSearchResultsResponse(result))
}

func parseGlobalSearchLimit(raw string) int {
	const defaultLimit = 8
	const maxLimit = 25

	value := strings.TrimSpace(raw)
	if value == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit <= 0 {
		return defaultLimit
	}

	if limit > maxLimit {
		return maxLimit
	}

	return limit
}
