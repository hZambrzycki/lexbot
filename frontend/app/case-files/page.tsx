import Link from "next/link";
import { getCaseFiles } from "@/lib/api";
import { CaseFilesTable } from "@/app/components/case-files-table";

export default async function CaseFilesPage() {
  const items = await getCaseFiles();

  return (
    <main className="space-y-6">
      <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <Link
            href="/"
            className="text-sm text-neutral-400 underline-offset-4 hover:underline"
          >
            ← Volver al inicio
          </Link>

          <h1 className="mt-3 text-3xl font-bold">Expedientes</h1>

          <p className="mt-2 text-neutral-400">
            Vista operativa de los asuntos activos, configuración procesal y
            acceso rápido a cada expediente.
          </p>
        </div>

        <div className="flex flex-wrap gap-2">
        <Link
          href="/case-files/new"
          className="inline-flex w-fit rounded-xl border border-neutral-700 bg-neutral-100 px-4 py-2 text-sm font-medium text-neutral-950 transition hover:bg-white"
        >
          Crear expediente
        </Link>

        <Link
          href="/events"
          className="inline-flex w-fit rounded-xl border border-red-900/70 bg-red-950/30 px-4 py-2 text-sm font-medium text-red-100 transition hover:bg-red-950/50"
        >
          Ver agenda global
        </Link>
      </div>
      </div>

      <CaseFilesTable items={items} />
    </main>
  );
}