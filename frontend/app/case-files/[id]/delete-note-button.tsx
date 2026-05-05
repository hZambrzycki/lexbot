"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { deleteNote } from "@/lib/api";

type Props = {
  caseFileId: string;
  noteId: string;
};

export function DeleteNoteButton({ caseFileId, noteId }: Props) {
  const router = useRouter();
  const [isDeleting, setIsDeleting] = useState(false);

  async function handleDelete() {
    if (!confirm("¿Eliminar esta nota?")) return;

    setIsDeleting(true);

    try {
      await deleteNote(caseFileId, noteId);
      router.refresh();
    } finally {
      setIsDeleting(false);
    }
  }

  return (
    <button
      type="button"
      onClick={handleDelete}
      disabled={isDeleting}
      className="rounded-xl border border-red-900/60 bg-red-950/20 px-3 py-1.5 text-xs font-medium text-red-200 transition hover:bg-red-950/40 disabled:cursor-not-allowed disabled:opacity-60"
    >
      {isDeleting ? "Eliminando..." : "Eliminar"}
    </button>
  );
}