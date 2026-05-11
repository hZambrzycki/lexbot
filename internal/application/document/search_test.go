package documentapp

import (
	"context"
	"testing"

	"lexbox/internal/application/querymodels"
)

type fakeDocumentSearchIndexRepository struct {
	results []querymodels.SearchDocumentResult

	searchQuery      string
	searchCaseFileID string
	searchFilters    querymodels.SearchDocumentFilters
	searchLimit      int
}

func (f *fakeDocumentSearchIndexRepository) UpsertDocument(
	ctx context.Context,
	documentID string,
	caseFileID string,
	originalName string,
	content string,
	documentType string,
	legalArea string,
) error {
	return nil
}

func (f *fakeDocumentSearchIndexRepository) DeleteDocument(
	ctx context.Context,
	documentID string,
) error {
	return nil
}

func (f *fakeDocumentSearchIndexRepository) Search(
	ctx context.Context,
	query string,
	caseFileID string,
	filters querymodels.SearchDocumentFilters,
	limit int,
) ([]querymodels.SearchDocumentResult, error) {
	f.searchQuery = query
	f.searchCaseFileID = caseFileID
	f.searchFilters = filters
	f.searchLimit = limit

	return f.results, nil
}

func TestSearchDocuments_Execute_GlobalSearch(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentSearchIndexRepository{
		results: []querymodels.SearchDocumentResult{
			{DocumentID: "doc-1"},
		},
	}

	uc := SearchDocuments{
		SearchIndex: repo,
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

	if repo.searchQuery != "despido" {
		t.Fatalf("unexpected search query: %q", repo.searchQuery)
	}

	if repo.searchCaseFileID != "" {
		t.Fatalf("expected empty case file id, got %q", repo.searchCaseFileID)
	}

	if repo.searchLimit != 10 {
		t.Fatalf("unexpected limit: %d", repo.searchLimit)
	}
}

func TestSearchDocuments_Execute_SearchByCaseFile(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentSearchIndexRepository{
		results: []querymodels.SearchDocumentResult{
			{DocumentID: "doc-2"},
		},
	}

	uc := SearchDocuments{
		SearchIndex: repo,
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

	if repo.searchCaseFileID != caseFileID {
		t.Fatalf("unexpected case file id: %q", repo.searchCaseFileID)
	}

	if repo.searchQuery != "juicio" {
		t.Fatalf("unexpected search query: %q", repo.searchQuery)
	}

	if repo.searchLimit != 5 {
		t.Fatalf("unexpected limit: %d", repo.searchLimit)
	}
}

func TestSearchDocuments_Execute_EmptyQueryReturnsEmptyResults(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentSearchIndexRepository{}

	uc := SearchDocuments{
		SearchIndex: repo,
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

	if repo.searchQuery != "" {
		t.Fatalf("expected search repository not to be called")
	}
}

func TestSearchDocuments_Execute_DefaultLimit(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentSearchIndexRepository{}

	uc := SearchDocuments{
		SearchIndex: repo,
	}

	_, err := uc.Execute(context.Background(), SearchDocumentsInput{
		Query: "demanda",
		Limit: 0,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if repo.searchLimit != 20 {
		t.Fatalf("expected default limit 20, got %d", repo.searchLimit)
	}
}

func TestSearchDocuments_Execute_ParsesDSLFilters(t *testing.T) {
	t.Parallel()

	repo := &fakeDocumentSearchIndexRepository{
		results: []querymodels.SearchDocumentResult{
			{DocumentID: "doc-3"},
		},
	}

	uc := SearchDocuments{
		SearchIndex: repo,
	}

	results, err := uc.Execute(context.Background(), SearchDocumentsInput{
		Query: "type:demanda area:laboral despido",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if repo.searchQuery != "despido" {
		t.Fatalf(
			"expected free query %q, got %q",
			"despido",
			repo.searchQuery,
		)
	}

	if repo.searchFilters.DocumentType != "demanda" {
		t.Fatalf(
			"expected document type demanda, got %q",
			repo.searchFilters.DocumentType,
		)
	}

	if repo.searchFilters.LegalArea != "laboral" {
		t.Fatalf(
			"expected legal area laboral, got %q",
			repo.searchFilters.LegalArea,
		)
	}

	if repo.searchLimit != 10 {
		t.Fatalf("unexpected limit: %d", repo.searchLimit)
	}
}
