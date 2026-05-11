import type {
  DocumentReviewStatus,
  DocumentSearchResult,
  DocumentSummary,
} from "@/lib/types";
import { DocumentListItem } from "./document-list-item";

type Props = {
  caseFileId: string;
  documents: DocumentSummary[];
  trimmedQuery: string;
  contentMatchMap: Map<string, DocumentSearchResult>;
  busyId: string | null;
  errorDraftId: string | null;
  deleteDraftId: string | null;
  errorNote: string;
  onReprocess: (documentId: string) => Promise<void>;
  onReview: (
    documentId: string,
    status: DocumentReviewStatus,
    note?: string,
  ) => Promise<void>;
  onDelete: (documentId: string) => Promise<void>;
  onOpenErrorDraft: (documentId: string, note: string) => void;
  onCloseErrorDraft: () => void;
  onErrorNoteChange: (note: string) => void;
  onOpenDeleteDraft: (documentId: string) => void;
  onCloseDeleteDraft: () => void;
};

export function DocumentResultsList({
  caseFileId,
  documents,
  trimmedQuery,
  contentMatchMap,
  busyId,
  errorDraftId,
  deleteDraftId,
  errorNote,
  onReprocess,
  onReview,
  onDelete,
  onOpenErrorDraft,
  onCloseErrorDraft,
  onErrorNoteChange,
  onOpenDeleteDraft,
  onCloseDeleteDraft,
}: Props) {
  return (
    <ul className="mt-4 divide-y divide-neutral-800 rounded-2xl border border-neutral-800 bg-neutral-950/40">
      {documents.map((summary) => (
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
          onReprocess={onReprocess}
          onReview={onReview}
          onDelete={onDelete}
          onOpenErrorDraft={onOpenErrorDraft}
          onCloseErrorDraft={onCloseErrorDraft}
          onErrorNoteChange={onErrorNoteChange}
          onOpenDeleteDraft={onOpenDeleteDraft}
          onCloseDeleteDraft={onCloseDeleteDraft}
        />
      ))}
    </ul>
  );
}