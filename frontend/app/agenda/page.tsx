import Link from "next/link";
import {
  getGlobalUpcomingEvents,
  getGlobalUpcomingICSUrl,
} from "@/lib/api";
import { EventsTable } from "@/app/components/events-table";
import { IcsExportActions } from "@/app/components/ics-export-actions";
import { EventItem } from "@/lib/types";

function countBy(
  items: EventItem[],
  predicate: (item: EventItem) => boolean,
): number {
  return items.filter(predicate).length;
}

function findTopAlert(items: EventItem[]): string {
  const criticalOverdue = items.find(
    (item) => item.status === "overdue" && item.priority === "critical",
  );

  if (criticalOverdue) {
    return `Plazo crítico vencido desde ${criticalOverdue.event_date}`;
  }

  const today = items.find((item) => item.status === "today");

  if (today) {
    return `Actuación para hoy: ${today.event_date}`;
  }

  const upcomingHigh = items.find(
    (item) => item.status === "upcoming" && item.priority === "high",
  );

  if (upcomingHigh) {
    return `Próximo hito prioritario: ${upcomingHigh.event_date}`;
  }

  return "Sin alertas procesales críticas";
}

export default async function AgendaPage() {
  const events = await getGlobalUpcomingEvents({
    reviewStatus: "pending",
  });

  const criticalOverdueCount = countBy(
    events,
    (item) => item.status === "overdue" && item.priority === "critical",
  );

  const overdueCount = countBy(events, (item) => item.status === "overdue");
  const todayCount = countBy(events, (item) => item.status === "today");

  const upcomingPriorityCount = countBy(
    events,
    (item) =>
      item.status === "upcoming" &&
      (item.priority === "critical" || item.priority === "high"),
  );

  const upcomingCount = countBy(events, (item) => item.status === "upcoming");

  const resolvedCount = countBy(
    events,
    (item) => item.review_status === "resolved",
  );

  return (
    <main className="space-y-8">
      <div className="space-y-3">
        <Link
          href="/case-files"
          className="text-sm text-neutral-400 underline-offset-4 hover:underline"
        >
          ← Volver a expedientes
        </Link>

        <div>
          <h1 className="text-3xl font-bold">Agenda procesal</h1>
          <p className="mt-2 max-w-3xl text-neutral-400">
            Vista global de hitos procesales detectados en todos los
            expedientes. Prioriza plazos vencidos, actuaciones de hoy y
            próximos eventos relevantes.
          </p>
        </div>
      </div>

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <p className="text-xs font-medium uppercase tracking-wide text-neutral-500">
          Alerta principal
        </p>
        <p className="mt-2 text-lg font-semibold text-neutral-100">
          {findTopAlert(events)}
        </p>

        <div className="mt-5 grid gap-3 md:grid-cols-3 xl:grid-cols-6">
          <div className="rounded-2xl border border-red-900/60 bg-red-950/20 p-4">
            <p className="text-2xl font-bold text-red-100">
              {criticalOverdueCount}
            </p>
            <p className="mt-1 text-xs text-red-200/80">
              Críticos vencidos
            </p>
          </div>

          <div className="rounded-2xl border border-red-900/50 bg-neutral-900 p-4">
            <p className="text-2xl font-bold text-neutral-100">
              {overdueCount}
            </p>
            <p className="mt-1 text-xs text-neutral-400">Vencidos</p>
          </div>

          <div className="rounded-2xl border border-yellow-900/60 bg-yellow-950/20 p-4">
            <p className="text-2xl font-bold text-yellow-100">
              {todayCount}
            </p>
            <p className="mt-1 text-xs text-yellow-200/80">Hoy</p>
          </div>

          <div className="rounded-2xl border border-orange-900/60 bg-orange-950/20 p-4">
            <p className="text-2xl font-bold text-orange-100">
              {upcomingPriorityCount}
            </p>
            <p className="mt-1 text-xs text-orange-200/80">
              Próximos prioritarios
            </p>
          </div>

          <div className="rounded-2xl border border-emerald-900/60 bg-emerald-950/20 p-4">
            <p className="text-2xl font-bold text-emerald-100">
              {upcomingCount}
            </p>
            <p className="mt-1 text-xs text-emerald-200/80">Próximos</p>
          </div>

          <div className="rounded-2xl border border-blue-900/60 bg-blue-950/20 p-4">
            <p className="text-2xl font-bold text-blue-100">
              {resolvedCount}
            </p>
            <p className="mt-1 text-xs text-blue-200/80">Resueltos</p>
          </div>
        </div>
      </section>

      <IcsExportActions
        title="Exportar agenda global"
        description="Descarga todos los hitos procesales pendientes de todos los expedientes."
        links={[
          {
            href: getGlobalUpcomingICSUrl(),
            label: "Exportar agenda global",
            description:
              "Incluye todos los eventos detectados en todos los expedientes.",
          },
          {
            href: getGlobalUpcomingICSUrl({ type: "deadline" }),
            label: "Exportar solo plazos",
            description:
              "Incluye únicamente vencimientos y plazos procesales globales.",
          },
        ]}
      />

      <section className="space-y-3">
        <div>
          <h2 className="text-2xl font-semibold">Eventos pendientes</h2>
          <p className="mt-1 text-sm text-neutral-400">
            Mostrando eventos pendientes de revisión o resolución en todos los
            expedientes.
          </p>
        </div>

    <EventsTable items={events} showCaseFileInfo />
      </section>
    </main>
  );
}