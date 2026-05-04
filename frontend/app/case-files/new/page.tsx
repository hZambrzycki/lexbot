import { getCaseFiles } from "@/lib/api";
import { NewCaseFileForm } from "./new-case-file-form";

function nextReference(existingReferences: string[]) {
  const year = new Date().getFullYear();

  const numbers = existingReferences
    .map((reference) => {
      const match = reference.match(new RegExp(`^EXP-${year}-(\\d+)$`));
      return match ? Number(match[1]) : 0;
    })
    .filter((value) => Number.isFinite(value));

  const next = Math.max(0, ...numbers) + 1;

  return `EXP-${year}-${String(next).padStart(4, "0")}`;
}

export default async function NewCaseFilePage() {
  const caseFiles = await getCaseFiles();

  const suggestedReference = nextReference(
    caseFiles.map((caseFile) => caseFile.reference),
  );

  return <NewCaseFileForm suggestedReference={suggestedReference} />;
}