"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";

import { globalSearch } from "@/lib/api";
import type { GlobalSearchResult, GlobalSearchResultType } from "@/lib/types";

const GROUP_LABELS: Record<GlobalSearchResultType, string> = {
  case_file: "Expedientes",
  document: "Documentos",
  event: "Eventos y plazos",
  note: "Notas",
  action: "Acciones rápidas",
  navigation: "Navegación",
};

function renderSnippet(snippet: string) {
  const parts = snippet.split(/(\[[^\]]+\])/g);

  return parts.map((part, index) => {
    const isHighlighted = part.startsWith("[") && part.endsWith("]");

    if (!isHighlighted) {
      return <span key={`${part}-${index}`}>{part}</span>;
    }

    return (
      <mark
        key={`${part}-${index}`}
        className="rounded bg-yellow-400/10 px-1 font-medium text-yellow-100"
      >
        {part.slice(1, -1)}
      </mark>
    );
  });
}

function groupResults(results: GlobalSearchResult[]) {
  const groups: {
    type: GlobalSearchResultType;
    label: string;
    items: { result: GlobalSearchResult; index: number }[];
  }[] = [];

  for (const [index, result] of results.entries()) {
    let group = groups.find((item) => item.type === result.type);

    if (!group) {
      group = {
        type: result.type,
        label: GROUP_LABELS[result.type],
        items: [],
      };

      groups.push(group);
    }

    group.items.push({ result, index });
  }

  return groups;
}

export function CommandPalette() {
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement | null>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<GlobalSearchResult[]>([]);
  const [activeIndex, setActiveIndex] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    function handleShortcut(event: KeyboardEvent) {
      const isMac = navigator.platform.toLowerCase().includes("mac");
      const modifierPressed = isMac ? event.metaKey : event.ctrlKey;

      if (modifierPressed && event.key.toLowerCase() === "k") {
        event.preventDefault();
        setIsOpen(true);
      }
    }

    window.addEventListener("keydown", handleShortcut);

    return () => {
      window.removeEventListener("keydown", handleShortcut);
    };
  }, []);

  useEffect(() => {
    if (!isOpen) return;

    const timeout = window.setTimeout(() => {
      inputRef.current?.focus();
    }, 0);

    return () => window.clearTimeout(timeout);
  }, [isOpen]);

  useEffect(() => {
    const trimmedQuery = query.trim();

    if (!isOpen || trimmedQuery.length < 2) {
      return;
    }

    let cancelled = false;

    const timeout = window.setTimeout(async () => {
      setIsLoading(true);

      try {
        const nextResults = await globalSearch(trimmedQuery);

        if (!cancelled) {
          setResults(nextResults);
          setActiveIndex(0);
        }
      } catch {
        if (!cancelled) {
          setResults([]);
          setActiveIndex(0);
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }, 250);

    return () => {
      cancelled = true;
      window.clearTimeout(timeout);
    };
  }, [isOpen, query]);

  function closePalette() {
    setIsOpen(false);
    setQuery("");
    setResults([]);
    setActiveIndex(0);
  }

  const trimmedQuery = query.trim();
  const canSearch = trimmedQuery.length >= 2;
  const visibleResults = canSearch ? results : [];
  const groupedResults = groupResults(visibleResults);
  const activeResult = visibleResults[activeIndex] ?? null;

  function openActiveResult() {
    if (!activeResult) return;

    router.push(activeResult.href);
    closePalette();
  }

  if (!isOpen) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center bg-black/60 px-4 pt-24 backdrop-blur-sm">
      <div className="w-full max-w-2xl overflow-hidden rounded-2xl border border-neutral-800 bg-neutral-950 shadow-2xl">
        <div className="border-b border-neutral-800 p-3">
          <input
            ref={inputRef}
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Escape") {
                event.preventDefault();
                closePalette();
                return;
              }

              if (event.key === "ArrowDown") {
                event.preventDefault();

                setActiveIndex((prev) =>
                  visibleResults.length === 0
                    ? 0
                    : prev === visibleResults.length - 1
                      ? 0
                      : prev + 1,
                );

                return;
              }

              if (event.key === "ArrowUp") {
                event.preventDefault();

                setActiveIndex((prev) =>
                  visibleResults.length === 0
                    ? 0
                    : prev === 0
                      ? visibleResults.length - 1
                      : prev - 1,
                );

                return;
              }

              if (event.key === "Enter") {
                event.preventDefault();
                openActiveResult();
              }
            }}
            placeholder="Buscar en LEXBOX..."
            autoComplete="off"
            spellCheck={false}
            className="w-full bg-transparent px-2 py-2 text-base text-neutral-100 outline-none placeholder:text-neutral-600"
          />
        </div>

        <div className="border-b border-neutral-800 px-4 py-2 text-[11px] uppercase tracking-wide text-neutral-500">
          Ctrl+K · Escape para cerrar · Enter para abrir
        </div>

        <div className="max-h-96 overflow-y-auto p-2">
          {!canSearch ? (
            <p className="p-4 text-sm text-neutral-500">
             Escribe al menos 2 caracteres para buscar expedientes, documentos, eventos y notas.
            </p>
          ) : null}

          {canSearch && isLoading ? (
            <p className="p-4 text-sm text-neutral-500">Buscando...</p>
          ) : null}

          {canSearch && !isLoading && visibleResults.length === 0 ? (
            <p className="p-4 text-sm text-neutral-500">No hay resultados.</p>
          ) : null}

          {groupedResults.map((group) => (
            <div key={group.type} className="py-1">
              <div className="px-3 py-2 text-[11px] font-medium uppercase tracking-wide text-neutral-600">
                {group.label}
              </div>

              <div className="space-y-1">
                {group.items.map(({ result, index }) => {
                  const isActive = index === activeIndex;

                  return (
                    <button
                      key={`${result.type}-${result.id}`}
                      type="button"
                      onMouseEnter={() => setActiveIndex(index)}
                      onMouseDown={(event) => {
                        event.preventDefault();

                        router.push(result.href);
                        closePalette();
                      }}
                      className={`block w-full rounded-xl px-3 py-3 text-left transition ${
                        isActive
                          ? "bg-red-950/40 text-red-100"
                          : "text-neutral-300 hover:bg-neutral-900 hover:text-neutral-100"
                      }`}
                    >
                      <div className="flex items-center justify-between gap-3">
                        <span className="truncate text-sm font-medium">
                          {result.title}
                        </span>

                        <span className="shrink-0 rounded-full border border-neutral-700 px-2 py-0.5 text-[11px] text-neutral-500">
                          {result.subtitle}
                        </span>
                      </div>

                      {result.snippet ? (
                        <p className="mt-1 line-clamp-2 text-xs text-neutral-500">
                          {renderSnippet(result.snippet)}
                        </p>
                      ) : null}
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}