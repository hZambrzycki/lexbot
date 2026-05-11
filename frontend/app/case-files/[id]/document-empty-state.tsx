type Props = {
  canSearchContent: boolean;
  contentResultsCount: number;
};

export function DocumentEmptyState({
  canSearchContent,
  contentResultsCount,
}: Props) {
  return (
    <p className="mt-4 rounded-xl border border-neutral-800 bg-neutral-900 p-4 text-sm text-neutral-400">
      {canSearchContent && contentResultsCount > 0
        ? "No hay coincidencias por nombre, tipo o estado, pero sí hay resultados dentro del contenido."
        : "No hay documentos para este filtro, búsqueda u ordenación."}
    </p>
  );
}