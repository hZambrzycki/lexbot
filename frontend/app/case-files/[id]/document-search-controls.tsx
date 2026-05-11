import type { RefObject } from "react";
import { useEffect, useMemo, useState, useSyncExternalStore } from "react";

import { SearchAutocomplete } from "@/app/components/search-autocomplete";
import { displayDslChipValue } from "@/lib/document-display";

import type { DocumentFilter, DocumentSort } from "./document-list-utils";
import { filterButtonClass } from "./document-list-utils";
import { useSearchAutocomplete } from "./use-search-autocomplete";

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

type DslChip = {
  key: string;
  label: string;
  token: string;
};

type SavedSearch = {
  label: string;
  query: string;
};

const RECENT_SEARCHES_KEY = "lexbox:document-search:recent";
const SAVED_SEARCHES_KEY = "lexbox:document-search:saved";
const SEARCH_STORAGE_EVENT = "lexbox:document-search-storage";

const MAX_RECENT_SEARCHES = 6;

const DSL_LABELS: Record<string, string> = {
  type: "Tipo",
  area: "Área",
  review: "Revisión",
  has: "Tiene",
  doc: "Formato",
};

const QUICK_SEARCHES: SavedSearch[] = [
  { label: "Pendientes", query: "review:pending_review" },
  { label: "PDF sin texto", query: "doc:pdf has:no_text" },
  { label: "Laboral con hitos", query: "area:labor has:events" },
  { label: "Sin clasificar", query: "type:unknown area:unknown" },
  { label: "Errores", query: "review:error" },
];

function readStoredSearches(key: string): SavedSearch[] {
  if (typeof window === "undefined") return [];

  try {
    const raw = window.localStorage.getItem(key);
    if (!raw) return [];

    const parsed = JSON.parse(raw);

    if (!Array.isArray(parsed)) return [];

    return parsed.filter(
      (item): item is SavedSearch =>
        typeof item?.label === "string" && typeof item?.query === "string",
    );
  } catch {
    return [];
  }
}

function writeStoredSearches(key: string, searches: SavedSearch[]) {
  if (typeof window === "undefined") return;

  window.localStorage.setItem(key, JSON.stringify(searches));
}

function emitStoredSearchesChange() {
  if (typeof window === "undefined") return;

  window.dispatchEvent(new Event(SEARCH_STORAGE_EVENT));
}

function subscribeStoredSearches(onStoreChange: () => void) {
  window.addEventListener("storage", onStoreChange);
  window.addEventListener(SEARCH_STORAGE_EVENT, onStoreChange);

  return () => {
    window.removeEventListener("storage", onStoreChange);
    window.removeEventListener(SEARCH_STORAGE_EVENT, onStoreChange);
  };
}

function getEmptySearchesSnapshot() {
  return "[]";
}

function useStoredSearches(key: string): SavedSearch[] {
  const snapshot = useSyncExternalStore(
    subscribeStoredSearches,
    () => JSON.stringify(readStoredSearches(key)),
    getEmptySearchesSnapshot,
  );

  return useMemo(() => {
    try {
      const parsed = JSON.parse(snapshot);

      if (!Array.isArray(parsed)) return [];

      return parsed.filter(
        (item): item is SavedSearch =>
          typeof item?.label === "string" && typeof item?.query === "string",
      );
    } catch {
      return [];
    }
  }, [snapshot]);
}

function buildRecentSearch(query: string): SavedSearch {
  return {
    label: query,
    query,
  };
}

function getDslChips(query: string): DslChip[] {
  return query
    .trim()
    .split(/\s+/)
    .map((token) => {
      const separatorIndex = token.indexOf(":");

      if (separatorIndex === -1) return null;

      const key = token.slice(0, separatorIndex);
      const value = token.slice(separatorIndex + 1);

      if (!key || !value) return null;

      const normalizedKey = key.toLowerCase();

      if (!DSL_LABELS[normalizedKey]) return null;

      return {
        key: normalizedKey,
        label: `${DSL_LABELS[normalizedKey]}: ${displayDslChipValue(
          normalizedKey,
          value,
        )}`,
        token,
      };
    })
    .filter((chip): chip is DslChip => chip !== null);
}

