package bootstrap

import (
	"context"

	casefileapp "lexbox/internal/application/casefile"
	clientapp "lexbox/internal/application/client"
	documentapp "lexbox/internal/application/document"
	noteapp "lexbox/internal/application/note"
	localextractor "lexbox/internal/infrastructure/extraction/local"
	sha256hasher "lexbox/internal/infrastructure/hash/sha256"
	idgen "lexbox/internal/infrastructure/idgen/uuid"
	reposqlite "lexbox/internal/infrastructure/persistence/sqlite"
	localstorage "lexbox/internal/infrastructure/storage/local"
)

type App struct {
	DB interface {
		Close() error
	}

	CreateClient               clientapp.CreateClient
	GetClient                  clientapp.GetClient
	ListClients                clientapp.ListClients
	CreateCaseFile             casefileapp.CreateCaseFile
	UpdateCaseFileConfig       casefileapp.UpdateCaseFileConfig
	ListCaseFiles              casefileapp.ListCaseFiles
	ListCaseFilesByClient      casefileapp.ListCaseFilesByClient
	GetCaseFileDashboard       casefileapp.GetCaseFileDashboard
	AddNote                    noteapp.AddNote
	ListNotesByCaseFile        noteapp.ListNotesByCaseFile
	AttachDocument             documentapp.AttachDocument
	ImportDocument             documentapp.ImportDocument
	ListDocumentsByCaseFile    documentapp.ListDocumentsByCaseFile
	DeleteDocument             documentapp.DeleteDocument
	UpdateDocumentReviewState  documentapp.UpdateDocumentReviewState
	ExtractDocumentText        documentapp.ExtractDocumentText
	SearchDocuments            documentapp.SearchDocuments
	GetDocumentDetail          documentapp.GetDocumentDetail
	GetCaseFileDetail          casefileapp.GetCaseFileDetail
	StorageAudit               documentapp.StorageAudit
	StorageCleanOrphans        documentapp.StorageCleanOrphans
	MimeNormalize              documentapp.MimeNormalize
	VerifyDocument             documentapp.VerifyDocument
	VerifyAllDocuments         documentapp.VerifyAllDocuments
	ReindexDocument            documentapp.ReindexDocument
	ReindexAllDocuments        documentapp.ReindexAllDocuments
	AnalyzeDocumentMetadata    documentapp.AnalyzeDocumentMetadata
	AnalyzeAllDocumentMetadata documentapp.AnalyzeAllDocumentMetadata
	AnalyzeDocumentEvents      documentapp.AnalyzeDocumentEvents
	AnalyzeAllDocumentEvents   documentapp.AnalyzeAllDocumentEvents
	ListDocumentEvents         documentapp.ListDocumentEvents
	ListEventsByCaseFile       documentapp.ListEventsByCaseFile
	ListUpcomingEvents         documentapp.ListUpcomingEvents
	ExportUpcomingEventsICS    documentapp.ExportUpcomingEventsICS
	MarkEventReviewed          documentapp.MarkEventReviewed
	MarkEventResolved          documentapp.MarkEventResolved
	ReopenEvent                documentapp.ReopenEvent
	FixCaseFile                casefileapp.FixCaseFile
	AuditCaseFile              casefileapp.AuditCaseFile
	GetEvent                   documentapp.GetEvent
	ReprocessDocument          documentapp.ReprocessDocument
}

