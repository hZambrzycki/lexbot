"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { reprocessDocument } from "@/lib/api";

type Props = {
  documentId: string;
};

export function ReprocessDocumentButton({ documentId }: Props) {
  const router = useRouter();

  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  async function handleReprocess() {
    setLoading(true);
    setMessage("");
    setError("");

    try {
      const result = await reprocessDocument(documentId);

      setMessage(
        `Documento reanalizado. Texto: ${result.extracted_length} caracteres. Eventos detectados: ${result.events_detected}.`,
      );

      router.refresh();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Error al reanalizar documento.",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-2">
      <button
        type="button"
        onClick={handleReprocess}
        disabled={loading}
        className="inline-flex rounded-xl border border-blue-900/70 bg-blue-950/30 px-4 py-2 text-sm font-medium text-blue-100 transition hover:bg-blue-950/50 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {loading ? "Reanalizando..." : "Reanalizar documento"}
      </button>

      {message ? <p className="text-sm text-green-300">{message}</p> : null}
      {error ? <p className="text-sm text-red-300">{error}</p> : null}
    </div>
  );
}