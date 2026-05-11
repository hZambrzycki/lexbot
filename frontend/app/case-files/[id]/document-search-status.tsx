type Props = {
  visibleCount: number;
  totalCount: number;
  canSearchContent: boolean;
  contentSearchLoading: boolean;
  contentResultsCount: number;
};

export function DocumentSearchStatus({
  visibleCount,
  totalCount,
  canSearchContent,
  contentSearchLoading,
  contentResultsCount,
}: Props) {
  return (
    <>
      <p className="mt-4 text-xs text-neutral-500">
        Mostrando {visibleCount} de {totalCount} documentos por filtros
        visibles.
      </p>

      {canSearchContent && contentSearchLoading ? (
        <p className="mt-2 text-xs text-neutral-500">
          Buscando dentro del texto extraído...
        </p>
      ) : null}

      {canSearchContent && contentResultsCount > 0 ? (
        <p className="mt-2 text-xs text-yellow-200">
          {contentResultsCount === 1
            ? "1 documento contiene coincidencias dentro del texto extraído."
            : `${contentResultsCount} documentos contienen coincidencias dentro del texto extraído.`}
        </p>
      ) : null}
    </>
  );
}