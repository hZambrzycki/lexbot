"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { reopenEvent, resolveEvent, reviewEvent } from "@/lib/api";

type ReviewStatus = "pending" | "reviewed" | "resolved" | string;

type Props = {
  eventIds: string[];
  reviewStatus?: ReviewStatus;
  isGroup?: boolean;
};

type ActionState =
  | "idle"
  | "reviewing"
  | "resolving"
  | "quick-resolving"
  | "reopening";

const QUICK_NOTES = {
  noAction:
    "Revisado el expediente. El hito no requiere actuación procesal adicional en este momento.",
  detectionError:
    "Error de detección — el hito detectado no corresponde a una actuación procesal pendiente.",
};

const RESOLUTION_TEMPLATES = [
  {
    label: "Presentado por LexNET",
    note: "Presentado escrito por LexNET. Se deja constancia de que el hito queda atendido.",
  },
  {
    label: "Ya estaba atendido",
    note: "Revisado el expediente. La actuación correspondiente ya consta realizada.",
  },
  {
    label: "No requiere actuación",
    note: QUICK_NOTES.noAction,
  },
  {
    label: "Vista preparada",
    note: "Vista revisada y agendada. Preparada la documentación necesaria.",
  },
  {
    label: "Detección incorrecta",
    note: QUICK_NOTES.detectionError,
  },
  {
    label: "Duplicado",
    note: "Hito duplicado o ya contemplado en otra actuación del expediente.",
  },
];

