package documentapp

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"lexbox/internal/domain/document"
	localextractor "lexbox/internal/infrastructure/extraction/local"
	sha256hasher "lexbox/internal/infrastructure/hash/sha256"
	idgen "lexbox/internal/infrastructure/idgen/uuid"
	reposqlite "lexbox/internal/infrastructure/persistence/sqlite"
	localstorage "lexbox/internal/infrastructure/storage/local"
)

func TestImportDocument_ExtractsMetadataAndEventsAutomatically(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tmpDir := t.TempDir()

	dbPath := filepath.Join(tmpDir, "lexbox-test.db")
	db, err := reposqlite.NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	runTestMigrations(t, db)

	now := time.Now().Format(time.RFC3339)

	_, err = db.ExecContext(ctx, `
		INSERT INTO clients (id, name, email, phone, identifier, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		"client-1",
		"Cliente Test",
		"cliente@test.local",
		"600000000",
		"TEST-001",
		now,
		now,
	)
	if err != nil {
		t.Fatalf("insert client: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO case_files (id, client_id, reference, title, type, status, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		"case-1",
		"client-1",
		"EXP-TEST-001",
		"Expediente de prueba",
		"civil",
		"open",
		"Expediente para test de import automático",
		now,
		now,
	)
	if err != nil {
		t.Fatalf("insert case file: %v", err)
	}

	sourcePath := filepath.Join(tmpDir, "eventos_import_auto.txt")
	sourceContent := strings.TrimSpace(`
DILIGENCIA DE ORDENACIÓN de 09/04/2026.

Notifíquese la presente resolución a las partes el 11/04/2026.

Se concede plazo de 5 días para formular alegaciones.

Se señala juicio para el día 10/05/2026.

Requiérase a la parte demandada para que aporte la documentación dentro de 3 días hábiles.

La comparecencia tendrá lugar el 15/05/2026.
`)
	if err := os.WriteFile(sourcePath, []byte(sourceContent), 0644); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	documentRepo := reposqlite.NewDocumentRepository(db)
	documentContentRepo := reposqlite.NewDocumentContentRepository(db)
	documentMetadataRepo := reposqlite.NewDocumentMetadataRepository(db)
	documentEventRepo := reposqlite.NewDocumentEventRepository(db)
	caseFileRepo := reposqlite.NewCaseFileRepository(db)

	idGenerator := idgen.NewIDGenerator()
	storage := localstorage.NewStorage(filepath.Join(tmpDir, ".lexbox", "storage"))
	textExtractor := localextractor.NewTextExtractor()
	fileHasher := sha256hasher.NewFileHasher()

	analyzeDocumentEvents := AnalyzeDocumentEvents{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		Events:           documentEventRepo,
		IDs:              idGenerator,
	}

	uc := ImportDocument{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		Metadata:         documentMetadataRepo,
		CaseFiles:        caseFileRepo,
		Storage:          storage,
		Extractor:        textExtractor,
		Hasher:           fileHasher,
		IDs:              idGenerator,
		AnalyzeEvents:    analyzeDocumentEvents,
	}

	result, err := uc.Execute(ctx, ImportDocumentInput{
		CaseFileID: "case-1",
		SourcePath: sourcePath,
		MimeType:   "text/plain",
	})
	if err != nil {
		t.Fatalf("ImportDocument.Execute: %v", err)
	}

	if result.Document.ID == "" {
		t.Fatal("expected imported document ID")
	}
	if !result.TextExtracted {
		t.Fatal("expected text to be extracted")
	}
	if !result.MetadataAnalyzed {
		t.Fatal("expected metadata to be analyzed")
	}
	if !result.EventsAnalyzed {
		t.Fatal("expected events to be analyzed")
	}
	if result.EventsDetected != 5 {
		t.Fatalf("expected 5 detected events, got %d", result.EventsDetected)
	}

	savedDoc, err := documentRepo.GetByID(ctx, result.Document.ID)
	if err != nil {
		t.Fatalf("DocumentRepository.GetByID: %v", err)
	}
	if savedDoc.OriginalName != "eventos_import_auto.txt" {
		t.Fatalf("unexpected original name: %s", savedDoc.OriginalName)
	}
	if savedDoc.CaseFileID.String() != "case-1" {
		t.Fatalf("unexpected case_file_id: %s", savedDoc.CaseFileID)
	}

	savedContent, err := documentContentRepo.GetByDocumentID(ctx, result.Document.ID.String())
	if err != nil {
		t.Fatalf("DocumentContentRepository.GetByDocumentID: %v", err)
	}
	if !strings.Contains(savedContent, "Notifíquese la presente resolución") {
		t.Fatalf("expected extracted content to contain notification line, got: %q", savedContent)
	}

	metadata, err := documentMetadataRepo.GetByDocumentID(ctx, result.Document.ID)
	if err != nil {
		t.Fatalf("DocumentMetadataRepository.GetByDocumentID: %v", err)
	}
	if metadata.DocumentType != "order" {
		t.Fatalf("expected document type %q, got %q", "order", metadata.DocumentType)
	}
	if metadata.LegalArea != "procedural" {
		t.Fatalf("expected legal area %q, got %q", "procedural", metadata.LegalArea)
	}

	events, err := documentEventRepo.ListByDocumentID(ctx, result.Document.ID)
	if err != nil {
		t.Fatalf("DocumentEventRepository.ListByDocumentID: %v", err)
	}
	if len(events) != 5 {
		t.Fatalf("expected 5 events, got %d", len(events))
	}

	assertHasPersistedEvent(t, events, "notification", "2026-04-11")
	assertHasPersistedEvent(t, events, "requirement", "2026-04-15")
	assertHasPersistedEvent(t, events, "deadline", "2026-04-16")
	assertHasPersistedEvent(t, events, "hearing", "2026-05-10")
	assertHasPersistedEvent(t, events, "appearance", "2026-05-15")
}

func runTestMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	paths := []string{
		filepath.Join("..", "..", "..", "migrations", "001_init.sql"),
		filepath.Join("..", "..", "..", "migrations", "002_document_contents.sql"),
		filepath.Join("..", "..", "..", "migrations", "003_document_metadata.sql"),
		filepath.Join("..", "..", "..", "migrations", "004_document_events.sql"),
		filepath.Join("..", "..", "..", "migrations", "005_document_event_context.sql"),
		filepath.Join("..", "..", "..", "migrations", "006_document_event_semantics.sql"),
	}
	for _, path := range paths {
		if err := reposqlite.RunMigrations(db, path); err != nil {
			t.Fatalf("RunMigrations(%s): %v", path, err)
		}
	}
}

func assertHasPersistedEvent(t *testing.T, events []document.Event, wantType, wantDate string) {
	t.Helper()

	for _, e := range events {
		if e.EventType == wantType && e.EventDate == wantDate {
			return
		}
	}

	t.Fatalf("expected persisted event type=%q date=%q not found", wantType, wantDate)
}
