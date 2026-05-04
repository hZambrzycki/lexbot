type ExportLink = {
  href: string;
  label: string;
  description?: string;
};

type Props = {
  title?: string;
  description?: string;
  links: ExportLink[];
};

export function IcsExportActions({
  title = "Exportar agenda",
  description = "Descarga los hitos procesales en formato .ics para importarlos en tu calendario.",
  links,
}: Props) {
  if (links.length === 0) return null;

  return (
    <section className="rounded-2xl border border-neutral-800 bg-neutral-950/70 p-5">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div className="space-y-1">
          <h2 className="text-lg font-semibold text-neutral-50">{title}</h2>
          <p className="max-w-2xl text-sm leading-6 text-neutral-400">
            {description}
          </p>
        </div>

        <div className="flex flex-wrap gap-2">
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="inline-flex items-center rounded-xl border border-neutral-700 bg-neutral-900 px-4 py-2 text-sm font-medium text-neutral-100 transition hover:border-neutral-500 hover:bg-neutral-800"
            >
              {link.label}
            </a>
          ))}
        </div>
      </div>

      {links.some((link) => link.description) ? (
        <div className="mt-4 grid gap-2 md:grid-cols-2">
          {links
            .filter((link) => link.description)
            .map((link) => (
              <p
                key={`${link.href}-description`}
                className="rounded-xl border border-neutral-800 bg-black/20 p-3 text-xs leading-5 text-neutral-400"
              >
                <span className="font-medium text-neutral-200">
                  {link.label}:
                </span>{" "}
                {link.description}
              </p>
            ))}
        </div>
      ) : null}
    </section>
  );
}