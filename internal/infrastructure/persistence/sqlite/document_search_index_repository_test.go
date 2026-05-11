package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"lexbox/internal/application/querymodels"
)

func TestDocumentSearchIndexRepository_UpsertAndSearch(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"automatizacion.txt",
		"text/plain",
		"pending_review",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"automatizacion.txt",
		"Documento sobre automatización procesal y ejecución jurídica.",
		"claim",
		"civil",
	)
	if err != nil {
		t.Fatalf("UpsertDocument: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"automatizacion",
		"",
		querymodels.SearchDocumentFilters{},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", results[0].DocumentID)
	}

	assertSearchIndexContains(t, results[0].Snippet, "[automatización]")
}

func TestDocumentSearchIndexRepository_SearchByCaseFile(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustUpsertSearchIndexFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"Demanda por despido improcedente.",
	)
	mustUpsertSearchIndexFixture(
		t,
		ctx,
		repo,
		"doc-2",
		"case-2",
		"Demanda por despido objetivo.",
	)

	results, err := repo.Search(
		ctx,
		"despido",
		"case-1",
		querymodels.SearchDocumentFilters{},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", results[0].DocumentID)
	}
}

func TestDocumentSearchIndexRepository_DeleteDocument(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustUpsertSearchIndexFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"Demanda por despido.",
	)

	if err := repo.DeleteDocument(ctx, "doc-1"); err != nil {
		t.Fatalf("DeleteDocument: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"despido",
		"",
		querymodels.SearchDocumentFilters{},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 results after delete, got %d", len(results))
	}
}

func TestDocumentSearchIndexRepository_UpsertReplacesExistingDocument(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustUpsertSearchIndexFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"Demanda por despido.",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"sentencia.txt",
		"Sentencia sobre reclamación de cantidad.",
		"judgment",
		"labor",
	)
	if err != nil {
		t.Fatalf("UpsertDocument replace: %v", err)
	}

	oldResults, err := repo.Search(
		ctx,
		"despido",
		"",
		querymodels.SearchDocumentFilters{},
		10,
	)
	if err != nil {
		t.Fatalf("Search old term: %v", err)
	}

	if len(oldResults) != 0 {
		t.Fatalf("expected 0 old results after replace, got %d", len(oldResults))
	}

	newResults, err := repo.Search(
		ctx,
		"sentencia",
		"",
		querymodels.SearchDocumentFilters{},
		10,
	)
	if err != nil {
		t.Fatalf("Search new term: %v", err)
	}

	if len(newResults) != 1 {
		t.Fatalf("expected 1 new result after replace, got %d", len(newResults))
	}

	if newResults[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", newResults[0].DocumentID)
	}
}

func TestDocumentSearchIndexRepository_SearchFindsMetadata(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"contrato_arrendamiento.pdf",
		"application/pdf",
		"pending_review",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"contrato_arrendamiento.pdf",
		"Documento auxiliar sin la palabra buscada en contenido.",
		"contract",
		"civil",
	)
	if err != nil {
		t.Fatalf("UpsertDocument: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"arrendamiento",
		"",
		querymodels.SearchDocumentFilters{},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result by original_name, got %d", len(results))
	}
}

func TestDocumentSearchIndexRepository_SearchByDocumentTypeFilter(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"demanda.txt",
		"text/plain",
		"pending_review",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"demanda.txt",
		"Documento sobre despido.",
		"claim",
		"labor",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-1: %v", err)
	}

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-2",
		"case-1",
		"sentencia.txt",
		"text/plain",
		"pending_review",
	)

	err = repo.UpsertDocument(
		ctx,
		"doc-2",
		"case-1",
		"sentencia.txt",
		"Documento sobre despido.",
		"judgment",
		"labor",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-2: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"despido",
		"",
		querymodels.SearchDocumentFilters{
			DocumentType: "claim",
		},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result by document type, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", results[0].DocumentID)
	}
}

func TestDocumentSearchIndexRepository_SearchByLegalAreaFilter(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"demanda_laboral.txt",
		"text/plain",
		"pending_review",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"demanda_laboral.txt",
		"Documento sobre recurso.",
		"claim",
		"labor",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-1: %v", err)
	}

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-2",
		"case-1",
		"demanda_civil.txt",
		"text/plain",
		"pending_review",
	)

	err = repo.UpsertDocument(
		ctx,
		"doc-2",
		"case-1",
		"demanda_civil.txt",
		"Documento sobre recurso.",
		"claim",
		"civil",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-2: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"recurso",
		"",
		querymodels.SearchDocumentFilters{
			LegalArea: "labor",
		},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result by legal area, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", results[0].DocumentID)
	}
}

func TestDocumentSearchIndexRepository_SearchByReviewStatusFilter(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"pendiente.txt",
		"text/plain",
		"pending_review",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"pendiente.txt",
		"Documento sobre ejecución.",
		"claim",
		"civil",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-1: %v", err)
	}

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-2",
		"case-1",
		"error.txt",
		"text/plain",
		"error",
	)

	err = repo.UpsertDocument(
		ctx,
		"doc-2",
		"case-1",
		"error.txt",
		"Documento sobre ejecución.",
		"claim",
		"civil",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-2: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"ejecución",
		"",
		querymodels.SearchDocumentFilters{
			ReviewStatus: "error",
		},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result by review status, got %d", len(results))
	}

	if results[0].DocumentID != "doc-2" {
		t.Fatalf("expected doc-2, got %s", results[0].DocumentID)
	}
}

