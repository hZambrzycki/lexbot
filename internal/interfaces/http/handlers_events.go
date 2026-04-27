package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	documentapp "lexbox/internal/application/document"
)

type EventHandler struct {
	ListEventsByCaseFile documentapp.ListEventsByCaseFile
	ListUpcomingEvents   documentapp.ListUpcomingEvents
	MarkEventReviewed    documentapp.MarkEventReviewed
	MarkEventResolved    documentapp.MarkEventResolved
	ReopenEvent          documentapp.ReopenEvent
}

type resolveEventRequest struct {
	ResolutionNote string `json:"resolution_note"`
}

func (h EventHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/events/")
		parts := strings.Split(path, "/")

		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}

		eventID := strings.TrimSpace(parts[0])
		action := strings.TrimSpace(parts[1])

		if eventID == "" {
			http.NotFound(w, r)
			return
		}

		switch action {
		case "review":
			if r.Method != http.MethodPost {
				writeMethodNotAllowed(w, http.MethodPost)
				return
			}
			h.handleReview(w, r, eventID)
			return

		case "resolve":
			if r.Method != http.MethodPost {
				writeMethodNotAllowed(w, http.MethodPost)
				return
			}
			h.handleResolve(w, r, eventID)
			return

		case "reopen":
			if r.Method != http.MethodPost {
				writeMethodNotAllowed(w, http.MethodPost)
				return
			}
			h.handleReopen(w, r, eventID)
			return
		}

		http.NotFound(w, r)
	})
}

func (h EventHandler) handleCaseFileEvents(w http.ResponseWriter, r *http.Request, caseFileID string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	reviewStatus := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("review_status")))

	result, err := h.ListEventsByCaseFile.Execute(r.Context(), documentapp.ListEventsByCaseFileInput{
		CaseFileID:   caseFileID,
		ReviewStatus: reviewStatus,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCaseFileEventsResponse(result))
}

func (h EventHandler) handleCaseFileUpcomingEvents(w http.ResponseWriter, r *http.Request, caseFileID string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	eventType := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("type")))
	reviewStatus := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("review_status")))
	relativeOnly := parseBoolQuery(r.URL.Query().Get("relative_only"))

	result, err := h.ListUpcomingEvents.Execute(r.Context(), documentapp.ListUpcomingEventsInput{
		CaseFileID:   caseFileID,
		EventType:    eventType,
		ReviewStatus: reviewStatus,
		RelativeOnly: relativeOnly,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUpcomingEventsResponse(result))
}

func (h EventHandler) handleReview(w http.ResponseWriter, r *http.Request, eventID string) {
	result, err := h.MarkEventReviewed.Execute(r.Context(), documentapp.MarkEventReviewedInput{
		EventID: eventID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEventActionResponse(
		result.EventID,
		result.ReviewStatus,
		result.ReviewedAt,
		result.ResolvedAt,
		result.ResolutionNote,
	))
}

func (h EventHandler) handleResolve(w http.ResponseWriter, r *http.Request, eventID string) {
	var req resolveEventRequest

	if r.Body != nil {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid json body",
			})
			return
		}
	}

	result, err := h.MarkEventResolved.Execute(r.Context(), documentapp.MarkEventResolvedInput{
		EventID:        eventID,
		ResolutionNote: req.ResolutionNote,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEventActionResponse(
		result.EventID,
		result.ReviewStatus,
		result.ReviewedAt,
		result.ResolvedAt,
		result.ResolutionNote,
	))
}

func (h EventHandler) handleReopen(w http.ResponseWriter, r *http.Request, eventID string) {
	result, err := h.ReopenEvent.Execute(r.Context(), documentapp.ReopenEventInput{
		EventID: eventID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEventActionResponse(
		result.EventID,
		result.ReviewStatus,
		result.ReviewedAt,
		result.ResolvedAt,
		result.ResolutionNote,
	))
}

func parseBoolQuery(value string) bool {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
