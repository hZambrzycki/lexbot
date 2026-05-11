package main

import (
	"context"
	"log"

	"lexbox/internal/infrastructure/bootstrap"
	httpapi "lexbox/internal/interfaces/http"
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

	server := httpapi.NewServer()
	server.RegisterRoutes(httpapi.RegisterHealthRoutes)

	eventHandler := httpapi.EventHandler{
		GetEvent:                appServices.GetEvent,
		ListEventsByCaseFile:    appServices.ListEventsByCaseFile,
		ListUpcomingEvents:      appServices.ListUpcomingEvents,
		ExportUpcomingEventsICS: appServices.ExportUpcomingEventsICS,
		MarkEventReviewed:       appServices.MarkEventReviewed,
		MarkEventResolved:       appServices.MarkEventResolved,
		ReopenEvent:             appServices.ReopenEvent,
	}

	caseFileHandler := httpapi.CaseFileHandler{
		CreateCaseFile:       appServices.CreateCaseFile,
		ListCaseFiles:        appServices.ListCaseFiles,
		GetCaseFileDetail:    appServices.GetCaseFileDetail,
		GetCaseFileDashboard: appServices.GetCaseFileDashboard,
		ImportDocument:       appServices.ImportDocument,
		EventHandler:         eventHandler,
		AddNote:              appServices.AddNote,
		DeleteNote:           appServices.DeleteNote,
	}

	documentHandler := httpapi.DocumentHandler{
		GetDocumentDetail: appServices.GetDocumentDetail,
		DeleteDocument:    appServices.DeleteDocument,
		ReprocessDocument: appServices.ReprocessDocument,
		ReviewDocument:    appServices.ReviewDocument,
		SearchDocuments:   appServices.SearchDocuments,
	}

	searchHandler := httpapi.SearchHandler{
		GlobalSearch: appServices.GlobalSearch,
	}

	server.RegisterRoutes(caseFileHandler.Register)
	server.RegisterRoutes(eventHandler.Register)
	server.RegisterRoutes(documentHandler.Register)
	server.RegisterRoutes(searchHandler.Register)

	log.Fatal(server.Start(":8080"))
}
