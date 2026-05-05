"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { DocumentReviewBadge } from "@/app/components/document-review-badge";
import type {
  DocumentReviewStatus,
  DocumentSearchResult,
  DocumentSummary,
} from "@/lib/types";
import {
  deleteDocument,
  reprocessDocument,
  reviewDocument,
  searchDocuments,
} from "@/lib/api";

import {
  displayDocumentType,
  displayLegalArea,
  displayMimeType,
  documentEventsLabel,
  documentExtractionLabel,
} from "@/lib/document-display";

type Props = {
  caseFileId: string;
  documents: DocumentSummary[];
};

type DocumentFilter =
  | "all"
  | "pending"
  | "reviewed"
  | "error"
  | "without_text"
  | "unknown"
  | "without_events"
  | "with_events";

type DocumentSort = "recent" | "name" | "events" | "pending" | "errors";
type HealthTone = "red" | "orange" | "blue" | "green" | "neutral";

type DocumentHealth = {
  label: string;
  tone: HealthTone;
};

type Toast = {
  text: string;
  type: "success" | "error";
};

const VALID_FILTERS: DocumentFilter[] = [
  "all",
  "pending",
  "reviewed",
  "error",
  "without_text",
  "unknown",
  "without_events",
  "with_events",
];

const VALID_SORTS: DocumentSort[] = [
  "recent",
  "name",
  "events",
  "pending",
  "errors",
];

function normalizeFilter(value: string | null): DocumentFilter {
  if (VALID_FILTERS.includes(value as DocumentFilter))
    return value as DocumentFilter;
  return "all";
}

function normalizeSort(value: string | null): DocumentSort {
  if (VALID_SORTS.includes(value as DocumentSort)) return value as DocumentSort;
  return "recent";
}

function documentHealth(summary: DocumentSummary): DocumentHealth {
  if (summary.document.review_status === "error") {
    return { label: "Requiere corrección", tone: "red" };
  }

  if (!summary.has_extracted_text) {
    return { label: "Revisar extracción", tone: "orange" };
  }

  if (summary.document.review_status === "reviewed") {
    return { label: "Revisado", tone: "green" };
  }

  if ((summary.event_count ?? 0) > 0) {
    return { label: "Contiene hitos", tone: "blue" };
  }

  if (summary.document_type === "unknown" && summary.legal_area === "unknown") {
    return { label: "Documento auxiliar", tone: "neutral" };
  }

  return { label: "Sin hitos detectados", tone: "neutral" };
}

function healthClass(tone: HealthTone) {
  switch (tone) {
    case "red":
      return "border-red-900/70 bg-red-950/40 text-red-100";
    case "orange":
      return "border-orange-900/70 bg-orange-950/30 text-orange-100";
    case "blue":
      return "border-blue-900/70 bg-blue-950/30 text-blue-100";
    case "green":
      return "border-emerald-900/70 bg-emerald-950/30 text-emerald-100";
    default:
      return "border-neutral-700 bg-neutral-950 text-neutral-300";
  }
}

function filterButtonClass(isActive: boolean) {
  return isActive
    ? "rounded-full border border-red-900/70 bg-red-950/40 px-3 py-1.5 text-xs font-medium text-red-100"
    : "rounded-full border border-neutral-800 bg-neutral-950 px-3 py-1.5 text-xs font-medium text-neutral-400 transition hover:border-neutral-700 hover:bg-neutral-900 hover:text-neutral-100";
}

function actionButtonClass(
  tone: "neutral" | "blue" | "green" | "red" = "neutral",
) {
  const map = {
    neutral:
      "border-neutral-700 bg-neutral-950 text-neutral-300 hover:bg-neutral-800 hover:text-neutral-100",
    blue: "border-blue-900/70 bg-blue-950/20 text-blue-100 hover:bg-blue-950/40",
    green:
      "border-emerald-900/70 bg-emerald-950/20 text-emerald-100 hover:bg-emerald-950/40",
    red: "border-red-900/70 bg-red-950/20 text-red-100 hover:bg-red-950/40",
  };

  return `inline-flex items-center rounded-lg border px-2.5 py-1.5 text-xs font-medium transition disabled:cursor-not-allowed disabled:opacity-50 ${map[tone]}`;
}

function normalizeText(value: unknown) {
  return String(value ?? "")
    .trim()
    .toLowerCase();
}

