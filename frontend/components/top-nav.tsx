"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { isAuthenticated } from "@/lib/auth";

const links = [
  { href: "/login", label: "Login" },
  { href: "/register", label: "Register" },
  { href: "/dashboard", label: "Dashboard" },
];

export function TopNav() {
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setLoggedIn(isAuthenticated());
  }, []);

  return (
    <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-6">
      <Link href="/" className="text-2xl font-bold tracking-tight">
        LinkHub
      </Link>
      <nav className="flex gap-3 text-sm">
        {links
          .filter((item) => !(loggedIn && item.href === "/login"))
          .map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="rounded-full border border-black/10 bg-white/70 px-4 py-2 transition hover:-translate-y-0.5"
            >
              {item.label}
            </Link>
          ))}
      </nav>
    </header>
  );
}
