"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { importCaseFileDocument } from "@/lib/api";

type Props = {
  caseFileId: string;
};

export function DocumentUploadForm({ caseFileId }: Props) {
  const router = useRouter();

  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<
    "idle" | "uploading" | "success" | "error"
  >("idle");
  const [message, setMessage] = useState("");

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!file) {
      setStatus("error");
      setMessage("Selecciona un documento antes de subirlo.");
      return;
    }

    setStatus("uploading");
    setMessage("Subiendo y analizando documento...");

    try {
      const result = await importCaseFileDocument(caseFileId, file);

      setStatus("success");
      setMessage(
        `Documento importado. Texto extraído: ${
          result.text_extracted ? "sí" : "no"
        }. Eventos detectados: ${result.events_detected}.`,
      );

      setFile(null);
      router.refresh();
    } catch (error) {
      setStatus("error");
      setMessage(
        error instanceof Error ? error.message : "Error al subir documento.",
      );
    }
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5"
    >
      <h3 className="text-lg font-semibold text-neutral-50">
        Subir documento
      </h3>

      <p className="mt-2 text-sm leading-6 text-neutral-400">
        Importa un documento al expediente. LEXBOX intentará extraer el texto y
        detectar hitos procesales automáticamente.
      </p>

      <div className="mt-4 flex flex-col gap-3 md:flex-row md:items-center">
        <input
          type="file"
          name="file"
          accept=".txt,.pdf,.docx,application/pdf,text/plain,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
          onChange={(event) => setFile(event.target.files?.[0] ?? null)}
          className="block w-full rounded-xl border border-neutral-700 bg-black/30 px-3 py-2 text-sm text-neutral-200 file:mr-4 file:rounded-lg file:border-0 file:bg-neutral-800 file:px-3 file:py-2 file:text-sm file:text-neutral-100 hover:file:bg-neutral-700"
        />

        <button
          type="submit"
          disabled={status === "uploading"}
          className="inline-flex justify-center rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/50 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {status === "uploading" ? "Analizando..." : "Subir documento"}
        </button>
      </div>

      {message ? (
        <p
          className={`mt-3 text-sm ${
            status === "error"
              ? "text-red-300"
              : status === "success"
                ? "text-green-300"
                : "text-neutral-400"
          }`}
        >
          {message}
        </p>
      ) : null}
    </form>
  );
}