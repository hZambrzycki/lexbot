package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"testing"
)

func TestDocumentSearchIndexRepository_UpsertAndSearch(t *testing.T) {
	t.Parallel()

	db := mustOpenFTS5TestDB(t)
	repo := NewDocumentSearchIndexRepository(db)

	ctx := context.Background()

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

	results, err := repo.Search(ctx, "automatizacion", "", 10)
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

	mustUpsertSearchIndexFixture(t, ctx, repo, "doc-1", "case-1", "Demanda por despido improcedente.")
	mustUpsertSearchIndexFixture(t, ctx, repo, "doc-2", "case-2", "Demanda por despido objetivo.")

	results, err := repo.Search(ctx, "despido", "case-1", 10)
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

	mustUpsertSearchIndexFixture(t, ctx, repo, "doc-1", "case-1", "Demanda por despido.")

	if err := repo.DeleteDocument(ctx, "doc-1"); err != nil {
		t.Fatalf("DeleteDocument: %v", err)
	}

	results, err := repo.Search(ctx, "despido", "", 10)
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

	mustUpsertSearchIndexFixture(t, ctx, repo, "doc-1", "case-1", "Demanda por despido.")

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

	oldResults, err := repo.Search(ctx, "despido", "", 10)
	if err != nil {
		t.Fatalf("Search old term: %v", err)
	}

	if len(oldResults) != 0 {
		t.Fatalf("expected 0 old results after replace, got %d", len(oldResults))
	}

	newResults, err := repo.Search(ctx, "sentencia", "", 10)
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

	results, err := repo.Search(ctx, "arrendamiento", "", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result by original_name, got %d", len(results))
	}
}

func mustOpenFTS5TestDB(t *testing.T) *sql.DB {
	t.Helper()

	db := mustOpenTestDB(t)

	_, err := db.Exec(`
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
		t.Fatalf("create document_search_index: %v", err)
	}

	return db
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
