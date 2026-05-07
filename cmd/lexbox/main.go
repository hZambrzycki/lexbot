package main

import (
	"context"
	"log"
	"os"

	"lexbox/internal/infrastructure/bootstrap"
	"lexbox/internal/interfaces/cli"
)

func main() {
	ctx := context.Background()

	appServices, err := bootstrap.BuildApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = appServices.Close()
	}()

	app := cli.App{
		Out: os.Stdout,

		CreateClient:                appServices.CreateClient,
		GetClient:                   appServices.GetClient,
		ListClients:                 appServices.ListClients,
		CreateCaseFile:              appServices.CreateCaseFile,
		UpdateCaseFileConfig:        appServices.UpdateCaseFileConfig,
		ListCaseFiles:               appServices.ListCaseFiles,
		ListCaseFilesByClient:       appServices.ListCaseFilesByClient,
		GetCaseFileDashboard:        appServices.GetCaseFileDashboard,
		AddNote:                     appServices.AddNote,
		ListNotesByCaseFile:         appServices.ListNotesByCaseFile,
		AttachDocument:              appServices.AttachDocument,
		ImportDocument:              appServices.ImportDocument,
		ListDocumentsByCaseFile:     appServices.ListDocumentsByCaseFile,
		ExtractDocumentText:         appServices.ExtractDocumentText,
		SearchDocuments:             appServices.SearchDocuments,
		GetDocumentDetail:           appServices.GetDocumentDetail,
		GetCaseFileDetail:           appServices.GetCaseFileDetail,
		StorageAudit:                appServices.StorageAudit,
		StorageCleanOrphans:         appServices.StorageCleanOrphans,
		MimeNormalize:               appServices.MimeNormalize,
		VerifyDocument:              appServices.VerifyDocument,
		VerifyAllDocuments:          appServices.VerifyAllDocuments,
		ReindexDocument:             appServices.ReindexDocument,
		ReindexAllDocuments:         appServices.ReindexAllDocuments,
		AnalyzeDocumentMetadata:     appServices.AnalyzeDocumentMetadata,
		AnalyzeAllDocumentMetadata:  appServices.AnalyzeAllDocumentMetadata,
		AnalyzeDocumentEvents:       appServices.AnalyzeDocumentEvents,
		AnalyzeAllDocumentEvents:    appServices.AnalyzeAllDocumentEvents,
		ListDocumentEvents:          appServices.ListDocumentEvents,
		ListEventsByCaseFile:        appServices.ListEventsByCaseFile,
		ListUpcomingEvents:          appServices.ListUpcomingEvents,
		ExportUpcomingEventsICS:     appServices.ExportUpcomingEventsICS,
		MarkEventReviewed:           appServices.MarkEventReviewed,
		MarkEventResolved:           appServices.MarkEventResolved,
		ReopenEvent:                 appServices.ReopenEvent,
		FixCaseFile:                 appServices.FixCaseFile,
		AuditCaseFile:               appServices.AuditCaseFile,
		BackfillDocumentSearchIndex: appServices.BackfillDocumentSearchIndex,
	}

	if err := app.Run(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