export function EventActions({
  eventIds,
  reviewStatus = "pending",
  isGroup = false,
}: Props) {
  const router = useRouter();

  const [actionState, setActionState] = useState<ActionState>("idle");
  const [error, setError] = useState<string | null>(null);
  const [showResolveForm, setShowResolveForm] = useState(false);
  const [resolutionNote, setResolutionNote] = useState("");

  const isPending = reviewStatus === "pending";
  const isReviewed = reviewStatus === "reviewed";
  const isResolved = reviewStatus === "resolved";

  const isLoading = actionState !== "idle";
  const hasEvents = eventIds.length > 0;
  const disabled = isLoading || !hasEvents;

  const affectedLabel =
    isGroup && eventIds.length > 1
      ? `${eventIds.length} localizaciones`
      : "1 localización";

  const trimmedResolutionNote = resolutionNote.trim();
  const canResolve = trimmedResolutionNote.length > 0 && !disabled;

  async function runAction(
    nextState: ActionState,
    action: (eventId: string) => Promise<unknown>,
    fallbackError: string,
  ) {
    if (disabled) return;

    setActionState(nextState);
    setError(null);

    try {
      await Promise.all(eventIds.map((eventId) => action(eventId)));
      router.refresh();
    } catch (error) {
      console.error(error);
      setError(error instanceof Error ? error.message : fallbackError);
    } finally {
      setActionState("idle");
    }
  }

  async function handleReview() {
    await runAction(
      "reviewing",
      reviewEvent,
      isGroup
        ? "No se pudo marcar el hito como revisado."
        : "No se pudo marcar el evento como revisado.",
    );
  }

  async function handleQuickResolve(note: string) {
    await runAction(
      "quick-resolving",
      (eventId) => resolveEvent(eventId, note),
      isGroup
        ? "No se pudo resolver este hito."
        : "No se pudo resolver el evento.",
    );
  }

  function openResolveForm() {
    setError(null);
    setShowResolveForm(true);
    setResolutionNote("");
  }

  function closeResolveForm() {
    if (isLoading) return;

    setShowResolveForm(false);
    setResolutionNote("");
    setError(null);
  }

  async function handleConfirmResolve() {
    if (!trimmedResolutionNote) {
      setError("Añade una nota antes de resolver.");
      return;
    }

    await runAction(
      "resolving",
      (eventId) => resolveEvent(eventId, trimmedResolutionNote),
      isGroup
        ? "No se pudo resolver este hito."
        : "No se pudo resolver el evento.",
    );

    setShowResolveForm(false);
    setResolutionNote("");
  }

  async function handleReopen() {
    await runAction(
      "reopening",
      reopenEvent,
      isGroup
        ? "No se pudo reabrir el hito."
        : "No se pudo reabrir el evento.",
    );
  }

  if (isResolved) {
    return (
      <div className="space-y-2">
        <div className="flex flex-wrap gap-2">
          <div className="inline-flex rounded-xl border border-blue-900/70 bg-blue-950/30 px-3 py-1.5 text-xs font-medium text-blue-100">
            Resuelto
          </div>

          <button
            type="button"
            onClick={handleReopen}
            disabled={disabled}
            className="inline-flex rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-1.5 text-xs font-medium text-neutral-100 transition hover:bg-neutral-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {actionState === "reopening" ? "Reabriendo..." : "Reabrir"}
          </button>
        </div>

        {error ? <p className="text-xs text-red-300">{error}</p> : null}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap gap-2">
        {isPending ? (
          <button
            type="button"
            onClick={handleReview}
            disabled={disabled}
            className="inline-flex rounded-xl border border-purple-900/70 bg-purple-950/30 px-3 py-1.5 text-xs font-medium text-purple-100 transition hover:bg-purple-950/50 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {actionState === "reviewing" ? "Revisando..." : "Revisar"}
          </button>
        ) : null}

        {isReviewed ? (
          <div className="inline-flex rounded-xl border border-purple-900/70 bg-purple-950/30 px-3 py-1.5 text-xs font-medium text-purple-100">
            Revisado
          </div>
        ) : null}

        {isPending || isReviewed ? (
          <>
            <button
              type="button"
              onClick={openResolveForm}
              disabled={disabled}
              className="inline-flex rounded-xl border border-blue-900/70 bg-blue-950/30 px-3 py-1.5 text-xs font-medium text-blue-100 transition hover:bg-blue-950/50 disabled:cursor-not-allowed disabled:opacity-50"
            >
              Resolver
            </button>

            <button
              type="button"
              onClick={() => handleQuickResolve(QUICK_NOTES.noAction)}
              disabled={disabled}
              className="inline-flex rounded-xl border border-green-900/70 bg-green-950/30 px-3 py-1.5 text-xs font-medium text-green-100 transition hover:bg-green-950/50 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {actionState === "quick-resolving"
                ? "Resolviendo..."
                : "No requiere actuación"}
            </button>

            <button
              type="button"
              onClick={() => handleQuickResolve(QUICK_NOTES.detectionError)}
              disabled={disabled}
              className="inline-flex rounded-xl border border-orange-900/70 bg-orange-950/30 px-3 py-1.5 text-xs font-medium text-orange-100 transition hover:bg-orange-950/50 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {actionState === "quick-resolving"
                ? "Resolviendo..."
                : "Error de detección"}
            </button>
          </>
        ) : null}

        {isReviewed ? (
          <button
            type="button"
            onClick={handleReopen}
            disabled={disabled}
            className="inline-flex rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-1.5 text-xs font-medium text-neutral-100 transition hover:bg-neutral-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {actionState === "reopening" ? "Reabriendo..." : "Reabrir"}
          </button>
        ) : null}
      </div>

      {showResolveForm ? (
        <div className="rounded-2xl border border-blue-900/70 bg-blue-950/20 p-4">
          <div className="space-y-4">
            <div>
              <p className="text-sm font-medium text-blue-100">
                Resolver actuación
              </p>

              <p className="mt-1 text-xs leading-5 text-blue-200/80">
                Deja una nota sencilla para saber por qué este hito queda
                resuelto. Se guardará en el histórico del expediente.
              </p>

              {isGroup ? (
                <p className="mt-3 rounded-xl border border-blue-900/70 bg-blue-950/30 px-3 py-2 text-xs leading-5 text-blue-100/80">
                  Este hito aparece en varios documentos. La resolución se
                  aplicará a {affectedLabel}.
                </p>
              ) : null}
            </div>

            <div className="space-y-2">
              <p className="text-xs font-medium text-blue-100">
                Notas rápidas
              </p>

              <div className="flex flex-wrap gap-2">
                {RESOLUTION_TEMPLATES.map((template) => (
                  <button
                    key={template.label}
                    type="button"
                    onClick={() => {
                      setResolutionNote(template.note);
                      setError(null);
                    }}
                    disabled={isLoading}
                    className="inline-flex rounded-xl border border-blue-900/70 bg-neutral-950/80 px-3 py-1.5 text-xs font-medium text-blue-100 transition hover:border-blue-600 hover:bg-blue-950/40 disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    {template.label}
                  </button>
                ))}
              </div>
            </div>

            <label className="block space-y-1">
              <span className="text-xs font-medium text-blue-100">
                Nota de resolución
              </span>

              <textarea
                value={resolutionNote}
                onChange={(event) => {
                  setResolutionNote(event.target.value);
                  setError(null);
                }}
                disabled={isLoading}
                rows={4}
                className="w-full rounded-xl border border-blue-900/70 bg-neutral-950 px-3 py-2 text-sm text-neutral-100 outline-none transition placeholder:text-neutral-600 focus:border-blue-500 disabled:cursor-not-allowed disabled:opacity-60"
                placeholder="Ej.: Presentado escrito por LexNET el día..., no requiere actuación adicional..."
              />
            </label>

            <div className="flex flex-wrap gap-2">
              <button
                type="button"
                onClick={handleConfirmResolve}
                disabled={!canResolve}
                className="inline-flex rounded-xl border border-blue-700 bg-blue-700 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-blue-600 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {actionState === "resolving"
                  ? "Resolviendo..."
                  : "Guardar y resolver"}
              </button>

              <button
                type="button"
                onClick={closeResolveForm}
                disabled={isLoading}
                className="inline-flex rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-1.5 text-xs font-medium text-neutral-100 transition hover:bg-neutral-800 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Cancelar
              </button>
            </div>
          </div>
        </div>
      ) : null}

      {error ? <p className="text-xs text-red-300">{error}</p> : null}
    </div>
  );
}