import Link from "next/link";
import { CaseFile } from "@/lib/types";

type Props = {
  items: CaseFile[];
};

export function CaseFilesTable({ items }: Props) {
  return (
    <div className="overflow-hidden rounded-2xl border border-neutral-800 bg-neutral-900">
      <table className="min-w-full text-sm">
        <thead className="bg-neutral-800/70 text-neutral-300">
          <tr>
            <th className="px-4 py-3 text-left">Reference</th>
            <th className="px-4 py-3 text-left">Title</th>
            <th className="px-4 py-3 text-left">Type</th>
            <th className="px-4 py-3 text-left">Status</th>
            <th className="px-4 py-3 text-left">Scope</th>
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.id} className="border-t border-neutral-800">
              <td className="px-4 py-3">
                <Link
                  href={`/case-files/${item.id}`}
                  className="font-medium text-white underline-offset-4 hover:underline"
                >
                  {item.reference}
                </Link>
              </td>
              <td className="px-4 py-3">{item.title}</td>
              <td className="px-4 py-3">{item.type}</td>
              <td className="px-4 py-3">{item.status}</td>
              <td className="px-4 py-3">{item.calendar_scope}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}