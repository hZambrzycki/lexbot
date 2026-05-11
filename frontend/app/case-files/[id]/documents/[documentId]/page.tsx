import Link from "next/link";
import { getDocument } from "@/lib/api";
import { DeleteDocumentButton } from "./delete-document-button";
import { ReprocessDocumentButton } from "./reprocess-document-button";
import { HighlightedText } from "@/app/components/highlighted-text";
import { eventRelationLabel } from "@/lib/event-relations";
import { DocumentReviewBadge } from "@/app/components/document-review-badge";
import { DocumentReviewActions } from "@/app/components/document-review-actions";
import {
  displayDocumentType,
  displayLegalArea,
  displayMimeType,
} from "@/lib/document-display";

type Props = {
  params: Promise<{
    id: string;
    documentId: string;
  }>;
  searchParams?: Promise<{
    q?: string;
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

export default async function DocumentDetailPage({
  params,
  searchParams,
}: Props) {
  const { id, documentId } = await params;
  const resolvedSearchParams = searchParams ? await searchParams : {};
  const query = resolvedSearchParams.q?.trim() ?? "";

  const detail = await getDocument(documentId);
  const doc = detail.document;

  return (
    <main className="space-y-8">
      <section className="space-y-4">
        <Link
          href={`/case-files/${id}?tab=documentos${
            query ? `&q=${encodeURIComponent(query)}` : ""
          }`}
          className="text-sm text-neutral-400 underline-offset-4 hover:underline"
        >
          ← Volver al expediente
        </Link>

        <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight text-neutral-50">
              {doc.original_name}
            </h1>

            <p className="mt-2 text-sm text-neutral-400">
              Documento del expediente · {displayMimeType(doc.mime_type)}
            </p>
          </div>

          <DocumentReviewBadge status={doc.review_status} />
        </div>

        <div className="flex flex-wrap gap-2">
          <span className="rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
            Archivo físico:{" "}
            {detail.file_exists ? "localizado" : "no localizado"}
          </span>

          <span className="rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
            Texto extraído: {detail.has_extracted_text ? "sí" : "no"}
          </span>

          <span className="rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
            Eventos detectados: {detail.events.length}
          </span>

          {doc.reviewed_at ? (
            <span className="rounded-full border border-neutral-700 bg-neutral-900 px-3 py-1 text-xs font-medium text-neutral-300">
              Revisado: {new Date(doc.reviewed_at).toLocaleString("es-ES")}
            </span>
          ) : null}
        </div>

        {query ? (
          <div className="rounded-2xl border border-yellow-900/50 bg-yellow-950/20 p-4">
            <p className="text-sm font-medium text-yellow-100">
              Búsqueda activa: “{query}”
            </p>
            <p className="mt-1 text-xs text-yellow-200/70">
              Este documento se abrió desde una coincidencia encontrada en el
              contenido.
            </p>
          </div>
        ) : null}

        {doc.review_note ? (
          <div className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-4">
            <p className="text-xs font-medium uppercase tracking-wide text-neutral-500">
              Nota de revisión documental
            </p>

            <p className="mt-2 text-sm leading-6 text-neutral-300">
              {doc.review_note}
            </p>
          </div>
        ) : null}

        <div className="grid gap-4 lg:grid-cols-2">
          <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
            <h2 className="text-base font-semibold text-neutral-50">
              Acciones del documento
            </h2>

            <p className="mt-1 text-sm leading-6 text-neutral-400">
              Reprocesa el documento si has mejorado el extractor o elimina el
              archivo si fue subido por error.
            </p>

            <div className="mt-4 flex flex-wrap gap-3">
              <ReprocessDocumentButton documentId={doc.id} />
              <DeleteDocumentButton documentId={doc.id} caseFileId={id} />
            </div>
          </section>

          <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div>
                <h2 className="text-base font-semibold text-neutral-50">
                  Revisión documental
                </h2>

                <p className="mt-1 text-sm leading-6 text-neutral-400">
                  Marca el documento como revisado, pendiente o problemático.
                </p>
              </div>

              <DocumentReviewBadge status={doc.review_status} />
            </div>

            <div className="mt-4">
              <DocumentReviewActions
                documentId={doc.id}
                reviewStatus={doc.review_status}
                reviewNote={doc.review_note}
              />
            </div>
          </section>
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
          <div className="mt-3 rounded-xl border border-neutral-800 bg-black/20 p-4">
            <p className="text-sm text-neutral-300">
              No se han detectado hitos procesales en este documento.
            </p>

            <p className="mt-2 text-sm leading-6 text-neutral-500">
              Esto puede ser normal si el documento es un CV, justificante,
              contrato, nómina, anexo o documento auxiliar. Solo requiere
              revisión si esperabas encontrar un plazo, una vista, una
              notificación o un requerimiento.
            </p>
          </div>
        ) : (
          <ul className="mt-4 space-y-3">
            {detail.events.map((event) => {
              const relation = eventRelationLabel(event);

              return (
                <li
                  key={event.event_id}
                  className="rounded-xl border border-neutral-800 bg-neutral-900 p-4"
                >
                  <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                    <div>
                      <p className="font-medium text-neutral-100">
                        {displayEventType(event.event_type)} ·{" "}
                        {event.event_date}
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

                  {relation ? (
                    <div className="mt-3 rounded-xl border border-sky-900/60 bg-sky-950/20 p-3">
                      <p className="text-xs font-medium uppercase tracking-wide text-sky-300/80">
                        Relación procesal
                      </p>

                      <p className="mt-1 text-sm font-medium text-sky-100">
                        ↳ {relation}
                      </p>
                    </div>
                  ) : null}

                  <div className="mt-3 grid gap-2 text-xs text-neutral-500 md:grid-cols-2">
                    <div>Tipo de fecha: {displayDateKind(event.date_kind)}</div>
                    <div>Fecha base: {event.anchor_date || "—"}</div>
                    <div>Fuente: {event.anchor_source || "—"}</div>
                    <div>Días: {event.relative_days ?? 0}</div>
                    <div>
                      Cómputo:{" "}
                      {event.is_business_days
                        ? "días hábiles"
                        : "días naturales"}
                    </div>
                    <div>
                      Día adicional: {event.add_extra_day ? "sí" : "no"}
                    </div>
                    <div className="md:col-span-2">
                      Regla: {event.computation || "—"}
                    </div>
                  </div>
                </li>
              );
            })}
          </ul>
        )}
      </section>

      <details
        open={Boolean(query)}
        className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5"
      >
        <summary className="cursor-pointer text-xl font-semibold text-neutral-50">
          Texto extraído completo
        </summary>

        {!detail.has_extracted_text ? (
          <p className="mt-3 text-sm text-neutral-400">
            Este documento todavía no tiene texto extraído.
          </p>
        ) : (
          <pre className="mt-4 max-h-[560px] overflow-auto whitespace-pre-wrap rounded-xl border border-neutral-800 bg-black/30 p-4 text-sm leading-6 text-neutral-300">
            <HighlightedText text={detail.extracted_text} query={query} />
          </pre>
        )}
      </details>
    </main>
  );
}