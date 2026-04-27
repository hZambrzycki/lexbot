import Link from "next/link";
import { getCaseFile, getCaseFileDashboard } from "@/lib/api";
import { CaseDashboard } from "@/app/components/case-dashboard";
import { EventsTable } from "@/app/components/events-table";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function CaseFileDetailPage({ params }: Props) {
  const { id } = await params;

  const [detail, dashboard] = await Promise.all([
    getCaseFile(id),
    getCaseFileDashboard(id),
  ]);

  const events = dashboard.upcoming_events;

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
          <h1 className="text-3xl font-bold">{detail.case_file.title}</h1>
          <p className="mt-2 text-neutral-400">
            {detail.case_file.reference} · {detail.case_file.type} ·{" "}
            {detail.case_file.status}
          </p>
        </div>
      </div>

      <CaseDashboard dashboard={dashboard} />

      <section className="space-y-3">
<h2 className="text-2xl font-semibold">Eventos</h2>        
<EventsTable items={events} />
      </section>

      <section className="space-y-3">
<h2 className="text-2xl font-semibold">Documentos</h2>
        <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5">
          <ul className="space-y-2">
            {detail.documents.map((doc) => (
              <li key={doc.id} className="text-sm">
                {doc.original_name} · {doc.mime_type}
              </li>
            ))}
          </ul>
        </div>
      </section>

      <section className="space-y-3">
<h2 className="text-2xl font-semibold">Notas</h2>
        <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5">
          <ul className="space-y-3">
            {detail.notes.map((note) => (
              <li key={note.id}>
                <div className="font-medium">{note.title}</div>
                <div className="text-sm text-neutral-400">{note.content}</div>
              </li>
            ))}
          </ul>
        </div>
      </section>
    </main>
  );
}