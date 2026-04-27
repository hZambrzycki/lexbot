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
		ListEventsByCaseFile: appServices.ListEventsByCaseFile,
		ListUpcomingEvents:   appServices.ListUpcomingEvents,
		MarkEventReviewed:    appServices.MarkEventReviewed,
		MarkEventResolved:    appServices.MarkEventResolved,
		ReopenEvent:          appServices.ReopenEvent,
	}

	caseFileHandler := httpapi.CaseFileHandler{
		ListCaseFiles:        appServices.ListCaseFiles,
		GetCaseFileDetail:    appServices.GetCaseFileDetail,
		GetCaseFileDashboard: appServices.GetCaseFileDashboard,
		EventHandler:         eventHandler,
	}

	server.RegisterRoutes(caseFileHandler.Register)
	server.RegisterRoutes(eventHandler.Register)

	log.Fatal(server.Start(":8080"))
}
