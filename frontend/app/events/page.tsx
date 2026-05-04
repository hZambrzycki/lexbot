import Link from "next/link";
import { getGlobalUpcomingEvents } from "@/lib/api";
import { EventsTabs } from "./events-tabs";

function apiBaseUrl() {
  return process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
}

export default async function EventsPage() {
  const [pendingEvents, reviewedEvents, resolvedEvents] = await Promise.all([
    getGlobalUpcomingEvents({
      reviewStatus: "pending",
    }),
    getGlobalUpcomingEvents({
      reviewStatus: "reviewed",
    }),
    getGlobalUpcomingEvents({
      reviewStatus: "resolved",
    }),
  ]);

  const baseUrl = apiBaseUrl();

  return (
    <main className="space-y-8">
      <div className="space-y-3">
        <Link
          href="/case-files"
          className="text-sm text-neutral-400 underline-offset-4 hover:text-neutral-200 hover:underline"
        >
          ← Volver a expedientes
        </Link>

        <div>
          <p className="text-sm font-medium uppercase tracking-wide text-neutral-500">
            Agenda procesal
          </p>

          <h1 className="mt-2 text-3xl font-bold">
            Próximos eventos y vencimientos
          </h1>

          <p className="mt-2 max-w-3xl text-neutral-400">
            Vista global de hitos procesales detectados en todos los
            expedientes. La agenda separa lo pendiente, lo revisado y lo ya
            resuelto para que no se mezclen actuaciones abiertas con histórico.
          </p>
        </div>
      </div>

      <EventsTabs
        pendingEvents={pendingEvents}
        reviewedEvents={reviewedEvents}
        resolvedEvents={resolvedEvents}
        baseUrl={baseUrl}
      />
    </main>
  );
}