import Link from "next/link";
import { getDocument } from "@/lib/api";
import { DeleteDocumentButton } from "./delete-document-button";
import { ReprocessDocumentButton } from "./reprocess-document-button";


type Props = {
  params: Promise<{
    id: string;
    documentId: string;
  }>;
};

function displayEventType(value: string) {
  switch (value) {
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
      return value || "Evento";
  }
}

function displayReviewStatus(value: string) {
  switch (value) {
    case "pending":
      return "Pendiente";
    case "reviewed":
      return "Revisado";
    case "resolved":
      return "Resuelto";
    default:
      return value || "Sin estado";
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
      return value || "Sin tipo";
  }
}

function displayDocumentType(value: string) {
  switch (value) {
    case "order":
      return "Resolución / acto procesal";
    case "claim":
      return "Demanda / reclamación";
    case "contract":
      return "Contrato";
    case "notice":
      return "Notificación";
    case "evidence":
      return "Documento probatorio";
    case "unknown":
    case "":
      return "Sin clasificar";
    default:
      return value;
  }
}

function displayLegalArea(value: string) {
  switch (value) {
    case "procedural":
      return "Procesal";
    case "civil":
      return "Civil";
    case "labor":
    case "laboral":
      return "Laboral";
    case "administrative":
    case "administrativo":
      return "Administrativo";
    case "immigration":
    case "extranjeria":
      return "Extranjería";
    case "commercial":
    case "mercantil":
      return "Mercantil";
    case "unknown":
    case "":
      return "Sin clasificar";
    default:
      return value;
  }
}

function displayDateKind(value?: string) {
  switch (value) {
    case "absolute":
      return "Fecha absoluta";
    case "relative":
      return "Fecha relativa";
    default:
      return "—";
  }
}

export default async function DocumentDetailPage({ params }: Props) {
  const { id, documentId } = await params;
  const detail = await getDocument(documentId);
  const doc = detail.document;

  return (
    <main className="space-y-8">
      <section className="space-y-4">
        <Link
          href={`/case-files/${id}`}
          className="text-sm text-neutral-400 underline-offset-4 hover:underline"
        >
          ← Volver al expediente
        </Link>

        <div>
          <h1 className="text-3xl font-bold tracking-tight text-neutral-50">
            {doc.original_name}
          </h1>

          <p className="mt-2 text-sm text-neutral-400">
            Documento del expediente · {displayMimeType(doc.mime_type)}
          </p>
        </div>

        <div className="flex flex-wrap gap-2">
          <span className="rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
            Archivo físico: {detail.file_exists ? "localizado" : "no localizado"}
          </span>

          <span className="rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
            Texto extraído: {detail.has_extracted_text ? "sí" : "no"}
          </span>

          <span className="rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
            Eventos detectados: {detail.events.length}
          </span>
        </div>
        
        <div className="flex gap-3">
        <DeleteDocumentButton documentId={doc.id} caseFileId={id} />
        <ReprocessDocumentButton documentId={doc.id} />
        </div>
      </section>

      <section className="grid gap-5 md:grid-cols-3">
        <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5">
          <p className="text-sm text-neutral-400">Longitud del texto</p>
          <p className="mt-2 text-2xl font-semibold text-neutral-50">
            {detail.extracted_text_length}
          </p>
          <p className="mt-1 text-xs text-neutral-500">caracteres extraídos</p>
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5">
          <p className="text-sm text-neutral-400">Tipo documental</p>
          <p className="mt-2 text-lg font-semibold text-neutral-50">
            {detail.has_metadata
              ? displayDocumentType(detail.document_type)
              : "Sin clasificar"}
          </p>
        </div>

        <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-5">
          <p className="text-sm text-neutral-400">Área jurídica</p>
          <p className="mt-2 text-lg font-semibold text-neutral-50">
            {detail.has_metadata
              ? displayLegalArea(detail.legal_area)
              : "Sin clasificar"}
          </p>
        </div>
      </section>

      {detail.extracted_text_preview ? (
        <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
          <h2 className="text-xl font-semibold text-neutral-50">
            Vista rápida del texto
          </h2>

          <p className="mt-3 text-sm leading-6 text-neutral-300">
            {detail.extracted_text_preview}
          </p>
        </section>
      ) : null}

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <h2 className="text-xl font-semibold text-neutral-50">
          Eventos detectados en este documento
        </h2>

        {detail.events.length === 0 ? (
          <p className="mt-3 text-sm text-neutral-400">
            No se han detectado hitos procesales en este documento.
          </p>
        ) : (
          <ul className="mt-4 space-y-3">
            {detail.events.map((event) => (
              <li
                key={event.event_id}
                className="rounded-xl border border-neutral-800 bg-neutral-900 p-4"
              >
                <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                  <div>
                    <p className="font-medium text-neutral-100">
                      {displayEventType(event.event_type)} · {event.event_date}
                    </p>

                    <p className="mt-1 text-xs text-neutral-500">
                      {displayReviewStatus(event.review_status)}
                    </p>
                  </div>

                  <Link
                    href={`/events/${event.event_id}`}
                    className="inline-flex w-fit rounded-xl border border-neutral-700 bg-neutral-950 px-3 py-2 text-xs font-medium text-neutral-100 transition hover:bg-neutral-800"
                  >
                    Ver hito
                  </Link>
                </div>

                <p className="mt-3 text-sm leading-6 text-neutral-300">
                  “{event.source_text}”
                </p>

                <div className="mt-3 grid gap-2 text-xs text-neutral-500 md:grid-cols-2">
                  <div>Tipo de fecha: {displayDateKind(event.date_kind)}</div>
                  <div>Fecha base: {event.anchor_date || "—"}</div>
                  <div>Fuente: {event.anchor_source || "—"}</div>
                  <div>Días: {event.relative_days ?? 0}</div>
                  <div>
                    Cómputo:{" "}
                    {event.is_business_days ? "días hábiles" : "días naturales"}
                  </div>
                  <div>
                    Día adicional: {event.add_extra_day ? "sí" : "no"}
                  </div>
                  <div className="md:col-span-2">
                    Regla: {event.computation || "—"}
                  </div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>

      <details className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
        <summary className="cursor-pointer text-xl font-semibold text-neutral-50">
          Texto extraído completo
        </summary>

        {!detail.has_extracted_text ? (
          <p className="mt-3 text-sm text-neutral-400">
            Este documento todavía no tiene texto extraído.
          </p>
        ) : (
          <pre className="mt-4 max-h-[560px] overflow-auto whitespace-pre-wrap rounded-xl border border-neutral-800 bg-black/30 p-4 text-sm leading-6 text-neutral-300">
            {detail.extracted_text}
          </pre>
        )}
      </details>
    </main>
  );
}