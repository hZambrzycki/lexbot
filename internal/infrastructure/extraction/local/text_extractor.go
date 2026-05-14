package local

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrUnsupportedFileType = errors.New("unsupported file type for extraction")
var ErrEmptyContent = errors.New("extracted content is empty")

type TextExtractor struct{}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

func (e *TextExtractor) ExtractText(ctx context.Context, path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".txt", ".md":
		return e.extractPlainTextFile(path)
	case ".pdf":
		return e.extractPDFFile(ctx, path)
	case ".docx":
		return e.extractDOCXFile(path)
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedFileType, ext)
	}
}

func (e *TextExtractor) extractPlainTextFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read file: %w", err)
	}

	content := strings.TrimSpace(string(data))
	if content == "" {
		return "", ErrEmptyContent
	}

	return content, nil
}

func (e *TextExtractor) extractPDFFile(ctx context.Context, path string) (string, error) {
	cmd := exec.CommandContext(ctx, "pdftotext", "-layout", path, "-")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())

		fmt.Println("OCR fallback triggered:", path)

		ocrText, ocrErr := e.extractPDFWithOCR(ctx, path)
		if ocrErr == nil {
			fmt.Println("OCR extraction succeeded:", path)
			return ocrText, nil
		}

		fmt.Println("OCR extraction failed:", ocrErr)

		if errMsg != "" {
			return "", fmt.Errorf(
				"pdf extraction failed: %s: %w",
				errMsg,
				err,
			)
		}

		return "", fmt.Errorf("pdf extraction failed: %w", err)
	}

	content := normalizeWhitespace(stdout.String())

	// PDF con texto normal → no usar OCR
	if content != "" && len(content) >= 50 {
		return content, nil
	}

	fmt.Println("Low text PDF detected, trying OCR:", path)

	ocrText, ocrErr := e.extractPDFWithOCR(ctx, path)
	if ocrErr == nil {
		fmt.Println("OCR extraction succeeded:", path)

		// Preferimos OCR si recupera más contenido
		if len(ocrText) > len(content) {
			return ocrText, nil
		}
	}

	if ocrErr != nil {
		fmt.Println("OCR extraction failed:", ocrErr)
	}

	if content == "" {
		return "", ErrEmptyContent
	}

	return content, nil
}
func (e *TextExtractor) extractDOCXFile(path string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("cannot open docx file: %w", err)
	}
	defer r.Close()

	var documentXML *zip.File
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			documentXML = f
			break
		}
	}

	if documentXML == nil {
		return "", fmt.Errorf("docx extraction failed: word/document.xml not found")
	}

	rc, err := documentXML.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open docx document.xml: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("cannot read docx document.xml: %w", err)
	}

	content, err := extractTextFromWordXML(data)
	if err != nil {
		return "", fmt.Errorf("docx extraction failed: %w", err)
	}

	content = normalizeWhitespace(content)
	if content == "" {
		return "", ErrEmptyContent
	}

	return content, nil
}

func extractTextFromWordXML(data []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	var b strings.Builder
	pendingNewline := false
	pendingTab := false

	for {
		tok, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "tab":
				pendingTab = true
			case "br":
				b.WriteString("\n")
			}

		case xml.EndElement:
			switch tok.Name.Local {
			case "p":
				pendingNewline = true
			case "tr":
				pendingNewline = true
			}

		case xml.CharData:
			text := string(tok)
			if strings.TrimSpace(text) == "" {
				continue
			}

			if pendingNewline && b.Len() > 0 {
				b.WriteString("\n")
				pendingNewline = false
			}
			if pendingTab {
				b.WriteString("\t")
				pendingTab = false
			}

			b.WriteString(text)
		}
	}

	return b.String(), nil
}

func normalizeWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cleaned = append(cleaned, line)
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}
