package httpapi

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	casefileapp "lexbox/internal/application/casefile"
	documentapp "lexbox/internal/application/document"
	noteapp "lexbox/internal/application/note"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/client"
	"lexbox/internal/domain/document"
	idgen "lexbox/internal/infrastructure/idgen/uuid"
	reposqlite "lexbox/internal/infrastructure/persistence/sqlite"

	_ "modernc.org/sqlite"
)

type httpTestDeps struct {
	db                   *sql.DB
	ids                  *idgen.IDGenerator
	clientRepo           *reposqlite.ClientRepository
	caseFileRepo         *reposqlite.CaseFileRepository
	noteRepo             *reposqlite.NoteRepository
	documentRepo         *reposqlite.DocumentRepository
	documentContentRepo  *reposqlite.DocumentContentRepository
	documentMetadataRepo *reposqlite.DocumentMetadataRepository
	documentEventRepo    *reposqlite.DocumentEventRepository
}

func setupHTTPTestDeps(t *testing.T) httpTestDeps {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	runHTTPTestMigrations(t, db)

	return httpTestDeps{
		db:                   db,
		ids:                  idgen.NewIDGenerator(),
		clientRepo:           reposqlite.NewClientRepository(db),
		caseFileRepo:         reposqlite.NewCaseFileRepository(db),
		noteRepo:             reposqlite.NewNoteRepository(db),
		documentRepo:         reposqlite.NewDocumentRepository(db),
		documentContentRepo:  reposqlite.NewDocumentContentRepository(db),
		documentMetadataRepo: reposqlite.NewDocumentMetadataRepository(db),
		documentEventRepo:    reposqlite.NewDocumentEventRepository(db),
	}
}

func runHTTPTestMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	paths := []string{
		filepath.Join("..", "..", "..", "migrations", "001_init.sql"),
		filepath.Join("..", "..", "..", "migrations", "002_document_contents.sql"),
		filepath.Join("..", "..", "..", "migrations", "003_document_metadata.sql"),
		filepath.Join("..", "..", "..", "migrations", "004_document_events.sql"),
		filepath.Join("..", "..", "..", "migrations", "005_document_event_context.sql"),
		filepath.Join("..", "..", "..", "migrations", "006_document_event_semantics.sql"),
		filepath.Join("..", "..", "..", "migrations", "007_case_file_event_config.sql"),
		filepath.Join("..", "..", "..", "migrations", "008_document_event_review_state.sql"),
		filepath.Join("..", "..", "..", "migrations", "009_document_review_state.sql"),
	}

	for _, path := range paths {
		if err := reposqlite.RunMigrations(db, path); err != nil {
			t.Fatalf("RunMigrations(%s): %v", path, err)
		}
	}
}

//
// ==========================
// TEST 1: DASHBOARD
// ==========================
//

func TestCaseFileDashboardEndpoint(t *testing.T) {
	ctx := context.Background()
	deps := setupHTTPTestDeps(t)
	clientID := deps.ids.NewID()
	cl, _ := client.NewClient(clientID, "Cliente HTTP")

	_ = deps.clientRepo.Save(ctx, cl)

	cf, _ := casefile.NewCaseFile(deps.ids.NewID(), clientID, "Expediente HTTP")
	cf.Reference = "EXP-HTTP-1"
	cf.Type = casefile.TypeCivil
	cf.CalendarScope = "madrid"

	_ = deps.caseFileRepo.Save(ctx, cf)

	doc, _ := document.NewDocument(
		deps.ids.NewID(),
		cf.ID,
		"doc.txt",
		filepath.Join(".lexbox", "doc.txt"),
	)

	doc = doc.WithUpdatedMetadata("text/plain", "hash")
	_ = deps.documentRepo.Save(ctx, doc)

	_ = deps.documentContentRepo.Save(ctx, doc.ID.String(), "contenido")

	eventDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	ev, _ := document.NewEvent(
		deps.ids.NewID(),
		doc.ID,
		document.EventTypeDeadline,
		eventDate,
		"evento",
		time.Now().Format(time.RFC3339),
		eventDate,
		"absolute",
		"inline",
		0,
		false,
		false,
		"madrid",
		"",
		"absolute date",
	)

	_ = deps.documentEventRepo.ReplaceByDocumentID(ctx, doc.ID, []document.Event{ev})

	handler := CaseFileHandler{
		GetCaseFileDashboard: casefileapp.GetCaseFileDashboard{
			GetCaseFileDetail: casefileapp.GetCaseFileDetail{
				CaseFiles: deps.caseFileRepo,
				Notes:     deps.noteRepo,
				Documents: deps.documentRepo,
			},
			ListUpcomingEvents: documentapp.ListUpcomingEvents{
				Events: deps.documentEventRepo,
			},
			DocumentContents: deps.documentContentRepo,
			Metadata:         deps.documentMetadataRepo,
			Events:           deps.documentEventRepo,
		},
		AddNote: noteapp.AddNote{
			Notes:     deps.noteRepo,
			CaseFiles: deps.caseFileRepo,
			IDs:       deps.ids,
		},

		DeleteNote: noteapp.DeleteNote{
			Notes: deps.noteRepo,
		},
	}

	mux := http.NewServeMux()
	handler.Register(mux)

	req := httptest.NewRequest("GET", "/case-files/"+cf.ID.String()+"/dashboard", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status != 200")
	}

	if !strings.Contains(rec.Body.String(), eventDate) {
		t.Fatalf("dashboard no contiene evento")
	}
}

