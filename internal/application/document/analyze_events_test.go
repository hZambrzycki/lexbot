package documentapp

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	domaincasefile "lexbox/internal/domain/casefile"
	domainclient "lexbox/internal/domain/client"
	domaindocument "lexbox/internal/domain/document"
	idgen "lexbox/internal/infrastructure/idgen/uuid"
	reposqlite "lexbox/internal/infrastructure/persistence/sqlite"

	_ "modernc.org/sqlite"
)

func TestAnalyzeDocumentEvents_UsesCaseFileComputationConfig(t *testing.T) {
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
	documentContentRepo := reposqlite.NewDocumentContentRepository(db)
	eventRepo := reposqlite.NewDocumentEventRepository(db)
	ids := idgen.NewIDGenerator()
	// Cliente
	clientID := ids.NewID()
	clientEntity, err := domainclient.NewClient(clientID, "Cliente Test")
	if err != nil {
		t.Fatalf("client.NewClient: %v", err)
	}

	if err := clientRepo.Save(ctx, clientEntity); err != nil {
		t.Fatalf("ClientRepository.Save: %v", err)
	}

	// Expediente con config explícita
	cf, err := domaincasefile.NewCaseFile(ids.NewID(), clientID, "Expediente prueba")
	if err != nil {
		t.Fatalf("casefile.NewCaseFile: %v", err)
	}
	cf.Reference = "EXP-TEST-001"
	cf.Type = domaincasefile.TypeCivil
	cf.Description = "Config de calendario"
	cf.CalendarScope = "madrid"
	cf.AugustNonBusiness = true

	if err := caseFileRepo.Save(ctx, cf); err != nil {
		t.Fatalf("CaseFileRepository.Save: %v", err)
	}

	// Documento asociado al expediente
	doc, err := domaindocument.NewDocument(
		ids.NewID(),
		cf.ID,
		"auto.txt",
		filepath.Join(".lexbox", "storage", "documents", "auto.txt"),
	)
	if err != nil {
		t.Fatalf("document.NewDocument: %v", err)
	}
	doc = doc.WithUpdatedMetadata("text/plain", "hash-test")

	if err := documentRepo.Save(ctx, doc); err != nil {
		t.Fatalf("DocumentRepository.Save: %v", err)
	}

	content := "Notifíquese la resolución el 31/07/2026.\nDentro de 1 día hábil deberá aportar la documentación."
	if err := documentContentRepo.Save(ctx, doc.ID.String(), content); err != nil {
		t.Fatalf("DocumentContentRepository.Save: %v", err)
	}

	uc := AnalyzeDocumentEvents{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		CaseFiles:        caseFileRepo,
		Events:           eventRepo,
		IDs:              ids,
	}

	result, err := uc.Execute(ctx, AnalyzeDocumentEventsInput{
		DocumentID: doc.ID.String(),
	})
	if err != nil {
		t.Fatalf("AnalyzeDocumentEvents.Execute: %v", err)
	}

	if result.Detected != 2 {
		t.Fatalf("expected 2 detected events, got %d", result.Detected)
	}

	events, err := eventRepo.ListByDocumentID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("DocumentEventRepository.ListByDocumentID: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 persisted events, got %d", len(events))
	}

	var deadline *domaindocument.Event
	for i := range events {
		if events[i].EventType == domaindocument.EventTypeDeadline {
			deadline = &events[i]
			break
		}
	}

	if deadline == nil {
		t.Fatal("expected deadline event")
	}

	if deadline.EventDate != "2026-09-01" {
		t.Fatalf("expected deadline date 2026-09-01, got %s", deadline.EventDate)
	}

	if deadline.CalendarScope != "madrid" {
		t.Fatalf("expected calendar scope madrid, got %s", deadline.CalendarScope)
	}

	if !deadline.IsBusinessDays {
		t.Fatal("expected business days deadline")
	}

	if deadline.AnchorDate != "2026-07-31" {
		t.Fatalf("expected anchor date 2026-07-31, got %s", deadline.AnchorDate)
	}
}

func TestAnalyzeDocumentEvents_UsesCaseFileConfig_StateScope_And_AugustBusiness(t *testing.T) {
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
	documentContentRepo := reposqlite.NewDocumentContentRepository(db)
	eventRepo := reposqlite.NewDocumentEventRepository(db)
	ids := idgen.NewIDGenerator()

	// Cliente
	clientID := ids.NewID()
	clientEntity, err := domainclient.NewClient(clientID, "Cliente Test")
	if err != nil {
		t.Fatalf("client.NewClient: %v", err)
	}

	if err := clientRepo.Save(ctx, clientEntity); err != nil {
		t.Fatalf("ClientRepository.Save: %v", err)
	}

	// Expediente con config explícita distinta del default
	cf, err := domaincasefile.NewCaseFile(ids.NewID(), clientID, "Expediente state agosto hábil")
	if err != nil {
		t.Fatalf("casefile.NewCaseFile: %v", err)
	}
	cf.Reference = "EXP-TEST-STATE-001"
	cf.Type = domaincasefile.TypeCivil
	cf.Description = "Config state + agosto hábil"
	cf.CalendarScope = "state"
	cf.AugustNonBusiness = false

	if err := caseFileRepo.Save(ctx, cf); err != nil {
		t.Fatalf("CaseFileRepository.Save: %v", err)
	}

	// Documento asociado
	doc, err := domaindocument.NewDocument(
		ids.NewID(),
		cf.ID,
		"auto_state.txt",
		filepath.Join(".lexbox", "storage", "documents", "auto_state.txt"),
	)
	if err != nil {
		t.Fatalf("document.NewDocument: %v", err)
	}
	doc = doc.WithUpdatedMetadata("text/plain", "hash-state-test")

	if err := documentRepo.Save(ctx, doc); err != nil {
		t.Fatalf("DocumentRepository.Save: %v", err)
	}

	content := "Notifíquese la resolución el 31/07/2026.\nDentro de 1 día hábil deberá aportar la documentación."
	if err := documentContentRepo.Save(ctx, doc.ID.String(), content); err != nil {
		t.Fatalf("DocumentContentRepository.Save: %v", err)
	}

	uc := AnalyzeDocumentEvents{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		CaseFiles:        caseFileRepo,
		Events:           eventRepo,
		IDs:              ids,
	}

	result, err := uc.Execute(ctx, AnalyzeDocumentEventsInput{
		DocumentID: doc.ID.String(),
	})
	if err != nil {
		t.Fatalf("AnalyzeDocumentEvents.Execute: %v", err)
	}

	if result.Detected != 2 {
		t.Fatalf("expected 2 detected events, got %d", result.Detected)
	}

	events, err := eventRepo.ListByDocumentID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("DocumentEventRepository.ListByDocumentID: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 persisted events, got %d", len(events))
	}

	var deadline *domaindocument.Event
	for i := range events {
		if events[i].EventType == domaindocument.EventTypeDeadline {
			deadline = &events[i]
			break
		}
	}

	if deadline == nil {
		t.Fatal("expected deadline event")
	}

	if deadline.EventDate != "2026-08-03" {
		t.Fatalf("expected deadline date 2026-08-03, got %s", deadline.EventDate)
	}

	if deadline.CalendarScope != "state" {
		t.Fatalf("expected calendar scope state, got %s", deadline.CalendarScope)
	}

	if !deadline.IsBusinessDays {
		t.Fatal("expected business days deadline")
	}

	if deadline.AnchorDate != "2026-07-31" {
		t.Fatalf("expected anchor date 2026-07-31, got %s", deadline.AnchorDate)
	}
}
