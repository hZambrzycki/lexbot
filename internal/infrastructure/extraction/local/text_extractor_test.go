package local

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTextExtractor_ExtractText_TXT(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	err := os.WriteFile(path, []byte("Hola mundo\n\nEsto es una prueba."), 0o644)
	if err != nil {
		t.Fatalf("write txt file: %v", err)
	}

	extractor := NewTextExtractor()

	got, err := extractor.ExtractText(context.Background(), path)
	if err != nil {
		t.Fatalf("ExtractText returned error: %v", err)
	}

	want := "Hola mundo\n\nEsto es una prueba."
	if got != want {
		t.Fatalf("unexpected extracted text\nwant: %q\ngot:  %q", want, got)
	}
}

func TestTextExtractor_ExtractText_TXT_EmptyContent(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	err := os.WriteFile(path, []byte("   \n \n\t"), 0o644)
	if err != nil {
		t.Fatalf("write txt file: %v", err)
	}

	extractor := NewTextExtractor()

	_, err = extractor.ExtractText(context.Background(), path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != ErrEmptyContent {
		t.Fatalf("expected ErrEmptyContent, got: %v", err)
	}
}

func TestTextExtractor_ExtractText_DOCX(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.docx")

	err := createMinimalDOCX(path, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:r><w:t>Demanda por despido</w:t></w:r></w:p>
    <w:p><w:r><w:t>Juzgado de lo Social</w:t></w:r></w:p>
  </w:body>
</w:document>`)
	if err != nil {
		t.Fatalf("create docx: %v", err)
	}

	extractor := NewTextExtractor()

	got, err := extractor.ExtractText(context.Background(), path)
	if err != nil {
		t.Fatalf("ExtractText returned error: %v", err)
	}

	if !strings.Contains(got, "Demanda por despido") {
		t.Fatalf("expected extracted text to contain first paragraph, got: %q", got)
	}

	if !strings.Contains(got, "Juzgado de lo Social") {
		t.Fatalf("expected extracted text to contain second paragraph, got: %q", got)
	}
}

func TestTextExtractor_ExtractText_UnsupportedType(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.bin")

	err := os.WriteFile(path, []byte("whatever"), 0o644)
	if err != nil {
		t.Fatalf("write bin file: %v", err)
	}

	extractor := NewTextExtractor()

	_, err = extractor.ExtractText(context.Background(), path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported file type") {
		t.Fatalf("expected unsupported file type error, got: %v", err)
	}
}

func createMinimalDOCX(path string, documentXML string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	zw := zip.NewWriter(file)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`,
		"word/document.xml": documentXML,
	}

	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(content)); err != nil {
			return err
		}
	}

	return zw.Close()
}
