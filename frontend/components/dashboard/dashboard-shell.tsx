"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { DashboardCard } from "@/components/dashboard-card";
import { fetchAPI } from "@/lib/api";
import { clearSession, getAccessToken } from "@/lib/auth";

type Overview = {
  total_short_links: number;
  total_clicks: number;
  top_links: { id: string; code: string; title: string; clicks: number }[];
};

const sections = [
  { href: "/dashboard/profile", label: "Profile setup", desc: "Username, bio, avatar, theme" },
  { href: "/dashboard/links", label: "Profile links", desc: "Add, sort, enable, disable" },
  { href: "/dashboard/short-links", label: "Short links", desc: "Custom code and redirects" },
  { href: "/dashboard/analytics", label: "Analytics", desc: "Clicks, top links, trends" },
  { href: "/dashboard/billing", label: "Billing", desc: "Free vs Pro roadmap" },
];

export function DashboardShell() {
  const router = useRouter();
  const [overview, setOverview] = useState<Overview | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    const token = getAccessToken();
    if (!token) {
      router.push("/login");
      return;
    }

    fetchAPI<Overview>("/api/v1/analytics/overview", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(setOverview)
      .catch((err) => setError(err instanceof Error ? err.message : "Could not load dashboard"));
  }, [router]);

  function logout() {
    clearSession();
    router.push("/login");
  }

  const cards = [
    {
      title: "Short Links",
      value: String(overview?.total_short_links ?? 0).padStart(2, "0"),
      hint: "Live from /analytics/overview",
    },
    {
      title: "Total Clicks",
      value: String(overview?.total_clicks ?? 0),
      hint: "Redirect events counted in memory",
    },
    {
      title: "Top Links",
      value: String(overview?.top_links?.length ?? 0),
      hint: "Ready for richer analytics",
    },
  ];

  return (
    <section className="mx-auto max-w-6xl px-6 py-8">
      <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.24em] text-black/55">Dashboard</div>
          <h1 className="mt-2 text-4xl font-bold">MVP control room</h1>
          <p className="mt-3 max-w-2xl text-black/65">
            Dashboard này đang đọc dữ liệu thật từ backend.
          </p>
        </div>
        <button
          onClick={logout}
          className="rounded-full border border-black/10 bg-white/70 px-5 py-3 text-sm font-semibold"
        >
          Logout
        </button>
      </div>

      {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}

      <div className="mt-8 grid gap-4 md:grid-cols-3">
        {cards.map((item) => (
          <DashboardCard key={item.title} {...item} />
        ))}
      </div>

      {overview?.top_links?.length ? (
        <div className="card mt-8 rounded-[28px] p-6 shadow-card">
          <div className="text-sm uppercase tracking-[0.24em] text-black/55">Top links</div>
          <div className="mt-4 grid gap-3">
            {overview.top_links.map((item) => (
              <div key={item.id} className="rounded-2xl bg-white px-4 py-4">
                <div className="font-semibold">{item.title || item.code}</div>
                <div className="mt-1 text-sm text-black/55">
                  /{item.code} • {item.clicks} clicks
                </div>
              </div>
            ))}
          </div>
        </div>
      ) : null}

      <div className="mt-8 grid gap-4 md:grid-cols-2">
        {sections.map((section) => (
          <Link
            key={section.href}
            href={section.href}
            className="card rounded-[28px] p-6 shadow-card transition hover:-translate-y-1"
          >
            <div className="text-2xl font-semibold">{section.label}</div>
            <div className="mt-3 text-black/60">{section.desc}</div>
          </Link>
        ))}
      </div>
    </section>
  );
}

