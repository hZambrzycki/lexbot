package documentapp

import (
	"context"

	"lexbox/internal/application/ports"
)

type StorageAuditMissingDocument struct {
	DocumentID   string
	OriginalName string
	StoragePath  string
}

type StorageAuditResult struct {
	OrphanFiles      []string
	MissingDocuments []StorageAuditMissingDocument
}

type StorageAudit struct {
	Documents ports.DocumentRepository
	Storage   ports.FileStorage
}

func (uc StorageAudit) Execute(ctx context.Context) (StorageAuditResult, error) {
	documents, err := uc.Documents.ListAll(ctx)
	if err != nil {
		return StorageAuditResult{}, err
	}

	storedFiles, err := uc.Storage.ListStoredDocuments(ctx)
	if err != nil {
		return StorageAuditResult{}, err
	}

	dbPaths := make(map[string]StorageAuditMissingDocument, len(documents))
	for _, doc := range documents {
		dbPaths[doc.StoragePath] = StorageAuditMissingDocument{
			DocumentID:   doc.ID.String(),
			OriginalName: doc.OriginalName,
			StoragePath:  doc.StoragePath,
		}
	}

	storagePaths := make(map[string]struct{}, len(storedFiles))
	for _, file := range storedFiles {
		storagePaths[file.StoragePath] = struct{}{}
	}

	result := StorageAuditResult{
		OrphanFiles:      make([]string, 0),
		MissingDocuments: make([]StorageAuditMissingDocument, 0),
	}

	for _, file := range storedFiles {
		if _, ok := dbPaths[file.StoragePath]; !ok {
			result.OrphanFiles = append(result.OrphanFiles, file.StoragePath)
		}
	}

	for _, doc := range documents {
		if _, ok := storagePaths[doc.StoragePath]; !ok {
			result.MissingDocuments = append(result.MissingDocuments, StorageAuditMissingDocument{
				DocumentID:   doc.ID.String(),
				OriginalName: doc.OriginalName,
				StoragePath:  doc.StoragePath,
			})
		}
	}

	return result, nil
}
