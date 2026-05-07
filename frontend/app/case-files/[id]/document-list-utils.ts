import type { DocumentSummary } from "@/lib/types";
import { normalizeSearchText } from "@/lib/search-normalization";
export type DocumentFilter =
  | "all"
  | "pending"
  | "reviewed"
  | "error"
  | "without_text"
  | "unknown"
  | "without_events"
  | "with_events";

export type DocumentSort = "recent" | "name" | "events" | "pending" | "errors";
export type HealthTone = "red" | "orange" | "blue" | "green" | "neutral";

export type DocumentHealth = {
  label: string;
  tone: HealthTone;
};

export const VALID_FILTERS: DocumentFilter[] = [
  "all",
  "pending",
  "reviewed",
  "error",
  "without_text",
  "unknown",
  "without_events",
  "with_events",
];

export const VALID_SORTS: DocumentSort[] = [
  "recent",
  "name",
  "events",
  "pending",
  "errors",
];

export function normalizeFilter(value: string | null): DocumentFilter {
  if (VALID_FILTERS.includes(value as DocumentFilter)) {
    return value as DocumentFilter;
  }

  return "all";
}

export function normalizeSort(value: string | null): DocumentSort {
  if (VALID_SORTS.includes(value as DocumentSort)) {
    return value as DocumentSort;
  }

  return "recent";
}

export function documentHealth(summary: DocumentSummary): DocumentHealth {
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

export function healthClass(tone: HealthTone) {
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

export function filterButtonClass(isActive: boolean) {
  return isActive
    ? "rounded-full border border-red-900/70 bg-red-950/40 px-3 py-1.5 text-xs font-medium text-red-100"
    : "rounded-full border border-neutral-800 bg-neutral-950 px-3 py-1.5 text-xs font-medium text-neutral-400 transition hover:border-neutral-700 hover:bg-neutral-900 hover:text-neutral-100";
}

export function actionButtonClass(
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

export function normalizeText(value: unknown) {
  return normalizeSearchText(value);
}

export function reviewRank(status?: string) {
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