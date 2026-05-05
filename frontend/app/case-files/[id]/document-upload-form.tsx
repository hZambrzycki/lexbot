"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { importCaseFileDocument } from "@/lib/api";

type Props = {
  caseFileId: string;
};

const MAX_FILE_SIZE_BYTES = 20 * 1024 * 1024;

const allowedExtensions = [".txt", ".md", ".pdf", ".docx"];

type UploadStatus = "idle" | "uploading" | "success" | "warning" | "error";

function getFileExtension(fileName: string): string {
  const index = fileName.lastIndexOf(".");
  if (index === -1) {
    return "";
  }

  return fileName.slice(index).toLowerCase();
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024 * 1024) {
    return `${Math.max(1, Math.round(bytes / 1024))} KB`;
  }

  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function validateFile(file: File): string | null {
  const extension = getFileExtension(file.name);

  if (!allowedExtensions.includes(extension)) {
    return "Formato no permitido. Sube un documento TXT, MD, PDF o DOCX.";
  }

  if (file.size <= 0) {
    return "El archivo está vacío.";
  }

  if (file.size > MAX_FILE_SIZE_BYTES) {
    return `El archivo es demasiado grande. Máximo ${formatFileSize(
      MAX_FILE_SIZE_BYTES,
    )}.`;
  }

  return null;
}

function buildSuccessMessage(
  result: Awaited<ReturnType<typeof importCaseFileDocument>>,
) {
  if (!result.text_extracted) {
    return "Documento importado, pero no se ha podido extraer texto. Revisa si es un PDF escaneado o si el archivo está vacío.";
  }

  if (result.events_detected === 0) {
    return "Documento importado y texto extraído, pero no se han detectado hitos procesales.";
  }

  return `Documento importado correctamente. Texto extraído: sí. Eventos detectados: ${result.events_detected}.`;
}

export function DocumentUploadForm({ caseFileId }: Props) {
  const router = useRouter();

  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<UploadStatus>("idle");
  const [message, setMessage] = useState("");

  function handleFileChange(event: React.ChangeEvent<HTMLInputElement>) {
    const selectedFile = event.target.files?.[0] ?? null;

    setFile(selectedFile);
    setStatus("idle");
    setMessage("");

    if (!selectedFile) {
      return;
    }

    const validationError = validateFile(selectedFile);
    if (validationError) {
      setStatus("error");
      setMessage(validationError);
      return;
    }

    setMessage(
      `Archivo seleccionado: ${selectedFile.name} (${formatFileSize(
        selectedFile.size,
      )}).`,
    );
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const form = event.currentTarget;

    if (!file) {
      setStatus("error");
      setMessage("Selecciona un documento antes de subirlo.");
      return;
    }

    const validationError = validateFile(file);
    if (validationError) {
      setStatus("error");
      setMessage(validationError);
      return;
    }

    setStatus("uploading");
    setMessage("Subiendo, extrayendo texto y buscando hitos procesales...");

    try {
      const result = await importCaseFileDocument(caseFileId, file);
      const successMessage = buildSuccessMessage(result);

      setStatus(
        !result.text_extracted || result.events_detected === 0
          ? "warning"
          : "success",
      );

      setMessage(successMessage);
      setFile(null);
      form.reset();

      router.refresh();
    } catch (error) {
      setStatus("error");
      setMessage(
        error instanceof Error ? error.message : "Error al subir documento.",
      );
    }
  }

  const messageClassName = {
    idle: "text-neutral-400",
    uploading: "text-neutral-400",
    success: "text-green-300",
    warning: "text-amber-300",
    error: "text-red-300",
  }[status];

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

      <p className="mt-2 text-xs leading-5 text-neutral-500">
        Formatos permitidos: TXT, MD, PDF y DOCX. Tamaño máximo:{" "}
        {formatFileSize(MAX_FILE_SIZE_BYTES)}.
      </p>

      <div className="mt-4 flex flex-col gap-3 md:flex-row md:items-center">
        <input
          type="file"
          name="file"
          accept=".txt,.md,.pdf,.docx,application/pdf,text/plain,text/markdown,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
          onChange={handleFileChange}
          disabled={status === "uploading"}
          className="block w-full rounded-xl border border-neutral-700 bg-black/30 px-3 py-2 text-sm text-neutral-200 file:mr-4 file:rounded-lg file:border-0 file:bg-neutral-800 file:px-3 file:py-2 file:text-sm file:text-neutral-100 hover:file:bg-neutral-700 disabled:cursor-not-allowed disabled:opacity-60"
        />

        <button
          type="submit"
          disabled={status === "uploading" || !file}
          className="inline-flex justify-center rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/50 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {status === "uploading" ? "Analizando..." : "Subir documento"}
        </button>
      </div>

      {message ? (
        <p className={`mt-3 text-sm leading-6 ${messageClassName}`}>
          {message}
        </p>
      ) : null}
    </form>
  );
}