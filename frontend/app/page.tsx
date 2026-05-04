import Link from "next/link";

export default function HomePage() {
  return (
    <main className="space-y-8">
      <section className="space-y-3">
        <p className="text-sm font-medium uppercase tracking-[0.3em] text-neutral-500">
          LEXBOX
        </p>

        <div className="max-w-3xl space-y-3">
          <h1 className="text-4xl font-bold tracking-tight text-neutral-50">
            Panel local de gestión procesal
          </h1>

          <p className="text-base leading-7 text-neutral-400">
            Controla expedientes, documentos, notas, plazos y eventos
            procesales detectados automáticamente desde una única vista
            local-first.
          </p>
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-2">
        <Link
          href="/case-files"
          className="group rounded-2xl border border-neutral-800 bg-neutral-900 p-6 transition hover:border-neutral-600 hover:bg-neutral-850"
        >
          <div className="space-y-3">
            <div className="text-sm font-medium text-neutral-400">
              Expedientes
            </div>

            <h2 className="text-2xl font-semibold text-neutral-50">
              Ver expedientes
            </h2>

            <p className="text-sm leading-6 text-neutral-400">
              Accede al listado de asuntos, revisa documentos, notas, alertas
              procesales y estado de cada expediente.
            </p>

            <div className="pt-2 text-sm font-medium text-neutral-200 group-hover:underline">
              Abrir expedientes →
            </div>
          </div>
        </Link>

        <Link
          href="/events"
          className="group rounded-2xl border border-red-900/50 bg-red-950/20 p-6 transition hover:border-red-700/70 hover:bg-red-950/30"
        >
          <div className="space-y-3">
            <div className="text-sm font-medium text-red-300/80">
              Agenda procesal
            </div>

            <h2 className="text-2xl font-semibold text-red-100">
              Ver agenda global
            </h2>

            <p className="text-sm leading-6 text-red-100/75">
              Consulta todos los plazos, vistas, comparecencias,
              requerimientos y notificaciones detectados en todos los
              expedientes.
            </p>

            <div className="pt-2 text-sm font-medium text-red-100 group-hover:underline">
              Abrir agenda →
            </div>
          </div>
        </Link>
      </section>

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <h2 className="text-lg font-semibold text-neutral-100">
          Flujo recomendado
        </h2>

        <div className="mt-4 grid gap-3 md:grid-cols-3">
          <div className="rounded-xl border border-neutral-800 bg-neutral-900 p-4">
            <div className="text-sm font-medium text-neutral-200">
              1. Importar documentos
            </div>
            <p className="mt-2 text-sm leading-6 text-neutral-400">
              El sistema extrae texto, clasifica metadatos y detecta hitos
              procesales.
            </p>
          </div>

          <div className="rounded-xl border border-neutral-800 bg-neutral-900 p-4">
            <div className="text-sm font-medium text-neutral-200">
              2. Revisar eventos
            </div>
            <p className="mt-2 text-sm leading-6 text-neutral-400">
              Los plazos vencidos y próximos aparecen priorizados según riesgo
              procesal.
            </p>
          </div>

          <div className="rounded-xl border border-neutral-800 bg-neutral-900 p-4">
            <div className="text-sm font-medium text-neutral-200">
              3. Resolver y exportar
            </div>
            <p className="mt-2 text-sm leading-6 text-neutral-400">
              Marca eventos como revisados o resueltos y exporta la agenda en
              formato ICS.
            </p>
          </div>
        </div>
      </section>
    </main>
  );
}