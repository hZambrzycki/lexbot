package casefileapp

import (
	"context"
	"database/sql"
	documentapp "lexbox/internal/application/document"
	domaincasefile "lexbox/internal/domain/casefile"
	domainclient "lexbox/internal/domain/client"
	domaindocument "lexbox/internal/domain/document"
	idgen "lexbox/internal/infrastructure/idgen/uuid"
	reposqlite "lexbox/internal/infrastructure/persistence/sqlite"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetCaseFileDashboard_Execute_IgnoresResolvedEventsForActiveAttention(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	dbPath := filepath.Join(tmpDir, "lexbox-dashboard-test.db")
	db, err := reposqlite.NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	runDashboardTestMigrations(t, db)

	clientRepo := reposqlite.NewClientRepository(db)
	caseFileRepo := reposqlite.NewCaseFileRepository(db)
	noteRepo := reposqlite.NewNoteRepository(db)
	documentRepo := reposqlite.NewDocumentRepository(db)
	documentContentRepo := reposqlite.NewDocumentContentRepository(db)
	documentMetadataRepo := reposqlite.NewDocumentMetadataRepository(db)
	documentEventRepo := reposqlite.NewDocumentEventRepository(db)

	ids := idgen.NewIDGenerator()

	clientID := ids.NewID()
	clientEntity, err := domainclient.NewClient(clientID, "Cliente Dashboard Test")
	if err != nil {
		t.Fatalf("client.NewClient: %v", err)
	}
	if err := clientRepo.Save(ctx, clientEntity); err != nil {
		t.Fatalf("ClientRepository.Save: %v", err)
	}

	cf, err := domaincasefile.NewCaseFile(ids.NewID(), clientID, "Expediente dashboard")
	if err != nil {
		t.Fatalf("casefile.NewCaseFile: %v", err)
	}
	cf.Reference = "EXP-DASH-001"
	cf.Type = domaincasefile.TypeCivil
	cf.Description = "Test integrado dashboard"
	cf.CalendarScope = "madrid"
	cf.AugustNonBusiness = true

	if err := caseFileRepo.Save(ctx, cf); err != nil {
		t.Fatalf("CaseFileRepository.Save: %v", err)
	}

	doc, err := domaindocument.NewDocument(
		ids.NewID(),
		cf.ID,
		"dashboard_test.txt",
		filepath.Join(".lexbox", "storage", "documents", "dashboard_test.txt"),
	)
	if err != nil {
		t.Fatalf("document.NewDocument: %v", err)
	}
	doc = doc.WithUpdatedMetadata("text/plain", "hash-dashboard-test")

	if err := documentRepo.Save(ctx, doc); err != nil {
		t.Fatalf("DocumentRepository.Save: %v", err)
	}

	if err := documentContentRepo.Save(ctx, doc.ID.String(), "Texto extraído de prueba"); err != nil {
		t.Fatalf("DocumentContentRepository.Save: %v", err)
	}

	if err := documentMetadataRepo.Save(ctx, domaindocument.Metadata{
		DocumentID:   doc.ID,
		DocumentType: "order",
		LegalArea:    "procedural",
		AnalyzedAt:   time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("DocumentMetadataRepository.Save: %v", err)
	}

	now := time.Now().UTC()
	createdAt := now.Format(time.RFC3339)

	pendingDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	resolvedDate := now.AddDate(0, 0, -2).Format("2006-01-02")

	pendingEvent, err := domaindocument.NewEvent(
		ids.NewID(),
		doc.ID,
		domaindocument.EventTypeDeadline,
		pendingDate,
		"Plazo pendiente activo",
		createdAt,
		pendingDate,
		"absolute",
		"inline",
		0,
		false,
		false,
		"madrid",
		"",
		"absolute date",
	)
	if err != nil {
		t.Fatalf("document.NewEvent pending: %v", err)
	}

	resolvedEvent, err := domaindocument.NewEvent(
		ids.NewID(),
		doc.ID,
		domaindocument.EventTypeDeadline,
		resolvedDate,
		"Plazo ya resuelto",
		createdAt,
		resolvedDate,
		"absolute",
		"inline",
		0,
		false,
		false,
		"madrid",
		"",
		"absolute date",
	)
	if err != nil {
		t.Fatalf("document.NewEvent resolved: %v", err)
	}
	resolvedEvent.ReviewStatus = domaindocument.ReviewStatusResolved
	resolvedEvent.ReviewedAt = createdAt
	resolvedEvent.ResolvedAt = createdAt
	resolvedEvent.ResolutionNote = "Escrito presentado"

	if err := documentEventRepo.ReplaceByDocumentID(ctx, doc.ID, []domaindocument.Event{
		pendingEvent,
		resolvedEvent,
	}); err != nil {
		t.Fatalf("DocumentEventRepository.ReplaceByDocumentID: %v", err)
	}

	listUpcomingEvents := documentapp.ListUpcomingEvents{
		Events: documentEventRepo,
	}

	uc := GetCaseFileDashboard{
		GetCaseFileDetail: GetCaseFileDetail{
			CaseFiles: caseFileRepo,
			Notes:     noteRepo,
			Documents: documentRepo,
		},
		ListUpcomingEvents: listUpcomingEvents,
		DocumentContents:   documentContentRepo,
		Metadata:           documentMetadataRepo,
		Events:             documentEventRepo,
	}

	result, err := uc.Execute(ctx, GetCaseFileDashboardInput{
		CaseFileID: cf.ID.String(),
	})
	if err != nil {
		t.Fatalf("GetCaseFileDashboard.Execute: %v", err)
	}

	if result.DocumentCount != 1 {
		t.Fatalf("expected DocumentCount=1, got %d", result.DocumentCount)
	}

	if result.DocumentsWithoutText != 0 {
		t.Fatalf("expected DocumentsWithoutText=0, got %d", result.DocumentsWithoutText)
	}

	if result.DocumentsWithUnknownMetadata != 0 {
		t.Fatalf("expected DocumentsWithUnknownMetadata=0, got %d", result.DocumentsWithUnknownMetadata)
	}

	if result.DocumentsWithoutEvents != 0 {
		t.Fatalf("expected DocumentsWithoutEvents=0, got %d", result.DocumentsWithoutEvents)
	}

	if result.PendingReviewCount != 1 {
		t.Fatalf("expected PendingReviewCount=1, got %d", result.PendingReviewCount)
	}

	if result.ReviewedCount != 0 {
		t.Fatalf("expected ReviewedCount=0, got %d", result.ReviewedCount)
	}

	if result.ResolvedCount != 1 {
		t.Fatalf("expected ResolvedCount=1, got %d", result.ResolvedCount)
	}

	if result.ActiveEventCount != 1 {
		t.Fatalf("expected ActiveEventCount=1, got %d", result.ActiveEventCount)
	}

	if result.ResolvedEventCount != 1 {
		t.Fatalf("expected ResolvedEventCount=1, got %d", result.ResolvedEventCount)
	}

	if !result.NeedsAttention {
		t.Fatal("expected NeedsAttention=true")
	}

	if !strings.Contains(result.TopAlert, pendingDate) {
		t.Fatalf("expected TopAlert to reference pending date %q, got %q", pendingDate, result.TopAlert)
	}

	if strings.Contains(result.TopAlert, resolvedDate) {
		t.Fatalf("expected TopAlert to ignore resolved date %q, got %q", resolvedDate, result.TopAlert)
	}

	if !strings.Contains(result.RecommendedNextAction, pendingDate) {
		t.Fatalf("expected RecommendedNextAction to reference pending date %q, got %q", pendingDate, result.RecommendedNextAction)
	}

	if strings.Contains(result.RecommendedNextAction, resolvedDate) {
		t.Fatalf("expected RecommendedNextAction to ignore resolved date %q, got %q", resolvedDate, result.RecommendedNextAction)
	}

	if result.ProceduralHint != "possible deadline breach" {
		t.Fatalf("expected ProceduralHint=%q, got %q", "possible deadline breach", result.ProceduralHint)
	}

	if len(result.UpcomingEvents) != 2 {
		t.Fatalf("expected 2 UpcomingEvents, got %d", len(result.UpcomingEvents))
	}
}

func runDashboardTestMigrations(t *testing.T, db *sql.DB) {
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
	}

	for _, path := range paths {
		if err := reposqlite.RunMigrations(db, path); err != nil {
			t.Fatalf("RunMigrations(%s): %v", path, err)
		}
	}
}