//
// ==========================
// TEST 2: RESOLVE
// ==========================
//

func TestEventResolveEndpoint(t *testing.T) {
	ctx := context.Background()
	deps := setupHTTPTestDeps(t)
	clientID := deps.ids.NewID()
	cl, _ := client.NewClient(clientID, "Cliente")
	_ = deps.clientRepo.Save(ctx, cl)

	cf, _ := casefile.NewCaseFile(deps.ids.NewID(), clientID, "Expediente")
	_ = deps.caseFileRepo.Save(ctx, cf)

	doc, _ := document.NewDocument(
		deps.ids.NewID(),
		cf.ID,
		"doc.txt",
		"/tmp/doc.txt",
	)

	doc = doc.WithUpdatedMetadata("text/plain", "hash")
	_ = deps.documentRepo.Save(ctx, doc)

	ev, _ := document.NewEvent(
		deps.ids.NewID(),
		doc.ID,
		document.EventTypeDeadline,
		"2026-04-30",
		"evento",
		time.Now().Format(time.RFC3339),
		"2026-04-30",
		"absolute",
		"inline",
		0,
		false,
		false,
		"madrid",
		"",
		"absolute date",
	)

	_ = deps.documentEventRepo.ReplaceByDocumentID(ctx, doc.ID, []document.Event{ev})

	handler := EventHandler{
		MarkEventResolved: documentapp.MarkEventResolved{
			Events: deps.documentEventRepo,
		},
	}

	mux := http.NewServeMux()
	handler.Register(mux)

	body := `{"resolution_note":"ok"}`
	req := httptest.NewRequest("POST", "/events/"+ev.ID.String()+"/resolve", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status != 200")
	}

	if !strings.Contains(rec.Body.String(), "resolved") {
		t.Fatalf("no se resolvió")
	}
}

//
// ==========================
// TEST 3: REOPEN
// ==========================
//

func TestEventReopenEndpoint(t *testing.T) {
	ctx := context.Background()
	deps := setupHTTPTestDeps(t)
	clientID := deps.ids.NewID()
	cl, _ := client.NewClient(clientID, "Cliente")
	_ = deps.clientRepo.Save(ctx, cl)

	cf, _ := casefile.NewCaseFile(deps.ids.NewID(), clientID, "Expediente")
	_ = deps.caseFileRepo.Save(ctx, cf)

	doc, _ := document.NewDocument(
		deps.ids.NewID(),
		cf.ID,
		"doc.txt",
		"/tmp/doc.txt",
	)

	doc = doc.WithUpdatedMetadata("text/plain", "hash")
	_ = deps.documentRepo.Save(ctx, doc)

	ev, _ := document.NewEvent(
		deps.ids.NewID(),
		doc.ID,
		document.EventTypeDeadline,
		"2026-04-30",
		"evento",
		time.Now().Format(time.RFC3339),
		"2026-04-30",
		"absolute",
		"inline",
		0,
		false,
		false,
		"madrid",
		"",
		"absolute date",
	)

	ev.ReviewStatus = document.ReviewStatusResolved
	ev.ResolvedAt = time.Now().Format(time.RFC3339)

	_ = deps.documentEventRepo.ReplaceByDocumentID(ctx, doc.ID, []document.Event{ev})

	handler := EventHandler{
		ReopenEvent: documentapp.ReopenEvent{
			Events: deps.documentEventRepo,
		},
	}

	mux := http.NewServeMux()
	handler.Register(mux)

	req := httptest.NewRequest("POST", "/events/"+ev.ID.String()+"/reopen", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status != 200")
	}

	if !strings.Contains(rec.Body.String(), "pending") {
		t.Fatalf("no se reabrió")
	}
}
