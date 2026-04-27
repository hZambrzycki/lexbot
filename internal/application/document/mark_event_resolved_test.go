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

func TestMarkEventResolved_UpdatesReviewStateAndSetsNote(t *testing.T) {
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

	cf, err := domaincasefile.NewCaseFile(ids.NewID(), clientID, "Expediente resolve")
	if err != nil {
		t.Fatalf("casefile.NewCaseFile: %v", err)
	}
	cf.Reference = "EXP-RESOLVE-001"
	cf.Type = domaincasefile.TypeCivil
	cf.Description = "Test mark resolved"

	if err := caseFileRepo.Save(ctx, cf); err != nil {
		t.Fatalf("CaseFileRepository.Save: %v", err)
	}

	doc, err := domaindocument.NewDocument(
		ids.NewID(),
		cf.ID,
		"resolve.txt",
		filepath.Join(".lexbox", "storage", "documents", "resolve.txt"),
	)
	if err != nil {
		t.Fatalf("document.NewDocument: %v", err)
	}
	doc = doc.WithUpdatedMetadata("text/plain", "hash-resolve-test")

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

	uc := MarkEventResolved{
		Events: eventRepo,
	}

	result, err := uc.Execute(ctx, MarkEventResolvedInput{
		EventID:        event.ID.String(),
		ResolutionNote: "Escrito presentado en plazo",
	})
	if err != nil {
		t.Fatalf("MarkEventResolved.Execute: %v", err)
	}

	if result.EventID != event.ID.String() {
		t.Fatalf("expected event id %q, got %q", event.ID.String(), result.EventID)
	}

	if result.ReviewStatus != domaindocument.ReviewStatusResolved {
		t.Fatalf("expected review status %q, got %q", domaindocument.ReviewStatusResolved, result.ReviewStatus)
	}

	if result.ReviewedAt == "" {
		t.Fatal("expected reviewed_at to be set")
	}

	if result.ResolvedAt == "" {
		t.Fatal("expected resolved_at to be set")
	}

	if result.ResolutionNote != "Escrito presentado en plazo" {
		t.Fatalf("expected resolution note %q, got %q", "Escrito presentado en plazo", result.ResolutionNote)
	}

	saved, err := eventRepo.GetByID(ctx, event.ID)
	if err != nil {
		t.Fatalf("DocumentEventRepository.GetByID: %v", err)
	}

	if saved.ReviewStatus != domaindocument.ReviewStatusResolved {
		t.Fatalf("expected persisted review status %q, got %q", domaindocument.ReviewStatusResolved, saved.ReviewStatus)
	}

	if saved.ReviewedAt == "" {
		t.Fatal("expected persisted reviewed_at to be set")
	}

	if saved.ResolvedAt == "" {
		t.Fatal("expected persisted resolved_at to be set")
	}

	if saved.ResolutionNote != "Escrito presentado en plazo" {
		t.Fatalf("expected persisted resolution note %q, got %q", "Escrito presentado en plazo", saved.ResolutionNote)
	}
}

func TestMarkEventResolved_PreservesExistingReviewedAt(t *testing.T) {
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

	cf, err := domaincasefile.NewCaseFile(ids.NewID(), clientID, "Expediente resolve reviewed")
	if err != nil {
		t.Fatalf("casefile.NewCaseFile: %v", err)
	}
	cf.Reference = "EXP-RESOLVE-002"
	cf.Type = domaincasefile.TypeCivil
	cf.Description = "Test preserve reviewed_at"

	if err := caseFileRepo.Save(ctx, cf); err != nil {
		t.Fatalf("CaseFileRepository.Save: %v", err)
	}

	doc, err := domaindocument.NewDocument(
		ids.NewID(),
		cf.ID,
		"resolve-reviewed.txt",
		filepath.Join(".lexbox", "storage", "documents", "resolve-reviewed.txt"),
	)
	if err != nil {
		t.Fatalf("document.NewDocument: %v", err)
	}
	doc = doc.WithUpdatedMetadata("text/plain", "hash-resolve-reviewed-test")

	if err := documentRepo.Save(ctx, doc); err != nil {
		t.Fatalf("DocumentRepository.Save: %v", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)

	event, err := domaindocument.NewEvent(
		ids.NewID(),
		doc.ID,
		domaindocument.EventTypeRequirement,
		"2026-05-03",
		"Aporte la documentación requerida",
		now,
		"2026-04-30",
		"relative",
		"notification_line",
		3,
		true,
		false,
		"madrid",
		"3 días hábiles",
		"anchor 2026-04-30 + 3 business days",
	)
	if err != nil {
		t.Fatalf("document.NewEvent: %v", err)
	}

	if err := eventRepo.ReplaceByDocumentID(ctx, doc.ID, []domaindocument.Event{event}); err != nil {
		t.Fatalf("DocumentEventRepository.ReplaceByDocumentID: %v", err)
	}

	existingReviewedAt := "2026-04-29T10:00:00Z"
	err = eventRepo.UpdateReviewState(
		ctx,
		event.ID,
		domaindocument.ReviewStatusReviewed,
		existingReviewedAt,
		"",
		"",
	)
	if err != nil {
		t.Fatalf("DocumentEventRepository.UpdateReviewState: %v", err)
	}

	uc := MarkEventResolved{
		Events: eventRepo,
	}

	result, err := uc.Execute(ctx, MarkEventResolvedInput{
		EventID:        event.ID.String(),
		ResolutionNote: "Documentación aportada",
	})
	if err != nil {
		t.Fatalf("MarkEventResolved.Execute: %v", err)
	}

	if result.ReviewStatus != domaindocument.ReviewStatusResolved {
		t.Fatalf("expected review status %q, got %q", domaindocument.ReviewStatusResolved, result.ReviewStatus)
	}

	if result.ReviewedAt != existingReviewedAt {
		t.Fatalf("expected reviewed_at %q, got %q", existingReviewedAt, result.ReviewedAt)
	}

	if result.ResolvedAt == "" {
		t.Fatal("expected resolved_at to be set")
	}

	if result.ResolutionNote != "Documentación aportada" {
		t.Fatalf("expected resolution note %q, got %q", "Documentación aportada", result.ResolutionNote)
	}

	saved, err := eventRepo.GetByID(ctx, event.ID)
	if err != nil {
		t.Fatalf("DocumentEventRepository.GetByID: %v", err)
	}

	if saved.ReviewStatus != domaindocument.ReviewStatusResolved {
		t.Fatalf("expected persisted review status %q, got %q", domaindocument.ReviewStatusResolved, saved.ReviewStatus)
	}

	if saved.ReviewedAt != existingReviewedAt {
		t.Fatalf("expected persisted reviewed_at %q, got %q", existingReviewedAt, saved.ReviewedAt)
	}

	if saved.ResolvedAt == "" {
		t.Fatal("expected persisted resolved_at to be set")
	}

	if saved.ResolutionNote != "Documentación aportada" {
		t.Fatalf("expected persisted resolution note %q, got %q", "Documentación aportada", saved.ResolutionNote)
	}
}
