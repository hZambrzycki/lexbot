package search

import (
	"strings"
	"testing"
)

func TestNormalizeText(t *testing.T) {
	t.Parallel()

	got := NormalizeText("Árbol ÉXTRANJERÍA Ñandú")

	want := "arbol extranjeria nandu"

	if got != want {
		t.Fatalf("unexpected normalized text\nwant: %q\ngot:  %q", want, got)
	}
}

func TestSplitTerms_DeduplicatesAccentInsensitiveTerms(t *testing.T) {
	t.Parallel()

	terms := SplitTerms("ejecución ejecucion EJECUCIÓN")

	if len(terms) != 1 {
		t.Fatalf("expected 1 deduplicated term, got %d", len(terms))
	}
}

func TestComputeScore_IsAccentInsensitive(t *testing.T) {
	t.Parallel()

	content := "Automatización de ejecución procesal"

	score := ComputeScore(content, []string{
		"automatizacion",
		"ejecucion",
	})

	if score <= 0 {
		t.Fatalf("expected positive score, got %d", score)
	}
}

func TestBuildSnippet_HighlightsMultipleTerms(t *testing.T) {
	t.Parallel()

	content := "La empresa interpuso recurso contra la sentencia dictada."

	snippet := BuildSnippet(
		content,
		"empresa sentencia",
		[]string{"empresa", "sentencia"},
		180,
	)

	assertContains(t, snippet, "[empresa]")
	assertContains(t, snippet, "[sentencia]")
}

func TestBuildSnippet_IsAccentInsensitive(t *testing.T) {
	t.Parallel()

	content := "Sistema de automatización para ejecución jurídica"

	snippet := BuildSnippet(
		content,
		"automatizacion ejecucion",
		[]string{"automatizacion", "ejecucion"},
		180,
	)

	assertContains(t, snippet, "[automatización]")
	assertContains(t, snippet, "[ejecución]")
}

func TestBuildSnippet_IsCaseInsensitive(t *testing.T) {
	t.Parallel()

	content := "Demanda por Despido improcedente"

	snippet := BuildSnippet(
		content,
		"DEMANDA DESPIDO",
		[]string{"DEMANDA", "DESPIDO"},
		180,
	)

	lower := strings.ToLower(snippet)

	assertContains(t, lower, "[demanda]")
	assertContains(t, lower, "[despido]")
}

func TestBuildSnippet_ReturnsTruncatedSnippet(t *testing.T) {
	t.Parallel()

	content := strings.Repeat("texto ", 100)

	snippet := BuildSnippet(
		content,
		"texto",
		[]string{"texto"},
		50,
	)

	if len(snippet) == 0 {
		t.Fatal("expected non-empty snippet")
	}
}

func TestNormalizeTextWithIndex(t *testing.T) {
	t.Parallel()

	normalized, indexMap := NormalizeTextWithIndex("acción")

	if normalized != "accion" {
		t.Fatalf("unexpected normalized value: %q", normalized)
	}

	if len(indexMap) != len(normalized) {
		t.Fatalf(
			"expected index map length %d, got %d",
			len(normalized),
			len(indexMap),
		)
	}
}

func assertContains(t *testing.T, value string, expected string) {
	t.Helper()

	if !strings.Contains(value, expected) {
		t.Fatalf("expected %q to contain %q", value, expected)
	}
}
