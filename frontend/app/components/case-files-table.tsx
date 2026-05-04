import Link from "next/link";
import { CaseFile } from "@/lib/types";

type Props = {
  items: CaseFile[];
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
    case "criminal":
      return "Penal";
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

function statusClass(value: string) {
  switch (value) {
    case "open":
      return "border-emerald-800 bg-emerald-950/50 text-emerald-200";
    case "closed":
      return "border-blue-800 bg-blue-950/50 text-blue-200";
    case "archived":
      return "border-neutral-700 bg-neutral-800 text-neutral-300";
    default:
      return "border-neutral-700 bg-neutral-800 text-neutral-300";
  }
}

export function CaseFilesTable({ items }: Props) {
  if (items.length === 0) {
    return (
      <div className="rounded-2xl border border-neutral-800 bg-neutral-950 p-6 text-sm text-neutral-400">
        No hay expedientes registrados.
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-neutral-800 bg-neutral-900">
      <table className="min-w-full text-sm">
        <thead className="bg-neutral-800/70 text-neutral-300">
          <tr>
            <th className="px-4 py-3 text-left">Referencia</th>
            <th className="px-4 py-3 text-left">Expediente</th>
            <th className="px-4 py-3 text-left">Tipo</th>
            <th className="px-4 py-3 text-left">Estado</th>
            <th className="px-4 py-3 text-left">Calendario</th>
            <th className="px-4 py-3 text-right">Acción</th>
          </tr>
        </thead>

        <tbody>
          {items.map((item) => (
            <tr
              key={item.id}
              className="border-t border-neutral-800 transition hover:bg-neutral-800/40"
            >
              <td className="px-4 py-4 align-top">
                <Link
                  href={`/case-files/${item.id}`}
                  className="font-semibold text-white underline-offset-4 hover:underline"
                >
                  {item.reference}
                </Link>
              </td>

              <td className="px-4 py-4 align-top">
                <div className="font-medium text-neutral-100">
                  {item.title}
                </div>

                {item.description ? (
                  <div className="mt-1 max-w-md truncate text-xs text-neutral-500">
                    {item.description}
                  </div>
                ) : null}
              </td>

              <td className="px-4 py-4 align-top text-neutral-300">
                {displayCaseType(item.type)}
              </td>

              <td className="px-4 py-4 align-top">
                <span
                  className={`inline-flex rounded-full border px-2.5 py-1 text-xs font-medium ${statusClass(
                    item.status,
                  )}`}
                >
                  {displayStatus(item.status)}
                </span>
              </td>

              <td className="px-4 py-4 align-top text-neutral-300">
                <div>{item.calendar_scope || "madrid"}</div>

                {item.august_non_business ? (
                  <div className="mt-1 text-xs text-neutral-500">
                    Agosto inhábil
                  </div>
                ) : null}
              </td>

              <td className="px-4 py-4 text-right align-top">
                <Link
                  href={`/case-files/${item.id}`}
                  className="inline-flex rounded-xl border border-neutral-700 px-3 py-1.5 text-xs font-medium text-neutral-200 transition hover:bg-neutral-800"
                >
                  Abrir
                </Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}