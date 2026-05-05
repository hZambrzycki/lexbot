"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { createNote } from "@/lib/api";

type Props = {
  caseFileId: string;
};

export function NoteForm({ caseFileId }: Props) {
  const router = useRouter();
  const [status, setStatus] = useState<"idle" | "saving" | "success" | "error">(
    "idle",
  );
  const [message, setMessage] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const form = event.currentTarget;
    const formData = new FormData(form);

    const title = String(formData.get("title") ?? "").trim();
    const content = String(formData.get("content") ?? "").trim();

    if (!title || !content) {
      setStatus("error");
      setMessage("Indica título y contenido.");
      return;
    }

    setStatus("saving");
    setMessage("");

    try {
      await createNote(caseFileId, { title, content });
      form.reset();
      setStatus("success");
      setMessage("Nota guardada correctamente.");
      router.refresh();
    } catch (error) {
      setStatus("error");
      setMessage(
        error instanceof Error ? error.message : "No se ha podido guardar la nota.",
      );
    }
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5"
    >
      <h3 className="text-lg font-semibold text-neutral-50">Añadir nota</h3>

      <div className="mt-4 grid gap-3">
        <input
          name="title"
          placeholder="Título de la nota"
          className="rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-900/70"
        />

        <textarea
          name="content"
          placeholder="Contenido de la nota..."
          rows={4}
          className="rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-900/70"
        />

        <button
          type="submit"
          disabled={status === "saving"}
          className="w-fit rounded-2xl border border-red-900/70 bg-red-950/40 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/60 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {status === "saving" ? "Guardando..." : "Guardar nota"}
        </button>

        {message ? (
          <p
            className={
              status === "error"
                ? "text-sm text-red-300"
                : "text-sm text-emerald-300"
            }
          >
            {message}
          </p>
        ) : null}
      </div>
    </form>
  );
}