"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { reopenEvent, resolveEvent, reviewEvent } from "@/lib/api";

type Props = {
  eventId: string;
};

type LoadingAction = null | "review" | "resolve" | "reopen";

export function EventActions({ eventId }: Props) {
  const router = useRouter();

  const [loading, setLoading] = useState<LoadingAction>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [resolutionNote, setResolutionNote] = useState("");

  async function handleReview() {
    try {
      setLoading("review");
      await reviewEvent(eventId);
      router.refresh();
    } catch (e) {
      console.error(e);
      alert("No se pudo marcar el evento como revisado.");
    } finally {
      setLoading(null);
    }
  }

  async function handleResolve() {
    try {
      setLoading("resolve");
      await resolveEvent(eventId, resolutionNote.trim());
      setModalOpen(false);
      setResolutionNote("");
      router.refresh();
    } catch (e) {
      console.error(e);
      alert("No se pudo resolver el evento.");
    } finally {
      setLoading(null);
    }
  }

  async function handleReopen() {
    try {
      setLoading("reopen");
      await reopenEvent(eventId);
      router.refresh();
    } catch (e) {
      console.error(e);
      alert("No se pudo reabrir el evento.");
    } finally {
      setLoading(null);
    }
  }

  const disabled = loading !== null;

  return (
    <>
      <div className="flex flex-wrap gap-2">
        <button
          onClick={handleReview}
          disabled={disabled}
          className="rounded-xl border border-purple-700 bg-purple-950/40 px-3 py-1.5 text-xs font-medium text-purple-200 transition hover:bg-purple-900/60 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading === "review" ? "Revisando..." : "Revisar"}
        </button>

        <button
          onClick={() => setModalOpen(true)}
          disabled={disabled}
          className="rounded-xl border border-emerald-700 bg-emerald-950/40 px-3 py-1.5 text-xs font-medium text-emerald-200 transition hover:bg-emerald-900/60 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Resolver
        </button>

        <button
          onClick={handleReopen}
          disabled={disabled}
          className="rounded-xl border border-neutral-700 bg-neutral-800 px-3 py-1.5 text-xs font-medium text-neutral-200 transition hover:bg-neutral-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading === "reopen" ? "Reabriendo..." : "Reabrir"}
        </button>
      </div>

      {modalOpen ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 px-4">
          <div className="w-full max-w-lg rounded-2xl border border-neutral-800 bg-neutral-950 p-6 shadow-2xl">
            <div className="space-y-2">
              <h3 className="text-xl font-semibold">Resolver evento</h3>
              <p className="text-sm text-neutral-400">
                Añade una nota breve dejando constancia de la actuación realizada.
              </p>
            </div>

            <textarea
              value={resolutionNote}
              onChange={(e) => setResolutionNote(e.target.value)}
              placeholder="Ej.: Documentación aportada, escrito presentado, plazo controlado..."
              className="mt-5 min-h-32 w-full resize-none rounded-xl border border-neutral-800 bg-neutral-900 p-3 text-sm text-neutral-100 outline-none ring-0 placeholder:text-neutral-500 focus:border-neutral-600"
            />

            <div className="mt-5 flex justify-end gap-2">
              <button
                onClick={() => {
                  setModalOpen(false);
                  setResolutionNote("");
                }}
                disabled={disabled}
                className="rounded-xl border border-neutral-700 px-4 py-2 text-sm text-neutral-200 transition hover:bg-neutral-800 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Cancelar
              </button>

              <button
                onClick={handleResolve}
                disabled={disabled}
                className="rounded-xl border border-emerald-700 bg-emerald-950/60 px-4 py-2 text-sm font-medium text-emerald-100 transition hover:bg-emerald-900/70 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {loading === "resolve" ? "Resolviendo..." : "Resolver evento"}
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}