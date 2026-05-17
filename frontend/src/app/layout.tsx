import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Lunch Order",
  description: "Daily lunch order submission portal",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-gray-50 min-h-screen flex items-center justify-center">
        {children}
      </body>
    </html>
  );
}
