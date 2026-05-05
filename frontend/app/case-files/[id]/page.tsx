import Link from "next/link";
import { getCaseFile, getCaseFileDashboard } from "@/lib/api";
import { CaseDashboard } from "@/app/components/case-dashboard";
import { EventsTable } from "@/app/components/events-table";
import { DocumentUploadForm } from "./document-upload-form";
import { DocumentsList } from "./documents-list";
import { NoteForm } from "./note-form";
import { DeleteNoteButton } from "./delete-note-button";
import { displayCaseType, displayStatus } from "@/lib/document-display";

type Props = {
  params: Promise<{ id: string }>;
  searchParams?: Promise<{ tab?: string }>;
};

type CaseFileTab = "resumen" | "eventos" | "documentos" | "notas";

function normalizeTab(value?: string): CaseFileTab {
  switch (value) {
    case "eventos":
    case "documentos":
    case "notas":
      return value;
    default:
      return "resumen";
  }
}

function tabClass(isActive: boolean) {
  return isActive
    ? "rounded-2xl border border-red-900/70 bg-red-950/40 p-4 shadow-sm shadow-red-950/20"
    : "rounded-2xl border border-neutral-800 bg-neutral-950/70 p-4 transition hover:border-neutral-700 hover:bg-neutral-900";
}

function tabTitleClass(isActive: boolean) {
  return isActive
    ? "text-sm font-semibold text-red-100"
    : "text-sm font-semibold text-neutral-100";
}

function tabHintClass(isActive: boolean) {
  return isActive ? "text-xs text-red-200/70" : "text-xs text-neutral-500";
}

function tabCountClass(isActive: boolean) {
  return isActive
    ? "rounded-full border border-red-800 bg-red-950 px-2 py-0.5 text-xs font-semibold text-red-100"
    : "rounded-full border border-neutral-700 bg-black/30 px-2 py-0.5 text-xs font-semibold text-neutral-400";
}

