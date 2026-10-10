import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Research Workflow Platform",
  description: "Submit and monitor asynchronous research workflows.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