function reviewRank(status?: string) {
  switch (status) {
    case "error":
      return 0;
    case "pending_review":
    case "":
    case undefined:
      return 1;
    case "reviewed":
      return 2;
    default:
      return 3;
  }
}

export function DocumentsList({ caseFileId, documents }: Props) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const [busyId, setBusyId] = useState<string | null>(null);
  const [toast, setToast] = useState<Toast | null>(null);
  const [errorDraftId, setErrorDraftId] = useState<string | null>(null);
  const [errorNote, setErrorNote] = useState("");
  const [deleteDraftId, setDeleteDraftId] = useState<string | null>(null);
  const [contentResults, setContentResults] = useState<DocumentSearchResult[]>([]);
  const [contentSearchLoading, setContentSearchLoading] = useState(false);
  const [contentSearchError, setContentSearchError] = useState("");
  const filter = normalizeFilter(searchParams.get("docFilter"));
  const query = searchParams.get("q") ?? "";
  const sort = normalizeSort(searchParams.get("sort"));

  useEffect(() => {
    if (!toast) return;

    const timeout = window.setTimeout(() => {
      setToast(null);
    }, 3000);

    return () => window.clearTimeout(timeout);
  }, [toast]);

  useEffect(() => {
    const normalizedQuery = query.trim();

    if (normalizedQuery.length < 3) {
      return;
    }

    let cancelled = false;

    const timeout = window.setTimeout(async () => {
      setContentSearchLoading(true);
      setContentSearchError("");

      try {
        const results = await searchDocuments(caseFileId, normalizedQuery);

        if (!cancelled) {
          setContentResults(results);
        }
      } catch (error) {
        if (!cancelled) {
          setContentResults([]);
          setContentSearchError(
            error instanceof Error
              ? error.message
              : "No se pudo buscar dentro de los documentos.",
          );
        }
      } finally {
        if (!cancelled) {
          setContentSearchLoading(false);
        }
      }
    }, 300);

    return () => {
      cancelled = true;
      window.clearTimeout(timeout);
    };
  }, [caseFileId, query]);

  function updateUrl(next: {
    filter?: DocumentFilter;
    query?: string;
    sort?: DocumentSort;
  }) {
    const params = new URLSearchParams(searchParams.toString());

    params.set("tab", "documentos");

    const nextFilter = next.filter ?? filter;
    const nextQuery = next.query ?? query;
    const nextSort = next.sort ?? sort;

    if (nextFilter === "all") params.delete("docFilter");
    else params.set("docFilter", nextFilter);

    if (nextQuery.trim()) params.set("q", nextQuery.trim());
    else params.delete("q");

    if (nextSort === "recent") params.delete("sort");
    else params.set("sort", nextSort);

    router.replace(`${pathname}?${params.toString()}`, { scroll: false });
  }

  function clearFilters() {
    const params = new URLSearchParams(searchParams.toString());

    params.set("tab", "documentos");
    params.delete("docFilter");
    params.delete("q");
    params.delete("sort");

    router.replace(`${pathname}?${params.toString()}`, { scroll: false });
  }

  async function runAction(
    documentId: string,
    action: () => Promise<unknown>,
    ok: string,
  ) {
    setBusyId(documentId);
    setToast(null);

    try {
      await action();
      setToast({ text: ok, type: "success" });
      router.refresh();
    } catch (error) {
      setToast({
        text:
          error instanceof Error
            ? error.message
            : "No se pudo completar la acción.",
        type: "error",
      });
    } finally {
      setBusyId(null);
    }
  }

  async function handleReview(
    documentId: string,
    status: DocumentReviewStatus,
    note = "",
  ) {
    await runAction(
      documentId,
      () => reviewDocument(documentId, status, note),
      "Estado documental actualizado.",
    );

    setErrorDraftId(null);
    setErrorNote("");
  }

  async function handleReprocess(documentId: string) {
    await runAction(
      documentId,
      () => reprocessDocument(documentId),
      "Documento reanalizado correctamente.",
    );
  }

  async function handleDelete(documentId: string) {
    await runAction(
      documentId,
      () => deleteDocument(documentId),
      "Documento eliminado correctamente.",
    );

    setDeleteDraftId(null);
  }

  const counts = useMemo(() => {
    return {
      all: documents.length,
      pending: documents.filter(
        (summary) =>
          !summary.document.review_status ||
          summary.document.review_status === "pending_review",
      ).length,
      reviewed: documents.filter(
        (summary) => summary.document.review_status === "reviewed",
      ).length,
      error: documents.filter(
        (summary) => summary.document.review_status === "error",
      ).length,
      without_text: documents.filter((summary) => !summary.has_extracted_text)
        .length,
      unknown: documents.filter(
        (summary) =>
          summary.has_extracted_text &&
          (summary.event_count ?? 0) === 0 &&
          summary.document_type === "unknown" &&
          summary.legal_area === "unknown",
      ).length,
      without_events: documents.filter(
        (summary) =>
          summary.has_extracted_text && (summary.event_count ?? 0) === 0,
      ).length,
      with_events: documents.filter((summary) => (summary.event_count ?? 0) > 0)
        .length,
    };
  }, [documents]);

  const byFilter = (() => {
    switch (filter) {
      case "pending":
        return documents.filter(
          (summary) =>
            !summary.document.review_status ||
            summary.document.review_status === "pending_review",
        );
      case "reviewed":
        return documents.filter(
          (summary) => summary.document.review_status === "reviewed",
        );
      case "error":
        return documents.filter(
          (summary) => summary.document.review_status === "error",
        );
      case "without_text":
        return documents.filter((summary) => !summary.has_extracted_text);
      case "unknown":
        return documents.filter(
          (summary) =>
            summary.has_extracted_text &&
            (summary.event_count ?? 0) === 0 &&
            summary.document_type === "unknown" &&
            summary.legal_area === "unknown",
        );
      case "without_events":
        return documents.filter(
          (summary) =>
            summary.has_extracted_text && (summary.event_count ?? 0) === 0,
        );
      case "with_events":
        return documents.filter((summary) => (summary.event_count ?? 0) > 0);
      default:
        return documents;
    }
  })();

  const normalizedQuery = normalizeText(query);

  const byQuery = !normalizedQuery
    ? byFilter
    : byFilter.filter((summary) => {
        const doc = summary.document;
        const health = documentHealth(summary);

        const haystack = [
          doc.original_name,
          doc.mime_type,
          doc.review_status,
          summary.document_type,
          summary.legal_area,
          displayMimeType(doc.mime_type),
          displayDocumentType(summary.document_type),
          displayLegalArea(summary.legal_area),
          health.label,
          summary.has_extracted_text ? "texto extraído" : "sin texto",
          (summary.event_count ?? 0) > 0 ? "con hitos" : "sin hitos",
        ]
          .map(normalizeText)
          .join(" ");

        return haystack.includes(normalizedQuery);
      });

  const filteredDocuments = [...byQuery].sort((a, b) => {
    switch (sort) {
      case "name":
        return a.document.original_name.localeCompare(
          b.document.original_name,
          "es",
          { sensitivity: "base" },
        );
      case "events":
        return (b.event_count ?? 0) - (a.event_count ?? 0);
      case "pending":
        return (
          reviewRank(a.document.review_status) -
          reviewRank(b.document.review_status)
        );
      case "errors":
        return (
          Number(b.document.review_status === "error") -
          Number(a.document.review_status === "error")
        );
      case "recent":
      default:
        return (
          new Date(b.document.created_at).getTime() -
          new Date(a.document.created_at).getTime()
        );
    }
  });

  const filters: { id: DocumentFilter; label: string; count: number }[] = [
    { id: "all", label: "Todos", count: counts.all },
    { id: "pending", label: "Pendientes", count: counts.pending },
    { id: "reviewed", label: "Revisados", count: counts.reviewed },
    { id: "error", label: "Error", count: counts.error },
    { id: "without_text", label: "Sin texto", count: counts.without_text },
    { id: "unknown", label: "Sin clasificar", count: counts.unknown },
    { id: "without_events", label: "Sin hitos", count: counts.without_events },
    { id: "with_events", label: "Con hitos", count: counts.with_events },
  ];

  return (
    <>
      <div className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap gap-2">
            {filters.map((item) => (
              <button
                key={item.id}
                type="button"
                onClick={() => updateUrl({ filter: item.id })}
                className={filterButtonClass(filter === item.id)}
              >
                {item.label}
                <span className="ml-1 text-neutral-500">{item.count}</span>
              </button>
            ))}
          </div>

          <div className="grid gap-2 lg:grid-cols-[1fr_auto_auto] lg:items-center">
            <input
              value={query}
              onChange={(event) => updateUrl({ query: event.target.value })}
              placeholder="Buscar por nombre, tipo, área, estado..."
              className="w-full rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-900/70"
            />

            <select
              value={sort}
              onChange={(event) =>
                updateUrl({ sort: event.target.value as DocumentSort })
              }
              className="rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm text-neutral-300 outline-none focus:border-red-900/70"
            >
              <option value="recent">Más recientes</option>
              <option value="name">Nombre A-Z</option>
              <option value="events">Más hitos</option>
              <option value="pending">Pendientes primero</option>
              <option value="errors">Errores primero</option>
            </select>

            {filter !== "all" || query.trim() || sort !== "recent" ? (
              <button
                type="button"
                onClick={clearFilters}
                className="whitespace-nowrap rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm font-medium text-neutral-400 transition hover:border-neutral-700 hover:bg-neutral-900 hover:text-neutral-100"
              >
                Limpiar filtros
              </button>
            ) : null}
          </div>

          <p className="text-xs text-neutral-500">
            Mostrando {filteredDocuments.length} de {documents.length}{" "}
            documentos.
          </p>
          {query.trim().length >= 3 ? (
            <div className="rounded-2xl border border-neutral-800 bg-neutral-950/80 p-4">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-sm font-semibold text-neutral-100">
                    Resultados dentro del contenido
                  </p>
                  <p className="mt-1 text-xs text-neutral-500">
                    Coincidencias encontradas en el texto extraído de los documentos.
                  </p>
                </div>

                {contentSearchLoading ? (
                  <span className="text-xs text-neutral-500">Buscando...</span>
                ) : (
                  <span className="rounded-full border border-neutral-700 bg-neutral-950 px-2.5 py-1 text-xs text-neutral-400">
                    {query.trim().length >= 3 ? contentResults.length : 0}
                  </span>
                )}
              </div>

              {contentSearchError ? (
                <p className="mt-3 rounded-xl border border-red-900/60 bg-red-950/20 p-3 text-xs text-red-100">
                  {contentSearchError}
                </p>
              ) : null}

              {!contentSearchLoading &&
              !contentSearchError &&
              contentResults.length === 0 ? (
                <p className="mt-3 text-sm text-neutral-500">
                  No hay coincidencias dentro del texto extraído.
                </p>
              ) : null}

              {contentResults.length > 0 ? (
                <ul className="mt-3 space-y-2">
                  {contentResults.map((result) => (
                    <li
                      key={result.document_id}
                      className="rounded-xl border border-neutral-800 bg-neutral-900/80 p-3"
                    >
                      <Link
                        href={`/case-files/${caseFileId}/documents/${result.document_id}`}
                        className="text-sm font-medium text-neutral-100 underline-offset-4 hover:underline"
                      >
                        {result.original_name}
                      </Link>

                      <p className="mt-2 text-xs leading-5 text-neutral-400">
                        {result.snippet}
                      </p>
                    </li>
                  ))}
                </ul>
              ) : null}
            </div>
          ) : null}
        </div>

        {filteredDocuments.length === 0 ? (
          <p className="mt-4 rounded-xl border border-neutral-800 bg-neutral-900 p-4 text-sm text-neutral-400">
            No hay documentos para este filtro, búsqueda u ordenación.
          </p>
        ) : (
          <ul className="mt-4 divide-y divide-neutral-800 rounded-2xl border border-neutral-800 bg-neutral-950/40">
            {filteredDocuments.map((summary) => {
              const doc = summary.document;
              const health = documentHealth(summary);
              const isBusy = busyId === doc.id;

              return (
                <li
                  key={doc.id}
                  className="grid gap-4 p-4 transition hover:bg-neutral-900/70 xl:grid-cols-[1fr_auto]"
                >
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <Link
                        href={`/case-files/${caseFileId}/documents/${doc.id}`}
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
                    </div>

                    <div className="mt-2 flex flex-wrap gap-2 text-xs text-neutral-400">
                      <span>{displayMimeType(doc.mime_type)}</span>
                      <span>·</span>
                      <span>
                        {documentExtractionLabel(summary.has_extracted_text)}
                      </span>
                      <span>·</span>
                      <span>
                        Tipo: {displayDocumentType(summary.document_type)}
                      </span>
                      <span>·</span>
                      <span>Área: {displayLegalArea(summary.legal_area)}</span>
                      <span>·</span>
                      <span>{documentEventsLabel(summary.event_count)}</span>
                    </div>

                    {doc.review_note ? (
                      <p className="mt-2 line-clamp-1 text-xs text-neutral-500">
                        Nota: {doc.review_note}
                      </p>
                    ) : null}

                    {!summary.has_extracted_text ? (
                      <p className="mt-2 text-xs text-orange-300">
                        Sin texto extraído. Puede ser un PDF escaneado o
                        requiere revisión.
                      </p>
                    ) : null}
                  </div>

                  <div className="flex flex-wrap items-center gap-2 xl:justify-end">
                    <Link
                      href={`/case-files/${caseFileId}/documents/${doc.id}`}
                      className={actionButtonClass("neutral")}
                    >
                      Ver
                    </Link>

                    <button
                      type="button"
                      disabled={isBusy}
                      onClick={() => handleReprocess(doc.id)}
                      className={actionButtonClass("blue")}
                    >
                      {isBusy ? "..." : "Reanalizar"}
                    </button>

                    {doc.review_status !== "reviewed" ? (
                      <button
                        type="button"
                        disabled={isBusy}
                        onClick={() => handleReview(doc.id, "reviewed")}
                        className={actionButtonClass("green")}
                      >
                        {isBusy ? "..." : "Revisado"}
                      </button>
                    ) : (
                      <button
                        type="button"
                        disabled={isBusy}
                        onClick={() => handleReview(doc.id, "pending_review")}
                        className={actionButtonClass("neutral")}
                      >
                        {isBusy ? "..." : "Reabrir"}
                      </button>
                    )}

                    {doc.review_status !== "error" ? (
                      <button
                        type="button"
                        disabled={isBusy}
                        onClick={() => {
                          setErrorDraftId(doc.id);
                          setDeleteDraftId(null);
                          setErrorNote(doc.review_note ?? "");
                        }}
                        className={actionButtonClass("red")}
                      >
                        Error
                      </button>
                    ) : null}

                    <button
                      type="button"
                      disabled={isBusy}
                      onClick={() => {
                        setDeleteDraftId(doc.id);
                        setErrorDraftId(null);
                        setErrorNote("");
                      }}
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
                        onChange={(event) => setErrorNote(event.target.value)}
                        rows={3}
                        className="mt-2 w-full rounded-xl border border-red-900/60 bg-neutral-950 px-3 py-2 text-sm text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-700"
                        placeholder="Ej.: PDF escaneado, documento ilegible, extracción incorrecta..."
                      />

                      <div className="mt-3 flex flex-wrap gap-2">
                        <button
                          type="button"
                          disabled={isBusy || !errorNote.trim()}
                          onClick={() =>
                            handleReview(doc.id, "error", errorNote.trim())
                          }
                          className={actionButtonClass("red")}
                        >
                          {isBusy ? "..." : "Guardar error"}
                        </button>

                        <button
                          type="button"
                          disabled={isBusy}
                          onClick={() => {
                            setErrorDraftId(null);
                            setErrorNote("");
                          }}
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
                        Se eliminará el documento y sus datos asociados del
                        expediente.
                      </p>

                      <div className="mt-3 flex flex-wrap gap-2">
                        <button
                          type="button"
                          disabled={isBusy}
                          onClick={() => handleDelete(doc.id)}
                          className={actionButtonClass("red")}
                        >
                          {isBusy ? "..." : "Confirmar borrado"}
                        </button>

                        <button
                          type="button"
                          disabled={isBusy}
                          onClick={() => setDeleteDraftId(null)}
                          className={actionButtonClass("neutral")}
                        >
                          Cancelar
                        </button>
                      </div>
                    </div>
                  ) : null}
                </li>
              );
            })}
          </ul>
        )}
      </div>

      {toast ? (
        <div className="fixed bottom-6 right-6 z-50">
          <div
            className={`rounded-xl border px-4 py-3 text-sm shadow-lg ${
              toast.type === "success"
                ? "border-emerald-800 bg-emerald-950/90 text-emerald-100"
                : "border-red-800 bg-red-950/90 text-red-100"
            }`}
          >
            {toast.text}
          </div>
        </div>
      ) : null}
    </>
  );
}
