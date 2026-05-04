"use client";

import { useMemo, useState } from "react";
import { EventItem } from "@/lib/types";
import { EventActions } from "./event-actions";

type Props = {
  items: EventItem[];
};

type Filter = "all" | "pending" | "critical" | "resolved" | "overdue";

function displayLabel(value?: string) {
  switch (value) {
    case "overdue":
      return "Vencido";
    case "today":
      return "Hoy";
    case "upcoming":
      return "Próximo";
    case "critical":
      return "Crítico";
    case "high":
      return "Alta";
    case "medium":
      return "Media";
    case "low":
      return "Baja";
    case "pending":
      return "Pendiente";
    case "reviewed":
      return "Revisado";
    case "resolved":
      return "Resuelto";
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
      return value ?? "-";
  }
}

function badgeClass(value?: string) {
  switch (value) {
    case "critical":
    case "overdue":
      return "border-red-800 bg-red-950/60 text-red-200";
    case "high":
      return "border-orange-800 bg-orange-950/60 text-orange-200";
    case "medium":
    case "today":
      return "border-yellow-800 bg-yellow-950/60 text-yellow-200";
    case "low":
    case "upcoming":
      return "border-emerald-800 bg-emerald-950/60 text-emerald-200";
    case "resolved":
      return "border-blue-800 bg-blue-950/60 text-blue-200";
    case "reviewed":
      return "border-purple-800 bg-purple-950/60 text-purple-200";
    case "pending":
      return "border-neutral-700 bg-neutral-800 text-neutral-200";
    default:
      return "border-neutral-800 bg-neutral-900 text-neutral-300";
  }
}

function Badge({ value }: { value?: string }) {
  return (
    <span
      className={`inline-flex rounded-full border px-2.5 py-1 text-xs font-medium ${badgeClass(
        value,
      )}`}
    >
      {displayLabel(value)}
    </span>
  );
}

function formatDocuments(item: EventItem) {
  if (item.document_names && item.document_names.length > 0) {
    return item.document_names.join(", ");
  }

  return item.original_name ?? "-";
}

function formatRemaining(item: EventItem) {
  if (typeof item.days_remaining !== "number") return "";

  if (item.days_remaining < 0) {
    return `Hace ${Math.abs(item.days_remaining)} días`;
  }

  if (item.days_remaining === 0) {
    return "Hoy";
  }

  return `En ${item.days_remaining} días`;
}

function displayComputation(value?: string) {
  switch (value) {
    case "absolute date":
      return "Fecha absoluta";
    case "relative date":
      return "Fecha relativa";
    case "anchor + 1 business days":
      return "Fecha base + 1 día hábil";
    case "anchor + 1 natural days":
      return "Fecha base + 1 día natural";
    case "anchor + next_day + 1 business days":
      return "Fecha base + día siguiente + 1 día hábil";
    case "anchor + next_day + 1 natural days":
      return "Fecha base + día siguiente + 1 día natural";
    default:
      if (!value) return "-";

      return value
        .replace("anchor", "Fecha base")
        .replace("next_day", "día siguiente")
        .replace("business days", "días hábiles")
        .replace("natural days", "días naturales")
        .replaceAll(" + ", " + ");
  }
}

function cardClass(item: EventItem) {
  const isCritical = item.priority === "critical";
  const isOverdue = item.status === "overdue";
  const isResolved = item.review_status === "resolved";

  if (isResolved) return "border-blue-900/70 bg-blue-950/20";
  if (isCritical && isOverdue) {
    return "border-red-800/80 bg-red-950/20 shadow-red-950/20";
  }
  if (isOverdue) return "border-red-900/60 bg-neutral-900";

  return "border-neutral-800 bg-neutral-900";
}

function filterItems(items: EventItem[], filter: Filter) {
  switch (filter) {
    case "pending":
      return items.filter((item) => item.review_status === "pending");
    case "critical":
      return items.filter((item) => item.priority === "critical");
    case "resolved":
      return items.filter((item) => item.review_status === "resolved");
    case "overdue":
      return items.filter((item) => item.status === "overdue");
    default:
      return items;
  }
}

