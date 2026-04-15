package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	casefileapp "lexbox/internal/application/casefile"
	clientapp "lexbox/internal/application/client"
	documentapp "lexbox/internal/application/document"
	noteapp "lexbox/internal/application/note"
	"lexbox/internal/domain/casefile"
)

type App struct {
	Out io.Writer

	CreateClient               clientapp.CreateClient
	GetClient                  clientapp.GetClient
	ListClients                clientapp.ListClients
	CreateCaseFile             casefileapp.CreateCaseFile
	ListCaseFiles              casefileapp.ListCaseFiles
	ListCaseFilesByClient      casefileapp.ListCaseFilesByClient
	AddNote                    noteapp.AddNote
	ListNotesByCaseFile        noteapp.ListNotesByCaseFile
	AttachDocument             documentapp.AttachDocument
	ImportDocument             documentapp.ImportDocument
	ListDocumentsByCaseFile    documentapp.ListDocumentsByCaseFile
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
}

func (a App) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return nil
	}

	switch args[0] {
	case "help":
		a.printHelp()
		return nil
	case "demo":
		return a.runDemo(ctx)
	case "client-create":
		return a.runCreateClient(ctx, args[1:])
	case "client-get":
		return a.runGetClient(ctx, args[1:])
	case "client-list":
		return a.runListClients(ctx)
	case "case-create":
		return a.runCreateCaseFile(ctx, args[1:])
	case "case-list":
		return a.runListCaseFiles(ctx)
	case "case-list-by-client":
		return a.runListCaseFilesByClient(ctx, args[1:])
	case "note-add":
		return a.runAddNote(ctx, args[1:])
	case "note-list-by-case":
		return a.runListNotesByCaseFile(ctx, args[1:])
	case "document-attach":
		return a.runAttachDocument(ctx, args[1:])
	case "document-import":
		return a.runImportDocument(ctx, args[1:])
	case "document-list-by-case":
		return a.runListDocumentsByCaseFile(ctx, args[1:])
	case "document-extract-text":
		return a.runExtractDocumentText(ctx, args[1:])
	case "document-reindex":
		return a.runReindexDocument(ctx, args[1:])
	case "document-reindex-all":
		return a.runReindexAllDocuments(ctx, args[1:])
	case "document-analyze-metadata":
		return a.runAnalyzeDocumentMetadata(ctx, args[1:])
	case "document-analyze-metadata-all":
		return a.runAnalyzeAllDocumentMetadata(ctx, args[1:])
	case "document-analyze-events":
		return a.runAnalyzeDocumentEvents(ctx, args[1:])
	case "document-analyze-events-all":
		return a.runAnalyzeAllDocumentEvents(ctx, args[1:])
	case "document-get":
		return a.runGetDocument(ctx, args[1:])
	case "document-verify":
		return a.runVerifyDocument(ctx, args[1:])
	case "document-verify-all":
		return a.runVerifyAllDocuments(ctx, args[1:])
	case "search":
		return a.runSearchDocuments(ctx, args[1:])
	case "case-get":
		return a.runGetCaseFile(ctx, args[1:])
	case "storage-audit":
		return a.runStorageAudit(ctx)
	case "storage-clean-orphans":
		return a.runStorageCleanOrphans(ctx)
	case "mime-normalize":
		return a.runMimeNormalize(ctx)
	case "document-events":
		return a.runListDocumentEvents(ctx, args[1:])
	case "events-list":
		return a.runListEventsByCaseFile(ctx, args[1:])
	case "events-upcoming":
		return a.runListUpcomingEvents(ctx, args[1:])
	case "events-export-ics":
		return a.runExportUpcomingEventsICS(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a App) printHelp() {
	fmt.Fprintln(a.Out, "LEXBOX CLI")
	fmt.Fprintln(a.Out)
	fmt.Fprintln(a.Out, "Commands:")
	fmt.Fprintln(a.Out, "  help")
	fmt.Fprintln(a.Out, "  demo")
	fmt.Fprintln(a.Out, `  client-create "<name>" "<email>" "<phone>" "<identifier>"`)
	fmt.Fprintln(a.Out, `  client-get "<clientID>"`)
	fmt.Fprintln(a.Out, "  client-list")
	fmt.Fprintln(a.Out, `  case-create "<clientID>" "<reference>" "<title>" "<type>" "<description>"`)
	fmt.Fprintln(a.Out, "  case-list")
	fmt.Fprintln(a.Out, `  case-list-by-client "<clientID>"`)
	fmt.Fprintln(a.Out, `  note-add "<caseFileID>" "<title>" "<content>"`)
	fmt.Fprintln(a.Out, `  note-list-by-case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  document-attach "<caseFileID>" "<originalName>" "<storagePath>" "<mimeType>" "<fileHash>"`)
	fmt.Fprintln(a.Out, `  document-import "<caseFileID>" "<sourcePath>" "<mimeType>"`)
	fmt.Fprintln(a.Out, `  document-import "<caseFileID>" "<sourcePath>" "<mimeType>" "<fileHash>"`)
	fmt.Fprintln(a.Out, `  document-list-by-case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  document-extract-text "<documentID>"`)
	fmt.Fprintln(a.Out, `  document-reindex "<documentID>"`)
	fmt.Fprintln(a.Out, `  document-reindex-all`)
	fmt.Fprintln(a.Out, `  document-reindex-all --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  document-analyze-metadata "<documentID>"`)
	fmt.Fprintln(a.Out, `  document-analyze-metadata-all`)
	fmt.Fprintln(a.Out, `  document-analyze-metadata-all --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  document-analyze-events "<documentID>"`)
	fmt.Fprintln(a.Out, `  document-analyze-events-all`)
	fmt.Fprintln(a.Out, `  document-analyze-events-all --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  document-get "<documentID>"`)
	fmt.Fprintln(a.Out, `  document-verify "<documentID>"`)
	fmt.Fprintln(a.Out, `  document-verify-all`)
	fmt.Fprintln(a.Out, `  document-verify-all --verbose`)
	fmt.Fprintln(a.Out, `  document-verify-all --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  document-verify-all --case "<caseFileID>" --verbose`)
	fmt.Fprintln(a.Out, `  search "<query>"`)
	fmt.Fprintln(a.Out, `  search "<query>" --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  case-get "<caseFileID>"`)
	fmt.Fprintln(a.Out, "  storage-audit")
	fmt.Fprintln(a.Out, "  storage-clean-orphans")
	fmt.Fprintln(a.Out, "  mime-normalize")
	fmt.Fprintln(a.Out, `  document-events "<documentID>"`)
	fmt.Fprintln(a.Out, `  events-list --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  events-upcoming`)
	fmt.Fprintln(a.Out, `  events-upcoming --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  events-upcoming --type "<eventType>"`)
	fmt.Fprintln(a.Out, `  events-upcoming --case "<caseFileID>" --type "<eventType>"`)
	fmt.Fprintln(a.Out, `  events-upcoming --relative-only`)
	fmt.Fprintln(a.Out, `  events-upcoming --verbose`)
	fmt.Fprintln(a.Out, `  events-upcoming --case "<caseFileID>" --relative-only --verbose`)
	fmt.Fprintln(a.Out, `  events-export-ics`)
	fmt.Fprintln(a.Out, `  events-export-ics --case "<caseFileID>"`)
	fmt.Fprintln(a.Out, `  events-export-ics --type "<eventType>"`)
	fmt.Fprintln(a.Out, `  events-export-ics --case "<caseFileID>" --type "<eventType>" --output "<filePath>"`)
}

func (a App) runDemo(ctx context.Context) error {
	clientEntity, err := a.CreateClient.Execute(ctx, clientapp.CreateClientInput{
		Name:       "María García Jiménez",
		Email:      "maria@example.com",
		Phone:      "600123123",
		Identifier: "12345678A",
	})
	if err != nil {
		return err
	}

	caseFileEntity, err := a.CreateCaseFile.Execute(ctx, casefileapp.CreateCaseFileInput{
		ClientID:    clientEntity.ID.String(),
		Reference:   "EXP-2026-0001",
		Title:       "Divorcio contencioso y medidas",
		Type:        casefile.TypeCivil,
		Description: "Procedimiento principal con medidas respecto de menores",
	})
	if err != nil {
		return err
	}

	noteEntity, err := a.AddNote.Execute(ctx, noteapp.AddNoteInput{
		CaseFileID: caseFileEntity.ID.String(),
		Title:      "Primera reunión",
		Content:    "La clienta aporta documentación inicial y explica antecedentes del procedimiento.",
	})
	if err != nil {
		return err
	}

	detail, err := a.GetCaseFileDetail.Execute(ctx, casefileapp.GetCaseFileDetailInput{
		ID: caseFileEntity.ID.String(),
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "=== LEXBOX :: DEMO SQLITE ===")
	fmt.Fprintln(a.Out)

	fmt.Fprintln(a.Out, "Cliente:")
	fmt.Fprintf(a.Out, "  ID: %s\n", clientEntity.ID)
	fmt.Fprintf(a.Out, "  Nombre: %s\n", clientEntity.Name)
	fmt.Fprintf(a.Out, "  Email: %s\n", clientEntity.Email)
	fmt.Fprintf(a.Out, "  Teléfono: %s\n", clientEntity.Phone)
	fmt.Fprintln(a.Out)

	fmt.Fprintln(a.Out, "Expediente:")
	fmt.Fprintf(a.Out, "  ID: %s\n", detail.CaseFile.ID)
	fmt.Fprintf(a.Out, "  ClienteID: %s\n", detail.CaseFile.ClientID)
	fmt.Fprintf(a.Out, "  Referencia: %s\n", detail.CaseFile.Reference)
	fmt.Fprintf(a.Out, "  Título: %s\n", detail.CaseFile.Title)
	fmt.Fprintf(a.Out, "  Tipo: %s\n", detail.CaseFile.Type)
	fmt.Fprintf(a.Out, "  Estado: %s\n", detail.CaseFile.Status)
	fmt.Fprintf(a.Out, "  Descripción: %s\n", detail.CaseFile.Description)
	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Notas (%d):\n", len(detail.Notes))
	for i, n := range detail.Notes {
		title := strings.TrimSpace(n.Title)
		if title == "" {
			title = "(untitled)"
		}
		fmt.Fprintf(a.Out, "  %d. [%s] %s\n", i+1, title, n.Content)
	}
	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Documentos (%d):\n", len(detail.Documents))
	for i, d := range detail.Documents {
		fmt.Fprintf(a.Out, "  %d. %s (%s)\n", i+1, d.OriginalName, d.StoragePath)
	}
	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Nota creada correctamente con ID: %s\n", noteEntity.ID)
	return nil
}

func (a App) runCreateClient(ctx context.Context, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf(`usage: client-create "<name>" "<email>" "<phone>" "<identifier>"`)
	}

	clientEntity, err := a.CreateClient.Execute(ctx, clientapp.CreateClientInput{
		Name:       args[0],
		Email:      args[1],
		Phone:      args[2],
		Identifier: args[3],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Client created successfully")
	fmt.Fprintf(a.Out, "ID: %s\n", clientEntity.ID)
	fmt.Fprintf(a.Out, "Name: %s\n", clientEntity.Name)
	return nil
}

func (a App) runGetClient(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: client-get "<clientID>"`)
	}

	clientEntity, err := a.GetClient.Execute(ctx, clientapp.GetClientInput{
		ID: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Client detail")
	fmt.Fprintf(a.Out, "ID: %s\n", clientEntity.ID)
	fmt.Fprintf(a.Out, "Name: %s\n", clientEntity.Name)
	fmt.Fprintf(a.Out, "Email: %s\n", clientEntity.Email)
	fmt.Fprintf(a.Out, "Phone: %s\n", clientEntity.Phone)
	fmt.Fprintf(a.Out, "Identifier: %s\n", clientEntity.Identifier)
	return nil
}

func (a App) runListClients(ctx context.Context) error {
	clients, err := a.ListClients.Execute(ctx)
	if err != nil {
		return err
	}

	if len(clients) == 0 {
		fmt.Fprintln(a.Out, "No clients found")
		return nil
	}

	fmt.Fprintf(a.Out, "Clients (%d):\n", len(clients))
	for i, c := range clients {
		fmt.Fprintf(a.Out, "  %d. %s | %s | %s | %s\n", i+1, c.ID, c.Name, c.Email, c.Phone)
	}
	return nil
}

func (a App) runCreateCaseFile(ctx context.Context, args []string) error {
	if len(args) < 5 {
		return fmt.Errorf(`usage: case-create "<clientID>" "<reference>" "<title>" "<type>" "<description>"`)
	}

	cfType, err := parseCaseFileType(args[3])
	if err != nil {
		return err
	}

	caseFileEntity, err := a.CreateCaseFile.Execute(ctx, casefileapp.CreateCaseFileInput{
		ClientID:    args[0],
		Reference:   args[1],
		Title:       args[2],
		Type:        cfType,
		Description: args[4],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Case file created successfully")
	fmt.Fprintf(a.Out, "ID: %s\n", caseFileEntity.ID)
	fmt.Fprintf(a.Out, "Title: %s\n", caseFileEntity.Title)
	return nil
}

func (a App) runListCaseFiles(ctx context.Context) error {
	caseFiles, err := a.ListCaseFiles.Execute(ctx)
	if err != nil {
		return err
	}

	if len(caseFiles) == 0 {
		fmt.Fprintln(a.Out, "No case files found")
		return nil
	}

	fmt.Fprintf(a.Out, "Case files (%d):\n", len(caseFiles))
	for i, cf := range caseFiles {
		fmt.Fprintf(
			a.Out,
			"  %d. %s | client=%s | ref=%s | title=%s | type=%s | status=%s\n",
			i+1, cf.ID, cf.ClientID, cf.Reference, cf.Title, cf.Type, cf.Status,
		)
	}
	return nil
}

func (a App) runListCaseFilesByClient(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: case-list-by-client "<clientID>"`)
	}

	caseFiles, err := a.ListCaseFilesByClient.Execute(ctx, casefileapp.ListCaseFilesByClientInput{
		ClientID: args[0],
	})
	if err != nil {
		return err
	}

	if len(caseFiles) == 0 {
		fmt.Fprintln(a.Out, "No case files found for this client")
		return nil
	}

	fmt.Fprintf(a.Out, "Case files for client %s (%d):\n", args[0], len(caseFiles))
	for i, cf := range caseFiles {
		fmt.Fprintf(
			a.Out,
			"  %d. %s | ref=%s | title=%s | type=%s | status=%s\n",
			i+1, cf.ID, cf.Reference, cf.Title, cf.Type, cf.Status,
		)
	}
	return nil
}

func (a App) runAddNote(ctx context.Context, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf(`usage: note-add "<caseFileID>" "<title>" "<content>"`)
	}

	noteEntity, err := a.AddNote.Execute(ctx, noteapp.AddNoteInput{
		CaseFileID: args[0],
		Title:      args[1],
		Content:    args[2],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Note added successfully")
	fmt.Fprintf(a.Out, "ID: %s\n", noteEntity.ID)
	fmt.Fprintf(a.Out, "Title: %s\n", noteEntity.Title)
	return nil
}

func (a App) runListNotesByCaseFile(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: note-list-by-case "<caseFileID>"`)
	}

	notes, err := a.ListNotesByCaseFile.Execute(ctx, noteapp.ListNotesByCaseFileInput{
		CaseFileID: args[0],
	})
	if err != nil {
		return err
	}

	if len(notes) == 0 {
		fmt.Fprintln(a.Out, "No notes found for this case file")
		return nil
	}

	fmt.Fprintf(a.Out, "Notes for case file %s (%d):\n", args[0], len(notes))
	for i, n := range notes {
		title := strings.TrimSpace(n.Title)
		if title == "" {
			title = "(untitled)"
		}
		fmt.Fprintf(a.Out, "  %d. %s | %s\n", i+1, title, n.Content)
	}
	return nil
}

func (a App) runAttachDocument(ctx context.Context, args []string) error {
	if len(args) < 5 {
		return fmt.Errorf(`usage: document-attach "<caseFileID>" "<originalName>" "<storagePath>" "<mimeType>" "<fileHash>"`)
	}

	doc, err := a.AttachDocument.Execute(ctx, documentapp.AttachDocumentInput{
		CaseFileID:   args[0],
		OriginalName: args[1],
		StoragePath:  args[2],
		MimeType:     args[3],
		FileHash:     args[4],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Document attached successfully")
	fmt.Fprintf(a.Out, "ID: %s\n", doc.ID)
	fmt.Fprintf(a.Out, "Original name: %s\n", doc.OriginalName)
	fmt.Fprintf(a.Out, "Storage path: %s\n", doc.StoragePath)
	return nil
}

func (a App) runImportDocument(ctx context.Context, args []string) error {
	if len(args) < 3 || len(args) > 4 {
		return fmt.Errorf(`usage: document-import "<caseFileID>" "<sourcePath>" "<mimeType>" ["<fileHash>"]`)
	}

	fileHash := ""
	if len(args) == 4 {
		fileHash = args[3]
	}

	result, err := a.ImportDocument.Execute(ctx, documentapp.ImportDocumentInput{
		CaseFileID: args[0],
		SourcePath: args[1],
		MimeType:   args[2],
		FileHash:   fileHash,
	})
	if err != nil {
		var dupErr documentapp.DuplicateDocumentError
		if errors.As(err, &dupErr) {
			fmt.Fprintln(a.Out, "Document already exists in this case file")
			fmt.Fprintf(a.Out, "Existing document_id=%s\n", dupErr.Existing.ID)
			fmt.Fprintf(a.Out, "Original name=%s\n", dupErr.Existing.OriginalName)
			fmt.Fprintf(a.Out, "Storage path=%s\n", dupErr.Existing.StoragePath)
			return nil
		}
		return err
	}

	doc := result.Document

	fmt.Fprintln(a.Out, "Document imported successfully")
	fmt.Fprintf(a.Out, "ID: %s\n", doc.ID)
	fmt.Fprintf(a.Out, "Original name: %s\n", doc.OriginalName)
	fmt.Fprintf(a.Out, "Storage path: %s\n", doc.StoragePath)
	fmt.Fprintf(a.Out, "File hash: %s\n", doc.FileHash)

	if result.TextExtracted {
		fmt.Fprintln(a.Out, "Text indexed: yes")
	} else {
		fmt.Fprintln(a.Out, "Text indexed: no")
	}

	if result.MetadataAnalyzed {
		fmt.Fprintln(a.Out, "Metadata analyzed: yes")
	} else {
		fmt.Fprintln(a.Out, "Metadata analyzed: no")
	}

	if result.EventsAnalyzed {
		fmt.Fprintln(a.Out, "Events analyzed: yes")
		fmt.Fprintf(a.Out, "Events detected: %d\n", result.EventsDetected)
	} else {
		fmt.Fprintln(a.Out, "Events analyzed: no")
		fmt.Fprintln(a.Out, "Events detected: 0")
	}

	return nil
}

func (a App) runListDocumentsByCaseFile(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-list-by-case "<caseFileID>"`)
	}

	documents, err := a.ListDocumentsByCaseFile.Execute(ctx, documentapp.ListDocumentsByCaseFileInput{
		CaseFileID: args[0],
	})
	if err != nil {
		return err
	}

	if len(documents) == 0 {
		fmt.Fprintln(a.Out, "No documents found for this case file")
		return nil
	}

	fmt.Fprintf(a.Out, "Documents for case file %s (%d):\n", args[0], len(documents))
	for i, item := range documents {
		d := item.Document

		fmt.Fprintf(
			a.Out,
			"  %d. id=%s | name=%s | path=%s | mime=%s | type=%s | area=%s\n",
			i+1,
			d.ID,
			d.OriginalName,
			d.StoragePath,
			d.MimeType,
			item.DocumentType,
			item.LegalArea,
		)
	}

	return nil
}

