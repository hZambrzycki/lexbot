import Link from "next/link";

import { DocumentReviewBadge } from "@/app/components/document-review-badge";
import {
  displayDocumentType,
  displayLegalArea,
  displayMimeType,
  documentEventsLabel,
  documentExtractionLabel,
} from "@/lib/document-display";
import type {
  DocumentReviewStatus,
  DocumentSearchResult,
  DocumentSummary,
} from "@/lib/types";

import {
  actionButtonClass,
  documentHealth,
  healthClass,
} from "./document-list-utils";

type Props = {
  caseFileId: string;
  summary: DocumentSummary;
  trimmedQuery: string;
  contentMatch?: DocumentSearchResult;
  isBusy: boolean;
  errorDraftId: string | null;
  deleteDraftId: string | null;
  errorNote: string;
  onReprocess: (documentId: string) => void;
  onReview: (
    documentId: string,
    status: DocumentReviewStatus,
    note?: string,
  ) => void;
  onDelete: (documentId: string) => void;
  onOpenErrorDraft: (documentId: string, note: string) => void;
  onCloseErrorDraft: () => void;
  onErrorNoteChange: (note: string) => void;
  onOpenDeleteDraft: (documentId: string) => void;
  onCloseDeleteDraft: () => void;
};

function renderHighlightedSnippet(snippet: string) {
  const parts = snippet.split(/(\[[^\]]+\])/g);

  return parts.map((part, index) => {
    const isHighlighted = part.startsWith("[") && part.endsWith("]");

    if (!isHighlighted) {
      return <span key={`${part}-${index}`}>{part}</span>;
    }

    return (
      <mark
        key={`${part}-${index}`}
        className="rounded-md border border-yellow-700/50 bg-yellow-400/10 px-1 font-medium text-yellow-100"
      >
        {part.slice(1, -1)}
      </mark>
    );
  });
}

