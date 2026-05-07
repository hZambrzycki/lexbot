package sqlite

import (
	"context"
	"testing"
)

func TestSQLiteSupportsFTS5(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
		CREATE VIRTUAL TABLE test_fts5 USING fts5(
			content,
			tokenize = 'unicode61 remove_diacritics 2'
		)
	`)
	if err != nil {
		t.Fatalf("expected SQLite to support FTS5: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO test_fts5(content)
		VALUES ('automatización procesal')
	`)
	if err != nil {
		t.Fatalf("insert fts5 content: %v", err)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT content
		FROM test_fts5
		WHERE test_fts5 MATCH 'automatizacion'
	`)
	if err != nil {
		t.Fatalf("query fts5 content: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("expected accent-insensitive FTS5 match")
	}
}
