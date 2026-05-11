import type {
  CaseFile,
  CaseFileDetail,
  CreateCaseFileInput,
  Dashboard,
  DocumentDetail,
  DocumentItem,
  DocumentReviewResponse,
  DocumentReviewStatus,
  EventActionResponse,
  EventItem,
  Note,
  GlobalSearchResult,
  DocumentSearchResult,
} from "@/lib/types";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

function buildApiUrl(path: string): string {
  return `${API_BASE_URL}${path}`;
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(buildApiUrl(path), {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
    cache: "no-store",
  });

  if (!response.ok) {
    let message = `HTTP ${response.status}`;

    try {
      const data = (await response.json()) as { error?: string };
      if (data?.error) {
        message = data.error;
      }
    } catch {
      // ignore non-json error body
    }

    throw new Error(message);
  }

  return response.json() as Promise<T>;
}

export function getGlobalUpcomingICSUrl(params?: { type?: string }): string {
  const searchParams = new URLSearchParams();

  if (params?.type) {
    searchParams.set("type", params.type);
  }

  const query = searchParams.toString();

  return buildApiUrl(`/events/upcoming.ics${query ? `?${query}` : ""}`);
}

export async function deleteNote(caseFileId: string, noteId: string) {
  return apiFetch<{ id: string; deleted: boolean }>(
    `/case-files/${caseFileId}/notes/${noteId}`,
    {
      method: "DELETE",
    },
  );
}
export async function createNote(caseFileId: string, input: {
  title: string;
  content: string;
}) {
  return apiFetch<Note>(`/case-files/${caseFileId}/notes`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function getCaseFileUpcomingICSUrl(
  id: string,
  params?: { type?: string },
): string {
  const searchParams = new URLSearchParams();

  if (params?.type) {
    searchParams.set("type", params.type);
  }

  const query = searchParams.toString();

  return buildApiUrl(
    `/case-files/${id}/events/upcoming.ics${query ? `?${query}` : ""}`,
  );
}

export async function getCaseFiles(): Promise<CaseFile[]> {
  return apiFetch<CaseFile[]>("/case-files");
}

export async function createCaseFile(
  input: CreateCaseFileInput,
): Promise<CaseFile> {
  return apiFetch<CaseFile>("/case-files", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function getCaseFile(id: string): Promise<CaseFileDetail> {
  return apiFetch<CaseFileDetail>(`/case-files/${id}`);
}

export async function getCaseFileDashboard(id: string): Promise<Dashboard> {
  return apiFetch<Dashboard>(`/case-files/${id}/dashboard`);
}

export async function getCaseFileEvents(id: string): Promise<EventItem[]> {
  return apiFetch<EventItem[]>(`/case-files/${id}/events`);
}

export async function getGlobalUpcomingEvents(params?: {
  type?: string;
  reviewStatus?: string;
  relativeOnly?: boolean;
}): Promise<EventItem[]> {
  const searchParams = new URLSearchParams();

  if (params?.type) {
    searchParams.set("type", params.type);
  }

  if (params?.reviewStatus) {
    searchParams.set("review_status", params.reviewStatus);
  }

  if (params?.relativeOnly) {
    searchParams.set("relative_only", "true");
  }

  const query = searchParams.toString();

  return apiFetch<EventItem[]>(`/events/upcoming${query ? `?${query}` : ""}`);
}

export async function reviewEvent(id: string): Promise<EventActionResponse> {
  return apiFetch<EventActionResponse>(`/events/${id}/review`, {
    method: "POST",
  });
}

export async function resolveEvent(
  id: string,
  resolutionNote: string,
): Promise<EventActionResponse> {
  return apiFetch<EventActionResponse>(`/events/${id}/resolve`, {
    method: "POST",
    body: JSON.stringify({ resolution_note: resolutionNote }),
  });
}

export async function reopenEvent(id: string): Promise<EventActionResponse> {
  return apiFetch<EventActionResponse>(`/events/${id}/reopen`, {
    method: "POST",
  });
}

export async function getEvent(id: string): Promise<EventItem> {
  return apiFetch<EventItem>(`/events/${id}`);
}

export async function getUpcomingEvents(params?: {
  type?: string;
  reviewStatus?: string;
  relativeOnly?: boolean;
}): Promise<EventItem[]> {
  return getGlobalUpcomingEvents(params);
}

export type ImportDocumentResponse = {
  document: DocumentItem;
  text_extracted: boolean;
  metadata_analyzed: boolean;
  events_analyzed: boolean;
  events_detected: number;
};

export async function importCaseFileDocument(
  caseFileId: string,
  file: File,
): Promise<ImportDocumentResponse> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(
    buildApiUrl(`/case-files/${caseFileId}/documents`),
    {
      method: "POST",
      body: formData,
      cache: "no-store",
    },
  );

  if (!response.ok) {
    let message = `HTTP ${response.status}`;

    try {
      const data = (await response.json()) as { error?: string };
      if (data?.error) {
        message = data.error;
      }
    } catch {
      // ignore non-json error body
    }

    throw new Error(message);
  }

  return response.json() as Promise<ImportDocumentResponse>;
}

export async function getDocument(id: string): Promise<DocumentDetail> {
  return apiFetch<DocumentDetail>(`/documents/${id}`);
}

export type DeleteDocumentResponse = {
  document_id: string;
  case_file_id: string;
  storage_path: string;
  deleted: boolean;
};

export async function deleteDocument(
  id: string,
): Promise<DeleteDocumentResponse> {
  return apiFetch<DeleteDocumentResponse>(`/documents/${id}`, {
    method: "DELETE",
  });
}

export type ReprocessDocumentResponse = {
  document_id: string;
  text_reindexed: boolean;
  extracted_length: number;
  metadata_analyzed: boolean;
  events_analyzed: boolean;
  events_detected: number;
};

export async function reprocessDocument(
  id: string,
): Promise<ReprocessDocumentResponse> {
  return apiFetch<ReprocessDocumentResponse>(`/documents/${id}/reprocess`, {
    method: "POST",
  });
}

export async function reviewDocument(
  id: string,
  reviewStatus: DocumentReviewStatus,
  reviewNote = "",
): Promise<DocumentReviewResponse> {
  return apiFetch<DocumentReviewResponse>(`/documents/${id}/review`, {
    method: "POST",
    body: JSON.stringify({
      review_status: reviewStatus,
      review_note: reviewNote,
    }),
  });
}

export async function searchDocuments(
  caseFileId: string,
  query: string,
  limit = 20,
): Promise<DocumentSearchResult[]> {
  const params = new URLSearchParams();

  if (caseFileId.trim().length > 0) {
    params.set("case_file_id", caseFileId);
  }

  params.set("q", query);
  params.set("limit", String(limit));

  return apiFetch<DocumentSearchResult[]>(
    `/documents/search?${params.toString()}`,
  );
}

const GLOBAL_NAVIGATION_RESULTS: GlobalSearchResult[] = [
  {
    type: "navigation",
    id: "home",
    title: "Ir al inicio",
    subtitle: "Navegación",
    href: "/",
    score: 10000,
  },
  {
    type: "navigation",
    id: "agenda",
    title: "Abrir agenda",
    subtitle: "Navegación",
    href: "/agenda",
    score: 10000,
  },
  {
    type: "navigation",
    id: "case-files",
    title: "Abrir expedientes",
    subtitle: "Navegación",
    href: "/case-files",
    score: 10000,
  },
];

const GLOBAL_ACTION_RESULTS: GlobalSearchResult[] = [
  {
    type: "action",
    id: "create-case-file",
    title: "Crear expediente",
    subtitle: "Acción rápida",
    href: "/case-files/new",
    score: 12000,
  },
  {
    type: "action",
    id: "open-agenda",
    title: "Ver agenda procesal",
    subtitle: "Acción rápida",
    href: "/agenda",
    score: 12000,
  },
  {
    type: "action",
    id: "pending-documents",
    title: "Ver documentos pendientes de revisión",
    subtitle: "Acción rápida",
    href: "/case-files",
    score: 11500,
  },
];

function normalizeGlobalSearchText(value: string): string {
  return value
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .trim();
}

function getNavigationResults(query: string): GlobalSearchResult[] {
  const normalizedQuery = normalizeGlobalSearchText(query);

  if (normalizedQuery.length < 2) {
    return [];
  }

  return GLOBAL_NAVIGATION_RESULTS.filter((item) => {
    const haystack = normalizeGlobalSearchText(
      `${item.title} ${item.subtitle} ${item.id}`,
    );

    return haystack.includes(normalizedQuery);
  });
}

export async function globalSearch(
  query: string,
): Promise<GlobalSearchResult[]> {
  const actionResults = getActionResults(query);
  const navigationResults = getNavigationResults(query);

  const backendResults = await apiFetch<GlobalSearchResult[]>(
    `/search/global?q=${encodeURIComponent(query)}&limit=8`,
  );

  return [
    ...actionResults,
    ...navigationResults,
    ...backendResults,
  ].sort((a, b) => b.score - a.score);
}

function getActionResults(query: string): GlobalSearchResult[] {
  const normalizedQuery = normalizeGlobalSearchText(query);

  if (normalizedQuery.length < 2) {
    return [];
  }

  return GLOBAL_ACTION_RESULTS.filter((item) => {
    const haystack = normalizeGlobalSearchText(
      `${item.title} ${item.subtitle} ${item.id}`,
    );

    return haystack.includes(normalizedQuery);
  });
}