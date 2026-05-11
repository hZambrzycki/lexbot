import { useMemo, useState } from "react";

import {
  LEGAL_SEARCH_TERMS,
  type LegalSearchTerm,
} from "@/lib/legal-search-terms";

import { normalizeSearchText } from "@/lib/search-normalization";

type Props = {
  query: string;
  onQueryChange: (query: string) => void;
};

type Result = {
  activeToken: string;
  suggestions: LegalSearchTerm[];
  activeIndex: number;
  activeSuggestion: LegalSearchTerm | null;
  setActiveIndex: (index: number) => void;
  applySuggestion: (suggestion?: LegalSearchTerm) => void;
  handleKeyDown: (event: React.KeyboardEvent<HTMLInputElement>) => void;
  closeSuggestions: () => void;
};

type RankedSuggestion = {
  item: LegalSearchTerm;
  rank: number;
};

const MIN_TOKEN_LENGTH = 2;
const MIN_KEYWORD_TOKEN_LENGTH = 6;
const MAX_SUGGESTIONS = 6;

const DSL_SEARCH_TERMS: LegalSearchTerm[] = [
  { trigger: "type", value: "type:claim", keywords: ["demanda"] },
  { trigger: "type", value: "type:judgment", keywords: ["sentencia"] },
  { trigger: "type", value: "type:answer", keywords: ["contestación"] },
  { trigger: "type", value: "type:contract", keywords: ["contrato"] },
  { trigger: "type", value: "type:order", keywords: ["diligencia", "providencia", "decreto"] },
  { trigger: "type", value: "type:order_decision", keywords: ["auto"] },
  { trigger: "type", value: "type:appeal_motion", keywords: ["reposición"] },
  { trigger: "type", value: "type:appeal_brief", keywords: ["apelación"] },
  { trigger: "type", value: "type:conciliation_filing", keywords: ["smac", "conciliación"] },
  { trigger: "type", value: "type:dismissal_letter", keywords: ["carta de despido"] },
  { trigger: "type", value: "type:divorce_petition", keywords: ["divorcio", "paternofiliales"] },
  { trigger: "type", value: "type:residence_decision", keywords: ["residencia", "extranjería"] },
  { trigger: "type", value: "type:administrative_resolution", keywords: ["resolución administrativa"] },
  { trigger: "type", value: "type:enforcement_filing", keywords: ["ejecución"] },
  { trigger: "type", value: "type:monitorio_filing", keywords: ["monitorio"] },
  { trigger: "type", value: "type:payroll", keywords: ["nómina"] },
  { trigger: "type", value: "type:settlement", keywords: ["finiquito"] },
  { trigger: "type", value: "type:non_legal", keywords: ["no jurídico"] },
  { trigger: "type", value: "type:unknown", keywords: ["sin clasificar"] },

  { trigger: "area", value: "area:labor", keywords: ["laboral"] },
  { trigger: "area", value: "area:immigration", keywords: ["extranjería"] },
  { trigger: "area", value: "area:family", keywords: ["familia"] },
  { trigger: "area", value: "area:commercial", keywords: ["mercantil"] },
  { trigger: "area", value: "area:civil" },
  { trigger: "area", value: "area:procedural", keywords: ["procesal"] },
  { trigger: "area", value: "area:non_legal", keywords: ["no jurídico"] },
  { trigger: "area", value: "area:unknown", keywords: ["sin clasificar"] },

  { trigger: "review", value: "review:pending_review", keywords: ["pendiente"] },
  { trigger: "review", value: "review:reviewed", keywords: ["revisado"] },
  { trigger: "review", value: "review:error" },

  { trigger: "has", value: "has:events", keywords: ["con hitos"] },
  { trigger: "has", value: "has:no_events", keywords: ["sin hitos"] },
  { trigger: "has", value: "has:text", keywords: ["con texto"] },
  { trigger: "has", value: "has:no_text", keywords: ["sin texto"] },
  { trigger: "doc", value: "doc:pdf" },
  { trigger: "doc", value: "doc:docx" },
  { trigger: "doc", value: "doc:txt" },
  { trigger: "doc", value: "doc:md" },
];

function getLastToken(query: string): string {
  const parts = query.trimEnd().split(/\s+/);
  return parts[parts.length - 1] ?? "";
}

function startsWithToken(values: string[], token: string): boolean {
  return values.some((value) => value.startsWith(token));
}

function wordStartsWithToken(values: string[], token: string): boolean {
  return values.some((value) =>
    value.split(/\s+/).some((word) => word.startsWith(token)),
  );
}

