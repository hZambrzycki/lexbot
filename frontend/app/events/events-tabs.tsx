"use client";

import { useMemo, useState } from "react";
import { EventsTable } from "@/app/components/events-table";
import { EventItem } from "@/lib/types";

type Tab = "pending" | "reviewed" | "resolved";

type Props = {
  pendingEvents: EventItem[];
  reviewedEvents: EventItem[];
  resolvedEvents: EventItem[];
  baseUrl: string;
};

function countByStatus(items: EventItem[], status: string) {
  return items.filter((item) => item.status === status).length;
}

function countUrgentOverdue(items: EventItem[]) {
  return items.filter(
    (item) => item.status === "overdue" && item.priority === "critical",
  ).length;
}

function countPriorityUpcoming(items: EventItem[]) {
  return items.filter(
    (item) =>
      item.status === "upcoming" &&
      (item.priority === "critical" || item.priority === "high"),
  ).length;
}

function nextMainAlert(items: EventItem[]) {
  const urgent = items.find(
    (item) => item.status === "overdue" && item.priority === "critical",
  );

  if (urgent) {
    return `Plazo urgente vencido desde ${urgent.event_date}`;
  }

  const today = items.find((item) => item.status === "today");

  if (today) {
    return `Actuación con vencimiento hoy: ${today.event_date}`;
  }

  const upcoming = items.find(
    (item) =>
      item.status === "upcoming" &&
      (item.priority === "critical" || item.priority === "high"),
  );

  if (upcoming) {
    return `Próximo hito prioritario: ${upcoming.event_date}`;
  }

  return "No hay alertas procesales prioritarias pendientes.";
}

function tabClass(active: boolean, color: "neutral" | "purple" | "blue") {
  if (active) {
    switch (color) {
      case "purple":
        return "border-purple-500 bg-purple-950/50 text-purple-100 shadow-lg shadow-purple-950/20";
      case "blue":
        return "border-blue-500 bg-blue-950/50 text-blue-100 shadow-lg shadow-blue-950/20";
      default:
        return "border-white bg-white text-black shadow-lg shadow-white/10";
    }
  }

  switch (color) {
    case "purple":
      return "border-purple-900/70 bg-purple-950/20 text-purple-200 hover:bg-purple-950/40";
    case "blue":
      return "border-blue-900/70 bg-blue-950/20 text-blue-200 hover:bg-blue-950/40";
    default:
      return "border-neutral-700 bg-neutral-900 text-neutral-300 hover:bg-neutral-800";
  }
}

function tabContent(tab: Tab) {
  switch (tab) {
    case "pending":
      return {
        title: "Eventos pendientes",
        subtitle: "Requieren revisión o actuación.",
        body: "Eventos que todavía requieren revisión, actuación o resolución.",
        intro:
          "Esta es la bandeja principal de trabajo. Aquí queda lo que todavía hay que revisar, atender o resolver.",
        empty: "No hay eventos pendientes.",
      };

    case "reviewed":
      return {
        title: "Eventos revisados",
        subtitle: "Comprobados, pendientes de resolver.",
        body: "Eventos ya comprobados, pendientes de resolver definitivamente o reabrir.",
        intro:
          "Aquí quedan los hitos ya comprobados. Es una zona intermedia antes de resolverlos definitivamente.",
        empty: "No hay eventos revisados pendientes de resolver.",
      };

    case "resolved":
      return {
        title: "Eventos resueltos",
        subtitle: "Histórico cerrado y reabrible.",
        body: "Histórico de hitos cerrados. Puedes reabrirlos si necesitan nueva revisión.",
        intro:
          "Aquí queda el histórico cerrado, con posibilidad de reapertura si el hito vuelve a necesitar revisión.",
        empty: "Todavía no hay eventos resueltos.",
      };
  }
}

function metricLabels(tab: Tab) {
  switch (tab) {
    case "pending":
      return {
        urgent: "Urgentes vencidos",
        overdue: "Con fecha pasada",
        today: "Hoy",
        priority: "Próximos prioritarios",
        upcoming: "Próximos",
      };

    case "reviewed":
      return {
        urgent: "Urgentes revisados",
        overdue: "Revisados con fecha pasada",
        today: "Revisados para hoy",
        priority: "Próximos revisados",
        upcoming: "Próximos revisados",
      };

    case "resolved":
      return {
        urgent: "Urgentes resueltos",
        overdue: "Resueltos con fecha pasada",
        today: "Resueltos de hoy",
        priority: "Próximos resueltos",
        upcoming: "Próximos resueltos",
      };
  }
}

