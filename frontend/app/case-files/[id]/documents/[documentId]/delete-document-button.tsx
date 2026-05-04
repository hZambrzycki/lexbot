"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { deleteDocument } from "@/lib/api";

type Props = {
  documentId: string;
  caseFileId: string;
};

export function DeleteDocumentButton({ documentId, caseFileId }: Props) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleDelete() {
    const confirmed = window.confirm(
      "¿Seguro que quieres borrar este documento? Se eliminarán también el texto extraído y los eventos detectados.",
    );

    if (!confirmed) return;

    setLoading(true);
    setError("");

    try {
      await deleteDocument(documentId);
      router.push(`/case-files/${caseFileId}`);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error al borrar documento.");
      setLoading(false);
    }
  }

  return (
    <div className="space-y-2">
      <button
        type="button"
        onClick={handleDelete}
        disabled={loading}
        className="inline-flex rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/50 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {loading ? "Borrando..." : "Borrar documento"}
      </button>

      {error ? <p className="text-sm text-red-300">{error}</p> : null}
    </div>
  );
}