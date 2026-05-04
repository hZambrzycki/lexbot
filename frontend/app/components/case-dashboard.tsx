import { Dashboard } from "@/lib/types";

type Props = {
  dashboard: Dashboard;
};

function StatCard({
  label,
  value,
  tone = "neutral",
}: {
  label: string;
  value: number | string;
  tone?: "neutral" | "red" | "orange" | "blue";
}) {
  const toneClass = {
    neutral: "border-neutral-800 bg-neutral-900 text-neutral-100",
    red: "border-red-900/60 bg-red-950/20 text-red-100",
    orange: "border-orange-900/60 bg-orange-950/20 text-orange-100",
    blue: "border-blue-900/60 bg-blue-950/20 text-blue-100",
  }[tone];

  return (
    <div className={`rounded-2xl border p-4 ${toneClass}`}>
      <div className="text-sm opacity-70">{label}</div>
      <div className="mt-2 text-2xl font-semibold">{value}</div>
    </div>
  );
}
function formatAlert(text: string) {
  if (!text || text === "no alerts" || text === "no immediate alerts") {
    return "Sin alertas inmediatas";
  }

  const deadlineMatch = text.match(
    /critical deadline on (\d{4}-\d{2}-\d{2}) \((\d+) days ago\)/,
  );

  if (deadlineMatch) {
    const [, date, days] = deadlineMatch;
    return `Plazo crítico vencido el ${date} (${days} días atrás)`;
  }

  return text
    .replaceAll("critical deadline", "plazo crítico")
    .replaceAll("deadline", "plazo")
    .replaceAll("notification", "notificación")
    .replaceAll("requirement", "requerimiento")
    .replaceAll("hearing", "vista")
    .replaceAll("appearance", "comparecencia")
    .replaceAll(" on ", " el ")
    .replaceAll("days ago", "días atrás");
}

function formatAction(text: string) {
  if (
    !text ||
    text === "no action required" ||
    text === "no immediate procedural action"
  ) {
    return "Sin actuación procesal inmediata";
  }

  const match = text.match(
    /review overdue critical deadline \((.+)\) from (\d{4}-\d{2}-\d{2}) immediately/,
  );

  if (match) {
    const [, documentName, date] = match;
    return `Revisar de inmediato el plazo crítico vencido desde ${date} (${documentName})`;
  }

  return text
    .replaceAll("review overdue", "revisar vencido")
    .replaceAll("critical deadline", "plazo crítico")
    .replaceAll("deadline", "plazo")
    .replaceAll("notification", "notificación")
    .replaceAll("requirement", "requerimiento")
    .replaceAll("hearing", "vista")
    .replaceAll("appearance", "comparecencia")
    .replaceAll(" from ", " desde ")
    .replaceAll("immediately", "de inmediato");
}

function formatHint(text: string) {
  if (
    !text ||
    text === "none" ||
    text === "no immediate procedural concerns"
  ) {
    return "Sin incidencias procesales inmediatas";
  }

  return text
    .replaceAll("possible deadline breach", "Posible incumplimiento de plazo")
    .replaceAll("deadline", "plazo")
    .replaceAll("review required", "Revisión necesaria");
}

export function CaseDashboard({ dashboard }: Props) {
  const hasAlert = dashboard.needs_attention || dashboard.overdue_count > 0;

  return (
    <section className="space-y-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-3 xl:grid-cols-5">
        <StatCard label="Documentos" value={dashboard.document_count} />
        <StatCard label="Notas" value={dashboard.note_count} />
        <StatCard
          label="Pendientes de revisión"
          value={dashboard.pending_review_count}
          tone={dashboard.pending_review_count > 0 ? "orange" : "neutral"}
        />
        <StatCard
          label="Eventos activos"
          value={dashboard.active_event_count}
          tone={dashboard.active_event_count > 0 ? "red" : "neutral"}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div
          className={`rounded-2xl border p-5 ${
            hasAlert
              ? "border-red-900/50 bg-red-950/20"
              : "border-neutral-800 bg-neutral-900"
          }`}
        >
          <div
            className={`text-sm ${
              hasAlert ? "text-red-300/80" : "text-neutral-400"
            }`}
          >
            Alerta principal
          </div>

          <div
            className={`mt-2 text-lg font-semibold ${
              hasAlert ? "text-red-200" : "text-neutral-100"
            }`}
          >
            {formatAlert(dashboard.top_alert)}
          </div>
        </div>

        <div className="rounded-2xl border border-orange-900/40 bg-orange-950/20 p-5">
          <div className="text-sm text-orange-300/80">
            Siguiente acción recomendada
          </div>

          <div className="mt-2 text-lg font-semibold leading-7 text-orange-200">
            {formatAction(dashboard.recommended_next_action)}
          </div>
        </div>

        <div className="rounded-2xl border border-yellow-900/40 bg-yellow-950/20 p-5">
          <div className="text-sm text-yellow-300/80">
            Advertencia procesal
          </div>

          <div className="mt-2 text-lg font-semibold leading-7 text-yellow-200">
            {formatHint(dashboard.procedural_hint)}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
        <div className="rounded-2xl border border-red-900/50 bg-red-950/10 p-4">
          <p className="text-xl font-semibold text-red-100">
            {dashboard.critical_count}
          </p>
          <p className="mt-1 text-xs text-red-200/70">Críticos</p>
        </div>

        <div className="rounded-2xl border border-orange-900/50 bg-orange-950/10 p-4">
          <p className="text-xl font-semibold text-orange-100">
            {dashboard.overdue_count}
          </p>
          <p className="mt-1 text-xs text-orange-200/70">Vencidos</p>
        </div>

        <div className="rounded-2xl border border-yellow-900/50 bg-yellow-950/10 p-4">
          <p className="text-xl font-semibold text-yellow-100">
            {dashboard.today_count}
          </p>
          <p className="mt-1 text-xs text-yellow-200/70">Hoy</p>
        </div>

        <div className="rounded-2xl border border-emerald-900/50 bg-emerald-950/10 p-4">
          <p className="text-xl font-semibold text-emerald-100">
            {dashboard.upcoming_count}
          </p>
          <p className="mt-1 text-xs text-emerald-200/70">Próximos</p>
        </div>
      </div>
    </section>
  );
}