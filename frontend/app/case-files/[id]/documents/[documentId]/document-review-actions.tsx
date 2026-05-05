"use client";

import { useState } from "react";
import type { DocumentReviewStatus } from "@/lib/types";

type Props = {
  documentId: string;
  reviewStatus?: DocumentReviewStatus;
  reviewNote?: string;
};

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const baseActionButtonClass =
  "inline-flex items-center justify-center gap-2 rounded-xl border px-3.5 py-2 text-xs font-semibold shadow-sm transition active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50";

const reviewedButtonClass = `${baseActionButtonClass} border-emerald-700/70 bg-emerald-950/40 text-emerald-100 hover:border-emerald-600 hover:bg-emerald-900/40 focus:outline-none focus:ring-2 focus:ring-emerald-700/40`;

const reopenButtonClass = `${baseActionButtonClass} border-sky-800/70 bg-sky-950/30 text-sky-100 hover:border-sky-700 hover:bg-sky-900/30 focus:outline-none focus:ring-2 focus:ring-sky-700/40`;

const dangerButtonClass = `${baseActionButtonClass} border-red-800/70 bg-red-950/35 text-red-100 hover:border-red-700 hover:bg-red-900/40 focus:outline-none focus:ring-2 focus:ring-red-700/40`;

const neutralButtonClass = `${baseActionButtonClass} border-neutral-700 bg-neutral-950 text-neutral-100 hover:bg-neutral-800 focus:outline-none focus:ring-2 focus:ring-neutral-700/40`;

export function DocumentReviewActions({
  documentId,
  reviewStatus = "pending_review",
  reviewNote = "",
}: Props) {
  const [note, setNote] = useState(reviewNote);
  const [loading, setLoading] = useState<DocumentReviewStatus | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [showErrorNote, setShowErrorNote] = useState(false);

  async function updateDocumentReview(nextStatus: DocumentReviewStatus) {
    const cleanNote = note.trim();

    setError(null);

    if (nextStatus === "error" && cleanNote.length === 0) {
      setError("Para marcar como error debes escribir una nota.");
      setShowErrorNote(true);
      return;
    }

    try {
      setLoading(nextStatus);

      const response = await fetch(
        `${API_BASE_URL}/documents/${documentId}/review`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            review_status: nextStatus,
            review_note: nextStatus === "error" ? cleanNote : "",
          }),
        },
      );

      if (!response.ok) {
        let message = `HTTP ${response.status}`;

        try {
          const data = (await response.json()) as { error?: string };
          if (data.error) {
            message = data.error;
          }
        } catch {
          // ignore non-json error body
        }

        throw new Error(message);
      }

      window.location.reload();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "No se pudo actualizar el documento",
      );
    } finally {
      setLoading(null);
    }
  }

  function handleErrorClick() {
    setError(null);
    setShowErrorNote(true);
  }

  function cancelErrorNote() {
    setError(null);
    setShowErrorNote(false);
  }

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap gap-3">
        <button
          type="button"
          disabled={loading !== null || reviewStatus === "reviewed"}
          onClick={() => updateDocumentReview("reviewed")}
          className={reviewedButtonClass}
        >
          <span className="text-sm leading-none">✓</span>
          <span>
            {loading === "reviewed" ? "Guardando..." : "Marcar revisado"}
          </span>
        </button>

        <button
          type="button"
          disabled={loading !== null || reviewStatus === "pending_review"}
          onClick={() => updateDocumentReview("pending_review")}
          className={reopenButtonClass}
        >
          <span className="text-sm leading-none">↺</span>
          <span>
            {loading === "pending_review"
              ? "Guardando..."
              : "Reabrir revisión"}
          </span>
        </button>

        <button
          type="button"
          disabled={loading !== null || reviewStatus === "error"}
          onClick={handleErrorClick}
          className={dangerButtonClass}
        >
          <span className="text-sm leading-none">!</span>
          <span>Marcar error</span>
        </button>
      </div>

      {showErrorNote ? (
        <div className="max-w-2xl rounded-2xl border border-red-900/40 bg-red-950/10 p-4">
          <label className="mb-1 block text-xs font-medium text-red-100">
            Nota del error
          </label>

          <textarea
            value={note}
            onChange={(event) => setNote(event.target.value)}
            placeholder="Ej.: OCR incompleto, falta una página, documento duplicado..."
            className="min-h-20 w-full rounded-xl border border-red-900/50 bg-neutral-950 px-3 py-2 text-xs leading-5 text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-700"
          />

          <div className="mt-3 flex flex-wrap gap-2">
            <button
              type="button"
              disabled={loading !== null}
              onClick={() => updateDocumentReview("error")}
              className={dangerButtonClass}
            >
              <span className="text-sm leading-none">!</span>
              <span>
                {loading === "error" ? "Guardando..." : "Confirmar error"}
              </span>
            </button>

            <button
              type="button"
              disabled={loading !== null}
              onClick={cancelErrorNote}
              className={neutralButtonClass}
            >
              <span className="text-sm leading-none">×</span>
              <span>Cancelar</span>
            </button>
          </div>
        </div>
      ) : null}


      {error ? (
        <p className="w-fit rounded-xl border border-red-900 bg-red-950/30 px-3 py-2 text-xs text-red-100">
          {error}
        </p>
      ) : null}
    </div>
  );
}