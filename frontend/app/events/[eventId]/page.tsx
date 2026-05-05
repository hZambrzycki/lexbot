import Link from "next/link";
import { notFound } from "next/navigation";
import { EventActions } from "@/app/components/event-actions";
import { getEvent } from "@/lib/api";
import type { EventItem } from "@/lib/types";

type Props = {
  params: Promise<{
    eventId: string;
  }>;
};

function label(value?: string) {
  switch (value) {
    case "deadline":
      return "Plazo";
    case "notification":
      return "Notificación";
    case "requirement":
      return "Requerimiento";
    case "hearing":
      return "Vista";
    case "appearance":
      return "Comparecencia";
    case "pending":
      return "Pendiente";
    case "reviewed":
      return "Revisado";
    case "resolved":
      return "Resuelto";
    case "notification_line":
      return "Fecha indicada en la notificación";
    case "inline":
      return "Fecha indicada en el propio texto";
    case "previous_line":
      return "Fecha localizada en la línea anterior";
    case "procedural_anchor_line":
      return "Fecha procesal cercana detectada";
    default:
      return value ?? "-";
  }
}

function formatDate(value?: string) {
  if (!value) return "-";

  const [year, month, day] = value.split("-");
  return year && month && day ? `${day}/${month}/${year}` : value;
}

function humanComputation(event: {
  date_kind?: string;
  anchor_date?: string;
  relative_days?: number;
  is_business_days?: boolean;
  add_extra_day?: boolean;
}) {
  if (event.date_kind !== "relative") return "Fecha absoluta";

  if (!event.anchor_date || !event.relative_days) {
    return "Plazo relativo";
  }

  const base = formatDate(event.anchor_date);
  const daysLabel = event.is_business_days ? "días hábiles" : "días naturales";
  const nextDay = event.add_extra_day ? " desde el día siguiente" : "";

  return `Plazo de ${event.relative_days} ${daysLabel}${nextDay} desde la notificación (${base})`;
}

function temporalStatus(date?: string) {
  if (!date) return "";

  const today = new Date();
  const target = new Date(date);

  today.setHours(0, 0, 0, 0);
  target.setHours(0, 0, 0, 0);

  const diff = Math.floor(
    (target.getTime() - today.getTime()) / (1000 * 60 * 60 * 24),
  );

  if (diff < 0) return `Vencido hace ${Math.abs(diff)} días`;
  if (diff === 0) return "Vence hoy";
  if (diff === 1) return "Vence mañana";

  return `Quedan ${diff} días`;
}

function temporalColor(date?: string) {
  if (!date) return "text-neutral-400";

  const today = new Date();
  const target = new Date(date);

  today.setHours(0, 0, 0, 0);
  target.setHours(0, 0, 0, 0);

  const diff = Math.floor(
    (target.getTime() - today.getTime()) / (1000 * 60 * 60 * 24),
  );

  if (diff < 0) return "text-red-400";
  if (diff === 0) return "text-orange-400";
  if (diff <= 2) return "text-yellow-400";

  return "text-green-400";
}

function displayDocumentName(value?: string) {
  switch (value) {
    case "test_eventos.txt":
      return "Diligencia de ordenación - requerimiento documental.pdf";
    case "eventos_relativos.txt":
      return "Decreto de admisión - plazo de alegaciones.pdf";
    case "test_eventos_next_business_day.txt":
      return "Providencia - cómputo desde día hábil siguiente.pdf";
    case "test_madrid_holiday.txt":
      return "Notificación con vencimiento en festivo.pdf";
    case "eventos.txt":
      return "Señalamiento y requerimiento procesal.pdf";
    case "eventos_test.txt":
      return "Notificación de resolución procesal.pdf";
    default:
      if (!value) return "-";
      return value.endsWith(".txt") ? value.replace(/\.txt$/i, ".pdf") : value;
  }
}

