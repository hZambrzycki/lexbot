package documentapp

import (
	"context"

	"lexbox/internal/application/ports"
)

type StorageCleanOrphansResult struct {
	DeletedFiles []string
}

type StorageCleanOrphans struct {
	Documents ports.DocumentRepository
	Storage   ports.FileStorage
}

func (uc StorageCleanOrphans) Execute(ctx context.Context) (StorageCleanOrphansResult, error) {
	documents, err := uc.Documents.ListAll(ctx)
	if err != nil {
		return StorageCleanOrphansResult{}, err
	}

	storedFiles, err := uc.Storage.ListStoredDocuments(ctx)
	if err != nil {
		return StorageCleanOrphansResult{}, err
	}

	dbPaths := make(map[string]struct{}, len(documents))
	for _, doc := range documents {
		dbPaths[doc.StoragePath] = struct{}{}
	}

	result := StorageCleanOrphansResult{
		DeletedFiles: make([]string, 0),
	}

	for _, file := range storedFiles {
		if _, ok := dbPaths[file.StoragePath]; ok {
			continue
		}

		if err := uc.Storage.DeleteStoredFile(ctx, file.StoragePath); err != nil {
			return StorageCleanOrphansResult{}, err
		}

		result.DeletedFiles = append(result.DeletedFiles, file.StoragePath)
	}

	return result, nil
}
