package ports

import "context"

type StoredFile struct {
	OriginalName string
	StoragePath  string
}

type FileStorage interface {
	StoreDocument(ctx context.Context, documentID string, sourcePath string) (StoredFile, error)
	DeleteDocument(ctx context.Context, documentID string) error
	DeleteStoredFile(ctx context.Context, storagePath string) error
	ListStoredDocuments(ctx context.Context) ([]StoredFile, error)
}
