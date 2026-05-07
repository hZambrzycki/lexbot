"use client";

import { useEffect, useMemo } from "react";
import {
  usePathname,
  useRouter,
  useSearchParams,
} from "next/navigation";
import type {
  DocumentSearchResult,
  DocumentSummary,
} from "@/lib/types";
import { searchDocuments } from "@/lib/api";
import {
  displayDocumentType,
  displayLegalArea,
  displayMimeType,
} from "@/lib/document-display";
import {
  type DocumentFilter,
  type DocumentSort,
  documentHealth,
  normalizeFilter,
  normalizeSort,
  normalizeText,
  reviewRank,
} from "./document-list-utils";
import { useState } from "react";

type UseDocumentSearchArgs = {
  caseFileId: string;
  documents: DocumentSummary[];
};

type UpdateUrlArgs = {
  filter?: DocumentFilter;
  query?: string;
  sort?: DocumentSort;
};

export function useDocumentSearch({
  caseFileId,
  documents,
}: UseDocumentSearchArgs) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const [contentResults, setContentResults] = useState<DocumentSearchResult[]>(
    [],
  );
  const [contentSearchLoading, setContentSearchLoading] = useState(false);

  const filter = normalizeFilter(searchParams.get("docFilter"));
  const query = searchParams.get("q") ?? "";
  const sort = normalizeSort(searchParams.get("sort"));

  const trimmedQuery = query.trim();
  const normalizedQuery = normalizeText(query);
  const canSearchContent = trimmedQuery.length >= 3;

  useEffect(() => {
    const normalizedContentQuery = query.trim();

    if (normalizedContentQuery.length < 3) {
      const timeout = window.setTimeout(() => {
        setContentResults([]);
        setContentSearchLoading(false);
      }, 0);

      return () => window.clearTimeout(timeout);
    }

    let cancelled = false;

    const timeout = window.setTimeout(async () => {
      setContentSearchLoading(true);
      setContentResults([]);

      try {
        const results = await searchDocuments(caseFileId, normalizedContentQuery);

        if (!cancelled) {
          setContentResults(results);
        }
      } catch {
        if (!cancelled) {
          setContentResults([]);
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

  function updateUrl(next: UpdateUrlArgs) {
    const params = new URLSearchParams(searchParams.toString());

    params.set("tab", "documentos");

    const nextFilter = next.filter ?? filter;
    const nextQuery = next.query ?? query;
    const nextSort = next.sort ?? sort;

    if (nextFilter === "all") params.delete("docFilter");
    else params.set("docFilter", nextFilter);

    if (nextQuery.length > 0) params.set("q", nextQuery);
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

  const byFilter = useMemo(() => {
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
  }, [documents, filter]);

  const contentMatchMap = useMemo(() => {
    const map = new Map<string, DocumentSearchResult>();

    for (const result of contentResults) {
      map.set(result.document_id, result);
    }

    return map;
  }, [contentResults]);

  const byQuery = useMemo(() => {
    if (!normalizedQuery) return byFilter;

    return byFilter.filter((summary) => {
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

      const matchesMetadata = haystack.includes(normalizedQuery);
      const matchesContent = contentMatchMap.has(doc.id);

      return matchesMetadata || matchesContent;
    });
  }, [byFilter, contentMatchMap, normalizedQuery]);

  const filteredDocuments = useMemo(() => {
    return [...byQuery].sort((a, b) => {
      const aContentMatch = contentMatchMap.has(a.document.id);
      const bContentMatch = contentMatchMap.has(b.document.id);

      if (aContentMatch !== bContentMatch) {
        return Number(bContentMatch) - Number(aContentMatch);
      }

      if (aContentMatch && bContentMatch) {
        const aScore = contentMatchMap.get(a.document.id)?.score ?? 0;
        const bScore = contentMatchMap.get(b.document.id)?.score ?? 0;

        if (aScore !== bScore) {
          return bScore - aScore;
        }
      }

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
  }, [byQuery, contentMatchMap, sort]);

  return {
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
  };
}