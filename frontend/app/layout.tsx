import "./globals.css";
import type { Metadata } from "next";

import { CommandPalette } from "@/app/components/command-palette";

export const metadata: Metadata = {
  title: "LEXBOX",
  description: "Legal case operations dashboard",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-neutral-950 text-neutral-100">
        <CommandPalette />

        <div className="mx-auto max-w-7xl px-6 py-8">
          {children}
        </div>
      </body>
    </html>
  );
}