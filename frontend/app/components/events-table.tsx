"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import { EventItem } from "@/lib/types";
import { EventActions } from "./event-actions";

export type EventsTableMode = "pending" | "reviewed" | "resolved" | "mixed";

type Props = {
  items: EventItem[];
  showCaseFileInfo?: boolean;
  showResolvedInfo?: boolean;
  mode?: EventsTableMode;
};

type Filter = "all" | "urgent" | "overdue" | "upcoming" | "resolved";

type EventGroup = {
  id: string;
  representative: EventItem;
  items: EventItem[];
  documentNames: string[];
  documentIds: string[];
  occurrenceCount: number;
};

type EventInsight = {
  title: string;
  body: string;
  className: string;
};

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

function genderedLabel(
  eventType: string | undefined,
  masculine: string,
  feminine: string,
) {
  return eventType === "hearing" || eventType === "appearance"
    ? feminine
    : masculine;
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
      className={`inline-flex items-center rounded-full border px-2.5 py-1 text-xs font-medium ${badgeClass(
        value,
      )}`}
    >
      {displayLabel(value)}
    </span>
  );
}

function normalizeText(value?: string) {
  return (value ?? "")
    .trim()
    .replace(/\s+/g, " ")
    .toLowerCase();
}

function groupKey(item: EventItem) {
  return [
    item.case_file_id ?? "",
    item.event_type,
    item.event_date,
    normalizeText(item.source_text),
    normalizeText(item.computation),
  ].join("|");
}

function unique(values: Array<string | undefined>) {
  return Array.from(
    new Set(
      values.filter((value): value is string =>
        Boolean(value && value.trim()),
      ),
    ),
  );
}

function getDocumentNames(item: EventItem) {
  if (item.document_names && item.document_names.length > 0) {
    return item.document_names.map(displayDocumentName);
  }

  if (item.original_name) {
    return [displayDocumentName(item.original_name)];
  }

  return [];
}

function statusRank(status?: string) {
  switch (status) {
    case "overdue":
      return 0;
    case "today":
      return 1;
    case "upcoming":
      return 2;
    default:
      return 3;
  }
}

function priorityRank(priority?: string) {
  switch (priority) {
    case "critical":
      return 0;
    case "high":
      return 1;
    case "medium":
      return 2;
    case "low":
      return 3;
    default:
      return 4;
  }
}

function reviewRank(reviewStatus?: string) {
  switch (reviewStatus) {
    case "pending":
      return 0;
    case "reviewed":
      return 1;
    case "resolved":
      return 2;
    default:
      return 3;
  }
}

function shouldReplaceRepresentative(current: EventItem, next: EventItem) {
  const nextReviewRank = reviewRank(next.review_status);
  const currentReviewRank = reviewRank(current.review_status);

  if (nextReviewRank !== currentReviewRank) {
    return nextReviewRank < currentReviewRank;
  }

  const nextPriorityRank = priorityRank(next.priority);
  const currentPriorityRank = priorityRank(current.priority);

  if (nextPriorityRank !== currentPriorityRank) {
    return nextPriorityRank < currentPriorityRank;
  }

  const nextStatusRank = statusRank(next.status);
  const currentStatusRank = statusRank(current.status);

  return nextStatusRank < currentStatusRank;
}