export function DocumentListItem({
  caseFileId,
  summary,
  trimmedQuery,
  contentMatch,
  isBusy,
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
  const doc = summary.document;
  const health = documentHealth(summary);
  const detailHref = `/case-files/${caseFileId}/documents/${doc.id}${
    trimmedQuery ? `?q=${encodeURIComponent(trimmedQuery)}` : ""
  }`;

  return (
    <li className="grid gap-4 p-4 transition hover:bg-neutral-900/70 xl:grid-cols-[1fr_auto]">
      <div className="min-w-0">
        <div className="flex flex-wrap items-center gap-2">
          <Link
            href={detailHref}
            className="truncate text-base font-semibold text-neutral-100 underline-offset-4 hover:underline"
          >
            {doc.original_name}
          </Link>

          <span
            className={`rounded-full border px-2.5 py-1 text-xs font-medium ${healthClass(
              health.tone,
            )}`}
          >
            {health.label}
          </span>

          <DocumentReviewBadge status={doc.review_status} />

          {contentMatch ? (
            <span className="rounded-full border border-yellow-900/70 bg-yellow-950/30 px-2.5 py-1 text-xs font-medium text-yellow-100">
              Coincidencia en contenido
            </span>
          ) : null}
        </div>

        <div className="mt-2 flex flex-wrap gap-2 text-xs text-neutral-400">
          <span>{displayMimeType(doc.mime_type)}</span>
          <span>·</span>
          <span>{documentExtractionLabel(summary.has_extracted_text)}</span>
          <span>·</span>
          <span>Tipo: {displayDocumentType(summary.document_type)}</span>
          <span>·</span>
          <span>Área: {displayLegalArea(summary.legal_area)}</span>
          <span>·</span>
          <span>{documentEventsLabel(summary.event_count)}</span>
        </div>

        {contentMatch ? (
          <div className="mt-3 rounded-2xl border border-yellow-900/40 bg-gradient-to-br from-yellow-950/30 to-neutral-950 p-4 text-sm leading-6 text-yellow-50/90 shadow-inner">
            <div className="mb-2 flex items-center gap-2 text-[11px] font-medium uppercase tracking-wide text-yellow-300/80">
              <span>Coincidencia documental</span>

              {contentMatch.score ? (
                <span className="rounded-full border border-yellow-900/40 px-2 py-0.5 text-[10px] normal-case tracking-normal text-yellow-200/70">
                  relevancia {contentMatch.score}
                </span>
              ) : null}
            </div>

            <p>{renderHighlightedSnippet(contentMatch.snippet)}</p>
          </div>
        ) : null}

        {doc.review_note ? (
          <p className="mt-2 line-clamp-1 text-xs text-neutral-500">
            Nota: {doc.review_note}
          </p>
        ) : null}

        {!summary.has_extracted_text ? (
          <p className="mt-2 text-xs text-orange-300">
            Sin texto extraído. Puede ser un PDF escaneado o requiere revisión.
          </p>
        ) : null}
      </div>

      <div className="flex flex-wrap items-center gap-2 xl:justify-end">
        <Link href={detailHref} className={actionButtonClass("neutral")}>
          Ver
        </Link>

        <button
          type="button"
          disabled={isBusy}
          onClick={() => onReprocess(doc.id)}
          className={actionButtonClass("blue")}
        >
          {isBusy ? "..." : "Reanalizar"}
        </button>

        {doc.review_status !== "reviewed" ? (
          <button
            type="button"
            disabled={isBusy}
            onClick={() => onReview(doc.id, "reviewed")}
            className={actionButtonClass("green")}
          >
            {isBusy ? "..." : "Revisado"}
          </button>
        ) : (
          <button
            type="button"
            disabled={isBusy}
            onClick={() => onReview(doc.id, "pending_review")}
            className={actionButtonClass("neutral")}
          >
            {isBusy ? "..." : "Reabrir"}
          </button>
        )}

        {doc.review_status !== "error" ? (
          <button
            type="button"
            disabled={isBusy}
            onClick={() => onOpenErrorDraft(doc.id, doc.review_note ?? "")}
            className={actionButtonClass("red")}
          >
            Error
          </button>
        ) : null}

        <button
          type="button"
          disabled={isBusy}
          onClick={() => onOpenDeleteDraft(doc.id)}
          className={actionButtonClass("red")}
        >
          Borrar
        </button>
      </div>

      {errorDraftId === doc.id ? (
        <div className="rounded-xl border border-red-900/60 bg-red-950/20 p-3 xl:col-span-2">
          <label className="text-xs font-medium text-red-100">
            Motivo del error documental
          </label>

          <textarea
            value={errorNote}
            onChange={(event) => onErrorNoteChange(event.target.value)}
            rows={3}
            className="mt-2 w-full rounded-xl border border-red-900/60 bg-neutral-950 px-3 py-2 text-sm text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-700"
            placeholder="Ej.: PDF escaneado, documento ilegible, extracción incorrecta..."
          />

          <div className="mt-3 flex flex-wrap gap-2">
            <button
              type="button"
              disabled={isBusy || !errorNote.trim()}
              onClick={() => onReview(doc.id, "error", errorNote.trim())}
              className={actionButtonClass("red")}
            >
              {isBusy ? "..." : "Guardar error"}
            </button>

            <button
              type="button"
              disabled={isBusy}
              onClick={onCloseErrorDraft}
              className={actionButtonClass("neutral")}
            >
              Cancelar
            </button>
          </div>
        </div>
      ) : null}

      {deleteDraftId === doc.id ? (
        <div className="rounded-xl border border-red-900/60 bg-red-950/20 p-3 xl:col-span-2">
          <p className="text-sm font-medium text-red-100">
            ¿Seguro que quieres borrar este documento?
          </p>

          <p className="mt-1 text-xs text-red-200/70">
            Se eliminará el documento y sus datos asociados del expediente.
          </p>

          <div className="mt-3 flex flex-wrap gap-2">
            <button
              type="button"
              disabled={isBusy}
              onClick={() => onDelete(doc.id)}
              className={actionButtonClass("red")}
            >
              {isBusy ? "..." : "Confirmar borrado"}
            </button>

            <button
              type="button"
              disabled={isBusy}
              onClick={onCloseDeleteDraft}
              className={actionButtonClass("neutral")}
            >
              Cancelar
            </button>
          </div>
        </div>
      ) : null}
    </li>
  );
}