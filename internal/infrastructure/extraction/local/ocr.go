package local

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func (e *TextExtractor) extractPDFWithOCR(ctx context.Context, path string) (string, error) {
	tempDir, err := os.MkdirTemp("", "lexbox-ocr-*")
	if err != nil {
		return "", fmt.Errorf("cannot create ocr temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	outputPrefix := filepath.Join(tempDir, "page")

	convertCmd := exec.CommandContext(
		ctx,
		"pdftoppm",
		"-png",
		"-r",
		"200",
		path,
		outputPrefix,
	)

	var convertErr bytes.Buffer
	convertCmd.Stderr = &convertErr

	if err := convertCmd.Run(); err != nil {
		errMsg := strings.TrimSpace(convertErr.String())
		if errMsg != "" {
			return "", fmt.Errorf("ocr pdf image conversion failed: %s: %w", errMsg, err)
		}
		return "", fmt.Errorf("ocr pdf image conversion failed: %w", err)
	}

	imageFiles, err := filepath.Glob(filepath.Join(tempDir, "*.png"))
	if err != nil {
		return "", fmt.Errorf("cannot list ocr images: %w", err)
	}

	sort.Strings(imageFiles)

	if len(imageFiles) == 0 {
		return "", ErrEmptyContent
	}

	var result strings.Builder

	for _, imagePath := range imageFiles {
		ocrCmd := exec.CommandContext(
			ctx,
			"tesseract",
			imagePath,
			"stdout",
			"-l",
			"spa",
		)

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		ocrCmd.Stdout = &stdout
		ocrCmd.Stderr = &stderr

		if err := ocrCmd.Run(); err != nil {
			continue
		}

		text := strings.TrimSpace(stdout.String())
		if text == "" {
			continue
		}

		result.WriteString(text)
		result.WriteString("\n\n")
	}

	content := normalizeWhitespace(result.String())
	if content == "" {
		return "", ErrEmptyContent
	}

	return content, nil
}