function groupEvents(items: EventItem[]): EventGroup[] {
  const groups = new Map<string, EventGroup>();

  for (const item of items) {
    const key = groupKey(item);
    const current = groups.get(key);

    const itemOccurrenceCount =
      item.duplicate_count && item.duplicate_count > 1
        ? item.duplicate_count
        : 1;

    if (!current) {
      groups.set(key, {
        id: key,
        representative: item,
        items: [item],
        documentNames: getDocumentNames(item),
        documentIds: unique([item.document_id, ...(item.document_ids ?? [])]),
        occurrenceCount: itemOccurrenceCount,
      });

      continue;
    }

    current.items.push(item);

    current.documentNames = unique([
      ...current.documentNames,
      ...getDocumentNames(item),
    ]);

    current.documentIds = unique([
      ...current.documentIds,
      item.document_id,
      ...(item.document_ids ?? []),
    ]);

    current.occurrenceCount += itemOccurrenceCount;

    if (shouldReplaceRepresentative(current.representative, item)) {
      current.representative = item;
    }
  }

  return Array.from(groups.values()).sort((a, b) => {
    const aItem = a.representative;
    const bItem = b.representative;

    const statusCompare = statusRank(aItem.status) - statusRank(bItem.status);
    if (statusCompare !== 0) return statusCompare;

    const priorityCompare =
      priorityRank(aItem.priority) - priorityRank(bItem.priority);
    if (priorityCompare !== 0) return priorityCompare;

    const dateCompare = aItem.event_date.localeCompare(bItem.event_date);
    if (dateCompare !== 0) return dateCompare;

    const typeCompare = aItem.event_type.localeCompare(bItem.event_type);
    if (typeCompare !== 0) return typeCompare;

    return normalizeText(aItem.source_text).localeCompare(
      normalizeText(bItem.source_text),
    );
  });
}

function groupMatchesFilter(group: EventGroup, filter: Filter) {
  switch (filter) {
    case "urgent":
      return group.items.some((item) => item.priority === "critical");
    case "resolved":
      return group.items.every((item) => item.review_status === "resolved");
    case "overdue":
      return group.items.some((item) => item.status === "overdue");
    case "upcoming":
      return group.items.some((item) => item.status === "upcoming");
    default:
      return true;
  }
}

function filterGroups(groups: EventGroup[], filter: Filter) {
  return groups.filter((group) => groupMatchesFilter(group, filter));
}

function filterGroupCount(groups: EventGroup[], filter: Filter) {
  return filterGroups(groups, filter).length;
}

function detectionCount(groups: EventGroup[]) {
  return groups.reduce((total, group) => total + group.occurrenceCount, 0);
}

function formatDate(value?: string) {
  if (!value) return "-";

  const [year, month, day] = value.split("-");

  if (!year || !month || !day) {
    return value;
  }

  return `${day}/${month}/${year}`;
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
        .replaceAll("anchor", "Fecha base")
        .replaceAll("next_day", "día siguiente")
        .replaceAll("business days", "días hábiles")
        .replaceAll("natural days", "días naturales");
  }
}

function displayAnchorSource(value?: string) {
  switch (value) {
    case "inline":
      return "Fecha indicada en el propio texto";
    case "notification_line":
      return "Fecha indicada en la notificación";
    case "previous_line":
      return "Fecha localizada en la línea anterior";
    case "procedural_anchor_line":
      return "Fecha procesal cercana detectada";
    default:
      return value ? displayTechnicalValue(value) : "-";
  }
}

function displayDateKind(value?: string) {
  switch (value) {
    case "absolute":
      return "Fecha exacta";
    case "relative":
      return "Fecha calculada desde otra fecha";
    default:
      return value ? displayTechnicalValue(value) : "-";
  }
}

function displayCalendarScope(value?: string) {
  switch (value) {
    case "madrid":
      return "Madrid";
    case "state":
    case "national":
      return "Estatal";
    default:
      return value ? displayTechnicalValue(value) : "-";
  }
}

function displayBusinessDays(value?: boolean) {
  if (typeof value !== "boolean") return "-";
  return value ? "Días hábiles" : "Días naturales";
}

function displayExtraDay(value?: boolean) {
  if (typeof value !== "boolean") {
    return "Desde la fecha base indicada";
  }

  return value
    ? "Desde el día siguiente a la fecha base"
    : "Desde la fecha base indicada";
}

function displayTriggerText(value?: string) {
  if (!value) return "-";

  return value
    .replaceAll("_", " ")
    .replaceAll("dias", "días")
    .replaceAll("habiles", "hábiles")
    .replaceAll("dia", "día")
    .replaceAll("habil", "hábil");
}

function displayTechnicalValue(value: string) {
  return value
    .replaceAll("_", " ")
    .replaceAll("-", " ")
    .replace(/\b\w/g, (letter) => letter.toUpperCase());
}

