import {
  CaseFile,
  CaseFileDetail,
  Dashboard,
  EventActionResponse,
  EventItem,
} from "@/lib/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

if (!API_BASE_URL) {
  throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
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
      if (data?.error) message = data.error;
    } catch {}
    throw new Error(message);
  }

  return response.json() as Promise<T>;
}

export async function getCaseFiles(): Promise<CaseFile[]> {
  return apiFetch<CaseFile[]>("/case-files");
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