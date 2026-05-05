import type {
  CaseFile,
  CaseFileDetail,
  CreateCaseFileInput,
  Dashboard,
  DocumentDetail,
  DocumentReviewResponse,
  DocumentReviewStatus,
  EventActionResponse,
  EventItem,
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
  document: {
    id: string;
    case_file_id: string;
    original_name: string;
    storage_path: string;
    mime_type: string;
    file_hash: string;
    created_at: string;
    updated_at: string;
  };
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