func TestDocumentSearchIndexRepository_SearchByDocTypeFilter(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"demanda.pdf",
		"application/pdf",
		"pending_review",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"demanda.pdf",
		"Documento sobre alimentos.",
		"claim",
		"family",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-1: %v", err)
	}

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-2",
		"case-1",
		"demanda.txt",
		"text/plain",
		"pending_review",
	)

	err = repo.UpsertDocument(
		ctx,
		"doc-2",
		"case-1",
		"demanda.txt",
		"Documento sobre alimentos.",
		"claim",
		"family",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-2: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"alimentos",
		"",
		querymodels.SearchDocumentFilters{
			DocType: "pdf",
		},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result by doc type, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", results[0].DocumentID)
	}
}

func TestDocumentSearchIndexRepository_SearchByHasEventsFilter(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustUpsertSearchIndexFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"Documento con requerimiento.",
	)
	mustUpsertSearchIndexFixture(
		t,
		ctx,
		repo,
		"doc-2",
		"case-1",
		"Documento con requerimiento.",
	)

	mustInsertEventFixture(t, ctx, repo, "event-1", "doc-1")

	results, err := repo.Search(
		ctx,
		"requerimiento",
		"",
		querymodels.SearchDocumentFilters{
			Has: "events",
		},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result with events, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", results[0].DocumentID)
	}
}

func TestDocumentSearchIndexRepository_SearchByHasNoTextFilter(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		"doc-1",
		"case-1",
		"sin_texto.txt",
		"text/plain",
		"pending_review",
	)

	err := repo.UpsertDocument(
		ctx,
		"doc-1",
		"case-1",
		"sin_texto.txt",
		"Documento indexado pero sin contenido extraído.",
		"unknown",
		"unknown",
	)
	if err != nil {
		t.Fatalf("UpsertDocument doc-1: %v", err)
	}

	results, err := repo.Search(
		ctx,
		"documento",
		"",
		querymodels.SearchDocumentFilters{
			Has: "no_text",
		},
		10,
	)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result without text, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected doc-1, got %s", results[0].DocumentID)
	}
}

func mustOpenFTS5TestDB(t *testing.T) *sql.DB {
	t.Helper()

	db := mustOpenTestDB(t)

	_, err := db.Exec(`
		CREATE TABLE documents (
			id TEXT PRIMARY KEY,
			case_file_id TEXT NOT NULL,
			original_name TEXT NOT NULL,
			storage_path TEXT NOT NULL,
			mime_type TEXT,
			file_hash TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			review_status TEXT NOT NULL DEFAULT 'pending_review',
			reviewed_at TEXT NOT NULL DEFAULT '',
			review_note TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE document_contents (
			document_id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			extracted_at TEXT NOT NULL
		);

		CREATE TABLE document_events (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			event_date TEXT NOT NULL,
			source_text TEXT NOT NULL,
			created_at TEXT NOT NULL
		);

		CREATE VIRTUAL TABLE document_search_index USING fts5(
			document_id UNINDEXED,
			case_file_id UNINDEXED,
			original_name,
			content,
			document_type,
			legal_area,
			tokenize = 'unicode61 remove_diacritics 2'
		)
	`)
	if err != nil {
		t.Fatalf("create test search schema: %v", err)
	}

	return db
}

func mustInsertDocumentFixture(
	t *testing.T,
	ctx context.Context,
	repo *DocumentSearchIndexRepository,
	documentID string,
	caseFileID string,
	originalName string,
	mimeType string,
	reviewStatus string,
) {
	t.Helper()

	if reviewStatus == "" {
		reviewStatus = "pending_review"
	}

	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO documents (
			id,
			case_file_id,
			original_name,
			storage_path,
			mime_type,
			file_hash,
			created_at,
			updated_at,
			review_status,
			reviewed_at,
			review_note
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		documentID,
		caseFileID,
		originalName,
		"/tmp/"+originalName,
		mimeType,
		"hash-"+documentID,
		"2026-01-01T00:00:00Z",
		"2026-01-01T00:00:00Z",
		reviewStatus,
		"",
		"",
	)
	if err != nil {
		t.Fatalf("insert document fixture: %v", err)
	}
}

func mustInsertDocumentContentFixture(
	t *testing.T,
	ctx context.Context,
	repo *DocumentSearchIndexRepository,
	documentID string,
	content string,
) {
	t.Helper()

	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO document_contents (
			document_id,
			content,
			extracted_at
		)
		VALUES (?, ?, ?)
	`,
		documentID,
		content,
		"2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("insert document content fixture: %v", err)
	}
}

func mustInsertEventFixture(
	t *testing.T,
	ctx context.Context,
	repo *DocumentSearchIndexRepository,
	eventID string,
	documentID string,
) {
	t.Helper()

	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO document_events (
			id,
			document_id,
			event_type,
			event_date,
			source_text,
			created_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		eventID,
		documentID,
		"deadline",
		"2026-01-10",
		"Requerimiento de subsanación",
		"2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("insert event fixture: %v", err)
	}
}

func mustUpsertSearchIndexFixture(
	t *testing.T,
	ctx context.Context,
	repo *DocumentSearchIndexRepository,
	documentID string,
	caseFileID string,
	content string,
) {
	t.Helper()

	mustInsertDocumentFixture(
		t,
		ctx,
		repo,
		documentID,
		caseFileID,
		documentID+".txt",
		"text/plain",
		"pending_review",
	)

	mustInsertDocumentContentFixture(t, ctx, repo, documentID, content)

	err := repo.UpsertDocument(
		ctx,
		documentID,
		caseFileID,
		documentID+".txt",
		content,
		"unknown",
		"unknown",
	)
	if err != nil {
		t.Fatalf("UpsertDocument fixture: %v", err)
	}
}

func assertSearchIndexContains(t *testing.T, value string, expected string) {
	t.Helper()

	if !strings.Contains(value, expected) {
		t.Fatalf("expected %q to contain %q", value, expected)
	}
}
