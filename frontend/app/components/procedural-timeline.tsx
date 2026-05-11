"use client";

import Link from "next/link";
import { useState } from "react";

import {
  eventRelationLabel,
  findRelatedEvent,
} from "@/lib/event-relations";
import type { EventItem } from "@/lib/types";

type Props = {
  events: EventItem[];
  caseFileId: string;
};

type DayGroup = {
  date: string;
  items: EventItem[];
};

function eventLabel(value?: string) {
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
    case "filing":
      return "Presentación";
    default:
      return value || "Hito";
  }
}

function proceduralPhase(item: EventItem) {
  if (item.event_type === "notification") {
    return {
      label: "Fecha base",
      description: "Puede iniciar el cómputo de plazos.",
      className: "border-sky-800 bg-sky-950/40 text-sky-100",
    };
  }

  if (item.event_type === "deadline" && item.date_kind === "relative") {
    return {
      label: "Plazo calculado",
      description: "Deriva de una fecha base detectada.",
      className: "border-red-800 bg-red-950/40 text-red-100",
    };
  }

  if (item.event_type === "deadline") {
    return {
      label: "Vencimiento",
      description: "Fecha límite detectada de forma directa.",
      className: "border-orange-800 bg-orange-950/40 text-orange-100",
    };
  }

  if (item.event_type === "requirement") {
    return {
      label: "Actuación requerida",
      description: "Puede exigir aportar documentación o presentar escrito.",
      className: "border-yellow-800 bg-yellow-950/40 text-yellow-100",
    };
  }

  if (item.event_type === "hearing" || item.event_type === "appearance") {
    return {
      label: "Señalamiento",
      description: "Actuación presencial o vista.",
      className: "border-purple-800 bg-purple-950/40 text-purple-100",
    };
  }

  return {
    label: "Hito procesal",
    description: "Evento relevante detectado en el expediente.",
    className: "border-neutral-700 bg-neutral-900 text-neutral-200",
  };
}

function statusLabel(value?: string) {
  switch (value) {
    case "overdue":
      return "Vencido";
    case "today":
      return "Hoy";
    case "upcoming":
      return "Próximo";
    default:
      return value || "Sin estado";
  }
}

function priorityLabel(value?: string) {
  switch (value) {
    case "critical":
      return "Crítico";
    case "high":
      return "Alta";
    case "medium":
      return "Media";
    case "low":
      return "Baja";
    default:
      return value || "Sin prioridad";
  }
}

function formatDate(value?: string) {
  if (!value) return "-";

  const [year, month, day] = value.split("-");

  if (!year || !month || !day) return value;

  return `${day}/${month}/${year}`;
}

function timelineTone(item: EventItem) {
  if (item.status === "overdue" && item.priority === "critical") {
    return {
      card: "border-red-900/70 bg-red-950/20",
      text: "text-red-100",
      badge: "border-red-800 bg-red-950/60 text-red-200",
    };
  }

  if (item.status === "overdue") {
    return {
      card: "border-orange-900/60 bg-orange-950/10",
      text: "text-orange-100",
      badge: "border-orange-800 bg-orange-950/60 text-orange-200",
    };
  }

  if (item.status === "today") {
    return {
      card: "border-yellow-900/60 bg-yellow-950/10",
      text: "text-yellow-100",
      badge: "border-yellow-800 bg-yellow-950/60 text-yellow-200",
    };
  }

  return {
    card: "border-neutral-800 bg-neutral-950/70",
    text: "text-neutral-100",
    badge: "border-neutral-700 bg-neutral-900 text-neutral-300",
  };
}

function shortText(value?: string, maxLength = 150) {
  const normalized = (value ?? "").trim().replace(/\s+/g, " ");

  if (!normalized) return "Sin texto de origen.";
  if (normalized.length <= maxLength) return normalized;

  return `${normalized.slice(0, maxLength).trim()}...`;
}

function phaseRank(eventType?: string) {
  switch (eventType) {
    case "notification":
      return 0;
    case "requirement":
      return 1;
    case "deadline":
      return 2;
    case "filing":
      return 3;
    case "hearing":
      return 4;
    case "appearance":
      return 5;
    default:
      return 9;
  }
}

