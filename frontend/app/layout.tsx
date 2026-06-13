import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "LinkHub MVP",
  description: "Link-in-bio and URL shortener SaaS scaffold",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="shell">{children}</body>
    </html>
  );
}