func (a App) runExtractDocumentText(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-extract-text "<documentID>"`)
	}

	content, err := a.ExtractDocumentText.Execute(ctx, documentapp.ExtractDocumentTextInput{
		DocumentID: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Document text extracted successfully")
	fmt.Fprintln(a.Out, "Content:")
	fmt.Fprintln(a.Out, content)
	return nil
}

func (a App) runReindexDocument(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-reindex "<documentID>"`)
	}

	result, err := a.ReindexDocument.Execute(ctx, documentapp.ReindexDocumentInput{
		DocumentID: args[0],
	})
	if err != nil {
		switch {
		case errors.Is(err, documentapp.ErrReindexFileNotFound):
			return fmt.Errorf("cannot reindex document: file not found")
		case errors.Is(err, documentapp.ErrReindexUnsupported):
			return fmt.Errorf("cannot reindex document: unsupported file type")
		case errors.Is(err, documentapp.ErrReindexEmptyContent):
			return fmt.Errorf("cannot reindex document: extracted content is empty")
		case errors.Is(err, documentapp.ErrReindexReadFailed):
			return fmt.Errorf("cannot reindex document: read error")
		default:
			return err
		}
	}

	fmt.Fprintln(a.Out, "Document reindexed successfully")
	fmt.Fprintf(a.Out, "Document ID: %s\n", result.DocumentID)
	fmt.Fprintf(a.Out, "Extracted text length: %d\n", result.ExtractedLength)
	return nil
}