function sortEvents(events: EventItem[]) {
  return [...events].sort((a, b) => {
    const dateCompare = a.event_date.localeCompare(b.event_date);

    if (dateCompare !== 0) {
      return dateCompare;
    }

    return phaseRank(a.event_type) - phaseRank(b.event_type);
  });
}

function groupByDay(events: EventItem[]): DayGroup[] {
  const groups = new Map<string, EventItem[]>();

  for (const event of events) {
    const current = groups.get(event.event_date) ?? [];

    current.push(event);
    groups.set(event.event_date, current);
  }

  return Array.from(groups.entries()).map(([date, items]) => ({
    date,
    items,
  }));
}

function computationLabel(item: EventItem) {
  if (item.date_kind === "relative") {
    const days =
      typeof item.relative_days === "number" && item.relative_days > 0
        ? `${item.relative_days} ${
            item.is_business_days ? "días hábiles" : "días naturales"
          }`
        : "plazo relativo";

    return item.anchor_date
      ? `${days} desde ${formatDate(item.anchor_date)}`
      : days;
  }

  return "Fecha directa";
}

function TimelineMetric({
  label,
  value,
}: {
  label: string;
  value: number;
}) {
  return (
    <span className="rounded-full border border-neutral-800 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
      {label}: {value}
    </span>
  );
}