export function EventsTabs({
  pendingEvents,
  reviewedEvents,
  resolvedEvents,
  baseUrl,
}: Props) {
  const [activeTab, setActiveTab] = useState<Tab>("pending");

  const activeEvents = useMemo(() => {
    switch (activeTab) {
      case "reviewed":
        return reviewedEvents;
      case "resolved":
        return resolvedEvents;
      default:
        return pendingEvents;
    }
  }, [activeTab, pendingEvents, reviewedEvents, resolvedEvents]);

  const activeContent = tabContent(activeTab);
  const labels = metricLabels(activeTab);

  const urgentOverdueCount = countUrgentOverdue(activeEvents);
  const overdueCount = countByStatus(activeEvents, "overdue");
  const todayCount = countByStatus(activeEvents, "today");
  const priorityUpcomingCount = countPriorityUpcoming(activeEvents);
  const upcomingCount = countByStatus(activeEvents, "upcoming");

  const mainAlert = nextMainAlert(pendingEvents);

  return (
    <div className="space-y-8">
      <section className="grid gap-3 md:grid-cols-3">
        {([
          ["pending", "Pendientes", pendingEvents.length, "neutral"],
          ["reviewed", "Revisados", reviewedEvents.length, "purple"],
          ["resolved", "Resueltos", resolvedEvents.length, "blue"],
        ] as const).map(([tab, label, count, color]) => {
          const content = tabContent(tab);

          return (
            <button
              key={tab}
              type="button"
              onClick={() => setActiveTab(tab)}
              className={`rounded-2xl border p-4 text-left transition ${tabClass(
                activeTab === tab,
                color,
              )}`}
            >
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-medium uppercase tracking-wide">
                    {label}
                  </p>

                  <p className="mt-2 text-xs leading-5 opacity-80">
                    {content.subtitle}
                  </p>
                </div>

                <span className="rounded-full bg-current/10 px-3 py-1 text-lg font-bold">
                  {count}
                </span>
              </div>
            </button>
          );
        })}
      </section>

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <p className="text-sm font-medium uppercase tracking-wide text-neutral-500">
          Alerta principal
        </p>

        <p className="mt-2 text-xl font-semibold text-neutral-100">
          {mainAlert}
        </p>

        <p className="mt-2 text-sm leading-6 text-neutral-400">
          Esta alerta solo mira los eventos pendientes. Así no se mezclan avisos
          activos con hitos ya revisados o resueltos.
        </p>
      </section>

      <section className="grid gap-3 md:grid-cols-2 xl:grid-cols-5">
        <MetricCard value={urgentOverdueCount} label={labels.urgent} tone="red" />
        <MetricCard value={overdueCount} label={labels.overdue} tone="redSoft" />
        <MetricCard value={todayCount} label={labels.today} tone="yellow" />
        <MetricCard
          value={priorityUpcomingCount}
          label={labels.priority}
          tone="orange"
        />
        <MetricCard value={upcomingCount} label={labels.upcoming} tone="green" />
      </section>

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
          <div>
            <h2 className="text-xl font-semibold text-neutral-50">
              Exportar agenda global
            </h2>

            <p className="mt-2 max-w-3xl text-sm leading-6 text-neutral-400">
              Descarga los hitos pendientes para importarlos en tu calendario.
            </p>
          </div>

          <div className="flex flex-wrap gap-2">
            <a
              href={`${baseUrl}/events/upcoming.ics`}
              className="inline-flex rounded-xl border border-neutral-700 bg-neutral-900 px-4 py-2 text-sm font-medium text-neutral-100 transition hover:bg-neutral-800"
            >
              Exportar agenda global
            </a>

            <a
              href={`${baseUrl}/events/upcoming.ics?type=deadline`}
              className="inline-flex rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/50"
            >
              Exportar solo plazos
            </a>
          </div>
        </div>

        <div className="mt-4 grid gap-3 md:grid-cols-2">
          <div className="rounded-xl border border-neutral-800 bg-neutral-900 p-4 text-sm leading-6 text-neutral-400">
            <span className="font-medium text-neutral-200">
              Agenda global:
            </span>{" "}
            incluye los eventos pendientes detectados en todos los expedientes.
          </div>

          <div className="rounded-xl border border-neutral-800 bg-neutral-900 p-4 text-sm leading-6 text-neutral-400">
            <span className="font-medium text-neutral-200">
              Solo plazos:
            </span>{" "}
            incluye únicamente vencimientos y plazos procesales pendientes.
          </div>
        </div>
      </section>

      <section className="space-y-3">
        <div className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
          <p className="text-sm font-medium uppercase tracking-wide text-neutral-500">
            Vista actual
          </p>

          <h2 className="mt-2 text-2xl font-semibold">
            {activeContent.title}
          </h2>

          <p className="mt-2 text-sm leading-6 text-neutral-400">
            {activeContent.body}
          </p>

          <p className="mt-3 rounded-xl border border-neutral-800 bg-neutral-900 px-3 py-2 text-xs leading-5 text-neutral-400">
            {activeContent.intro}
          </p>
        </div>

        {activeEvents.length > 0 ? (
          <EventsTable
            items={activeEvents}
            showCaseFileInfo
            showResolvedInfo={activeTab === "resolved"}
            mode={activeTab}
          />
        ) : (
          <div className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5 text-sm text-neutral-400">
            {activeContent.empty}
          </div>
        )}
      </section>
    </div>
  );
}

function MetricCard({
  value,
  label,
  tone,
}: {
  value: number;
  label: string;
  tone: "red" | "redSoft" | "yellow" | "orange" | "green";
}) {
  const className = {
    red: "border-red-800/70 bg-red-950/20 text-red-100 text-red-300/80",
    redSoft: "border-red-900/60 bg-red-950/10 text-red-100 text-red-300/80",
    yellow:
      "border-yellow-800/70 bg-yellow-950/20 text-yellow-100 text-yellow-300/80",
    orange:
      "border-orange-800/70 bg-orange-950/20 text-orange-100 text-orange-300/80",
    green:
      "border-emerald-800/70 bg-emerald-950/20 text-emerald-100 text-emerald-300/80",
  }[tone];

  const [borderBg, numberColor, labelColor] = className
    .split(" text-")
    .map((part, index) => (index === 0 ? part : `text-${part}`));

  return (
    <div className={`rounded-2xl border p-4 ${borderBg}`}>
      <p className={`text-2xl font-semibold ${numberColor}`}>{value}</p>
      <p className={`mt-1 text-xs uppercase tracking-wide ${labelColor}`}>
        {label}
      </p>
    </div>
  );
}