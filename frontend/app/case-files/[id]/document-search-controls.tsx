import type { RefObject } from "react";
import { useEffect, useState } from "react";
import type { DocumentFilter, DocumentSort } from "./document-list-utils";
import { filterButtonClass } from "./document-list-utils";

type FilterItem = {
  id: DocumentFilter;
  label: string;
  count: number;
};

type Props = {
  filters: FilterItem[];
  activeFilter: DocumentFilter;
  query: string;
  sort: DocumentSort;
  inputRef: RefObject<HTMLInputElement | null>;
  onFilterChange: (filter: DocumentFilter) => void;
  onQueryChange: (query: string) => void;
  onSortChange: (sort: DocumentSort) => void;
  onClear: () => void;
};

function DebouncedSearchInput({
  query,
  inputRef,
  onQueryChange,
}: {
  query: string;
  inputRef: RefObject<HTMLInputElement | null>;
  onQueryChange: (query: string) => void;
}) {
  const [draftQuery, setDraftQuery] = useState(query);

  useEffect(() => {
    if (draftQuery === query) return;

    const timeout = window.setTimeout(() => {
      onQueryChange(draftQuery);
    }, 700);

    return () => window.clearTimeout(timeout);
  }, [draftQuery, onQueryChange, query]);

  function clearQuery() {
    setDraftQuery("");
    onQueryChange("");
  }

  return (
    <div className="relative">
      <input
        ref={inputRef}
        value={draftQuery}
        onChange={(event) => setDraftQuery(event.target.value)}
        placeholder="Buscar por nombre, tipo, estado o contenido..."
        className="w-full rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 pr-11 text-sm text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-900/70"
      />

      {draftQuery.trim().length > 0 ? (
        <button
          type="button"
          onClick={clearQuery}
          aria-label="Limpiar búsqueda"
          className="absolute right-3 top-1/2 flex h-6 w-6 -translate-y-1/2 items-center justify-center rounded-full border border-neutral-700 bg-neutral-900 text-xs text-neutral-400 transition hover:border-red-800 hover:bg-red-950/40 hover:text-red-100"
        >
          ×
        </button>
      ) : null}
    </div>
  );
}

export function DocumentSearchControls({
  filters,
  activeFilter,
  query,
  sort,
  inputRef,
  onFilterChange,
  onQueryChange,
  onSortChange,
  onClear,
}: Props) {
  const hasActiveControls =
    activeFilter !== "all" || query.trim().length > 0 || sort !== "recent";

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-2">
        {filters.map((item) => (
          <button
            key={item.id}
            type="button"
            onClick={() => onFilterChange(item.id)}
            className={filterButtonClass(activeFilter === item.id)}
          >
            {item.label}
            <span className="ml-1 text-neutral-500">{item.count}</span>
          </button>
        ))}
      </div>

      <div className="grid gap-2 lg:grid-cols-[1fr_auto_auto] lg:items-center">
        <DebouncedSearchInput
          key={query}
          query={query}
          inputRef={inputRef}
          onQueryChange={onQueryChange}
        />

        <select
          value={sort}
          onChange={(event) => onSortChange(event.target.value as DocumentSort)}
          className="rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm text-neutral-300 outline-none focus:border-red-900/70"
        >
          <option value="recent">Más recientes</option>
          <option value="name">Nombre A-Z</option>
          <option value="events">Más hitos</option>
          <option value="pending">Pendientes primero</option>
          <option value="errors">Errores primero</option>
        </select>

        {hasActiveControls ? (
          <button
            type="button"
            onClick={onClear}
            className="whitespace-nowrap rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm font-medium text-neutral-400 transition hover:border-neutral-700 hover:bg-neutral-900 hover:text-neutral-100"
          >
            {query.trim() ? "Limpiar búsqueda" : "Limpiar filtros"}
          </button>
        ) : null}
      </div>
    </div>
  );
}