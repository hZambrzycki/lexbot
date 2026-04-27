import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "LEXBOX",
  description: "Legal case operations dashboard",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className="bg-neutral-950 text-neutral-100 min-h-screen">
        <div className="mx-auto max-w-7xl px-6 py-8">{children}</div>
      </body>
    </html>
  );
}