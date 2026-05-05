export type CaseFile = {
  id: string;
  client_id: string;
  reference: string;
  title: string;
  type: string;
  status: string;
  description: string;
  calendar_scope: string;
  august_non_business: boolean;
  created_at: string;
  updated_at: string;
};

export type Note = {
  id: string;
  case_file_id: string;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
};

export type DocumentReviewStatus = "pending_review" | "reviewed" | "error";

export type DocumentReviewResponse = {
  document_id: string;
  review_status: DocumentReviewStatus;
  reviewed_at: string;
  review_note: string;
};

export type DocumentItem = {
  id: string;
  case_file_id: string;
  original_name: string;
  storage_path: string;
  mime_type: string;
  file_hash: string;
  created_at: string;
  updated_at: string;

  review_status: DocumentReviewStatus;
  reviewed_at?: string;
  review_note?: string;
};

export type EventItem = {
  event_id: string;
  document_id?: string;
  original_name?: string;

  case_file_id?: string;
  case_file_reference?: string;
  case_file_title?: string;

  event_type: string;
  event_date: string;
  source_text: string;
  days_remaining?: number;
  status?: string;
  priority?: string;
  duplicate_count?: number;
  document_names?: string[];
  document_ids?: string[];
  anchor_date?: string;
  date_kind?: string;
  anchor_source?: string;
  relative_days?: number;
  is_business_days?: boolean;
  add_extra_day?: boolean;
  calendar_scope?: string;
  trigger_text?: string;
  computation?: string;
  review_status: string;
  reviewed_at?: string;
  resolved_at?: string;
  resolution_note?: string;
};

export type CaseFileDetail = {
  case_file: CaseFile;
  notes: Note[];
  documents: DocumentItem[];
};

export type Dashboard = {
  case_file: CaseFile;
  note_count: number;
  document_count: number;
  upcoming_events: EventItem[];
  documents_without_text: number;
  documents_without_text_list: string[];
  documents_with_unknown_metadata: number;
  documents_with_unknown_metadata_list: string[];
  documents_without_events: number;
  documents_without_events_list: string[];
  overdue_count: number;
  today_count: number;
  upcoming_count: number;
  critical_count: number;
  high_count: number;
  medium_count: number;
  low_count: number;
  pending_review_count: number;
  reviewed_count: number;
  resolved_count: number;
  active_event_count: number;
  resolved_event_count: number;
  needs_attention: boolean;
  top_alert: string;
  recommended_next_action: string;
  procedural_hint: string;
};

export type EventActionResponse = {
  event_id: string;
  review_status: string;
  reviewed_at?: string;
  resolved_at?: string;
  resolution_note?: string;
};

export type CreateCaseFileInput = {
  client_id: string;
  reference: string;
  title: string;
  type: string;
  description: string;
  calendar_scope: string;
  august_non_business: boolean;
};

export type DocumentEventDetail = {
  event_id: string;
  document_id: string;
  event_type: string;
  event_date: string;
  source_text: string;
  created_at: string;
  anchor_date?: string;
  date_kind?: string;
  anchor_source?: string;
  relative_days?: number;
  is_business_days?: boolean;
  add_extra_day?: boolean;
  calendar_scope?: string;
  trigger_text?: string;
  computation?: string;
  review_status: string;
  reviewed_at?: string;
  resolved_at?: string;
  resolution_note?: string;
};

export type DocumentDetail = {
  document: DocumentItem;
  file_exists: boolean;
  has_extracted_text: boolean;
  extracted_text: string;
  extracted_text_length: number;
  extracted_text_preview: string;
  has_metadata: boolean;
  document_type: string;
  legal_area: string;
  metadata_analyzed_at: string;
  has_events: boolean;
  events: DocumentEventDetail[];
};