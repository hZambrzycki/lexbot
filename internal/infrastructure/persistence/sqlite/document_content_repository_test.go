package sqlite

import (
	"context"
	"database/sql"
	"testing"
)

func TestDocumentContentRepository_SaveAndGetByDocumentID(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	repo := NewDocumentContentRepository(db)

	createDocumentsSchemaForSearchTests(t, db)

	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
		INSERT INTO documents (id, case_file_id, original_name)
		VALUES ('doc-1', 'case-1', 'demanda.txt')
	`)
	if err != nil {
		t.Fatalf("insert document: %v", err)
	}

	err = repo.Save(ctx, "doc-1", "contenido indexado de prueba")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.GetByDocumentID(ctx, "doc-1")
	if err != nil {
		t.Fatalf("GetByDocumentID: %v", err)
	}

	if got != "contenido indexado de prueba" {
		t.Fatalf("unexpected content\nwant: %q\ngot:  %q", "contenido indexado de prueba", got)
	}
}

func TestDocumentContentRepository_SearchByText(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	repo := NewDocumentContentRepository(db)

	createDocumentsSchemaForSearchTests(t, db)

	ctx := context.Background()

	insertSearchFixture(t, ctx, db,
		"doc-1", "case-1", "demanda_laboral.txt",
		"Demanda por despido improcedente. El despido fue comunicado sin preaviso.",
	)

	insertSearchFixture(t, ctx, db,
		"doc-2", "case-2", "resolucion_extranjeria.txt",
		"Resolución de extranjería sobre autorización de residencia temporal.",
	)

	results, err := repo.SearchByText(ctx, "despido", 10)
	if err != nil {
		t.Fatalf("SearchByText: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected document_id doc-1, got %s", results[0].DocumentID)
	}

	if results[0].OriginalName != "demanda_laboral.txt" {
		t.Fatalf("unexpected original_name: %s", results[0].OriginalName)
	}

	if results[0].CaseFileID != "case-1" {
		t.Fatalf("unexpected case_file_id: %s", results[0].CaseFileID)
	}

	if results[0].Snippet == "" {
		t.Fatal("expected non-empty snippet")
	}
}

func TestDocumentContentRepository_SearchByTextByCaseFile(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	repo := NewDocumentContentRepository(db)

	createDocumentsSchemaForSearchTests(t, db)

	ctx := context.Background()

	insertSearchFixture(t, ctx, db,
		"doc-1", "case-1", "demanda_1.txt",
		"Demanda por despido con referencia a juicio y conciliación.",
	)

	insertSearchFixture(t, ctx, db,
		"doc-2", "case-2", "demanda_2.txt",
		"Demanda por despido con hechos y fundamentos de derecho.",
	)

	results, err := repo.SearchByTextByCaseFile(ctx, "case-1", "despido", 10)
	if err != nil {
		t.Fatalf("SearchByTextByCaseFile: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].DocumentID != "doc-1" {
		t.Fatalf("expected document_id doc-1, got %s", results[0].DocumentID)
	}
}

func mustOpenTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func createDocumentsSchemaForSearchTests(t *testing.T, db *sql.DB) {
	t.Helper()

	ctx := context.Background()

	stmts := []string{
		`
		CREATE TABLE documents (
			id TEXT PRIMARY KEY,
			case_file_id TEXT NOT NULL,
			original_name TEXT NOT NULL
		)
		`,
		`
		CREATE TABLE document_contents (
			document_id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			extracted_at TEXT NOT NULL
		)
		`,
	}

	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("create schema: %v", err)
		}
	}
}

func insertSearchFixture(t *testing.T, ctx context.Context, db *sql.DB, documentID, caseFileID, originalName, content string) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO documents (id, case_file_id, original_name)
		VALUES (?, ?, ?)
	`, documentID, caseFileID, originalName)
	if err != nil {
		t.Fatalf("insert document: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO document_contents (document_id, content, extracted_at)
		VALUES (?, ?, '2026-04-09T10:00:00Z')
	`, documentID, content)
	if err != nil {
		t.Fatalf("insert document content: %v", err)
	}
}
