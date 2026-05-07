"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import type { DocumentReviewStatus } from "@/lib/types";
import {
  deleteDocument,
  reprocessDocument,
  reviewDocument,
} from "@/lib/api";

type Toast = {
  text: string;
  type: "success" | "error";
};

export function useDocumentActions() {
  const router = useRouter();

  const [busyId, setBusyId] = useState<string | null>(null);
  const [toast, setToast] = useState<Toast | null>(null);

  useEffect(() => {
    if (!toast) return;

    const timeout = window.setTimeout(() => {
      setToast(null);
    }, 3000);

    return () => window.clearTimeout(timeout);
  }, [toast]);

  async function runAction(
    documentId: string,
    action: () => Promise<unknown>,
    ok: string,
  ) {
    setBusyId(documentId);
    setToast(null);

    try {
      await action();
      setToast({ text: ok, type: "success" });
      router.refresh();
    } catch (error) {
      setToast({
        text:
          error instanceof Error
            ? error.message
            : "No se pudo completar la acción.",
        type: "error",
      });
    } finally {
      setBusyId(null);
    }
  }

  async function handleReview(
    documentId: string,
    status: DocumentReviewStatus,
    note = "",
  ) {
    await runAction(
      documentId,
      () => reviewDocument(documentId, status, note),
      "Estado documental actualizado.",
    );
  }

  async function handleReprocess(documentId: string) {
    await runAction(
      documentId,
      () => reprocessDocument(documentId),
      "Documento reanalizado correctamente.",
    );
  }

  async function handleDelete(documentId: string) {
    await runAction(
      documentId,
      () => deleteDocument(documentId),
      "Documento eliminado correctamente.",
    );
  }

  return {
    busyId,
    toast,
    handleReview,
    handleReprocess,
    handleDelete,
  };
}