func (a App) runReindexAllDocuments(ctx context.Context, args []string) error {
	caseFileID := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--case":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: document-reindex-all [--case "<caseFileID>"]`)
			}
			caseFileID = strings.TrimSpace(args[i+1])
			i++
		default:
			return fmt.Errorf(`usage: document-reindex-all [--case "<caseFileID>"]`)
		}
	}

	result, err := a.ReindexAllDocuments.Execute(ctx, documentapp.ReindexAllDocumentsInput{
		CaseFileID: caseFileID,
	})
	if err != nil {
		return err
	}

	if caseFileID != "" {
		fmt.Fprintf(a.Out, "Document reindex for case file %s\n", caseFileID)
	} else {
		fmt.Fprintln(a.Out, "Document reindex (all)")
	}

	fmt.Fprintf(a.Out, "Scanned: %d\n", result.Scanned)
	fmt.Fprintf(a.Out, "Reindexed: %d\n", result.Reindexed)
	fmt.Fprintf(a.Out, "Skipped: %d\n", result.Skipped)
	fmt.Fprintf(a.Out, "  Missing file: %d\n", result.SkippedMissingFile)
	fmt.Fprintf(a.Out, "  Already indexed: %d\n", result.SkippedAlreadyIndexed)
	fmt.Fprintf(a.Out, "  Unsupported: %d\n", result.SkippedUnsupported)
	fmt.Fprintf(a.Out, "  Empty content: %d\n", result.SkippedEmptyContent)
	fmt.Fprintf(a.Out, "Errors: %d\n", result.Errors)

	return nil
}