export default async function CaseFileDetailPage({
  params,
  searchParams,
}: Props) {
  const { id } = await params;
  const resolvedSearchParams = searchParams ? await searchParams : {};
  const activeTab = normalizeTab(resolvedSearchParams.tab);

  const [detail, dashboard] = await Promise.all([
    getCaseFile(id),
    getCaseFileDashboard(id),
  ]);

  const caseFile = detail.case_file;
  const events = dashboard.upcoming_events;

  const pendingReviewDocuments = detail.documents.filter(
    (summary) =>
      !summary.document.review_status ||
      summary.document.review_status === "pending_review",
  ).length;

  const reviewedDocuments = detail.documents.filter(
    (summary) => summary.document.review_status === "reviewed",
  ).length;

  const errorDocuments = detail.documents.filter(
    (summary) => summary.document.review_status === "error",
  ).length;

  return (
    <main className="space-y-8">
      {/* HEADER */}
      <section className="space-y-4">
        <Link
          href="/case-files"
          className="text-sm text-neutral-400 underline-offset-4 hover:underline"
        >
          ← Volver a expedientes
        </Link>

        <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
          <div className="space-y-3">
            <div>
              <h1 className="text-3xl font-bold tracking-tight text-neutral-50">
                {caseFile.title}
              </h1>

              <p className="mt-2 text-sm text-neutral-400">
                {caseFile.reference} · {displayCaseType(caseFile.type)} ·{" "}
                {displayStatus(caseFile.status)}
              </p>
            </div>

            {caseFile.description && (
              <p className="max-w-3xl text-sm leading-6 text-neutral-400">
                {caseFile.description}
              </p>
            )}

            <div className="flex flex-wrap gap-2">
              <span className="inline-flex rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs text-neutral-300">
                Calendario: {caseFile.calendar_scope || "madrid"}
              </span>

              {caseFile.august_non_business && (
                <span className="inline-flex rounded-full border border-yellow-800 bg-yellow-950/40 px-3 py-1 text-xs text-yellow-200">
                  Agosto inhábil
                </span>
              )}
            </div>
          </div>

          <Link
            href="/events"
            className="inline-flex rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm text-red-100 hover:bg-red-950/50"
          >
            Ver agenda global
          </Link>
        </div>
      </section>

      {/* TABS */}
      <nav className="rounded-3xl border border-neutral-800 bg-neutral-950/70 p-2">
        <div className="grid gap-2 md:grid-cols-4">
          {[
            ["resumen", "Resumen", "Estado general", "●"],
            ["eventos", "Eventos", "Agenda", events.length],
            ["documentos", "Documentos", "Gestión", detail.documents.length],
            ["notas", "Notas", "Seguimiento", detail.notes.length],
          ].map(([key, title, hint, count]) => (
            <Link
              key={key}
              href={`/case-files/${caseFile.id}?tab=${key}`}
              className={tabClass(activeTab === key)}
            >
              <div className="flex justify-between">
                <div>
                  <p className={tabTitleClass(activeTab === key)}>{title}</p>
                  <p className={tabHintClass(activeTab === key)}>{hint}</p>
                </div>
                <span className={tabCountClass(activeTab === key)}>
                  {count}
                </span>
              </div>
            </Link>
          ))}
        </div>
      </nav>

      {/* RESUMEN */}
      {activeTab === "resumen" && (
        <>
          <CaseDashboard dashboard={dashboard} />

          <section className="grid gap-3 md:grid-cols-3">
            <Stat label="Pendientes" value={pendingReviewDocuments} />
            <Stat label="Revisados" value={reviewedDocuments} tone="green" />
            <Stat label="Errores" value={errorDocuments} tone="red" />
          </section>
        </>
      )}

      {/* EVENTOS */}
      {activeTab === "eventos" && (
        <EventsTable items={events} showResolvedInfo={false} />
      )}

      {/* DOCUMENTOS */}
      {activeTab === "documentos" && (
        <section className="space-y-5">
          <DocumentUploadForm caseFileId={caseFile.id} />
          <DocumentsList
            caseFileId={caseFile.id}
            documents={detail.documents}
          />
        </section>
      )}

      {/* NOTAS */}
       {activeTab === "notas" && (
          <section className="space-y-5">
            <div>
              <h2 className="text-2xl font-semibold text-neutral-50">Notas</h2>
              <p className="mt-1 text-sm text-neutral-400">
                Seguimiento interno del expediente, observaciones y recordatorios.
              </p>
            </div>

            <NoteForm caseFileId={caseFile.id} />

            <div className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
              {detail.notes.length === 0 ? (
                <p className="text-sm text-neutral-400">
                  No hay notas todavía. Añade la primera y empieza a construir el expediente.
                </p>
              ) : (
                <div className="space-y-3">
                  {[...detail.notes]
                    .sort(
                      (a, b) =>
                        new Date(b.created_at).getTime() -
                        new Date(a.created_at).getTime(),
                    )
                    .map((note) => (
                      <div
                        key={note.id}
                        className="rounded-2xl border border-neutral-800 bg-neutral-900/80 p-4 transition hover:border-neutral-700"
                      >
                        <div className="flex items-start justify-between gap-4">
                          <div className="min-w-0">
                            <h3 className="text-sm font-semibold text-neutral-100">
                              {note.title}
                            </h3>

                            <p className="mt-2 text-sm leading-5 text-neutral-400">
                              {note.content}
                            </p>
                          </div>

                          <div className="flex shrink-0 flex-col items-end gap-2">
                            <span className="rounded-full border border-neutral-700 bg-neutral-950 px-2.5 py-1 text-xs text-neutral-400">
                              {new Date(note.created_at).toLocaleString("es-ES")}
                            </span>

                            <DeleteNoteButton caseFileId={caseFile.id} noteId={note.id} />
                          </div>
                        </div>
                      </div>
                    ))}
                </div>
              )}
            </div>
          </section>
        )}
    </main>
  );
}

/* mini componente */
function Stat({
  label,
  value,
  tone = "amber",
}: {
  label: string;
  value: number;
  tone?: "amber" | "green" | "red";
}) {
  const map = {
    amber: "border-amber-900/60 bg-amber-950/20 text-amber-50",
    green: "border-emerald-900/60 bg-emerald-950/20 text-emerald-50",
    red: "border-red-900/60 bg-red-950/20 text-red-50",
  };

  return (
    <div className={`rounded-2xl border p-4 ${map[tone]}`}>
      <p className="text-sm opacity-70">{label}</p>
      <p className="mt-2 text-2xl font-semibold">{value}</p>
    </div>
  );
}