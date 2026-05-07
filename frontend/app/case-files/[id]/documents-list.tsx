"use client";

import { useEffect, useRef, useState } from "react";
import type { DocumentReviewStatus, DocumentSummary } from "@/lib/types";
import { DocumentListItem } from "./document-list-item";
import { DocumentSearchControls } from "./document-search-controls";
import { useDocumentActions } from "./use-document-actions";
import { useDocumentSearch } from "./use-document-search";
import { DocumentToast } from "@/app/components/document-toast";


type Props = {
  caseFileId: string;
  documents: DocumentSummary[];
};

export function DocumentsList({ caseFileId, documents }: Props) {
  const searchInputRef = useRef<HTMLInputElement | null>(null);

  const [errorDraftId, setErrorDraftId] = useState<string | null>(null);
  const [errorNote, setErrorNote] = useState("");
  const [deleteDraftId, setDeleteDraftId] = useState<string | null>(null);

  const {
    filter,
    query,
    sort,
    trimmedQuery,
    canSearchContent,
    contentResults,
    contentSearchLoading,
    contentMatchMap,
    filters,
    filteredDocuments,
    updateUrl,
    clearFilters,
  } = useDocumentSearch({ caseFileId, documents });

  const {
    busyId,
    toast,
    handleReview: reviewDocumentItem,
    handleReprocess,
    handleDelete: deleteDocumentItem,
  } = useDocumentActions();

  useEffect(() => {
    searchInputRef.current?.focus();
  }, []);

  async function handleReview(
    documentId: string,
    status: DocumentReviewStatus,
    note = "",
  ) {
    await reviewDocumentItem(documentId, status, note);

    setErrorDraftId(null);
    setErrorNote("");
  }

  async function handleDelete(documentId: string) {
    await deleteDocumentItem(documentId);

    setDeleteDraftId(null);
  }

  return (
    <>
      <div className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <DocumentSearchControls
          filters={filters}
          activeFilter={filter}
          query={query}
          sort={sort}
          inputRef={searchInputRef}
          onFilterChange={(nextFilter) => updateUrl({ filter: nextFilter })}
          onQueryChange={(nextQuery) => updateUrl({ query: nextQuery })}
          onSortChange={(nextSort) => updateUrl({ sort: nextSort })}
          onClear={clearFilters}
        />

        <p className="mt-4 text-xs text-neutral-500">
          Mostrando {filteredDocuments.length} de {documents.length} documentos
          por filtros visibles.
        </p>

        {canSearchContent && contentSearchLoading ? (
          <p className="mt-2 text-xs text-neutral-500">
            Buscando dentro del texto extraído...
          </p>
        ) : null}

        {canSearchContent && contentResults.length > 0 ? (
          <p className="mt-2 text-xs text-yellow-200">
            {contentResults.length === 1
              ? "1 documento contiene coincidencias dentro del texto extraído."
              : `${contentResults.length} documentos contienen coincidencias dentro del texto extraído.`}
          </p>
        ) : null}

        {filteredDocuments.length === 0 ? (
          <p className="mt-4 rounded-xl border border-neutral-800 bg-neutral-900 p-4 text-sm text-neutral-400">
            {canSearchContent && contentResults.length > 0
              ? "No hay coincidencias por nombre, tipo o estado, pero sí hay resultados dentro del contenido."
              : "No hay documentos para este filtro, búsqueda u ordenación."}
          </p>
        ) : (
          <ul className="mt-4 divide-y divide-neutral-800 rounded-2xl border border-neutral-800 bg-neutral-950/40">
            {filteredDocuments.map((summary) => (
              <DocumentListItem
                key={summary.document.id}
                caseFileId={caseFileId}
                summary={summary}
                trimmedQuery={trimmedQuery}
                contentMatch={contentMatchMap.get(summary.document.id)}
                isBusy={busyId === summary.document.id}
                errorDraftId={errorDraftId}
                deleteDraftId={deleteDraftId}
                errorNote={errorNote}
                onReprocess={handleReprocess}
                onReview={handleReview}
                onDelete={handleDelete}
                onOpenErrorDraft={(documentId, note) => {
                  setErrorDraftId(documentId);
                  setDeleteDraftId(null);
                  setErrorNote(note);
                }}
                onCloseErrorDraft={() => {
                  setErrorDraftId(null);
                  setErrorNote("");
                }}
                onErrorNoteChange={setErrorNote}
                onOpenDeleteDraft={(documentId) => {
                  setDeleteDraftId(documentId);
                  setErrorDraftId(null);
                  setErrorNote("");
                }}
                onCloseDeleteDraft={() => setDeleteDraftId(null)}
              />
            ))}
          </ul>
        )}
      </div>

      <DocumentToast toast={toast} />
    </>
  );
}