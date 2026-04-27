package documentapp

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	domaincasefile "lexbox/internal/domain/casefile"
	domainclient "lexbox/internal/domain/client"
	domaindocument "lexbox/internal/domain/document"
	idgen "lexbox/internal/infrastructure/idgen/uuid"
	reposqlite "lexbox/internal/infrastructure/persistence/sqlite"

	_ "modernc.org/sqlite"
)

func TestReopenEvent_ResetsReviewStateToPending(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	defer db.Close()

	runTestMigrations(t, db)

	ctx := context.Background()

	clientRepo := reposqlite.NewClientRepository(db)
	caseFileRepo := reposqlite.NewCaseFileRepository(db)
	documentRepo := reposqlite.NewDocumentRepository(db)
	eventRepo := reposqlite.NewDocumentEventRepository(db)
	ids := idgen.NewIDGenerator()

	clientID := ids.NewID()
	clientEntity, err := domainclient.NewClient(clientID, "Cliente Test")
	if err != nil {
		t.Fatalf("client.NewClient: %v", err)
	}
	if err := clientRepo.Save(ctx, clientEntity); err != nil {
		t.Fatalf("ClientRepository.Save: %v", err)
	}

	cf, err := domaincasefile.NewCaseFile(ids.NewID(), clientID, "Expediente reopen")
	if err != nil {
		t.Fatalf("casefile.NewCaseFile: %v", err)
	}
	cf.Reference = "EXP-REOPEN-001"
	cf.Type = domaincasefile.TypeCivil
	cf.Description = "Test reopen event"

	if err := caseFileRepo.Save(ctx, cf); err != nil {
		t.Fatalf("CaseFileRepository.Save: %v", err)
	}

	doc, err := domaindocument.NewDocument(
		ids.NewID(),
		cf.ID,
		"reopen.txt",
		filepath.Join(".lexbox", "storage", "documents", "reopen.txt"),
	)
	if err != nil {
		t.Fatalf("document.NewDocument: %v", err)
	}
	doc = doc.WithUpdatedMetadata("text/plain", "hash-reopen-test")

	if err := documentRepo.Save(ctx, doc); err != nil {
		t.Fatalf("DocumentRepository.Save: %v", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)

	event, err := domaindocument.NewEvent(
		ids.NewID(),
		doc.ID,
		domaindocument.EventTypeDeadline,
		"2026-05-01",
		"Se concede plazo de 5 días",
		now,
		"2026-04-26",
		"relative",
		"notification_line",
		5,
		false,
		false,
		"madrid",
		"plazo de 5 días",
		"anchor 2026-04-26 + 5 natural days",
	)
	if err != nil {
		t.Fatalf("document.NewEvent: %v", err)
	}

	if err := eventRepo.ReplaceByDocumentID(ctx, doc.ID, []domaindocument.Event{event}); err != nil {
		t.Fatalf("DocumentEventRepository.ReplaceByDocumentID: %v", err)
	}

	err = eventRepo.UpdateReviewState(
		ctx,
		event.ID,
		domaindocument.ReviewStatusResolved,
		"2026-04-27T09:00:00Z",
		"2026-04-27T10:00:00Z",
		"Escrito presentado",
	)
	if err != nil {
		t.Fatalf("DocumentEventRepository.UpdateReviewState: %v", err)
	}

	uc := ReopenEvent{
		Events: eventRepo,
	}

	result, err := uc.Execute(ctx, ReopenEventInput{
		EventID: event.ID.String(),
	})
	if err != nil {
		t.Fatalf("ReopenEvent.Execute: %v", err)
	}

	if result.EventID != event.ID.String() {
		t.Fatalf("expected event id %q, got %q", event.ID.String(), result.EventID)
	}

	if result.ReviewStatus != domaindocument.ReviewStatusPending {
		t.Fatalf("expected review status %q, got %q", domaindocument.ReviewStatusPending, result.ReviewStatus)
	}

	if result.ReviewedAt != "" {
		t.Fatalf("expected empty reviewed_at, got %q", result.ReviewedAt)
	}

	if result.ResolvedAt != "" {
		t.Fatalf("expected empty resolved_at, got %q", result.ResolvedAt)
	}

	if result.ResolutionNote != "" {
		t.Fatalf("expected empty resolution_note, got %q", result.ResolutionNote)
	}

	saved, err := eventRepo.GetByID(ctx, event.ID)
	if err != nil {
		t.Fatalf("DocumentEventRepository.GetByID: %v", err)
	}

	if saved.ReviewStatus != domaindocument.ReviewStatusPending {
		t.Fatalf("expected persisted review status %q, got %q", domaindocument.ReviewStatusPending, saved.ReviewStatus)
	}

	if saved.ReviewedAt != "" {
		t.Fatalf("expected persisted reviewed_at empty, got %q", saved.ReviewedAt)
	}

	if saved.ResolvedAt != "" {
		t.Fatalf("expected persisted resolved_at empty, got %q", saved.ResolvedAt)
	}

	if saved.ResolutionNote != "" {
		t.Fatalf("expected persisted resolution_note empty, got %q", saved.ResolutionNote)
	}
}
