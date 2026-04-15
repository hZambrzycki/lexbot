package documentapp

import (
	"context"
	"testing"

	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/shared"
)

type fakeDocumentContentRepository struct {
	searchByTextResults           []querymodels.SearchDocumentResult
	searchByTextByCaseFileResults []querymodels.SearchDocumentResult

	searchByTextQuery           string
	searchByTextLimit           int
	searchByTextByCaseFileID    string
	searchByTextByCaseFileQuery string
	searchByTextByCaseFileLimit int
}

func (f *fakeDocumentContentRepository) Save(ctx context.Context, documentID string, content string) error {
	return nil
}

func (f *fakeDocumentContentRepository) GetByDocumentID(ctx context.Context, documentID string) (string, error) {
	return "", shared.ErrNotFound
}

func (f *fakeDocumentContentRepository) SearchByText(ctx context.Context, query string, limit int) ([]querymodels.SearchDocumentResult, error) {
	f.searchByTextQuery = query
	f.searchByTextLimit = limit
	return f.searchByTextResults, nil
}

func (f *fakeDocumentContentRepository) SearchByTextByCaseFile(ctx context.Context, caseFileID string, query string, limit int) ([]querymodels.SearchDocumentResult, error) {
	f.searchByTextByCaseFileID = caseFileID
	f.searchByTextByCaseFileQuery = query
	f.searchByTextByCaseFileLimit = limit
	return f.searchByTextByCaseFileResults, nil
}

func TestSearchDocuments_Execute_GlobalSearch(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentContentRepository{
		searchByTextResults: []querymodels.SearchDocumentResult{
			{DocumentID: "doc-1"},
		},
	}

	uc := SearchDocuments{
		DocumentContents: repo,
	}

	results, err := uc.Execute(context.Background(), SearchDocumentsInput{
		Query: "despido",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if repo.searchByTextQuery != "despido" {
		t.Fatalf("unexpected search query: %q", repo.searchByTextQuery)
	}

	if repo.searchByTextLimit != 10 {
		t.Fatalf("unexpected limit: %d", repo.searchByTextLimit)
	}
}

func TestSearchDocuments_Execute_SearchByCaseFile(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentContentRepository{
		searchByTextByCaseFileResults: []querymodels.SearchDocumentResult{
			{DocumentID: "doc-2"},
		},
	}

	uc := SearchDocuments{
		DocumentContents: repo,
	}

	caseFileID := "ddfcc249-3c3c-41ab-b823-b0583e8709e0"

	results, err := uc.Execute(context.Background(), SearchDocumentsInput{
		Query:      "juicio",
		CaseFileID: caseFileID,
		Limit:      5,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if repo.searchByTextByCaseFileID != caseFileID {
		t.Fatalf("unexpected case file id: %q", repo.searchByTextByCaseFileID)
	}

	if repo.searchByTextByCaseFileQuery != "juicio" {
		t.Fatalf("unexpected search query: %q", repo.searchByTextByCaseFileQuery)
	}

	if repo.searchByTextByCaseFileLimit != 5 {
		t.Fatalf("unexpected limit: %d", repo.searchByTextByCaseFileLimit)
	}
}

func TestSearchDocuments_Execute_EmptyQueryReturnsEmptyResults(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentContentRepository{}
	uc := SearchDocuments{
		DocumentContents: repo,
	}

	results, err := uc.Execute(context.Background(), SearchDocumentsInput{
		Query: "   ",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
