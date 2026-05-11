import { normalizeSearchText } from "@/lib/search-normalization";
import type { LegalSearchTerm } from "@/lib/legal-search-terms";

type Props = {
  suggestions: LegalSearchTerm[];
  activeIndex: number;
  activeToken: string;
  onSelect: (suggestion: LegalSearchTerm) => void;
  onHover: (index: number) => void;
};

function HighlightedSuggestion({
  text,
  token,
}: {
  text: string;
  token: string;
}) {
  const normalizedText = normalizeSearchText(text);
  const normalizedToken = normalizeSearchText(token);

  const cleanedToken = normalizedToken.includes(":")
    ? normalizedToken.split(":").at(-1) ?? normalizedToken
    : normalizedToken;

  const matchIndex = normalizedText.indexOf(cleanedToken);

  if (matchIndex === -1 || !cleanedToken) {
    return <>{text}</>;
  }

  const before = text.slice(0, matchIndex);
  const match = text.slice(matchIndex, matchIndex + cleanedToken.length);
  const after = text.slice(matchIndex + cleanedToken.length);

  return (
    <>
      {before}
      <span className="font-semibold text-red-200">{match}</span>
      {after}
    </>
  );
}

export function SearchAutocomplete({
  suggestions,
  activeIndex,
  activeToken,
  onSelect,
  onHover,
}: Props) {
  if (suggestions.length === 0) {
    return null;
  }

  return (
    <div className="absolute left-0 right-0 top-full z-30 mt-2 overflow-hidden rounded-2xl border border-neutral-800 bg-neutral-950 shadow-2xl">
      <div className="flex items-center justify-between border-b border-neutral-800 px-3 py-2 text-[11px] font-medium uppercase tracking-wide text-neutral-500">
        <span>Sugerencias jurídicas</span>
        <span>Tab / Enter</span>
      </div>

      <ul className="max-h-64 overflow-y-auto py-1">
        {suggestions.map((suggestion, index) => {
          const isActive = index === activeIndex;

          return (
            <li key={`${suggestion.trigger}-${suggestion.value}`}>
              <button
                type="button"
                onMouseEnter={() => onHover(index)}
                onMouseDown={(event) => {
                  event.preventDefault();
                  onSelect(suggestion);
                }}
                className={`flex w-full items-center justify-between gap-3 px-3 py-2 text-left text-sm transition ${
                  isActive
                    ? "bg-red-950/40 text-red-100"
                    : "text-neutral-300 hover:bg-neutral-900 hover:text-neutral-100"
                }`}
              >
                <span className="min-w-0 truncate">
                  <HighlightedSuggestion
                    text={suggestion.value}
                    token={activeToken}
                  />
                </span>

                <span
                  className={`shrink-0 rounded-full border px-2 py-0.5 text-[11px] ${
                    isActive
                      ? "border-red-800/70 text-red-200"
                      : "border-neutral-700 text-neutral-500"
                  }`}
                >
                  {suggestion.trigger}
                </span>
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}