func (a App) runAnalyzeDocumentMetadata(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-analyze-metadata "<documentID>"`)
	}

	result, err := a.AnalyzeDocumentMetadata.Execute(ctx, documentapp.AnalyzeDocumentMetadataInput{
		DocumentID: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Document metadata analyzed successfully")
	fmt.Fprintf(a.Out, "Document ID: %s\n", result.DocumentID)
	fmt.Fprintf(a.Out, "Document type: %s\n", result.DocumentType)
	fmt.Fprintf(a.Out, "Legal area: %s\n", result.LegalArea)
	fmt.Fprintf(a.Out, "Analyzed at: %s\n", result.AnalyzedAt)
	return nil
}

func (a App) runAnalyzeAllDocumentMetadata(ctx context.Context, args []string) error {
	caseFileID := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--case":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: document-analyze-metadata-all [--case "<caseFileID>"]`)
			}
			caseFileID = strings.TrimSpace(args[i+1])
			i++
		default:
			return fmt.Errorf(`usage: document-analyze-metadata-all [--case "<caseFileID>"]`)
		}
	}

	result, err := a.AnalyzeAllDocumentMetadata.Execute(ctx, documentapp.AnalyzeAllDocumentMetadataInput{
		CaseFileID: caseFileID,
	})
	if err != nil {
		return err
	}

	if caseFileID != "" {
		fmt.Fprintf(a.Out, "Document metadata analysis for case file %s\n", caseFileID)
	} else {
		fmt.Fprintln(a.Out, "Document metadata analysis (all)")
	}

	fmt.Fprintf(a.Out, "Scanned: %d\n", result.Scanned)
	fmt.Fprintf(a.Out, "Analyzed: %d\n", result.Analyzed)
	fmt.Fprintf(a.Out, "Skipped: %d\n", result.Skipped)
	fmt.Fprintf(a.Out, "  No extracted text: %d\n", result.SkippedNoExtractedText)
	fmt.Fprintf(a.Out, "Errors: %d\n", result.Errors)

	return nil
}