function removeTokenFromQuery(query: string, tokenToRemove: string): string {
  return query
    .trim()
    .split(/\s+/)
    .filter((token) => token !== tokenToRemove)
    .join(" ");
}

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

  const recentSearches = useStoredSearches(RECENT_SEARCHES_KEY);
  const savedSearches = useStoredSearches(SAVED_SEARCHES_KEY);

  const dslChips = useMemo(() => getDslChips(draftQuery), [draftQuery]);

  const {
    activeToken,
    suggestions,
    activeIndex,
    activeSuggestion,
    setActiveIndex,
    applySuggestion,
    handleKeyDown,
    closeSuggestions,
  } = useSearchAutocomplete({
    query: draftQuery,
    onQueryChange: setDraftQuery,
  });

  useEffect(() => {
    if (draftQuery === query) return;

    const timeout = window.setTimeout(() => {
      const normalizedQuery = draftQuery.trim();

      onQueryChange(draftQuery);

      if (normalizedQuery.length >= 3) {
        rememberRecentSearch(normalizedQuery);
      }
    }, 700);

    return () => window.clearTimeout(timeout);
  }, [draftQuery, onQueryChange, query]);

  function rememberRecentSearch(nextQuery: string) {
    const current = readStoredSearches(RECENT_SEARCHES_KEY);
    const nextSearch = buildRecentSearch(nextQuery);

    const next = [
      nextSearch,
      ...current.filter((item) => item.query !== nextQuery),
    ].slice(0, MAX_RECENT_SEARCHES);

    writeStoredSearches(RECENT_SEARCHES_KEY, next);
    emitStoredSearchesChange();
  }

  function clearQuery() {
    setDraftQuery("");
    onQueryChange("");
    closeSuggestions();
  }

  function removeChip(token: string) {
    const nextQuery = removeTokenFromQuery(draftQuery, token);

    setDraftQuery(nextQuery);
    onQueryChange(nextQuery);
    closeSuggestions();

    if (nextQuery.trim().length >= 3) {
      rememberRecentSearch(nextQuery.trim());
    }
  }

  function applySearch(nextQuery: string) {
    setDraftQuery(nextQuery);
    onQueryChange(nextQuery);
    closeSuggestions();

    if (nextQuery.trim().length >= 3) {
      rememberRecentSearch(nextQuery.trim());
    }
  }

  function saveCurrentSearch() {
    const normalizedQuery = draftQuery.trim();

    if (normalizedQuery.length < 3) return;

    const current = readStoredSearches(SAVED_SEARCHES_KEY);

    if (current.some((item) => item.query === normalizedQuery)) {
      return;
    }

    const next = [buildRecentSearch(normalizedQuery), ...current];

    writeStoredSearches(SAVED_SEARCHES_KEY, next);
    emitStoredSearchesChange();
  }

  function removeSavedSearch(queryToRemove: string) {
    const current = readStoredSearches(SAVED_SEARCHES_KEY);
    const next = current.filter((item) => item.query !== queryToRemove);

    writeStoredSearches(SAVED_SEARCHES_KEY, next);
    emitStoredSearchesChange();
  }

  function clearRecentSearches() {
    writeStoredSearches(RECENT_SEARCHES_KEY, []);
    emitStoredSearchesChange();
  }

  const lastToken = draftQuery.trim().split(/\s+/).pop() ?? "";

  const ghostSuffix =
    activeSuggestion && activeSuggestion.value !== lastToken
      ? activeSuggestion.value.slice(lastToken.length)
      : "";

  const canSaveCurrentSearch =
    draftQuery.trim().length >= 3 &&
    !savedSearches.some((item) => item.query === draftQuery.trim());

  return (
    <div className="flex flex-col gap-3">
      <div className="relative">
        <div className="relative">
          {activeSuggestion && ghostSuffix && draftQuery.trim().length > 0 ? (
            <div className="pointer-events-none absolute inset-0 flex items-center rounded-2xl px-4 py-3 text-sm">
              <span className="invisible whitespace-pre">{draftQuery}</span>
              <span className="text-neutral-600">{ghostSuffix}</span>
            </div>
          ) : null}

          <input
            ref={inputRef}
            value={draftQuery}
            onChange={(event) => setDraftQuery(event.target.value)}
            onKeyDown={(event) => {
              handleKeyDown(event);

              if (event.key === "Escape") {
                inputRef.current?.blur();
              }
            }}
            onBlur={() => {
              window.setTimeout(() => {
                closeSuggestions();
              }, 100);
            }}
            placeholder="Buscar: type:claim area:labor despido..."
            autoComplete="off"
            spellCheck={false}
            className="relative z-10 w-full rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 pr-11 text-sm text-neutral-100 outline-none placeholder:text-neutral-600 focus:border-red-900/70"
          />

          <SearchAutocomplete
            suggestions={suggestions}
            activeIndex={activeIndex}
            activeToken={activeToken}
            onHover={setActiveIndex}
            onSelect={applySuggestion}
          />
        </div>

        {draftQuery.trim().length > 0 ? (
          <button
            type="button"
            onClick={clearQuery}
            aria-label="Limpiar búsqueda"
            className="absolute right-3 top-1/2 z-40 flex h-6 w-6 -translate-y-1/2 items-center justify-center rounded-full border border-neutral-700 bg-neutral-900 text-xs text-neutral-400 transition hover:border-red-800 hover:bg-red-950/40 hover:text-red-100"
          >
            ×
          </button>
        ) : null}
      </div>

      {dslChips.length > 0 ? (
        <div className="flex flex-wrap gap-2">
          {dslChips.map((chip) => (
            <button
              key={`${chip.key}-${chip.token}`}
              type="button"
              onClick={() => removeChip(chip.token)}
              className="rounded-full border border-red-900/40 bg-red-950/20 px-3 py-1 text-xs font-medium text-red-100 transition hover:border-red-800 hover:bg-red-950/40"
              title="Quitar filtro"
            >
              {chip.label}
              <span className="ml-2 text-red-300/70">×</span>
            </button>
          ))}
        </div>
      ) : null}

      <div className="flex flex-wrap items-center gap-2">
        <span className="text-xs font-medium text-neutral-500">
          Búsquedas rápidas:
        </span>

        {QUICK_SEARCHES.map((item) => (
          <button
            key={item.query}
            type="button"
            onClick={() => applySearch(item.query)}
            className="rounded-full border border-neutral-800 bg-neutral-950 px-3 py-1 text-xs font-medium text-neutral-400 transition hover:border-red-900/60 hover:bg-red-950/20 hover:text-red-100"
          >
            {item.label}
          </button>
        ))}
      </div>

      {(savedSearches.length > 0 || canSaveCurrentSearch) && (
        <div className="flex flex-wrap items-center gap-2">
          <span className="text-xs font-medium text-neutral-500">
            Favoritas:
          </span>

          {canSaveCurrentSearch ? (
            <button
              type="button"
              onClick={saveCurrentSearch}
              className="rounded-full border border-amber-900/50 bg-amber-950/20 px-3 py-1 text-xs font-medium text-amber-100 transition hover:border-amber-700 hover:bg-amber-950/40"
            >
              Guardar búsqueda actual
            </button>
          ) : null}

          {savedSearches.map((item) => (
            <span
              key={item.query}
              className="inline-flex items-center overflow-hidden rounded-full border border-neutral-800 bg-neutral-950 text-xs font-medium text-neutral-300"
            >
              <button
                type="button"
                onClick={() => applySearch(item.query)}
                className="px-3 py-1 transition hover:bg-neutral-900 hover:text-neutral-100"
                title={item.query}
              >
                {item.label}
              </button>

              <button
                type="button"
                onClick={() => removeSavedSearch(item.query)}
                className="border-l border-neutral-800 px-2 py-1 text-neutral-500 transition hover:bg-red-950/40 hover:text-red-100"
                aria-label="Eliminar búsqueda favorita"
              >
                ×
              </button>
            </span>
          ))}
        </div>
      )}

      {recentSearches.length > 0 ? (
        <div className="flex flex-wrap items-center gap-2">
          <span className="text-xs font-medium text-neutral-500">
            Recientes:
          </span>

          {recentSearches.map((item) => (
            <button
              key={item.query}
              type="button"
              onClick={() => applySearch(item.query)}
              className="max-w-xs truncate rounded-full border border-neutral-800 bg-neutral-950 px-3 py-1 text-xs font-medium text-neutral-500 transition hover:border-neutral-700 hover:bg-neutral-900 hover:text-neutral-200"
              title={item.query}
            >
              {item.label}
            </button>
          ))}

          <button
            type="button"
            onClick={clearRecentSearches}
            className="rounded-full px-2 py-1 text-xs text-neutral-600 transition hover:text-red-200"
          >
            limpiar
          </button>
        </div>
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
  const [searchResetKey, setSearchResetKey] = useState(0);

  const hasActiveControls =
    activeFilter !== "all" || query.trim().length > 0 || sort !== "recent";

  function handleClearControls() {
    setSearchResetKey((prev) => prev + 1);
    onClear();
  }

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

      <div className="grid gap-2 lg:grid-cols-[1fr_auto_auto] lg:items-start">
        <DebouncedSearchInput
          key={searchResetKey}
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
            onClick={handleClearControls}
            className="whitespace-nowrap rounded-2xl border border-neutral-800 bg-neutral-950 px-4 py-3 text-sm font-medium text-neutral-400 transition hover:border-neutral-700 hover:bg-neutral-900 hover:text-neutral-100"
          >
            {query.trim() ? "Limpiar búsqueda" : "Limpiar filtros"}
          </button>
        ) : null}
      </div>
    </div>
  );
}