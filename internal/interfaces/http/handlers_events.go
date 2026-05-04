package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	documentapp "lexbox/internal/application/document"
)

type EventHandler struct {
	GetEvent                documentapp.GetEvent
	ListEventsByCaseFile    documentapp.ListEventsByCaseFile
	ListUpcomingEvents      documentapp.ListUpcomingEvents
	ExportUpcomingEventsICS documentapp.ExportUpcomingEventsICS
	MarkEventReviewed       documentapp.MarkEventReviewed
	MarkEventResolved       documentapp.MarkEventResolved
	ReopenEvent             documentapp.ReopenEvent
}

type reviewEventRequest struct {
	ReviewStatus   string `json:"review_status"`
	ResolutionNote string `json:"resolution_note"`
}

type resolveEventRequest struct {
	ResolutionNote string `json:"resolution_note"`
}

func (h EventHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/events/upcoming", h.handleGlobalUpcomingEvents)
	mux.HandleFunc("/events/upcoming.ics", h.handleGlobalUpcomingEventsICS)

	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/events/")
		parts := strings.Split(path, "/")

		if len(parts) == 1 {
			eventID := strings.TrimSpace(parts[0])
			if eventID == "" {
				http.NotFound(w, r)
				return
			}

			if r.Method != http.MethodGet {
				writeMethodNotAllowed(w, http.MethodGet)
				return
			}

			h.handleGetEvent(w, r, eventID)
			return
		}

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

func (h EventHandler) handleGlobalUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/events/upcoming" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	eventType := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("type")))
	reviewStatus := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("review_status")))
	relativeOnly := parseBoolQuery(r.URL.Query().Get("relative_only"))

	result, err := h.ListUpcomingEvents.Execute(r.Context(), documentapp.ListUpcomingEventsInput{
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

func (h EventHandler) handleGlobalUpcomingEventsICS(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/events/upcoming.ics" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	eventType := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("type")))

	result, err := h.ExportUpcomingEventsICS.Execute(r.Context(), documentapp.ExportUpcomingEventsICSInput{
		EventType: eventType,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeICS(w, http.StatusOK, "lexbox-upcoming-events.ics", result.Content)
}

func (h EventHandler) handleGetEvent(w http.ResponseWriter, r *http.Request, eventID string) {
	result, err := h.GetEvent.Execute(r.Context(), documentapp.GetEventInput{
		EventID: eventID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUpcomingEventResponse(result))
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

func (h EventHandler) handleCaseFileUpcomingEventsICS(w http.ResponseWriter, r *http.Request, caseFileID string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	eventType := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("type")))

	result, err := h.ExportUpcomingEventsICS.Execute(r.Context(), documentapp.ExportUpcomingEventsICSInput{
		CaseFileID: caseFileID,
		EventType:  eventType,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	filename := fmt.Sprintf("lexbox-case-file-%s-upcoming-events.ics", caseFileID)
	writeICS(w, http.StatusOK, filename, result.Content)
}

func (h EventHandler) handleReview(w http.ResponseWriter, r *http.Request, eventID string) {
	var req reviewEventRequest

	if r.Body != nil {
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil && !errors.Is(err, io.EOF) {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid json body",
			})
			return
		}
	}

	reviewStatus := strings.TrimSpace(strings.ToLower(req.ReviewStatus))

	if reviewStatus == "" {
		reviewStatus = "reviewed"
	}

	switch reviewStatus {
	case "reviewed":
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
		return

	case "resolved":
		result, err := h.MarkEventResolved.Execute(r.Context(), documentapp.MarkEventResolvedInput{
			EventID:        eventID,
			ResolutionNote: strings.TrimSpace(req.ResolutionNote),
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
		return

	case "pending":
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
		return

	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "review_status must be reviewed, resolved or pending",
		})
		return
	}
}

func (h EventHandler) handleResolve(w http.ResponseWriter, r *http.Request, eventID string) {
	var req resolveEventRequest

	if r.Body != nil {
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil && !errors.Is(err, io.EOF) {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid json body",
			})
			return
		}
	}

	result, err := h.MarkEventResolved.Execute(r.Context(), documentapp.MarkEventResolvedInput{
		EventID:        eventID,
		ResolutionNote: strings.TrimSpace(req.ResolutionNote),
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