func (a App) runGetDocument(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-get "<documentID>"`)
	}

	result, err := a.GetDocumentDetail.Execute(ctx, documentapp.GetDocumentDetailInput{
		DocumentID: args[0],
	})
	if err != nil {
		return err
	}

	doc := result.Document

	fmt.Fprintln(a.Out, "Document detail")
	fmt.Fprintf(a.Out, "ID: %s\n", doc.ID)
	fmt.Fprintf(a.Out, "CaseFileID: %s\n", doc.CaseFileID)
	fmt.Fprintf(a.Out, "Original name: %s\n", doc.OriginalName)
	fmt.Fprintf(a.Out, "Storage path: %s\n", doc.StoragePath)
	fmt.Fprintf(a.Out, "File exists: %v\n", result.FileExists)
	fmt.Fprintf(a.Out, "Mime type: %s\n", doc.MimeType)
	fmt.Fprintf(a.Out, "File hash: %s\n", doc.FileHash)
	fmt.Fprintf(a.Out, "Created at: %s\n", doc.CreatedAt.Time().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(a.Out, "Updated at: %s\n", doc.UpdatedAt.Time().Format("2006-01-02 15:04:05"))

	if result.HasExtractedText {
		fmt.Fprintln(a.Out, "Has extracted text: yes")
		fmt.Fprintf(a.Out, "Extracted text length: %d\n", result.ExtractedTextLength)
		if result.ExtractedTextPreview != "" {
			fmt.Fprintf(a.Out, "Extracted text preview: %s\n", result.ExtractedTextPreview)
		}
	} else {
		fmt.Fprintln(a.Out, "Has extracted text: no")
		fmt.Fprintln(a.Out, "Extracted text length: 0")
	}

	if result.HasMetadata {
		fmt.Fprintln(a.Out, "Has metadata: yes")
		fmt.Fprintf(a.Out, "Document type: %s\n", result.DocumentType)
		fmt.Fprintf(a.Out, "Legal area: %s\n", result.LegalArea)
		fmt.Fprintf(a.Out, "Metadata analyzed at: %s\n", result.MetadataAnalyzedAt)
	} else {
		fmt.Fprintln(a.Out, "Has metadata: no")
	}

	if result.HasEvents {
		fmt.Fprintf(a.Out, "Events (%d):\n", len(result.Events))
		for i, e := range result.Events {
			fmt.Fprintf(a.Out, "  %d. type=%s | date=%s\n", i+1, e.EventType, e.EventDate)
			fmt.Fprintf(a.Out, "     source=%s\n", e.SourceText)
		}
	} else {
		fmt.Fprintln(a.Out, "Events (0):")
	}

	if result.HasEventTrace {
		fmt.Fprintf(a.Out, "Event trace (%d):\n", len(result.EventTrace))
		for i, e := range result.EventTrace {
			fmt.Fprintf(a.Out, "  %d. type=%s | date=%s | kind=%s\n", i+1, e.EventType, e.EventDate, e.DateKind)
			fmt.Fprintf(a.Out, "     source=%s\n", e.SourceText)

			if strings.TrimSpace(e.AnchorDate) != "" {
				fmt.Fprintf(a.Out, "     anchor_date=%s\n", e.AnchorDate)
			}
			if strings.TrimSpace(e.AnchorSource) != "" {
				fmt.Fprintf(a.Out, "     anchor_source=%s\n", e.AnchorSource)
			}
			if e.RelativeDays > 0 {
				fmt.Fprintf(a.Out, "     relative_days=%d\n", e.RelativeDays)
			}
			if e.DateKind == "relative" {
				fmt.Fprintf(a.Out, "     business_days=%v\n", e.IsBusinessDays)
			}
			if strings.TrimSpace(e.TriggerText) != "" {
				fmt.Fprintf(a.Out, "     trigger=%s\n", e.TriggerText)
			}
		}
	} else {
		fmt.Fprintln(a.Out, "Event trace (0):")
	}

	return nil
}

func (a App) runVerifyDocument(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-verify "<documentID>"`)
	}

	result, err := a.VerifyDocument.Execute(ctx, documentapp.VerifyDocumentInput{
		DocumentID: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Document verification")
	fmt.Fprintf(a.Out, "ID: %s\n", result.DocumentID)
	fmt.Fprintf(a.Out, "Exists in DB: %v\n", result.ExistsInDB)
	fmt.Fprintf(a.Out, "File exists: %v\n", result.FileExists)
	fmt.Fprintf(a.Out, "Mime type: %s\n", result.MimeType)
	fmt.Fprintf(a.Out, "Normalized mime: %s\n", result.MimeNormalized)
	fmt.Fprintf(a.Out, "Mime OK: %v\n", result.IsMimeNormalized)
	fmt.Fprintf(a.Out, "Has hash: %v\n", result.HasHash)
	fmt.Fprintf(a.Out, "Has extracted text: %v\n", result.HasExtractedText)

	return nil
}

func (a App) runVerifyAllDocuments(ctx context.Context, args []string) error {
	verbose := false
	caseFileID := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--verbose":
			verbose = true
		case "--case":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: document-verify-all [--case "<caseFileID>"] [--verbose]`)
			}
			caseFileID = strings.TrimSpace(args[i+1])
			i++
		default:
			return fmt.Errorf(`usage: document-verify-all [--case "<caseFileID>"] [--verbose]`)
		}
	}

	result, err := a.VerifyAllDocuments.Execute(ctx, documentapp.VerifyAllDocumentsInput{
		CaseFileID: caseFileID,
	})
	if err != nil {
		return err
	}

	if caseFileID != "" {
		fmt.Fprintf(a.Out, "Document verification for case file %s\n", caseFileID)
	} else {
		fmt.Fprintln(a.Out, "Document verification (all)")
	}

	fmt.Fprintf(a.Out, "Total: %d\n", result.Total)
	fmt.Fprintf(a.Out, "OK: %d\n", result.OK)
	fmt.Fprintf(a.Out, "Missing file: %d\n", result.MissingFile)
	fmt.Fprintf(a.Out, "Missing hash: %d\n", result.MissingHash)
	fmt.Fprintf(a.Out, "Missing text: %d\n", result.MissingText)
	fmt.Fprintf(a.Out, "Invalid mime: %d\n", result.InvalidMime)

	if !verbose {
		return nil
	}

	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Missing file items (%d):\n", len(result.MissingFileItems))
	if len(result.MissingFileItems) == 0 {
		fmt.Fprintln(a.Out, "  none")
	} else {
		for i, item := range result.MissingFileItems {
			fmt.Fprintf(a.Out, "  %d. document_id=%s\n", i+1, item.DocumentID)
			fmt.Fprintf(a.Out, "     original_name=%s\n", item.OriginalName)
			fmt.Fprintf(a.Out, "     storage_path=%s\n", item.StoragePath)
		}
	}

	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Missing hash items (%d):\n", len(result.MissingHashItems))
	if len(result.MissingHashItems) == 0 {
		fmt.Fprintln(a.Out, "  none")
	} else {
		for i, item := range result.MissingHashItems {
			fmt.Fprintf(a.Out, "  %d. document_id=%s\n", i+1, item.DocumentID)
			fmt.Fprintf(a.Out, "     original_name=%s\n", item.OriginalName)
			fmt.Fprintf(a.Out, "     storage_path=%s\n", item.StoragePath)
		}
	}

	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Missing text items (%d):\n", len(result.MissingTextItems))
	if len(result.MissingTextItems) == 0 {
		fmt.Fprintln(a.Out, "  none")
	} else {
		for i, item := range result.MissingTextItems {
			fmt.Fprintf(a.Out, "  %d. document_id=%s\n", i+1, item.DocumentID)
			fmt.Fprintf(a.Out, "     original_name=%s\n", item.OriginalName)
			fmt.Fprintf(a.Out, "     storage_path=%s\n", item.StoragePath)
		}
	}

	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Invalid mime items (%d):\n", len(result.InvalidMimeItems))
	if len(result.InvalidMimeItems) == 0 {
		fmt.Fprintln(a.Out, "  none")
	} else {
		for i, item := range result.InvalidMimeItems {
			fmt.Fprintf(a.Out, "  %d. document_id=%s\n", i+1, item.DocumentID)
			fmt.Fprintf(a.Out, "     original_name=%s\n", item.OriginalName)
			fmt.Fprintf(a.Out, "     storage_path=%s\n", item.StoragePath)
			fmt.Fprintf(a.Out, "     mime_type=%s\n", item.MimeType)
			fmt.Fprintf(a.Out, "     normalized_mime=%s\n", item.NormalizedMime)
		}
	}

	return nil
}

func (a App) runSearchDocuments(ctx context.Context, args []string) error {
	input, err := parseSearchDocumentsInput(args)
	if err != nil {
		return err
	}

	results, err := a.SearchDocuments.Execute(ctx, input)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Fprintln(a.Out, "No matching documents found")
		return nil
	}

	if input.CaseFileID != "" {
		fmt.Fprintf(a.Out, "Search results for case file %s (%d):\n", input.CaseFileID, len(results))
	} else {
		fmt.Fprintf(a.Out, "Search results (%d):\n", len(results))
	}

	for i, r := range results {
		fmt.Fprintf(a.Out, "  %d. document_id=%s\n", i+1, r.DocumentID)
		fmt.Fprintf(a.Out, "     original_name=%s\n", r.OriginalName)
		fmt.Fprintf(a.Out, "     case_file_id=%s\n", r.CaseFileID)
		fmt.Fprintf(a.Out, "     snippet=%s\n", r.Snippet)
	}

	return nil
}

func (a App) runGetCaseFile(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: case-get "<caseFileID>"`)
	}

	detail, err := a.GetCaseFileDetail.Execute(ctx, casefileapp.GetCaseFileDetailInput{
		ID: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Case file detail")
	fmt.Fprintf(a.Out, "ID: %s\n", detail.CaseFile.ID)
	fmt.Fprintf(a.Out, "ClientID: %s\n", detail.CaseFile.ClientID)
	fmt.Fprintf(a.Out, "Reference: %s\n", detail.CaseFile.Reference)
	fmt.Fprintf(a.Out, "Title: %s\n", detail.CaseFile.Title)
	fmt.Fprintf(a.Out, "Type: %s\n", detail.CaseFile.Type)
	fmt.Fprintf(a.Out, "Status: %s\n", detail.CaseFile.Status)
	fmt.Fprintf(a.Out, "Description: %s\n", detail.CaseFile.Description)
	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Notes (%d):\n", len(detail.Notes))
	for i, n := range detail.Notes {
		title := strings.TrimSpace(n.Title)
		if title == "" {
			title = "(untitled)"
		}
		fmt.Fprintf(a.Out, "  %d. [%s] %s\n", i+1, title, n.Content)
	}

	fmt.Fprintf(a.Out, "Documents (%d):\n", len(detail.Documents))
	for i, d := range detail.Documents {
		fmt.Fprintf(a.Out, "  %d. %s (%s)\n", i+1, d.OriginalName, d.StoragePath)
	}
	return nil
}

func (a App) runStorageAudit(ctx context.Context) error {
	result, err := a.StorageAudit.Execute(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Storage audit")
	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Orphan files (%d):\n", len(result.OrphanFiles))
	if len(result.OrphanFiles) == 0 {
		fmt.Fprintln(a.Out, "  none")
	} else {
		for i, path := range result.OrphanFiles {
			fmt.Fprintf(a.Out, "  %d. %s\n", i+1, path)
		}
	}

	fmt.Fprintln(a.Out)

	fmt.Fprintf(a.Out, "Missing files (%d):\n", len(result.MissingDocuments))
	if len(result.MissingDocuments) == 0 {
		fmt.Fprintln(a.Out, "  none")
	} else {
		for i, doc := range result.MissingDocuments {
			fmt.Fprintf(a.Out, "  %d. document_id=%s\n", i+1, doc.DocumentID)
			fmt.Fprintf(a.Out, "     original_name=%s\n", doc.OriginalName)
			fmt.Fprintf(a.Out, "     storage_path=%s\n", doc.StoragePath)
		}
	}

	return nil
}

func (a App) runStorageCleanOrphans(ctx context.Context) error {
	result, err := a.StorageCleanOrphans.Execute(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Storage orphan cleanup")

	if len(result.DeletedFiles) == 0 {
		fmt.Fprintln(a.Out, "No orphan files deleted")
		return nil
	}

	fmt.Fprintf(a.Out, "Deleted orphan files (%d):\n", len(result.DeletedFiles))
	for i, path := range result.DeletedFiles {
		fmt.Fprintf(a.Out, "  %d. %s\n", i+1, path)
	}

	return nil
}

func (a App) runMimeNormalize(ctx context.Context) error {
	result, err := a.MimeNormalize.Execute(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Mime normalize")
	fmt.Fprintf(a.Out, "Scanned: %d\n", result.Scanned)
	fmt.Fprintf(a.Out, "Updated: %d\n", result.Updated)

	if len(result.Changed) == 0 {
		fmt.Fprintln(a.Out, "No mime types changed")
		return nil
	}

	fmt.Fprintln(a.Out, "Changed documents:")
	for i, item := range result.Changed {
		fmt.Fprintf(a.Out, "  %d. document_id=%s\n", i+1, item.DocumentID)
		fmt.Fprintf(a.Out, "     old=%s\n", item.OldMimeType)
		fmt.Fprintf(a.Out, "     new=%s\n", item.NewMimeType)
	}

	return nil
}

func parseSearchDocumentsInput(args []string) (documentapp.SearchDocumentsInput, error) {
	if len(args) < 1 {
		return documentapp.SearchDocumentsInput{}, fmt.Errorf(`usage: search "<query>" [--case "<caseFileID>"]`)
	}

	input := documentapp.SearchDocumentsInput{
		Query: strings.TrimSpace(args[0]),
		Limit: 20,
	}

	if input.Query == "" {
		return documentapp.SearchDocumentsInput{}, fmt.Errorf(`usage: search "<query>" [--case "<caseFileID>"]`)
	}

	rest := args[1:]
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--case":
			if i+1 >= len(rest) {
				return documentapp.SearchDocumentsInput{}, fmt.Errorf(`usage: search "<query>" [--case "<caseFileID>"]`)
			}
			input.CaseFileID = strings.TrimSpace(rest[i+1])
			i++
		default:
			return documentapp.SearchDocumentsInput{}, fmt.Errorf("unknown search option: %s", rest[i])
		}
	}

	return input, nil
}

func parseCaseFileType(value string) (casefile.Type, error) {
	t := casefile.Type(strings.TrimSpace(strings.ToLower(value)))
	if !casefile.IsValidType(t) {
		return "", fmt.Errorf("invalid case file type: %s", value)
	}
	return t, nil
}
func (a App) runAnalyzeDocumentEvents(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-analyze-events "<documentID>"`)
	}

	result, err := a.AnalyzeDocumentEvents.Execute(ctx, documentapp.AnalyzeDocumentEventsInput{
		DocumentID: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Document events analyzed successfully")
	fmt.Fprintf(a.Out, "Document ID: %s\n", result.DocumentID)
	fmt.Fprintf(a.Out, "Detected events: %d\n", result.Detected)

	return nil
}

func (a App) runAnalyzeAllDocumentEvents(ctx context.Context, args []string) error {
	caseFileID := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--case":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: document-analyze-events-all [--case "<caseFileID>"]`)
			}
			caseFileID = strings.TrimSpace(args[i+1])
			i++
		default:
			return fmt.Errorf(`usage: document-analyze-events-all [--case "<caseFileID>"]`)
		}
	}

	result, err := a.AnalyzeAllDocumentEvents.Execute(ctx, documentapp.AnalyzeAllDocumentEventsInput{
		CaseFileID: caseFileID,
	})
	if err != nil {
		return err
	}

	if caseFileID != "" {
		fmt.Fprintf(a.Out, "Document event analysis for case file %s\n", caseFileID)
	} else {
		fmt.Fprintln(a.Out, "Document event analysis (all)")
	}

	fmt.Fprintf(a.Out, "Scanned: %d\n", result.Scanned)
	fmt.Fprintf(a.Out, "Analyzed: %d\n", result.Analyzed)
	fmt.Fprintf(a.Out, "Skipped: %d\n", result.Skipped)
	fmt.Fprintf(a.Out, "  No extracted text: %d\n", result.SkippedNoExtractedText)
	fmt.Fprintf(a.Out, "Errors: %d\n", result.Errors)

	return nil
}

