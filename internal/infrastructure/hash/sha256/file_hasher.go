package sha256hasher

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type FileHasher struct{}

func NewFileHasher() *FileHasher {
	return &FileHasher{}
}

func (h *FileHasher) HashFile(ctx context.Context, path string) (string, error) {
	_ = ctx

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file for hashing: %w", err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", fmt.Errorf("cannot hash file: %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