func BuildApp(ctx context.Context) (*App, error) {
	_ = ctx

	db, err := reposqlite.NewDB("lexbox.db")
	if err != nil {
		return nil, err
	}

	if err := reposqlite.RunMigrations(db, "migrations/001_init.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/002_document_contents.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/003_document_metadata.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/004_document_events.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/005_document_event_context.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/006_document_event_semantics.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/007_case_file_event_config.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/008_document_event_review_state.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := reposqlite.RunMigrations(db, "migrations/009_document_review_state.sql"); err != nil {
		_ = db.Close()
		return nil, err
	}

	clientRepo := reposqlite.NewClientRepository(db)
	caseFileRepo := reposqlite.NewCaseFileRepository(db)
	noteRepo := reposqlite.NewNoteRepository(db)
	documentRepo := reposqlite.NewDocumentRepository(db)
	documentContentRepo := reposqlite.NewDocumentContentRepository(db)
	documentMetadataRepo := reposqlite.NewDocumentMetadataRepository(db)
	documentEventRepo := reposqlite.NewDocumentEventRepository(db)

	idGenerator := idgen.NewIDGenerator()
	storage := localstorage.NewStorage(".lexbox/storage")
	textExtractor := localextractor.NewTextExtractor()
	fileHasher := sha256hasher.NewFileHasher()

	extractDocumentText := documentapp.ExtractDocumentText{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		Extractor:        textExtractor,
	}

	reindexDocument := documentapp.ReindexDocument{
		ExtractDocumentText: extractDocumentText,
	}

	reindexAllDocuments := documentapp.ReindexAllDocuments{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		ReindexDocument:  reindexDocument,
	}

	analyzeDocumentMetadata := documentapp.AnalyzeDocumentMetadata{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		Metadata:         documentMetadataRepo,
	}

	analyzeAllDocumentMetadata := documentapp.AnalyzeAllDocumentMetadata{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		AnalyzeOne:       analyzeDocumentMetadata,
	}

	analyzeDocumentEvents := documentapp.AnalyzeDocumentEvents{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		CaseFiles:        caseFileRepo,
		Events:           documentEventRepo,
		IDs:              idGenerator,
	}

	analyzeAllDocumentEvents := documentapp.AnalyzeAllDocumentEvents{
		Documents:        documentRepo,
		DocumentContents: documentContentRepo,
		CaseFiles:        caseFileRepo,
		AnalyzeOne:       analyzeDocumentEvents,
	}

	reprocessDocument := documentapp.ReprocessDocument{
		ReindexDocument:         reindexDocument,
		AnalyzeDocumentMetadata: analyzeDocumentMetadata,
		AnalyzeDocumentEvents:   analyzeDocumentEvents,
	}

	listUpcomingEvents := documentapp.ListUpcomingEvents{
		Events: documentEventRepo,
	}

	exportUpcomingEventsICS := documentapp.ExportUpcomingEventsICS{
		ListUpcomingEvents: listUpcomingEvents,
	}

	getCaseFileDetail := casefileapp.GetCaseFileDetail{
		CaseFiles: caseFileRepo,
		Notes:     noteRepo,
		Documents: documentRepo,
	}

	getCaseFileDashboard := casefileapp.GetCaseFileDashboard{
		GetCaseFileDetail:  getCaseFileDetail,
		ListUpcomingEvents: listUpcomingEvents,
		DocumentContents:   documentContentRepo,
		Metadata:           documentMetadataRepo,
		Events:             documentEventRepo,
	}

	app := &App{
		DB: db,

		CreateClient: clientapp.CreateClient{
			Clients: clientRepo,
			IDs:     idGenerator,
		},

		GetClient: clientapp.GetClient{
			Clients: clientRepo,
		},

		ListClients: clientapp.ListClients{
			Clients: clientRepo,
		},

		CreateCaseFile: casefileapp.CreateCaseFile{
			CaseFiles: caseFileRepo,
			Clients:   clientRepo,
			IDs:       idGenerator,
		},

		UpdateCaseFileConfig: casefileapp.UpdateCaseFileConfig{
			CaseFiles: caseFileRepo,
		},

		ListCaseFiles: casefileapp.ListCaseFiles{
			CaseFiles: caseFileRepo,
		},

		ListCaseFilesByClient: casefileapp.ListCaseFilesByClient{
			CaseFiles: caseFileRepo,
		},

		GetCaseFileDashboard: getCaseFileDashboard,

		AddNote: noteapp.AddNote{
			Notes:     noteRepo,
			CaseFiles: caseFileRepo,
			IDs:       idGenerator,
		},

		ListNotesByCaseFile: noteapp.ListNotesByCaseFile{
			Notes: noteRepo,
		},
		ReprocessDocument: reprocessDocument,
		AttachDocument: documentapp.AttachDocument{
			Documents: documentRepo,
			CaseFiles: caseFileRepo,
			IDs:       idGenerator,
		},

		ImportDocument: documentapp.ImportDocument{
			Documents:        documentRepo,
			DocumentContents: documentContentRepo,
			Metadata:         documentMetadataRepo,
			CaseFiles:        caseFileRepo,
			Storage:          storage,
			Extractor:        textExtractor,
			Hasher:           fileHasher,
			IDs:              idGenerator,
			AnalyzeEvents:    analyzeDocumentEvents,
		},

		ListDocumentsByCaseFile: documentapp.ListDocumentsByCaseFile{
			Documents: documentRepo,
			Metadata:  documentMetadataRepo,
		},
		DeleteDocument: documentapp.DeleteDocument{
			Documents: documentRepo,
			Storage:   storage,
		},
		UpdateDocumentReviewState: documentapp.UpdateDocumentReviewState{
			Documents: documentRepo,
		},
		ExtractDocumentText: extractDocumentText,

		SearchDocuments: documentapp.SearchDocuments{
			DocumentContents: documentContentRepo,
		},

		GetDocumentDetail: documentapp.GetDocumentDetail{
			Documents:        documentRepo,
			DocumentContents: documentContentRepo,
			Metadata:         documentMetadataRepo,
			Events:           documentEventRepo,
		},

		GetCaseFileDetail: getCaseFileDetail,

		StorageAudit: documentapp.StorageAudit{
			Documents: documentRepo,
			Storage:   storage,
		},

		StorageCleanOrphans: documentapp.StorageCleanOrphans{
			Documents: documentRepo,
			Storage:   storage,
		},

		MimeNormalize: documentapp.MimeNormalize{
			Documents: documentRepo,
		},

		VerifyDocument: documentapp.VerifyDocument{
			Documents:        documentRepo,
			DocumentContents: documentContentRepo,
		},

		VerifyAllDocuments: documentapp.VerifyAllDocuments{
			Documents:        documentRepo,
			DocumentContents: documentContentRepo,
		},

		ReindexDocument: reindexDocument,

		ReindexAllDocuments: reindexAllDocuments,

		AnalyzeDocumentMetadata: analyzeDocumentMetadata,

		AnalyzeAllDocumentMetadata: analyzeAllDocumentMetadata,

		AnalyzeDocumentEvents: analyzeDocumentEvents,

		AnalyzeAllDocumentEvents: analyzeAllDocumentEvents,

		ListDocumentEvents: documentapp.ListDocumentEvents{
			Events: documentEventRepo,
		},

		ListEventsByCaseFile: documentapp.ListEventsByCaseFile{
			Events: documentEventRepo,
		},

		ListUpcomingEvents: listUpcomingEvents,

		ExportUpcomingEventsICS: exportUpcomingEventsICS,

		MarkEventReviewed: documentapp.MarkEventReviewed{
			Events: documentEventRepo,
		},

		MarkEventResolved: documentapp.MarkEventResolved{
			Events: documentEventRepo,
		},

		ReopenEvent: documentapp.ReopenEvent{
			Events: documentEventRepo,
		},

		FixCaseFile: casefileapp.FixCaseFile{
			ReindexAllDocuments: reindexAllDocuments,
			AnalyzeAllMetadata:  analyzeAllDocumentMetadata,
			AnalyzeAllEvents:    analyzeAllDocumentEvents,
			GetDashboard:        getCaseFileDashboard,
		},

		AuditCaseFile: casefileapp.AuditCaseFile{
			VerifyAllDocuments: documentapp.VerifyAllDocuments{
				Documents:        documentRepo,
				DocumentContents: documentContentRepo,
			},
			GetDashboard: getCaseFileDashboard,
		},

		GetEvent: documentapp.GetEvent{
			Events: documentEventRepo,
		},
	}

	return app, nil
}

func (a *App) Close() error {
	if a == nil || a.DB == nil {
		return nil
	}
	return a.DB.Close()
}