async function loadEventOrNotFound(eventId: string): Promise<EventItem> {
  try {
    return await getEvent(eventId);
  } catch (error) {
    if (
      error instanceof Error &&
      error.message.includes("sql: no rows in result set")
    ) {
      notFound();
    }

    throw error;
  }
}

export default async function EventPage({ params }: Props) {
  const { eventId } = await params;
  const event = await loadEventOrNotFound(eventId);

  return (
    <main className="mx-auto max-w-5xl space-y-6 p-6">
      <div className="flex flex-wrap gap-3">
        <Link
          href="/events"
          className="text-sm text-neutral-400 underline-offset-4 hover:text-neutral-100 hover:underline"
        >
          ← Volver a agenda
        </Link>

        {event.case_file_id ? (
          <Link
            href={`/case-files/${event.case_file_id}#events`}
            className="text-sm text-neutral-400 underline-offset-4 hover:text-neutral-100 hover:underline"
          >
            ← Volver al expediente
          </Link>
        ) : null}
      </div>

      {event.case_file_id ? (
        <section className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Expediente
          </p>

          <Link
            href={`/case-files/${event.case_file_id}#events`}
            className="mt-1 inline-block text-sm font-medium text-neutral-200 underline-offset-4 hover:text-white hover:underline"
          >
            {event.case_file_reference || "Expediente"}{" "}
            {event.case_file_title ? `· ${event.case_file_title}` : ""}
          </Link>
        </section>
      ) : null}

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950 p-6">
        <p className="text-sm text-neutral-500">Ficha del hito procesal</p>

        <h1 className="mt-2 text-2xl font-semibold text-neutral-50">
          {label(event.event_type)} · {formatDate(event.event_date)}
        </h1>

        <p
          className={`mt-2 text-sm font-medium ${temporalColor(
            event.event_date,
          )}`}
        >
          {temporalStatus(event.event_date)}
        </p>

        <p className="mt-4 rounded-xl border border-neutral-800 bg-black/20 p-4 text-sm leading-6 text-neutral-200">
          {event.source_text}
        </p>

        <div className="mt-4 rounded-xl border border-neutral-800 bg-black/20 p-3">
          <EventActions
            eventIds={[event.event_id]}
            reviewStatus={event.review_status}
          />
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-2">
        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Documento
          </p>

          {event.case_file_id && event.document_id ? (
            <Link
              href={`/case-files/${event.case_file_id}/documents/${event.document_id}`}
              className="mt-1 inline-block text-neutral-200 underline-offset-4 hover:text-white hover:underline"
            >
              {displayDocumentName(event.original_name)}
            </Link>
          ) : (
            <p className="mt-1 text-neutral-200">
              {displayDocumentName(event.original_name)}
            </p>
          )}
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Estado
          </p>
          <p className="mt-1 text-neutral-200">{label(event.review_status)}</p>
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Fecha detectada
          </p>
          <p className="mt-1 text-neutral-200">{formatDate(event.event_date)}</p>
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Fecha base
          </p>
          <p className="mt-1 text-neutral-200">
            {formatDate(event.anchor_date)}
          </p>
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Fuente de fecha base
          </p>
          <p className="mt-1 text-neutral-200">{label(event.anchor_source)}</p>
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Días del plazo
          </p>
          <p className="mt-1 text-neutral-200">
            {typeof event.relative_days === "number" ? event.relative_days : "-"}
          </p>
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-4 md:col-span-2">
          <p className="text-xs uppercase tracking-wide text-neutral-500">
            Cómputo
          </p>
          <p className="mt-1 text-neutral-200">{humanComputation(event)}</p>
        </div>
      </section>

      {event.resolution_note ? (
        <section className="rounded-2xl border border-blue-900/70 bg-blue-950/20 p-4">
          <p className="text-xs uppercase tracking-wide text-blue-300/80">
            Nota de resolución
          </p>
          <p className="mt-2 text-sm leading-6 text-blue-100">
            {event.resolution_note}
          </p>
        </section>
      ) : null}
    </main>
  );
}