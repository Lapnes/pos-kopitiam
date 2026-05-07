import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "POS KopiTiam — Kasir",
  description: "Point of Sale System for KopiTiam",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="id">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
      </head>
      <body>{children}</body>
    </html>
  );
}