function getSuggestionRank(
  item: LegalSearchTerm,
  activeToken: string,
): number | null {
  const trigger = normalizeSearchText(item.trigger);
  const value = normalizeSearchText(item.value);

  const aliases = (item.aliases ?? []).map((alias) =>
    normalizeSearchText(alias),
  );

  const keywords = (item.keywords ?? []).map((keyword) =>
    normalizeSearchText(keyword),
  );

  if (value === activeToken) return null;

  if (trigger === activeToken) return 1;
  if (startsWithToken(aliases, activeToken)) return 2;
  if (trigger.startsWith(activeToken)) return 3;
  if (value.startsWith(activeToken)) return 4;
  if (wordStartsWithToken([value], activeToken)) return 5;

  const canMatchKeyword = activeToken.length >= MIN_KEYWORD_TOKEN_LENGTH;

  if (canMatchKeyword && startsWithToken(keywords, activeToken)) {
    return 6;
  }

  if (canMatchKeyword && wordStartsWithToken(keywords, activeToken)) {
    return 7;
  }

  return null;
}

function getDslSuggestionRank(
  item: LegalSearchTerm,
  activeToken: string,
): number | null {
  const value = normalizeSearchText(item.value);
  const trigger = normalizeSearchText(item.trigger);

  if (value === activeToken) return null;

  if (value.startsWith(activeToken)) return 1;
  if (`${trigger}:`.startsWith(activeToken)) return 2;

  return null;
}

function isDslToken(token: string): boolean {
  return token.includes(":");
}

export function useSearchAutocomplete({
  query,
  onQueryChange,
}: Props): Result {
  const [activeIndex, setActiveIndex] = useState(0);
  const [closedToken, setClosedToken] = useState<string | null>(null);

  const rawActiveToken = useMemo(() => getLastToken(query), [query]);

  const activeToken = useMemo(
    () => normalizeSearchText(rawActiveToken),
    [rawActiveToken],
  );

  const rankedSuggestions = useMemo(() => {
    if (activeToken.length < MIN_TOKEN_LENGTH) {
      return [];
    }

    const seenValues = new Set<string>();
    const ranked: RankedSuggestion[] = [];

    const sourceTerms = isDslToken(activeToken)
      ? DSL_SEARCH_TERMS
      : [...DSL_SEARCH_TERMS, ...LEGAL_SEARCH_TERMS];

    for (const item of sourceTerms) {
      const normalizedValue = normalizeSearchText(item.value);

      if (seenValues.has(normalizedValue)) continue;

      const rank = isDslToken(activeToken)
        ? getDslSuggestionRank(item, activeToken)
        : getDslSuggestionRank(item, activeToken) ??
          getSuggestionRank(item, activeToken);

      if (rank === null) continue;

      seenValues.add(normalizedValue);
      ranked.push({ item, rank });
    }

    return ranked
      .sort(
        (a, b) =>
          a.rank - b.rank ||
          a.item.value.localeCompare(b.item.value, "es"),
      )
      .map((suggestion) => suggestion.item)
      .slice(0, MAX_SUGGESTIONS);
  }, [activeToken]);

  const isOpen = closedToken !== activeToken;
  const suggestions = isOpen ? rankedSuggestions : [];

  const safeActiveIndex =
    suggestions.length === 0
      ? 0
      : Math.min(activeIndex, suggestions.length - 1);

  const activeSuggestion = suggestions[safeActiveIndex] ?? null;

  function applySuggestion(selectedSuggestion?: LegalSearchTerm) {
    const target = selectedSuggestion ?? activeSuggestion;

    if (!target) return;

    const parts = query.trimEnd().split(/\s+/);
    parts[parts.length - 1] = target.value;

    onQueryChange(`${parts.join(" ")} `);
    setClosedToken(activeToken);
    setActiveIndex(0);
  }

  function closeSuggestions() {
    setActiveIndex(0);
    setClosedToken(activeToken);
  }

  function handleKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    if (suggestions.length === 0) return;

    if (event.key === "ArrowDown") {
      event.preventDefault();

      setActiveIndex((prev) =>
        prev === suggestions.length - 1 ? 0 : prev + 1,
      );

      return;
    }

    if (event.key === "ArrowUp") {
      event.preventDefault();

      setActiveIndex((prev) =>
        prev === 0 ? suggestions.length - 1 : prev - 1,
      );

      return;
    }

    if (event.key === "Tab" || event.key === "Enter") {
      event.preventDefault();
      applySuggestion();

      return;
    }

    if (event.key === "Escape") {
      event.preventDefault();
      closeSuggestions();
    }
  }

  return {
    activeToken,
    suggestions,
    activeIndex: safeActiveIndex,
    activeSuggestion,
    setActiveIndex,
    applySuggestion,
    handleKeyDown,
    closeSuggestions,
  };
}