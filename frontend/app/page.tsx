import Link from "next/link";

export default function HomePage() {
  return (
    <main className="space-y-6">
      <div>
        <h1 className="text-4xl font-bold">LEXBOX</h1>
        <p className="text-neutral-400 mt-2">
          Local-first legal operations dashboard.
        </p>
      </div>

      <div className="rounded-2xl border border-neutral-800 bg-neutral-900 p-6">
        <Link
          href="/case-files"
          className="inline-flex rounded-xl bg-white px-4 py-2 text-black font-medium hover:opacity-90"
        >
          Open case files
        </Link>
      </div>
    </main>
  );
}