export function ProceduralTimeline({ events, caseFileId }: Props) {
  const [showAll, setShowAll] = useState(false);

  const sortedEvents = sortEvents(events);
  const visibleEvents = showAll ? sortedEvents : sortedEvents.slice(0, 12);
  const hiddenCount = Math.max(0, sortedEvents.length - 12);
  const groups = groupByDay(visibleEvents);

  const relativeCount = sortedEvents.filter(
    (event) => event.date_kind === "relative",
  ).length;

  const linkedCount = sortedEvents.filter((event) =>
    findRelatedEvent(event, sortedEvents),
  ).length;

  if (sortedEvents.length === 0) {
    return (
      <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <p className="text-sm font-medium uppercase tracking-wide text-neutral-500">
          Timeline procesal
        </p>

        <p className="mt-2 text-sm text-neutral-400">
          Todavía no hay hitos procesales suficientes para construir una línea
          temporal.
        </p>
      </section>
    );
  }

  return (
    <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
      <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
        <div>
          <p className="text-sm font-medium uppercase tracking-wide text-neutral-500">
            Timeline procesal
          </p>

          <h2 className="mt-2 text-2xl font-semibold text-neutral-50">
            Secuencia de hitos del expediente
          </h2>

          <p className="mt-2 max-w-3xl text-sm leading-6 text-neutral-400">
            Vista cronológica con lectura procesal: fechas base,
            requerimientos, plazos calculados, vencimientos y señalamientos.
          </p>

          <div className="mt-4 flex flex-wrap gap-2">
            <TimelineMetric label="Hitos" value={sortedEvents.length} />
            <TimelineMetric label="Visibles" value={visibleEvents.length} />
            <TimelineMetric label="Plazos relativos" value={relativeCount} />
            <TimelineMetric label="Con origen" value={linkedCount} />
          </div>
        </div>

        <Link
          href={`/case-files/${caseFileId}?tab=eventos`}
          className="inline-flex rounded-xl border border-neutral-700 bg-neutral-900 px-4 py-2 text-sm font-medium text-neutral-100 transition hover:bg-neutral-800"
        >
          Ver agenda del expediente
        </Link>
      </div>

      <div className="mt-5 flex flex-wrap gap-2 text-xs">
        <LegendBadge
          label="Fecha base"
          className="border-sky-800 bg-sky-950/40 text-sky-100"
        />

        <LegendBadge
          label="Plazo calculado"
          className="border-red-800 bg-red-950/40 text-red-100"
        />

        <LegendBadge
          label="Actuación requerida"
          className="border-yellow-800 bg-yellow-950/40 text-yellow-100"
        />

        <LegendBadge
          label="Señalamiento"
          className="border-purple-800 bg-purple-950/40 text-purple-100"
        />
      </div>

      <div className="mt-6 space-y-6">
        {groups.map((group, groupIndex) => {
          const isLastGroup = groupIndex === groups.length - 1;

          return (
            <div key={group.date} className="grid grid-cols-[auto_1fr] gap-4">
              <div className="flex flex-col items-center">
                <div className="rounded-full border border-neutral-700 bg-neutral-900 px-2 py-1 text-[11px] font-medium text-neutral-300">
                  {formatDate(group.date)}
                </div>

                {!isLastGroup ? (
                  <div className="mt-2 h-full min-h-20 w-px bg-neutral-800" />
                ) : null}
              </div>

              <div className="space-y-3">
                {group.items.map((item) => {
                  const tone = timelineTone(item);
                  const phase = proceduralPhase(item);
                  const relation = eventRelationLabel(item, sortedEvents);
                  const relatedEvent = findRelatedEvent(item, sortedEvents);

                  return (
                    <article
                      key={item.event_id}
                      className={`rounded-2xl border p-4 ${tone.card}`}
                    >
                      <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                        <div className="min-w-0">
                          <div className="flex flex-wrap items-center gap-2">
                            <span
                              className={`rounded-full border px-2.5 py-1 text-xs font-medium ${phase.className}`}
                              title={phase.description}
                            >
                              {phase.label}
                            </span>

                            {item.date_kind === "relative" ? (
                              <span className="rounded-full border border-red-900/60 bg-red-950/20 px-2.5 py-1 text-xs font-medium text-red-100">
                                Plazo relativo
                              </span>
                            ) : null}

                            {item.anchor_date ? (
                              <span className="rounded-full border border-sky-900/60 bg-sky-950/20 px-2.5 py-1 text-xs font-medium text-sky-100">
                                Base: {formatDate(item.anchor_date)}
                              </span>
                            ) : null}
                          </div>

                          <Link
                            href={`/events/${item.event_id}`}
                            className={`mt-3 inline-block text-base font-semibold underline-offset-4 hover:underline ${tone.text}`}
                          >
                            {eventLabel(item.event_type)} ·{" "}
                            {formatDate(item.event_date)}
                          </Link>

                          <p className="mt-1 text-sm text-neutral-400">
                            {shortText(item.source_text)}
                          </p>

                          <p className="mt-2 text-xs text-neutral-500">
                            Cómputo: {computationLabel(item)}
                          </p>

                          {relation ? (
                            <div className="mt-1 space-y-1">
                              <p className="text-xs font-medium text-sky-300">
                                ↳ {relation}
                              </p>

                              {relatedEvent ? (
                                <Link
                                  href={`/events/${relatedEvent.event_id}`}
                                  className="inline-flex w-fit rounded-lg border border-sky-900/60 bg-sky-950/20 px-2.5 py-1 text-xs font-medium text-sky-100 transition hover:bg-sky-900/50"
                                >
                                  Ver evento origen
                                </Link>
                              ) : null}
                            </div>
                          ) : null}
                        </div>

                        <div className="flex shrink-0 flex-wrap gap-2">
                          <span
                            className={`rounded-full border px-2.5 py-1 text-xs font-medium ${tone.badge}`}
                          >
                            {statusLabel(item.status)}
                          </span>

                          <span className="rounded-full border border-neutral-700 bg-neutral-900 px-2.5 py-1 text-xs font-medium text-neutral-300">
                            {priorityLabel(item.priority)}
                          </span>
                        </div>
                      </div>

                      <div className="mt-3 flex flex-wrap gap-2 text-xs text-neutral-500">
                        {item.original_name ? (
                          <span>Documento: {item.original_name}</span>
                        ) : null}

                        {item.computation ? (
                          <span>· {item.computation}</span>
                        ) : null}
                      </div>
                    </article>
                  );
                })}
              </div>
            </div>
          );
        })}
      </div>

      {hiddenCount > 0 ? (
        <button
          type="button"
          onClick={() => setShowAll((value) => !value)}
          className="mt-4 w-full rounded-xl border border-neutral-800 bg-neutral-900 px-4 py-3 text-sm font-medium text-neutral-300 transition hover:bg-neutral-800 hover:text-neutral-100"
        >
          {showAll
            ? "Ver menos hitos"
            : `Ver ${hiddenCount} hitos adicionales`}
        </button>
      ) : null}
    </section>
  );
}

function LegendBadge({
  label,
  className,
}: {
  label: string;
  className: string;
}) {
  return (
    <span
      className={`rounded-full border px-2.5 py-1 font-medium ${className}`}
    >
      {label}
    </span>
  );
}