function cardClass(group: EventGroup) {
  const item = group.representative;

  const isCritical = item.priority === "critical";
  const isOverdue = item.status === "overdue";
  const isResolved = group.items.every(
    (event) => event.review_status === "resolved",
  );

  if (isResolved) {
    return "border-blue-900/70 bg-blue-950/20";
  }

  if (isCritical && isOverdue) {
    return "border-red-800/80 bg-red-950/20 shadow-red-950/20";
  }

  if (isOverdue) {
    return "border-red-900/60 bg-neutral-900";
  }

  return "border-neutral-800 bg-neutral-900";
}

function formatDocuments(group: EventGroup) {
  if (group.documentNames.length === 0) return "-";

  if (group.documentNames.length <= 3) {
    return group.documentNames.join(", ");
  }

  return `${group.documentNames.slice(0, 3).join(", ")} y ${
    group.documentNames.length - 3
  } más`;
}

function formatOrigin(group: EventGroup) {
  if (group.documentNames.length === 0) return "Documento no identificado";

  if (group.documentNames.length === 1) {
    return group.documentNames[0];
  }

  return `${group.documentNames[0]} y ${group.documentNames.length - 1} más`;
}

function shortText(value: string, maxLength = 170) {
  const normalized = value.trim().replace(/\s+/g, " ");

  if (normalized.length <= maxLength) {
    return normalized;
  }

  return `${normalized.slice(0, maxLength).trim()}...`;
}

function groupReviewSummary(group: EventGroup) {
  const resolved = group.items.filter(
    (item) => item.review_status === "resolved",
  ).length;

  const reviewed = group.items.filter(
    (item) => item.review_status === "reviewed",
  ).length;

  const pending = group.items.filter(
    (item) => item.review_status === "pending",
  ).length;

  if (group.items.length <= 1) {
    return displayLabel(group.representative.review_status);
  }

  return `${pending} pendientes · ${reviewed} revisados · ${resolved} resueltos`;
}

function eventTitle(item: EventItem, mode: EventsTableMode) {
  const type = displayLabel(item.event_type);
  const date = formatDate(item.event_date);

  if (mode === "resolved") {
    return `${type} ${genderedLabel(
      item.event_type,
      "resuelto",
      "resuelta",
    )} · ${date}`;
  }

  if (mode === "reviewed") {
    return `${type} ${genderedLabel(
      item.event_type,
      "revisado",
      "revisada",
    )} · ${date}`;
  }

  if (item.status === "overdue" && item.priority === "critical") {
    return `${type} urgente vencido · ${date}`;
  }

  if (item.status === "overdue") {
    return `${type} con fecha pasada · ${date}`;
  }

  if (item.status === "today") {
    return `${type} para hoy · ${date}`;
  }

  if (item.status === "upcoming") {
    return `${type} ${genderedLabel(
      item.event_type,
      "próximo",
      "próxima",
    )} · ${date}`;
  }

  return `${type} · ${date}`;
}

function eventSubtitle(item: EventItem) {
  const remaining = formatRemaining(item);

  if (!remaining) {
    return displayLabel(item.review_status);
  }

  return `${displayLabel(item.review_status)} · ${remaining}`;
}

