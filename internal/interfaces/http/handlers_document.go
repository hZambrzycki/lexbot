package httpapi

import (
	"net/http"
	"strings"

	documentapp "lexbox/internal/application/document"
)

type DocumentHandler struct {
	GetDocumentDetail         documentapp.GetDocumentDetail
	DeleteDocument            documentapp.DeleteDocument
	ReprocessDocument         documentapp.ReprocessDocument
	UpdateDocumentReviewState documentapp.UpdateDocumentReviewState
}

func (h DocumentHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/documents/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/documents/")
		parts := strings.Split(path, "/")

		// 🔹 /documents/{id}
		if len(parts) == 1 {
			documentID := strings.TrimSpace(parts[0])
			if documentID == "" {
				http.NotFound(w, r)
				return
			}

			switch r.Method {
			case http.MethodGet:
				h.handleGetDocument(w, r, documentID)
			case http.MethodDelete:
				h.handleDeleteDocument(w, r, documentID)
			default:
				writeMethodNotAllowed(w, http.MethodGet, http.MethodDelete)
			}
			return
		}

		// 🔹 /documents/{id}/reprocess
		if len(parts) == 2 && parts[1] == "reprocess" {
			documentID := strings.TrimSpace(parts[0])
			if documentID == "" {
				http.NotFound(w, r)
				return
			}

			if r.Method != http.MethodPost {
				writeMethodNotAllowed(w, http.MethodPost)
				return
			}

			h.handleReprocessDocument(w, r, documentID)
			return
		}

		// 🔹 /documents/{id}/review
		if len(parts) == 2 && parts[1] == "review" {
			documentID := strings.TrimSpace(parts[0])
			if documentID == "" {
				http.NotFound(w, r)
				return
			}

			if r.Method != http.MethodPost {
				writeMethodNotAllowed(w, http.MethodPost)
				return
			}

			h.handleReviewDocument(w, r, documentID)
			return
		}

		http.NotFound(w, r)
	})
}

func (h DocumentHandler) handleGetDocument(w http.ResponseWriter, r *http.Request, documentID string) {
	result, err := h.GetDocumentDetail.Execute(r.Context(), documentapp.GetDocumentDetailInput{
		DocumentID: documentID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDocumentDetailResponse(result))
}

func (h DocumentHandler) handleDeleteDocument(w http.ResponseWriter, r *http.Request, documentID string) {
	result, err := h.DeleteDocument.Execute(r.Context(), documentapp.DeleteDocumentInput{
		DocumentID: documentID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h DocumentHandler) handleReprocessDocument(w http.ResponseWriter, r *http.Request, documentID string) {
	result, err := h.ReprocessDocument.Execute(r.Context(), documentapp.ReprocessDocumentInput{
		DocumentID: documentID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h DocumentHandler) handleReviewDocument(w http.ResponseWriter, r *http.Request, documentID string) {
	var body struct {
		ReviewStatus string `json:"review_status"`
		ReviewNote   string `json:"review_note"`
	}

	if err := readJSON(r, &body); err != nil {
		writeError(w, err)
		return
	}

	result, err := h.UpdateDocumentReviewState.Execute(
		r.Context(),
		documentapp.UpdateDocumentReviewStateInput{
			DocumentID:   documentID,
			ReviewStatus: body.ReviewStatus,
			ReviewNote:   body.ReviewNote,
		},
	)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