function filterCount(items: EventItem[], filter: Filter) {
  return filterItems(items, filter).length;
}

export function EventsTable({ items }: Props) {
  const [filter, setFilter] = useState<Filter>("all");

  const filters: { value: Filter; label: string }[] = [
    { value: "all", label: "Todos" },
    { value: "pending", label: "Pendientes" },
    { value: "critical", label: "Críticos" },
    { value: "resolved", label: "Resueltos" },
    { value: "overdue", label: "Vencidos" },
  ];

  const filteredItems = useMemo(() => filterItems(items, filter), [items, filter]);

  if (items.length === 0) {
    return (
      <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-6 text-neutral-400">
        No hay eventos detectados.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-row flex-wrap items-center gap-2">
        {filters.map((item) => (
          <button
            key={item.value}
            onClick={() => setFilter(item.value)}
            className={`inline-flex flex-row items-center gap-2 whitespace-nowrap rounded-xl border px-3 py-1.5 text-sm transition ${
              filter === item.value
                ? "border-white bg-white text-black"
                : "border-neutral-700 bg-neutral-900 text-neutral-300 hover:bg-neutral-800"
            }`}
          >
            <span>{item.label}</span>
            <span className="inline-flex h-5 min-w-5 items-center justify-center rounded-full bg-black/20 px-1.5 text-xs opacity-80">
              {filterCount(items, item.value)}
            </span>
          </button>
        ))}
      </div>

      {filteredItems.length === 0 ? (
        <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-6 text-neutral-400">
          No hay eventos para este filtro.
        </div>
      ) : (
        <div className="space-y-3">
          {filteredItems.map((item) => (
            <article
              key={item.event_id}
              className={`rounded-2xl border p-5 shadow-sm transition hover:bg-neutral-800/30 ${cardClass(
                item,
              )}`}
            >
              <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div className="space-y-3">
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge value={item.status} />
                    <Badge value={item.priority} />
                    <Badge value={item.review_status} />

                    {item.duplicate_count && item.duplicate_count > 1 ? (
                      <span className="rounded-full border border-neutral-700 bg-neutral-900 px-2.5 py-1 text-xs text-neutral-300">
                        Detectado en {item.duplicate_count} documentos
                      </span>
                    ) : null}
                  </div>

                  <div>
                    <div className="flex flex-wrap items-baseline gap-3">
                      <h3 className="text-lg font-semibold">
                        {displayLabel(item.event_type)} · {item.event_date}
                      </h3>
                      <span className="text-sm text-neutral-400">
                        {formatRemaining(item)}
                      </span>
                    </div>

                    <p className="mt-2 max-w-3xl text-sm text-neutral-300">
                      {item.source_text}
                    </p>
                  </div>

                  <div className="grid gap-3 text-sm text-neutral-400 md:grid-cols-2">
                    <div>
                      <div className="text-xs uppercase tracking-wide text-neutral-500">
                        Documentos
                      </div>
                      <div className="mt-1 text-neutral-200">
                        {formatDocuments(item)}
                      </div>
                    </div>

                    <div>
                      <div className="text-xs uppercase tracking-wide text-neutral-500">
                        Cómputo
                      </div>
                      <div className="mt-1 text-neutral-200">
                        {displayComputation(item.computation)}
                      </div>
                    </div>
                  </div>

                  {item.resolution_note ? (
                    <div className="rounded-xl border border-blue-900/60 bg-blue-950/30 p-3 text-sm text-blue-100">
                      <div className="text-xs uppercase tracking-wide text-blue-300/80">
                        Nota de resolución
                      </div>
                      <div className="mt-1">{item.resolution_note}</div>
                    </div>
                  ) : null}
                </div>

                <div className="lg:min-w-48">
                  <EventActions eventId={item.event_id} />
                </div>
              </div>
            </article>
          ))}
        </div>
      )}
    </div>
  );
}