package main

import (
	"context"
	"log"
	"os"

	casefileapp "lexbox/internal/application/casefile"
	clientapp "lexbox/internal/application/client"
	documentapp "lexbox/internal/application/document"
	noteapp "lexbox/internal/application/note"
	localextractor "lexbox/internal/infrastructure/extraction/local"
	sha256hasher "lexbox/internal/infrastructure/hash/sha256"
	idgen "lexbox/internal/infrastructure/idgen/uuid"
	reposqlite "lexbox/internal/infrastructure/persistence/sqlite"
	localstorage "lexbox/internal/infrastructure/storage/local"
	"lexbox/internal/interfaces/cli"
)

func main() {
	ctx := context.Background()

	db, err := reposqlite.NewDB("lexbox.db")
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer db.Close()

	if err := reposqlite.RunMigrations(db, "migrations/001_init.sql"); err != nil {
		log.Fatalf("migration error 001: %v", err)
	}
	if err := reposqlite.RunMigrations(db, "migrations/002_document_contents.sql"); err != nil {
		log.Fatalf("migration error 002: %v", err)
	}
	if err := reposqlite.RunMigrations(db, "migrations/003_document_metadata.sql"); err != nil {
		log.Fatalf("migration error 003: %v", err)
	}
	if err := reposqlite.RunMigrations(db, "migrations/004_document_events.sql"); err != nil {
		log.Fatalf("migration error 004: %v", err)
	}
	if err := reposqlite.RunMigrations(db, "migrations/005_document_event_context.sql"); err != nil {
		log.Fatalf("migration error 005: %v", err)
	}
	if err := reposqlite.RunMigrations(db, "migrations/006_document_event_semantics.sql"); err != nil {
		log.Fatalf("migration error 006: %v", err)
	}
	if err := reposqlite.RunMigrations(db, "migrations/007_case_file_event_config.sql"); err != nil {
		log.Fatalf("migration error 007: %v", err)
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

	listUpcomingEvents := documentapp.ListUpcomingEvents{
		Events: documentEventRepo,
	}

	exportUpcomingEventsICS := documentapp.ExportUpcomingEventsICS{
		ListUpcomingEvents: listUpcomingEvents,
	}
	app := cli.App{
		Out: os.Stdout,

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

		AddNote: noteapp.AddNote{
			Notes:     noteRepo,
			CaseFiles: caseFileRepo,
			IDs:       idGenerator,
		},

		ListNotesByCaseFile: noteapp.ListNotesByCaseFile{
			Notes: noteRepo,
		},

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

		ExtractDocumentText: extractDocumentText,

		ReindexDocument: reindexDocument,

		ReindexAllDocuments: reindexAllDocuments,

		AnalyzeDocumentMetadata: analyzeDocumentMetadata,

		AnalyzeAllDocumentMetadata: analyzeAllDocumentMetadata,

		AnalyzeDocumentEvents: analyzeDocumentEvents,

		AnalyzeAllDocumentEvents: analyzeAllDocumentEvents,

		GetDocumentDetail: documentapp.GetDocumentDetail{
			Documents:        documentRepo,
			DocumentContents: documentContentRepo,
			Metadata:         documentMetadataRepo,
			Events:           documentEventRepo,
		},

		SearchDocuments: documentapp.SearchDocuments{
			DocumentContents: documentContentRepo,
		},

		GetCaseFileDetail: casefileapp.GetCaseFileDetail{
			CaseFiles: caseFileRepo,
			Notes:     noteRepo,
			Documents: documentRepo,
		},

		GetCaseFileDashboard: casefileapp.GetCaseFileDashboard{
			GetCaseFileDetail: casefileapp.GetCaseFileDetail{
				CaseFiles: caseFileRepo,
				Notes:     noteRepo,
				Documents: documentRepo,
			},
			ListUpcomingEvents: listUpcomingEvents,
			DocumentContents:   documentContentRepo,
			Metadata:           documentMetadataRepo,
			Events:             documentEventRepo,
		},

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

		ListDocumentEvents: documentapp.ListDocumentEvents{
			Events: documentEventRepo,
		},

		ListEventsByCaseFile: documentapp.ListEventsByCaseFile{
			Events: documentEventRepo,
		},

		ListUpcomingEvents: listUpcomingEvents,

		ExportUpcomingEventsICS: exportUpcomingEventsICS,

		FixCaseFile: casefileapp.FixCaseFile{
			ReindexAllDocuments: reindexAllDocuments,
			AnalyzeAllMetadata:  analyzeAllDocumentMetadata,
			AnalyzeAllEvents:    analyzeAllDocumentEvents,
			GetDashboard: casefileapp.GetCaseFileDashboard{
				GetCaseFileDetail: casefileapp.GetCaseFileDetail{
					CaseFiles: caseFileRepo,
					Notes:     noteRepo,
					Documents: documentRepo,
				},
				ListUpcomingEvents: listUpcomingEvents,
				DocumentContents:   documentContentRepo,
				Metadata:           documentMetadataRepo,
				Events:             documentEventRepo,
			},
		},

		AuditCaseFile: casefileapp.AuditCaseFile{
			VerifyAllDocuments: documentapp.VerifyAllDocuments{
				Documents:        documentRepo,
				DocumentContents: documentContentRepo,
			},
			GetDashboard: casefileapp.GetCaseFileDashboard{
				GetCaseFileDetail: casefileapp.GetCaseFileDetail{
					CaseFiles: caseFileRepo,
					Notes:     noteRepo,
					Documents: documentRepo,
				},
				ListUpcomingEvents: listUpcomingEvents,
				DocumentContents:   documentContentRepo,
				Metadata:           documentMetadataRepo,
				Events:             documentEventRepo,
			},
		},
	}

	if err := app.Run(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
