package httpapi

import (
	"net/http"
	"strings"

	casefileapp "lexbox/internal/application/casefile"
)

type CaseFileHandler struct {
	ListCaseFiles        casefileapp.ListCaseFiles
	GetCaseFileDetail    casefileapp.GetCaseFileDetail
	GetCaseFileDashboard casefileapp.GetCaseFileDashboard
	EventHandler         EventHandler
}

func (h CaseFileHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/case-files", h.handleList)

	mux.HandleFunc("/case-files/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/case-files/")
		parts := strings.Split(path, "/")

		if len(parts) == 0 {
			http.NotFound(w, r)
			return
		}

		id := strings.TrimSpace(parts[0])
		if id == "" {
			http.NotFound(w, r)
			return
		}

		if len(parts) == 1 {
			h.handleGet(w, r, id)
			return
		}

		if len(parts) == 2 && parts[1] == "dashboard" {
			h.handleDashboard(w, r, id)
			return
		}

		if len(parts) == 2 && parts[1] == "events" {
			h.EventHandler.handleCaseFileEvents(w, r, id)
			return
		}

		if len(parts) == 3 && parts[1] == "events" && parts[2] == "upcoming" {
			h.EventHandler.handleCaseFileUpcomingEvents(w, r, id)
			return
		}

		http.NotFound(w, r)
	})
}

func (h CaseFileHandler) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	result, err := h.ListCaseFiles.Execute(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCaseFileListResponse(result))
}

func (h CaseFileHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	result, err := h.GetCaseFileDetail.Execute(r.Context(), casefileapp.GetCaseFileDetailInput{
		ID: id,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCaseFileDetailResponse(result))
}

func (h CaseFileHandler) handleDashboard(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	result, err := h.GetCaseFileDashboard.Execute(r.Context(), casefileapp.GetCaseFileDashboardInput{
		CaseFileID: id,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDashboardResponse(result))
}