function insightFor(item: EventItem, mode: EventsTableMode): EventInsight {
  if (mode === "resolved") {
    return {
      title: "Estado",
      body: "Este hito ya consta como resuelto. Puede consultarse como histórico o reabrirse si necesita nueva revisión.",
      className: "border-blue-900/70 bg-blue-950/20 text-blue-100",
    };
  }

  if (mode === "reviewed") {
    return {
      title: "Siguiente paso",
      body: "Este hito ya fue revisado. Ahora puedes resolverlo definitivamente o reabrirlo si necesita nueva comprobación.",
      className: "border-purple-900/70 bg-purple-950/20 text-purple-100",
    };
  }

  if (
    item.status === "overdue" &&
    item.priority === "critical" &&
    item.event_type === "deadline"
  ) {
    return {
      title: "Recomendación",
      body: "Plazo crítico vencido. Comprueba si ya se presentó el escrito o si procede una actuación urgente.",
      className: "border-red-900/70 bg-red-950/25 text-red-100",
    };
  }

  if (item.status === "overdue" && item.event_type === "requirement") {
    return {
      title: "Recomendación",
      body: "Requerimiento con fecha pasada. Conviene comprobar si se atendió, si consta presentación o si hay que actuar de inmediato.",
      className: "border-orange-900/70 bg-orange-950/20 text-orange-100",
    };
  }

  if (item.status === "overdue" && item.event_type === "notification") {
    return {
      title: "Información útil",
      body: "Notificación ya pasada. Puede servir como fecha base para revisar el cómputo de otros plazos.",
      className: "border-neutral-800 bg-black/20 text-neutral-200",
    };
  }

  if (item.status === "overdue") {
    return {
      title: "Recomendación",
      body: "Este hito tiene fecha pasada. Conviene revisar si requería actuación o si solo debe quedar como referencia.",
      className: "border-orange-900/70 bg-orange-950/20 text-orange-100",
    };
  }

  if (
    item.status === "upcoming" &&
    (item.priority === "critical" || item.priority === "high") &&
    item.event_type === "deadline"
  ) {
    return {
      title: "Recomendación",
      body: "Plazo próximo con prioridad alta. Conviene preparar la actuación y dejarlo controlado en agenda.",
      className: "border-yellow-900/70 bg-yellow-950/20 text-yellow-100",
    };
  }

  if (
    item.status === "upcoming" &&
    (item.priority === "critical" || item.priority === "high") &&
    item.event_type === "hearing"
  ) {
    return {
      title: "Recomendación",
      body: "Vista próxima. Conviene revisar prueba, documentación, asistencia y estrategia del expediente.",
      className: "border-yellow-900/70 bg-yellow-950/20 text-yellow-100",
    };
  }

  if (item.status === "upcoming" && item.event_type === "appearance") {
    return {
      title: "Recomendación",
      body: "Comparecencia próxima. Conviene revisar asistencia, documentación necesaria e instrucciones al cliente.",
      className: "border-yellow-900/70 bg-yellow-950/20 text-yellow-100",
    };
  }

  if (item.status === "upcoming" && item.event_type === "requirement") {
    return {
      title: "Recomendación",
      body: "Requerimiento pendiente. Conviene preparar la documentación antes del vencimiento.",
      className: "border-yellow-900/70 bg-yellow-950/20 text-yellow-100",
    };
  }

  if (item.event_type === "notification") {
    return {
      title: "Información útil",
      body: "Notificación detectada como posible fecha base para el cómputo de plazos.",
      className: "border-neutral-800 bg-black/20 text-neutral-200",
    };
  }

  if (item.event_type === "hearing") {
    return {
      title: "Recomendación",
      body: "Vista judicial detectada. Conviene verificar señalamiento, sala, asistencia y prueba.",
      className: "border-neutral-800 bg-black/20 text-neutral-200",
    };
  }

  if (item.event_type === "appearance") {
    return {
      title: "Recomendación",
      body: "Comparecencia detectada. Conviene verificar fecha, objeto y documentación necesaria.",
      className: "border-neutral-800 bg-black/20 text-neutral-200",
    };
  }

  if (item.event_type === "requirement") {
    return {
      title: "Recomendación",
      body: "Requerimiento procesal detectado. Conviene revisar el texto completo y confirmar si exige actuación.",
      className: "border-neutral-800 bg-black/20 text-neutral-200",
    };
  }

  return {
    title: "Recomendación",
    body: "Revisa el texto de origen y confirma si requiere actuación en el expediente.",
    className: "border-neutral-800 bg-black/20 text-neutral-200",
  };
}

function groupStats(groups: EventGroup[]) {
  return {
    urgent: groups.filter((group) =>
      group.items.some((item) => item.priority === "critical"),
    ).length,

    overdue: groups.filter((group) =>
      group.items.some((item) => item.status === "overdue"),
    ).length,

    upcoming: groups.filter((group) =>
      group.items.some((item) => item.status === "upcoming"),
    ).length,

    resolved: groups.filter((group) =>
      group.items.every((item) => item.review_status === "resolved"),
    ).length,
  };
}

