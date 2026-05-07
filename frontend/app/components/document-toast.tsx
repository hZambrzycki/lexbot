type Toast = {
  text: string;
  type: "success" | "error";
};

type Props = {
  toast: Toast | null;
};

export function DocumentToast({ toast }: Props) {
  if (!toast) return null;

  return (
    <div className="fixed bottom-6 right-6 z-50">
      <div
        className={`rounded-xl border px-4 py-3 text-sm shadow-lg ${
          toast.type === "success"
            ? "border-emerald-800 bg-emerald-950/90 text-emerald-100"
            : "border-red-800 bg-red-950/90 text-red-100"
        }`}
      >
        {toast.text}
      </div>
    </div>
  );
}