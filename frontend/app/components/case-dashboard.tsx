import { Dashboard } from "@/lib/types";

type Props = {
  dashboard: Dashboard;
};

function StatCard({ label, value }: { label: string; value: number | string }) {
  return (
    <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-4">
      <div className="text-sm text-neutral-400">{label}</div>
      <div className="mt-2 text-2xl font-semibold">{value}</div>
    </div>
  );
}

function translateAlert(text: string) {
  return text
    .replace("critical deadline", "plazo crítico")
    .replace("deadline", "plazo")
    .replace("notification", "notificación")
    .replace(" on ", " el ")
    .replace("days ago", "días atrás");
}

function translateAction(text: string) {
  return text
    .replace("review overdue", "revisar vencido")
    .replace("critical deadline", "plazo crítico")
    .replace("deadline", "plazo")
    .replace(" from ", " desde ")
    .replace("immediately", "de inmediato");
}
function translateHint(text: string) {
  return text
    .replace("possible deadline breach", "posible incumplimiento de plazo");
}

export function CaseDashboard({ dashboard }: Props) {
  return (
    <section className="space-y-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-3 xl:grid-cols-5">
        <StatCard label="Documentos" value={dashboard.document_count} />
        <StatCard label="Notas" value={dashboard.note_count} />
        <StatCard label="Pendientes de revisión" value={dashboard.pending_review_count} />
        <StatCard label="Resueltos" value={dashboard.resolved_count} />
        <StatCard label="Eventos activos" value={dashboard.active_event_count} />
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <div className="rounded-2xl border border-red-900/40 bg-red-950/20 p-5">
          <div className="text-sm text-red-300/80">Alerta principal</div>
          <div className="mt-2 text-lg font-semibold text-red-200">
            {translateAlert(dashboard.top_alert)}
          </div>
        </div>

        <div className="rounded-2xl border border-orange-900/40 bg-orange-950/20 p-5">
          <div className="text-sm text-orange-300/80">Siguiente acción recomendada</div>
          <div className="mt-2 text-lg font-semibold text-orange-200">
            {translateAction(dashboard.recommended_next_action)}
          </div>
        </div>
      </div>

      <div className="rounded-2xl border border-yellow-900/40 bg-yellow-950/20 p-5">
        <div className="text-sm text-yellow-300/80">Advertencia procesal</div>
        <div className="mt-2 text-lg font-semibold text-yellow-200">
          {translateHint(dashboard.procedural_hint)}
        </div>
      </div>
    </section>
  );
}