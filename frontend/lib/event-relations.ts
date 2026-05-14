import type { EventItem } from "@/lib/types";

type RelationEvent = Pick<
  EventItem,
  | "event_id"
  | "event_type"
  | "event_date"
  | "date_kind"
  | "anchor_date"
  | "anchor_source"
  | "document_id"
  | "case_file_id"
>;

function sameCaseOrUnknown(event: RelationEvent, candidate: RelationEvent) {
  return (
    !event.case_file_id ||
    !candidate.case_file_id ||
    candidate.case_file_id === event.case_file_id
  );
}

function sameDocumentOrUnknown(event: RelationEvent, candidate: RelationEvent) {
  return (
    !event.document_id ||
    !candidate.document_id ||
    candidate.document_id === event.document_id
  );
}

function canHaveOrigin(event: RelationEvent) {
  return event.event_type === "deadline" || event.event_type === "requirement";
}

function canBeOrigin(event: RelationEvent) {
  return (
    event.event_type === "notification" ||
    event.event_type === "filing" ||
    event.event_type === "requirement" ||
    event.event_type === "hearing" ||
    event.event_type === "appearance"
  );
}

export function findRelatedEvent(
  event: RelationEvent,
  allEvents: RelationEvent[] = [],
) {
  if (!canHaveOrigin(event)) return undefined;
  if (!event.anchor_date) return undefined;

  const notification = allEvents.find(
    (candidate) =>
      candidate.event_id !== event.event_id &&
      candidate.event_type === "notification" &&
      candidate.event_date === event.anchor_date &&
      sameCaseOrUnknown(event, candidate) &&
      sameDocumentOrUnknown(event, candidate),
  );

  if (notification) return notification;

  return allEvents.find(
    (candidate) =>
      candidate.event_id !== event.event_id &&
      candidate.event_date === event.anchor_date &&
      candidate.event_type !== event.event_type &&
      sameCaseOrUnknown(event, candidate) &&
      sameDocumentOrUnknown(event, candidate),
  );
}

export function findDerivedEvents(
  event: RelationEvent,
  allEvents: RelationEvent[] = [],
) {
  if (!canBeOrigin(event)) return [];

  return allEvents
    .filter(
      (candidate) =>
        candidate.event_id !== event.event_id &&
        canHaveOrigin(candidate) &&
        candidate.anchor_date === event.event_date &&
        sameCaseOrUnknown(event, candidate) &&
        sameDocumentOrUnknown(event, candidate),
    )
    .sort((a, b) => {
      const dateCompare = a.event_date.localeCompare(b.event_date);
      if (dateCompare !== 0) return dateCompare;

      return a.event_type.localeCompare(b.event_type);
    });
}

export function eventRelationLabel(
  event: Pick<
    EventItem,
    "event_type" | "date_kind" | "anchor_date" | "anchor_source"
  >,
  allEvents: Array<Pick<EventItem, "event_type" | "event_date">> = [],
) {
  if (
    event.event_type === "deadline" &&
    event.date_kind === "relative" &&
    event.anchor_date
  ) {
    const anchorEvent = allEvents.find(
      (candidate) =>
        candidate.event_date === event.anchor_date &&
        candidate.event_type === "notification",
    );

    if (anchorEvent || event.anchor_source === "notification_line") {
      return "Generado por notificación";
    }

    if (event.anchor_source === "procedural_anchor_line") {
      return "Deriva de resolución/providencia";
    }

    if (event.anchor_source === "previous_line") {
      return "Deriva de fecha anterior";
    }

    return "Deriva de fecha base";
  }

  if (event.event_type === "requirement" && event.anchor_date) {
    return "Actuación vinculada a fecha base";
  }

  if (event.event_type === "hearing" || event.event_type === "appearance") {
    return "Señalamiento procesal";
  }

  return "";
}

export function derivedEventsLabel(count: number) {
  if (count === 0) return "";

  if (count === 1) {
    return "Tiene 1 hito derivado";
  }

  return `Tiene ${count} hitos derivados`;
}