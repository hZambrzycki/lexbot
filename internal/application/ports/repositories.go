package ports

import (
	"context"
	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/client"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/note"
	"lexbox/internal/domain/shared"
)

type ClientRepository interface {
	Save(ctx context.Context, c client.Client) error
	GetByID(ctx context.Context, id shared.ID) (client.Client, error)
	List(ctx context.Context) ([]client.Client, error)
}

type CaseFileRepository interface {
	Save(ctx context.Context, cf casefile.CaseFile) error
	Update(ctx context.Context, cf casefile.CaseFile) error
	GetByID(ctx context.Context, id shared.ID) (casefile.CaseFile, error)
	List(ctx context.Context) ([]casefile.CaseFile, error)
	ListByClientID(ctx context.Context, clientID shared.ID) ([]casefile.CaseFile, error)
}

type NoteRepository interface {
	Save(ctx context.Context, n note.Note) error
	ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]note.Note, error)
}

type DocumentRepository interface {
	Save(ctx context.Context, doc document.Document) error
	Update(ctx context.Context, doc document.Document) error
	GetByID(ctx context.Context, id shared.ID) (document.Document, error)
	ListAll(ctx context.Context) ([]document.Document, error)
	ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]document.Document, error)
	GetByCaseFileIDAndFileHash(ctx context.Context, caseFileID shared.ID, fileHash string) (document.Document, error)
}

type DocumentMetadataRepository interface {
	Save(ctx context.Context, metadata document.Metadata) error
	GetByDocumentID(ctx context.Context, documentID shared.ID) (document.Metadata, error)
}

type DocumentEventRepository interface {
	ReplaceByDocumentID(ctx context.Context, documentID shared.ID, events []document.Event) error
	ListByDocumentID(ctx context.Context, documentID shared.ID) ([]document.Event, error)
	ListByCaseFileID(ctx context.Context, caseFileID shared.ID) ([]querymodels.CaseFileEventResult, error)
	ListUpcoming(ctx context.Context, fromDate string, caseFileID shared.ID, eventType string) ([]querymodels.CaseFileEventResult, error)
}