function statsLabels(mode: EventsTableMode) {
  if (mode === "resolved") {
    return {
      urgent: "Urgentes cerrados",
      overdue: "Cerrados con fecha pasada",
      upcoming: "Próximos cerrados",
    };
  }

  if (mode === "reviewed") {
    return {
      urgent: "Urgentes revisados",
      overdue: "Revisados con fecha pasada",
      upcoming: "Próximos revisados",
    };
  }

  return {
    urgent: "Urgentes",
    overdue: "Vencidos",
    upcoming: "Próximos",
  };
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

      return value.endsWith(".txt")
        ? value.replace(/\.txt$/i, ".pdf")
        : value;
  }
}

function filterLabel(filter: Filter, mode: EventsTableMode) {
  switch (filter) {
    case "all":
      return "Todos";
    case "urgent":
      return mode === "resolved"
        ? "Urgentes cerrados"
        : mode === "reviewed"
          ? "Urgentes revisados"
          : "Urgentes";
    case "overdue":
      return mode === "resolved"
        ? "Fecha pasada"
        : mode === "reviewed"
          ? "Fecha pasada"
          : "Vencidos";
    case "upcoming":
      return mode === "resolved"
        ? "Próximos cerrados"
        : mode === "reviewed"
          ? "Próximos revisados"
          : "Próximos";
    case "resolved":
      return "Resueltos";
  }
}

function filterDescription(filter: Filter, mode: EventsTableMode) {
  if (filter === "all") {
    if (mode === "resolved") {
      return "Mostrando todos los hitos resueltos.";
    }

    if (mode === "reviewed") {
      return "Mostrando todos los hitos revisados pendientes de resolución.";
    }

    return "Mostrando todos los hitos pendientes de revisión, actuación o resolución.";
  }

  if (filter === "urgent") {
    if (mode === "resolved") {
      return "Mostrando solo hitos urgentes que ya constan como resueltos.";
    }

    if (mode === "reviewed") {
      return "Mostrando solo hitos urgentes ya revisados.";
    }

    return "Mostrando solo hitos pendientes con prioridad crítica.";
  }

  if (filter === "overdue") {
    if (mode === "resolved") {
      return "Mostrando solo hitos resueltos cuya fecha ya pasó.";
    }

    if (mode === "reviewed") {
      return "Mostrando solo hitos revisados cuya fecha ya pasó.";
    }

    return "Mostrando solo hitos pendientes con fecha pasada.";
  }

  if (filter === "upcoming") {
    if (mode === "resolved") {
      return "Mostrando solo hitos resueltos con fecha futura.";
    }

    if (mode === "reviewed") {
      return "Mostrando solo hitos revisados con fecha futura.";
    }

    return "Mostrando solo hitos pendientes con fecha próxima.";
  }

  return "Mostrando hitos resueltos.";
}

function summaryText(
  groups: EventGroup[],
  allGroups: EventGroup[],
  filter: Filter,
  mode: EventsTableMode,
) {
  const total = allGroups.length;
  const current = groups.length;

  if (filter === "all") {
    if (mode === "resolved") {
      return `Mostrando ${current} hitos resueltos.`;
    }

    if (mode === "reviewed") {
      return `Mostrando ${current} hitos revisados.`;
    }

    return `Mostrando ${current} hitos pendientes.`;
  }

  return `Mostrando ${current} de ${total} hitos.`;
}

function hasCaseFileInfo(item: EventItem) {
  return Boolean(item.case_file_reference || item.case_file_title);
}

function caseFileLabel(item: EventItem) {
  return [item.case_file_reference, item.case_file_title]
    .filter(Boolean)
    .join(" · ");
}

function CaseFileLink({ item }: { item: EventItem }) {
  if (!hasCaseFileInfo(item)) return null;

  const label = caseFileLabel(item);

  if (item.case_file_id) {
    return (
      <Link
        href={`/case-files/${item.case_file_id}`}
        className="text-sm font-medium text-neutral-300 underline-offset-4 transition hover:text-neutral-100 hover:underline"
      >
        {label}
      </Link>
    );
  }

  return <span className="text-sm font-medium text-neutral-300">{label}</span>;
}

