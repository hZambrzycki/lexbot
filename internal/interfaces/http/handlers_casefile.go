package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	casefileapp "lexbox/internal/application/casefile"
	documentapp "lexbox/internal/application/document"
	"lexbox/internal/domain/casefile"
)

type CaseFileHandler struct {
	CreateCaseFile       casefileapp.CreateCaseFile
	ListCaseFiles        casefileapp.ListCaseFiles
	GetCaseFileDetail    casefileapp.GetCaseFileDetail
	GetCaseFileDashboard casefileapp.GetCaseFileDashboard
	ImportDocument       documentapp.ImportDocument
	EventHandler         EventHandler
}

type createCaseFileRequest struct {
	ClientID          string `json:"client_id"`
	Reference         string `json:"reference"`
	Title             string `json:"title"`
	Type              string `json:"type"`
	Description       string `json:"description"`
	CalendarScope     string `json:"calendar_scope"`
	AugustNonBusiness bool   `json:"august_non_business"`
}

func (h CaseFileHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/case-files", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleList(w, r)
			return
		case http.MethodPost:
			h.handleCreate(w, r)
			return
		default:
			writeMethodNotAllowed(w, http.MethodGet, http.MethodPost)
			return
		}
	})

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

		if len(parts) == 2 && parts[1] == "documents" {
			h.handleImportDocument(w, r, id)
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

		if len(parts) == 3 && parts[1] == "events" && parts[2] == "upcoming.ics" {
			h.EventHandler.handleCaseFileUpcomingEventsICS(w, r, id)
			return
		}

		http.NotFound(w, r)
	})
}

func (h CaseFileHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req createCaseFileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
		return
	}

	augustNonBusiness := req.AugustNonBusiness

	result, err := h.CreateCaseFile.Execute(r.Context(), casefileapp.CreateCaseFileInput{
		ClientID:          strings.TrimSpace(req.ClientID),
		Reference:         strings.TrimSpace(req.Reference),
		Title:             strings.TrimSpace(req.Title),
		Type:              casefile.Type(strings.TrimSpace(req.Type)),
		Description:       strings.TrimSpace(req.Description),
		CalendarScope:     strings.TrimSpace(req.CalendarScope),
		AugustNonBusiness: &augustNonBusiness,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toCaseFileResponse(result))
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

func (h CaseFileHandler) handleImportDocument(w http.ResponseWriter, r *http.Request, caseFileID string) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	if err := r.ParseMultipartForm(64 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid multipart form",
		})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing file field",
		})
		return
	}
	defer file.Close()

	tmpDir, err := os.MkdirTemp("", "lexbox-upload-*")
	if err != nil {
		writeError(w, err)
		return
	}
	defer os.RemoveAll(tmpDir)

	originalName := filepath.Base(header.Filename)
	if strings.TrimSpace(originalName) == "" || originalName == "." {
		originalName = "documento"
	}

	tmpPath := filepath.Join(tmpDir, originalName)

	tmp, err := os.Create(tmpPath)
	if err != nil {
		writeError(w, err)
		return
	}

	if _, err := io.Copy(tmp, file); err != nil {
		_ = tmp.Close()
		writeError(w, err)
		return
	}

	if err := tmp.Close(); err != nil {
		writeError(w, err)
		return
	}

	mimeType := header.Header.Get("Content-Type")

	result, err := h.ImportDocument.Execute(r.Context(), documentapp.ImportDocumentInput{
		CaseFileID: strings.TrimSpace(caseFileID),
		SourcePath: tmpPath,
		MimeType:   strings.TrimSpace(mimeType),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toImportDocumentResponse(result))
}
