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

export function findRelatedEvent(
  event: RelationEvent,
  allEvents: RelationEvent[] = [],
) {
  if (
    event.event_type !== "deadline" &&
    event.event_type !== "requirement"
  ) {
    return undefined;
  }

  if (!event.anchor_date) {
    return undefined;
  }

  const sameCaseOrUnknown = (candidate: RelationEvent) =>
    !event.case_file_id ||
    !candidate.case_file_id ||
    candidate.case_file_id === event.case_file_id;

  const sameDocumentOrUnknown = (candidate: RelationEvent) =>
    !event.document_id ||
    !candidate.document_id ||
    candidate.document_id === event.document_id;

  const notification = allEvents.find(
    (candidate) =>
      candidate.event_id !== event.event_id &&
      candidate.event_type === "notification" &&
      candidate.event_date === event.anchor_date &&
      sameCaseOrUnknown(candidate) &&
      sameDocumentOrUnknown(candidate),
  );

  if (notification) {
    return notification;
  }

  const proceduralAnchor = allEvents.find(
    (candidate) =>
      candidate.event_id !== event.event_id &&
      candidate.event_date === event.anchor_date &&
      candidate.event_type !== event.event_type &&
      sameCaseOrUnknown(candidate) &&
      sameDocumentOrUnknown(candidate),
  );

  return proceduralAnchor;
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