function RepeatedEventNotice({ count }: { count: number }) {
  return (
    <div className="rounded-xl border border-sky-900/60 bg-sky-950/20 px-3 py-2 text-xs leading-5 text-sky-100/90">
      <p className="font-medium text-sky-100">Aparece en varios documentos</p>

      <p className="mt-1 text-sky-100/70">
        Este mismo hito se ha localizado {count} veces. Si lo marcas como
        revisado o resuelto, se actualizarán todas esas apariciones.
      </p>
    </div>
  );
}

function ActiveFilterNotice({
  filter,
  mode,
  groups,
  allGroups,
  onClear,
}: {
  filter: Filter;
  mode: EventsTableMode;
  groups: EventGroup[];
  allGroups: EventGroup[];
  onClear: () => void;
}) {
  return (
    <div className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-4">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <p className="text-sm font-medium text-neutral-100">
            {filter === "all"
              ? "Vista completa"
              : `Vista filtrada: ${filterLabel(filter, mode)}`}
          </p>

          <p className="mt-1 text-sm leading-6 text-neutral-400">
            {filterDescription(filter, mode)}
          </p>

          <p className="mt-1 text-xs text-neutral-500">
            {summaryText(groups, allGroups, filter, mode)}
            {detectionCount(groups) > groups.length
              ? ` Localizaciones documentales: ${detectionCount(groups)}.`
              : ""}
          </p>
        </div>

        {filter !== "all" ? (
          <button
            type="button"
            onClick={onClear}
            className="inline-flex w-fit rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-1.5 text-xs font-medium text-neutral-100 transition hover:bg-neutral-800"
          >
            Ver todos
          </button>
        ) : null}
      </div>
    </div>
  );
}

