import Link from "next/link";
import { getCaseFile, getCaseFileDashboard } from "@/lib/api";
import { CaseDashboard } from "@/app/components/case-dashboard";
import { EventsTable } from "@/app/components/events-table";
import { DocumentUploadForm } from "./document-upload-form";
type Props = {
  params: Promise<{ id: string }>;
};

function displayCaseType(value: string) {
  switch (value) {
    case "civil":
      return "Civil";
    case "labor":
    case "laboral":
      return "Laboral";
    case "extranjeria":
      return "Extranjería";
    case "mercantil":
      return "Mercantil";
    case "administrativo":
    case "administrative":
      return "Administrativo";
    case "otros":
      return "Otros";
    default:
      return value;
  }
}

function displayStatus(value: string) {
  switch (value) {
    case "open":
      return "Abierto";
    case "closed":
      return "Cerrado";
    case "archived":
      return "Archivado";
    default:
      return value;
  }
}

function displayMimeType(value: string) {
  switch (value) {
    case "text/plain":
      return "Texto plano";
    case "application/pdf":
      return "PDF";
    case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
      return "Word";
    default:
      return value;
  }
}

export default async function CaseFileDetailPage({ params }: Props) {
  const { id } = await params;

  const [detail, dashboard] = await Promise.all([
    getCaseFile(id),
    getCaseFileDashboard(id),
  ]);

  const events = dashboard.upcoming_events;
  const caseFile = detail.case_file;

  return (
    <main className="space-y-8">
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

            {caseFile.description ? (
              <p className="max-w-3xl text-sm leading-6 text-neutral-400">
                {caseFile.description}
              </p>
            ) : null}

            <div className="flex flex-wrap gap-2">
              <span className="inline-flex rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
                Calendario: {caseFile.calendar_scope || "madrid"}
              </span>

              {caseFile.august_non_business ? (
                <span className="inline-flex rounded-full border border-yellow-800 bg-yellow-950/40 px-3 py-1 text-xs font-medium text-yellow-200">
                  Agosto inhábil
                </span>
              ) : null}
            </div>
          </div>

          <Link
            href="/events"
            className="inline-flex w-fit rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/50"
          >
            Ver agenda global
          </Link>
        </div>
      </section>

      <CaseDashboard dashboard={dashboard} />

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
          <div>
            <h2 className="text-xl font-semibold text-neutral-50">
              Exportar agenda del expediente
            </h2>

            <p className="mt-2 max-w-3xl text-sm leading-6 text-neutral-400">
              Descarga los hitos procesales detectados en este expediente para
              importarlos en tu calendario.
            </p>
          </div>

          <div className="flex flex-wrap gap-2">
            <a
              href={`${process.env.NEXT_PUBLIC_API_BASE_URL}/case-files/${caseFile.id}/events/upcoming.ics`}
              className="inline-flex rounded-xl border border-neutral-700 bg-neutral-900 px-4 py-2 text-sm font-medium text-neutral-100 transition hover:bg-neutral-800"
            >
              Exportar agenda
            </a>

            <a
              href={`${process.env.NEXT_PUBLIC_API_BASE_URL}/case-files/${caseFile.id}/events/upcoming.ics?type=deadline`}
              className="inline-flex rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/50"
            >
              Exportar solo plazos
            </a>
          </div>
        </div>

        <div className="mt-4 grid gap-3 md:grid-cols-2">
          <div className="rounded-xl border border-neutral-800 bg-neutral-900 p-4 text-sm leading-6 text-neutral-400">
            <span className="font-medium text-neutral-200">
              Exportar agenda:
            </span>{" "}
            incluye notificaciones, plazos, vistas, comparecencias y demás hitos
            detectados.
          </div>

          <div className="rounded-xl border border-neutral-800 bg-neutral-900 p-4 text-sm leading-6 text-neutral-400">
            <span className="font-medium text-neutral-200">
              Exportar solo plazos:
            </span>{" "}
            incluye únicamente vencimientos y plazos procesales del expediente.
          </div>
        </div>
      </section>

      <section className="space-y-3">
        <div>
          <h2 className="text-2xl font-semibold text-neutral-50">Eventos</h2>
          <p className="mt-1 text-sm text-neutral-400">
            Hitos procesales detectados en los documentos de este expediente.
          </p>
        </div>

        <EventsTable items={events} showResolvedInfo={false} />
      </section>

      <section className="grid gap-5 lg:grid-cols-2">
        <div className="space-y-3">
          <h2 className="text-2xl font-semibold text-neutral-50">
            Documentos
          </h2>
            <DocumentUploadForm caseFileId={caseFile.id} />
          <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5">
            {detail.documents.length === 0 ? (
              <p className="text-sm text-neutral-400">
                No hay documentos asociados a este expediente.
              </p>
            ) : (
              <ul className="space-y-3">
                {detail.documents.map((doc) => (
                  <li
                    key={doc.id}
                    className="rounded-xl border border-neutral-800 bg-black/20 p-3"
                  >
                    <Link
                      href={`/case-files/${caseFile.id}/documents/${doc.id}`}
                      className="font-medium text-neutral-100 underline-offset-4 hover:underline"
                    >
                      {doc.original_name}
                    </Link>

                    <div className="mt-1 text-xs text-neutral-500">
                      {displayMimeType(doc.mime_type)}
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        <div className="space-y-3">
          <h2 className="text-2xl font-semibold text-neutral-50">Notas</h2>

          <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5">
            {detail.notes.length === 0 ? (
              <p className="text-sm text-neutral-400">
                No hay notas registradas en este expediente.
              </p>
            ) : (
              <ul className="space-y-3">
                {detail.notes.map((note) => (
                  <li
                    key={note.id}
                    className="rounded-xl border border-neutral-800 bg-black/20 p-3"
                  >
                    <div className="font-medium text-neutral-100">
                      {note.title}
                    </div>

                    <div className="mt-2 text-sm leading-6 text-neutral-400">
                      {note.content}
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </section>
    </main>
  );
}