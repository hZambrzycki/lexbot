import type { DocumentReviewStatus } from "@/lib/types";

type Props = {
  status?: DocumentReviewStatus;
};

export function DocumentReviewBadge({ status = "pending_review" }: Props) {
  const config = {
    pending_review: {
      label: "Pendiente de revisión",
      className: "border-amber-800 bg-amber-950/40 text-amber-100",
    },
    reviewed: {
      label: "Revisado",
      className: "border-emerald-800 bg-emerald-950/40 text-emerald-100",
    },
    error: {
      label: "Error",
      className: "border-red-800 bg-red-950/40 text-red-100",
    },
  }[status];

  return (
    <span
      className={`inline-flex rounded-full border px-2.5 py-1 text-xs font-medium ${config.className}`}
    >
      {config.label}
    </span>
  );
}