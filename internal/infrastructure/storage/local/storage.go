package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"lexbox/internal/application/ports"
)

type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) StoreDocument(ctx context.Context, documentID string, sourcePath string) (ports.StoredFile, error) {
	_ = ctx

	sourcePath = filepath.Clean(sourcePath)

	info, err := os.Stat(sourcePath)
	if err != nil {
		return ports.StoredFile{}, fmt.Errorf("cannot stat source file: %w", err)
	}

	if info.IsDir() {
		return ports.StoredFile{}, fmt.Errorf("source path is a directory")
	}

	originalName := filepath.Base(sourcePath)
	targetDir := filepath.Join(s.basePath, "documents", documentID)
	targetPath := filepath.Join(targetDir, originalName)

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return ports.StoredFile{}, fmt.Errorf("cannot create target directory: %w", err)
	}

	src, err := os.Open(sourcePath)
	if err != nil {
		return ports.StoredFile{}, fmt.Errorf("cannot open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return ports.StoredFile{}, fmt.Errorf("cannot create target file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return ports.StoredFile{}, fmt.Errorf("cannot copy file: %w", err)
	}

	return ports.StoredFile{
		OriginalName: originalName,
		StoragePath:  targetPath,
	}, nil
}

func (s *Storage) DeleteDocument(ctx context.Context, documentID string) error {
	_ = ctx

	targetDir := filepath.Join(s.basePath, "documents", documentID)
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("cannot delete stored document: %w", err)
	}

	return nil
}

func (s *Storage) DeleteStoredFile(ctx context.Context, storagePath string) error {
	_ = ctx

	cleanPath := filepath.Clean(storagePath)

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("cannot stat stored file: %w", err)
	}

	if info.IsDir() {
		if err := os.RemoveAll(cleanPath); err != nil {
			return fmt.Errorf("cannot delete stored directory: %w", err)
		}
		return nil
	}

	if err := os.Remove(cleanPath); err != nil {
		return fmt.Errorf("cannot delete stored file: %w", err)
	}

	parentDir := filepath.Dir(cleanPath)
	_ = os.Remove(parentDir)

	return nil
}

func (s *Storage) ListStoredDocuments(ctx context.Context) ([]ports.StoredFile, error) {
	_ = ctx

	root := filepath.Join(s.basePath, "documents")

	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return []ports.StoredFile{}, nil
		}
		return nil, fmt.Errorf("cannot stat storage root: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("storage root is not a directory")
	}

	var files []ports.StoredFile

	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			return nil
		}

		files = append(files, ports.StoredFile{
			OriginalName: filepath.Base(path),
			StoragePath:  path,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot walk storage documents: %w", err)
	}

	return files, nil
}