func (a App) runListDocumentEvents(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`usage: document-events "<documentID>"`)
	}

	events, err := a.ListDocumentEvents.Execute(ctx, documentapp.ListDocumentEventsInput{
		DocumentID: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "Document events")
	fmt.Fprintf(a.Out, "Document ID: %s\n", args[0])

	if len(events) == 0 {
		fmt.Fprintln(a.Out, "No events found")
		return nil
	}

	for i, e := range events {
		fmt.Fprintf(
			a.Out,
			"  %d. type=%s | date=%s\n     source=%s\n",
			i+1,
			e.EventType,
			e.EventDate,
			e.SourceText,
		)
	}

	return nil
}

func (a App) runListEventsByCaseFile(ctx context.Context, args []string) error {
	caseFileID := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--case":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: events-list --case "<caseFileID>"`)
			}
			caseFileID = strings.TrimSpace(args[i+1])
			i++
		default:
			return fmt.Errorf(`usage: events-list --case "<caseFileID>"`)
		}
	}

	if caseFileID == "" {
		return fmt.Errorf(`usage: events-list --case "<caseFileID>"`)
	}

	events, err := a.ListEventsByCaseFile.Execute(ctx, documentapp.ListEventsByCaseFileInput{
		CaseFileID: caseFileID,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(a.Out, "Events for case file %s\n", caseFileID)

	if len(events) == 0 {
		fmt.Fprintln(a.Out, "No events found")
		return nil
	}

	for i, e := range events {
		fmt.Fprintf(
			a.Out,
			"  %d. date=%s | type=%s | document=%s\n     document_id=%s\n     source=%s\n",
			i+1,
			e.EventDate,
			e.EventType,
			e.OriginalName,
			e.DocumentID,
			e.SourceText,
		)
	}

	return nil
}

func (a App) runListUpcomingEvents(ctx context.Context, args []string) error {
	caseFileID := ""
	eventType := ""
	relativeOnly := false
	verbose := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--case":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: events-upcoming [--case "<caseFileID>"] [--type "<eventType>"] [--relative-only] [--verbose]`)
			}
			caseFileID = strings.TrimSpace(args[i+1])
			i++
		case "--type":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: events-upcoming [--case "<caseFileID>"] [--type "<eventType>"] [--relative-only] [--verbose]`)
			}
			eventType = strings.TrimSpace(strings.ToLower(args[i+1]))
			i++
		case "--relative-only":
			relativeOnly = true
		case "--verbose":
			verbose = true
		default:
			return fmt.Errorf(`usage: events-upcoming [--case "<caseFileID>"] [--type "<eventType>"] [--relative-only] [--verbose]`)
		}
	}

	events, err := a.ListUpcomingEvents.Execute(ctx, documentapp.ListUpcomingEventsInput{
		CaseFileID:   caseFileID,
		EventType:    eventType,
		RelativeOnly: relativeOnly,
	})
	if err != nil {
		return err
	}

	title := "Upcoming events"
	if caseFileID != "" && eventType != "" {
		title = fmt.Sprintf("Upcoming events for case file %s (type=%s)", caseFileID, eventType)
	} else if caseFileID != "" {
		title = fmt.Sprintf("Upcoming events for case file %s", caseFileID)
	} else if eventType != "" {
		title = fmt.Sprintf("Upcoming events (type=%s)", eventType)
	}

	if relativeOnly {
		title += " [relative only]"
	}
	if verbose {
		title += " [verbose]"
	}

	fmt.Fprintln(a.Out, title)

	if len(events) == 0 {
		fmt.Fprintln(a.Out, "No upcoming events found")
		return nil
	}

	for i, e := range events {
		statusText := ""
		switch e.Status {
		case "overdue":
			statusText = fmt.Sprintf("OVERDUE (%d days ago)", -e.DaysRemaining)
		case "today":
			statusText = "TODAY"
		default:
			statusText = fmt.Sprintf("in %d days", e.DaysRemaining)
		}

		duplicateText := ""
		if e.DuplicateCount > 1 {
			duplicateText = fmt.Sprintf(" | duplicated in %d documents", e.DuplicateCount)
		}

		fmt.Fprintf(
			a.Out,
			"  %d. date=%s | type=%s | priority=%s | %s%s\n",
			i+1,
			e.EventDate,
			e.EventType,
			strings.ToUpper(e.Priority),
			statusText,
			duplicateText,
		)

		if len(e.DocumentNames) > 0 {
			fmt.Fprintf(a.Out, "     documents=%s\n", strings.Join(e.DocumentNames, ", "))
		}

		if len(e.DocumentIDs) > 0 {
			fmt.Fprintf(a.Out, "     document_ids=%s\n", strings.Join(e.DocumentIDs, ", "))
		}

		fmt.Fprintf(a.Out, "     source=%s\n", e.SourceText)

		if verbose {
			if e.DateKind == "relative" {
				fmt.Fprintf(
					a.Out,
					"     derived=relative | anchor=%s (%s) | relative_days=%d | business_days=%v\n",
					e.AnchorDate,
					e.AnchorSource,
					e.RelativeDays,
					e.IsBusinessDays,
				)
				if strings.TrimSpace(e.TriggerText) != "" {
					fmt.Fprintf(a.Out, "     trigger=%s\n", e.TriggerText)
				}
			} else if e.DateKind == "absolute" {
				fmt.Fprintln(a.Out, "     derived=absolute")
			}
		}
	}

	return nil
}

func (a App) runExportUpcomingEventsICS(ctx context.Context, args []string) error {
	caseFileID := ""
	eventType := ""
	outputPath := "lexbox-upcoming-events.ics"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--case":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: events-export-ics [--case "<caseFileID>"] [--type "<eventType>"] [--output "<filePath>"]`)
			}
			caseFileID = strings.TrimSpace(args[i+1])
			i++
		case "--type":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: events-export-ics [--case "<caseFileID>"] [--type "<eventType>"] [--output "<filePath>"]`)
			}
			eventType = strings.TrimSpace(strings.ToLower(args[i+1]))
			i++
		case "--output":
			if i+1 >= len(args) {
				return fmt.Errorf(`usage: events-export-ics [--case "<caseFileID>"] [--type "<eventType>"] [--output "<filePath>"]`)
			}
			outputPath = strings.TrimSpace(args[i+1])
			i++
		default:
			return fmt.Errorf(`usage: events-export-ics [--case "<caseFileID>"] [--type "<eventType>"] [--output "<filePath>"]`)
		}
	}

	result, err := a.ExportUpcomingEventsICS.Execute(ctx, documentapp.ExportUpcomingEventsICSInput{
		CaseFileID: caseFileID,
		EventType:  eventType,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, []byte(result.Content), 0644); err != nil {
		return err
	}

	fmt.Fprintln(a.Out, "ICS exported successfully")
	fmt.Fprintf(a.Out, "Events exported: %d\n", result.EventCount)
	fmt.Fprintf(a.Out, "Output: %s\n", outputPath)

	return nil
}
