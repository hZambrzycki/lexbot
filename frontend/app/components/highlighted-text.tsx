"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { normalizeSearchText } from "@/lib/search-normalization";

type Props = {
  text: string;
  query: string;
};

type TextPart = {
  value: string;
  isMatch: boolean;
  matchIndex?: number;
};

type MatchRange = {
  start: number;
  end: number;
};

function getSearchTerms(query: string): string[] {
  const terms = query
    .trim()
    .split(/\s+/)
    .map((term) => normalizeSearchText(term))
    .filter(Boolean);

  return Array.from(new Set(terms));
}

function buildNormalizedIndex(text: string) {
  let normalized = "";
  const indexMap: number[] = [];

  for (let index = 0; index < text.length; index += 1) {
    const normalizedChar = normalizeSearchText(text[index]);

    for (let innerIndex = 0; innerIndex < normalizedChar.length; innerIndex += 1) {
      normalized += normalizedChar[innerIndex];
      indexMap.push(index);
    }
  }

  return { normalized, indexMap };
}

function findMatchRanges(text: string, query: string): MatchRange[] {
  const terms = getSearchTerms(query);

  if (terms.length === 0) {
    return [];
  }

  const { normalized, indexMap } = buildNormalizedIndex(text);
  const ranges: MatchRange[] = [];

  for (const term of terms.sort((a, b) => b.length - a.length)) {
    let searchFrom = 0;

    while (searchFrom < normalized.length) {
      const foundAt = normalized.indexOf(term, searchFrom);

      if (foundAt === -1) break;

      const start = indexMap[foundAt];
      const end = (indexMap[foundAt + term.length - 1] ?? start) + 1;

      ranges.push({ start, end });
      searchFrom = foundAt + term.length;
    }
  }

  return ranges
    .sort((a, b) => a.start - b.start || b.end - a.end)
    .filter((range, index, sorted) => {
      const previous = sorted[index - 1];

      if (!previous) return true;

      return range.start >= previous.end;
    });
}

function buildParts(text: string, query: string): TextPart[] {
  const ranges = findMatchRanges(text, query);

  if (ranges.length === 0) {
    return [{ value: text, isMatch: false }];
  }

  const parts: TextPart[] = [];
  let cursor = 0;

  ranges.forEach((range, matchIndex) => {
    if (range.start > cursor) {
      parts.push({
        value: text.slice(cursor, range.start),
        isMatch: false,
      });
    }

    parts.push({
      value: text.slice(range.start, range.end),
      isMatch: true,
      matchIndex,
    });

    cursor = range.end;
  });

  if (cursor < text.length) {
    parts.push({
      value: text.slice(cursor),
      isMatch: false,
    });
  }

  return parts;
}

export function HighlightedText({ text, query }: Props) {
  const trimmedQuery = query.trim();
  const matchRefs = useRef<Map<number, HTMLElement>>(new Map());
  const [activeIndex, setActiveIndex] = useState(0);

  const parts = useMemo(
    () => buildParts(text, trimmedQuery),
    [text, trimmedQuery],
  );

  const totalMatches = useMemo(
    () => parts.filter((part) => part.isMatch).length,
    [parts],
  );

  useEffect(() => {
    const timeout = window.setTimeout(() => {
      setActiveIndex(0);
    }, 0);

    return () => window.clearTimeout(timeout);
  }, [trimmedQuery]);

  useEffect(() => {
    if (!trimmedQuery || totalMatches === 0) return;

    const el = matchRefs.current.get(activeIndex);

    if (el) {
      el.scrollIntoView({
        behavior: "smooth",
        block: "center",
      });
    }
  }, [activeIndex, trimmedQuery, totalMatches]);

  if (!trimmedQuery) {
    return <>{text}</>;
  }

  return (
    <span className="block">
      {totalMatches > 0 ? (
        <div className="sticky top-0 z-20 mb-3 flex items-center justify-between gap-3 rounded-xl border border-yellow-900/40 bg-yellow-950/95 p-3 text-xs text-yellow-100 shadow-lg backdrop-blur">
          <span>
            {totalMatches === 1
              ? "1 coincidencia"
              : `${totalMatches} coincidencias`}{" "}
            · {activeIndex + 1}/{totalMatches}
          </span>

          <div className="flex gap-2">
            <button
              type="button"
              onClick={() =>
                setActiveIndex((prev) =>
                  prev === 0 ? totalMatches - 1 : prev - 1,
                )
              }
              className="rounded border border-yellow-700/40 px-2 py-1 hover:bg-yellow-900/40"
            >
              ↑
            </button>

            <button
              type="button"
              onClick={() =>
                setActiveIndex((prev) =>
                  prev === totalMatches - 1 ? 0 : prev + 1,
                )
              }
              className="rounded border border-yellow-700/40 px-2 py-1 hover:bg-yellow-900/40"
            >
              ↓
            </button>
          </div>
        </div>
      ) : null}

      {parts.map((part, index) => {
        if (!part.isMatch) {
          return <span key={index}>{part.value}</span>;
        }

        const isActive = part.matchIndex === activeIndex;

        return (
          <mark
            key={index}
            ref={(el) => {
              if (el && part.matchIndex !== undefined) {
                matchRefs.current.set(part.matchIndex, el);
              }
            }}
            className={`rounded-md px-1 text-yellow-100 ${
              isActive
                ? "border border-yellow-400 bg-yellow-700/60"
                : "border border-yellow-700/40 bg-yellow-900/50"
            }`}
          >
            {part.value}
          </mark>
        );
      })}
    </span>
  );
}