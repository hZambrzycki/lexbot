"use client";

import { useEffect, useRef, useState } from "react";

import { DocumentToast } from "@/app/components/document-toast";
import type { DocumentReviewStatus, DocumentSummary } from "@/lib/types";

import { DocumentEmptyState } from "./document-empty-state";
import { DocumentResultsList } from "./document-results-list";
import { DocumentSearchControls } from "./document-search-controls";
import { DocumentSearchStatus } from "./document-search-status";
import { useDocumentActions } from "./use-document-actions";
import { useDocumentSearch } from "./use-document-search";

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

  useEffect(() => {
    function handleShortcut(event: KeyboardEvent) {
      const target = event.target as HTMLElement | null;

      const isTyping =
        target?.tagName === "INPUT" ||
        target?.tagName === "TEXTAREA" ||
        target?.isContentEditable;

      if (isTyping) return;

      if (event.key === "/") {
        event.preventDefault();
        searchInputRef.current?.focus();
      }
    }

    window.addEventListener("keydown", handleShortcut);

    return () => {
      window.removeEventListener("keydown", handleShortcut);
    };
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

        <DocumentSearchStatus
          visibleCount={filteredDocuments.length}
          totalCount={documents.length}
          canSearchContent={canSearchContent}
          contentSearchLoading={contentSearchLoading}
          contentResultsCount={contentResults.length}
        />

        {filteredDocuments.length === 0 ? (
          <DocumentEmptyState
            canSearchContent={canSearchContent}
            contentResultsCount={contentResults.length}
          />
        ) : (
          <DocumentResultsList
            caseFileId={caseFileId}
            documents={filteredDocuments}
            trimmedQuery={trimmedQuery}
            contentMatchMap={contentMatchMap}
            busyId={busyId}
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
        )}
      </div>

      <DocumentToast toast={toast} />
    </>
  );
}