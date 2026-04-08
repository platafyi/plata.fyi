import type { Metadata } from "next";
import { Unbounded, Manrope } from "next/font/google";
import "./globals.css";
import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import Script from "next/script";

const jost = Unbounded({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700", "800"],
  variable: "--font-jost",
  display: "swap",
});

const manrope = Manrope({
  subsets: ["latin", "cyrillic"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-manrope",
  display: "swap",
});

export const metadata: Metadata = {
  icons: { icon: "/favicon.svg" },
  title: {
    default: "Plata.fyi — Споделете ја вашата плата анонимно",
    template: "%s | Plata.fyi",
  },
  description:
    "Анонимна платформа за споделување на плати во Македонија. Споредете ги вашите примања со останатите вработени во истата индустрија.",
  keywords: ["плата", "македонија", "споредба", "примања", "работа"],
  openGraph: {
    type: "website",
    locale: "mk_MK",
    url: "https://plata.fyi",
    siteName: "Plata.fyi",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="mk" className={`${jost.variable} ${manrope.variable}`}>
      <body className="min-h-screen bg-white flex flex-col">
        <Header />
        <main className="flex-1 container mx-auto px-4 py-8 max-w-6xl">
          {children}
        </main>
        <Footer />
        <Script
          defer
          src="https://cloud.umami.is/script.js"
          data-website-id="7cd147f4-0fe0-4220-9ea5-94bde4eaf0ff"
          strategy="afterInteractive"
        />
      </body>
    </html>
  );
}
