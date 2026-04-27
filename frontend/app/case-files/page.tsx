import { getCaseFiles } from "@/lib/api";
import { CaseFilesTable } from "@/app/components/case-files-table";

export default async function CaseFilesPage() {
  const items = await getCaseFiles();

  return (
    <main className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Case files</h1>
        <p className="mt-2 text-neutral-400">
          Operational overview of all tracked legal matters.
        </p>
      </div>

      <CaseFilesTable items={items} />
    </main>
  );
}