export function EventsTable({
  items,
  showCaseFileInfo = false,
  showResolvedInfo = false,
  mode = "mixed",
}: Props) {
  const [filter, setFilter] = useState<Filter>("all");
  const [expandedGroupId, setExpandedGroupId] = useState<string | null>(null);

  const filters: Filter[] = [
    "all",
    "urgent",
    "overdue",
    "upcoming",
    ...(showResolvedInfo || mode === "resolved" ? (["resolved"] as const) : []),
  ];

  const allGroups = useMemo(() => groupEvents(items), [items]);

  const groups = useMemo(
    () => filterGroups(allGroups, filter),
    [allGroups, filter],
  );

  const stats = useMemo(() => groupStats(groups), [groups]);
  const labels = statsLabels(mode);

  function clearFilter() {
    setFilter("all");
    setExpandedGroupId(null);
  }

  if (items.length === 0) {
    return (
      <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-6 text-sm text-neutral-400">
        No hay eventos detectados.
      </div>
    );
  }

  return (
    <section className="space-y-5">
      <div className="flex flex-wrap gap-2">
        {filters.map((item) => (
          <button
            key={item}
            type="button"
            onClick={() => {
              setFilter(item);
              setExpandedGroupId(null);
            }}
            className={`inline-flex flex-row items-center gap-2 whitespace-nowrap rounded-xl border px-3 py-1.5 text-sm transition ${
              filter === item
                ? "border-white bg-white text-black"
                : "border-neutral-700 bg-neutral-900 text-neutral-300 hover:bg-neutral-800"
            }`}
          >
            <span>{filterLabel(item, mode)}</span>
            <span className="rounded-full bg-current/10 px-2 py-0.5 text-xs">
              {filterGroupCount(allGroups, item)}
            </span>
          </button>
        ))}
      </div>

      <ActiveFilterNotice
        filter={filter}
        mode={mode}
        groups={groups}
        allGroups={allGroups}
        onClear={clearFilter}
      />

      <div
        className={`grid gap-3 ${
          showResolvedInfo || mode === "resolved"
            ? "md:grid-cols-4"
            : "md:grid-cols-3"
        }`}
      >
        <div className="rounded-2xl border border-red-900/60 bg-red-950/20 p-4">
          <p className="text-2xl font-bold text-red-100">{stats.urgent}</p>
          <p className="mt-1 text-xs text-red-200/80">{labels.urgent}</p>
        </div>

        <div className="rounded-2xl border border-orange-900/60 bg-orange-950/20 p-4">
          <p className="text-2xl font-bold text-orange-100">
            {stats.overdue}
          </p>
          <p className="mt-1 text-xs text-orange-200/80">{labels.overdue}</p>
        </div>

        <div className="rounded-2xl border border-yellow-900/60 bg-yellow-950/20 p-4">
          <p className="text-2xl font-bold text-yellow-100">
            {stats.upcoming}
          </p>
          <p className="mt-1 text-xs text-yellow-200/80">{labels.upcoming}</p>
        </div>

        {showResolvedInfo || mode === "resolved" ? (
          <div className="rounded-2xl border border-blue-900/60 bg-blue-950/20 p-4">
            <p className="text-2xl font-bold text-blue-100">
              {stats.resolved}
            </p>
            <p className="mt-1 text-xs text-blue-200/80">
              Resueltos en vista actual
            </p>
          </div>
        ) : null}
      </div>

      {groups.length === 0 ? (
        <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-6 text-sm text-neutral-400">
          No hay eventos para este filtro.
        </div>
      ) : (
        <div className="space-y-4">
          {groups.map((group) => {
            const item = group.representative;
            const insight = insightFor(item, mode);
            const isExpanded = expandedGroupId === group.id;
            const isGrouped = group.occurrenceCount > 1;

            return (
              <article
                key={group.id}
                className={`rounded-2xl border p-5 shadow-lg ${cardClass(group)}`}
              >
                <div className="space-y-5">
                  <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                    <div className="min-w-0 flex-1 space-y-4">
                      <div className="flex flex-wrap items-center gap-2">
                        <Badge value={item.status} />
                        <Badge value={item.priority} />
                        <Badge value={item.review_status} />

                        {isGrouped ? (
                          <span className="inline-flex items-center rounded-full border border-sky-800 bg-sky-950/50 px-2.5 py-1 text-xs font-medium text-sky-200">
                            Localizado {group.occurrenceCount} veces
                          </span>
                        ) : null}
                      </div>

                      <div className="space-y-2">
                        <div>
                          <Link
                            href={`/events/${item.event_id}`}
                            className="inline-block text-xl font-semibold leading-7 text-neutral-50 underline-offset-4 transition hover:text-white hover:underline"
                          >
                            {eventTitle(item, mode)}
                          </Link>

                          <p className="mt-1 text-sm text-neutral-400">
                            {eventSubtitle(item)}
                          </p>
                        </div>

                        {showCaseFileInfo ? (
                          <div className="rounded-xl border border-neutral-800 bg-black/20 p-3">
                            <p className="text-xs font-medium uppercase tracking-wide text-neutral-500">
                              Expediente
                            </p>

                            <div className="mt-1">
                              <CaseFileLink item={item} />
                            </div>
                          </div>
                        ) : null}

                        <div className="grid gap-3 md:grid-cols-2">
                          <div className="rounded-xl border border-neutral-800 bg-black/20 p-3">
                            <p className="text-xs font-medium uppercase tracking-wide text-neutral-500">
                              Detectado en
                            </p>

                            <p className="mt-1 text-sm text-neutral-300">
                              {formatOrigin(group)}
                            </p>
                          </div>

                          <div className="rounded-xl border border-neutral-800 bg-black/20 p-3">
                            <p className="text-xs font-medium uppercase tracking-wide text-neutral-500">
                              Tipo de hito
                            </p>

                            <p className="mt-1 text-sm text-neutral-300">
                              {displayLabel(item.event_type)}
                            </p>
                          </div>
                        </div>

                        {item.source_text ? (
                          <div className="rounded-xl border border-neutral-800 bg-neutral-950/70 p-3">
                            <p className="text-xs font-medium uppercase tracking-wide text-neutral-500">
                              Texto localizado
                            </p>

                            <p className="mt-2 max-w-4xl text-sm leading-6 text-neutral-200">
                              “{shortText(item.source_text)}”
                            </p>
                          </div>
                        ) : null}
                      </div>
                    </div>
                  </div>

                  <div
                    className={`rounded-xl border p-3 text-sm ${insight.className}`}
                  >
                    <p className="font-medium">{insight.title}</p>
                    <p className="mt-1 leading-6 opacity-90">{insight.body}</p>
                  </div>

                  <div className="flex flex-col gap-4 rounded-xl border border-neutral-800 bg-black/20 p-3 md:flex-row md:items-start md:justify-between">
                    <div className="space-y-3">
                      <EventActions
                        eventIds={group.items.map((event) => event.event_id)}
                        reviewStatus={item.review_status}
                        isGroup={isGrouped}
                      />

                      {isGrouped ? (
                        <RepeatedEventNotice count={group.occurrenceCount} />
                      ) : null}
                    </div>

                    <button
                      type="button"
                      onClick={() =>
                        setExpandedGroupId(isExpanded ? null : group.id)
                      }
                      className="inline-flex w-fit rounded-xl border border-neutral-700 px-3 py-1.5 text-xs font-medium text-neutral-200 transition hover:bg-neutral-800"
                    >
                      {isExpanded ? "Ocultar detalles" : "Ver detalles"}
                    </button>
                  </div>

                  {isExpanded ? (
                    <div className="mt-5 space-y-5 border-t border-neutral-800 pt-5">
                      <div>
                        <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                          Texto completo detectado
                        </p>

                        <p className="rounded-xl border border-neutral-800 bg-black/20 p-3 text-sm leading-6 text-neutral-200">
                          {item.source_text}
                        </p>
                      </div>

                      <div className="grid gap-4 md:grid-cols-2">
                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Documentos
                          </p>

                          <p className="text-sm text-neutral-300">
                            {formatDocuments(group)}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Estado del hito
                          </p>

                          <p className="text-sm text-neutral-300">
                            {groupReviewSummary(group)}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Fecha detectada
                          </p>

                          <p className="text-sm text-neutral-300">
                            {formatDate(item.event_date)}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Cómputo
                          </p>

                          <p className="text-sm text-neutral-300">
                            {displayComputation(item.computation)}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Fecha base
                          </p>

                          <p className="text-sm text-neutral-300">
                            {item.anchor_date
                              ? formatDate(item.anchor_date)
                              : "-"}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Fuente de fecha base
                          </p>

                          <p className="text-sm text-neutral-300">
                            {displayAnchorSource(item.anchor_source)}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Tipo de fecha
                          </p>

                          <p className="text-sm text-neutral-300">
                            {displayDateKind(item.date_kind)}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Criterio de días
                          </p>

                          <p className="text-sm text-neutral-300">
                            {displayBusinessDays(item.is_business_days)}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Inicio del cómputo
                          </p>

                          <p className="text-sm text-neutral-300">
                            {displayExtraDay(item.add_extra_day)}
                          </p>
                        </div>

                        {item.trigger_text ? (
                          <div>
                            <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                              Expresión detectada
                            </p>

                            <p className="text-sm text-neutral-300">
                              {displayTriggerText(item.trigger_text)}
                            </p>
                          </div>
                        ) : null}

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Días del plazo
                          </p>

                          <p className="text-sm text-neutral-300">
                            {typeof item.relative_days === "number"
                              ? item.relative_days
                              : "-"}
                          </p>
                        </div>

                        <div>
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-neutral-500">
                            Ámbito calendario
                          </p>

                          <p className="text-sm text-neutral-300">
                            {displayCalendarScope(item.calendar_scope)}
                          </p>
                        </div>
                      </div>

                      {item.resolution_note ? (
                        <div className="rounded-xl border border-blue-900/70 bg-blue-950/20 p-3">
                          <p className="mb-1 text-xs font-medium uppercase tracking-wide text-blue-300/80">
                            Nota de resolución
                          </p>

                          <p className="text-sm leading-6 text-blue-100">
                            {item.resolution_note}
                          </p>
                        </div>
                      ) : null}
                    </div>
                  ) : null}
                </div>
              </article>
            );
          })}
        </div>
      )}
    </section>
  );
}