"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { createCaseFile } from "@/lib/api";

type Props = {
  suggestedReference: string;
};

const TEMP_CLIENT_ID = "5b0222e8-e3ed-45d4-a0d4-12cb1a672db4";

export function NewCaseFileForm({ suggestedReference }: Props) {
  const router = useRouter();

  const [reference, setReference] = useState(suggestedReference);
  const [title, setTitle] = useState("");
  const [type, setType] = useState("otros");
  const [description, setDescription] = useState("");
  const [calendarScope, setCalendarScope] = useState("madrid");
  const [augustNonBusiness, setAugustNonBusiness] = useState(true);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setIsSubmitting(true);
    setError(null);

    try {
      const created = await createCaseFile({
        client_id: TEMP_CLIENT_ID,
        reference: reference.trim(),
        title: title.trim(),
        type,
        description: description.trim(),
        calendar_scope: calendarScope,
        august_non_business: augustNonBusiness,
      });

      router.push(`/case-files/${created.id}`);
      router.refresh();
    } catch (error) {
      console.error(error);
      setError(
        error instanceof Error
          ? error.message
          : "No se pudo crear el expediente.",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="mx-auto max-w-3xl space-y-6 p-6">
      <Link
        href="/case-files"
        className="text-sm text-neutral-400 underline-offset-4 hover:text-neutral-100 hover:underline"
      >
        ← Volver a expedientes
      </Link>

      <section className="rounded-2xl border border-neutral-800 bg-neutral-950 p-6">
        <p className="text-sm text-neutral-500">Nuevo expediente</p>

        <h1 className="mt-2 text-3xl font-bold text-neutral-50">
          Crear expediente
        </h1>

        <p className="mt-2 text-sm leading-6 text-neutral-400">
          Crea un expediente base para poder asociar documentos, eventos
          procesales y agenda.
        </p>

        <form onSubmit={handleSubmit} className="mt-6 space-y-5">
          <label className="block space-y-1">
            <span className="text-sm font-medium text-neutral-200">
              Referencia
            </span>
            <input
              value={reference}
              onChange={(event) => setReference(event.target.value)}
              required
              className="w-full rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-2 text-sm text-neutral-100 outline-none transition focus:border-neutral-400"
            />
          </label>

          <label className="block space-y-1">
            <span className="text-sm font-medium text-neutral-200">
              Título
            </span>
            <input
              value={title}
              onChange={(event) => setTitle(event.target.value)}
              required
              placeholder="Ej.: Despido disciplinario"
              className="w-full rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-2 text-sm text-neutral-100 outline-none transition placeholder:text-neutral-600 focus:border-neutral-400"
            />
          </label>

          <label className="block space-y-1">
            <span className="text-sm font-medium text-neutral-200">
              Tipo
            </span>
            <select
              value={type}
              onChange={(event) => setType(event.target.value)}
              className="w-full rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-2 text-sm text-neutral-100 outline-none transition focus:border-neutral-400"
            >
              <option value="civil">Civil</option>
              <option value="laboral">Laboral</option>
              <option value="extranjeria">Extranjería</option>
              <option value="mercantil">Mercantil</option>
              <option value="administrativo">Administrativo</option>
              <option value="otros">Otros</option>
            </select>
          </label>

          <label className="block space-y-1">
            <span className="text-sm font-medium text-neutral-200">
              Descripción
            </span>
            <textarea
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              rows={4}
              placeholder="Resumen breve del asunto..."
              className="w-full rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-2 text-sm text-neutral-100 outline-none transition placeholder:text-neutral-600 focus:border-neutral-400"
            />
          </label>

          <div className="grid gap-4 md:grid-cols-2">
            <label className="block space-y-1">
              <span className="text-sm font-medium text-neutral-200">
                Calendario
              </span>
              <select
                value={calendarScope}
                onChange={(event) => setCalendarScope(event.target.value)}
                className="w-full rounded-xl border border-neutral-700 bg-neutral-900 px-3 py-2 text-sm text-neutral-100 outline-none transition focus:border-neutral-400"
              >
                <option value="madrid">Madrid</option>
                <option value="state">Estatal</option>
              </select>
            </label>

            <label className="flex items-center gap-3 rounded-xl border border-neutral-800 bg-neutral-900 px-3 py-2">
              <input
                type="checkbox"
                checked={augustNonBusiness}
                onChange={(event) =>
                  setAugustNonBusiness(event.target.checked)
                }
                className="h-4 w-4"
              />
              <span className="text-sm text-neutral-200">
                Considerar agosto inhábil
              </span>
            </label>
          </div>

          {error ? (
            <div className="rounded-xl border border-red-900/70 bg-red-950/30 px-3 py-2 text-sm text-red-100">
              {error}
            </div>
          ) : null}

          <div className="flex flex-wrap gap-2">
            <button
              type="submit"
              disabled={isSubmitting}
              className="inline-flex rounded-xl border border-neutral-100 bg-neutral-100 px-4 py-2 text-sm font-medium text-neutral-950 transition hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
            >
              {isSubmitting ? "Creando..." : "Crear expediente"}
            </button>

            <Link
              href="/case-files"
              className="inline-flex rounded-xl border border-neutral-700 bg-neutral-900 px-4 py-2 text-sm font-medium text-neutral-100 transition hover:bg-neutral-800"
            >
              Cancelar
            </Link>
          </div>
        </form>
      </section>
    </main>
  );
}