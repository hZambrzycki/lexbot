"use client";

import { useEffect, useMemo, useRef, useState } from "react";

type Props = {
  text: string;
  query: string;
};

type TextPart = {
  value: string;
  isMatch: boolean;
  matchIndex?: number;
};

function escapeRegExp(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

function getSearchTerms(query: string): string[] {
  const terms = query
    .trim()
    .split(/\s+/)
    .map((term) => term.trim())
    .filter(Boolean);

  return Array.from(new Set(terms));
}

function buildParts(text: string, query: string): TextPart[] {
  const terms = getSearchTerms(query);

  if (terms.length === 0) {
    return [{ value: text, isMatch: false }];
  }

  const pattern = terms
    .sort((a, b) => b.length - a.length)
    .map(escapeRegExp)
    .join("|");

  const regex = new RegExp(`(${pattern})`, "gi");
  const rawParts = text.split(regex);

  let matchIndex = 0;
  const lowerTerms = terms.map((term) => term.toLowerCase());

  return rawParts.map((part) => {
    const isMatch = lowerTerms.includes(part.toLowerCase());

    if (!isMatch) {
      return { value: part, isMatch: false };
    }

    const currentIndex = matchIndex;
    matchIndex += 1;

    return {
      value: part,
      isMatch: true,
      matchIndex: currentIndex,
    };
  });
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