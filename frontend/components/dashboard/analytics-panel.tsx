"use client";

import { useEffect, useState } from "react";

import { fetchAPI } from "@/lib/api";
import { getAccessToken } from "@/lib/auth";

type Overview = {
  total_short_links: number;
  total_clicks: number;
  top_links: { id: string; code: string; title: string; clicks: number }[];
};

export function AnalyticsPanel() {
  const [data, setData] = useState<Overview | null>(null);

  useEffect(() => {
    const token = getAccessToken();
    if (!token) return;
    fetchAPI<Overview>("/api/v1/analytics/overview", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(setData)
      .catch(() => {});
  }, []);

  return (
    <div className="card rounded-[28px] p-8 shadow-card">
      <h1 className="text-3xl font-bold">Analytics</h1>
      <div className="mt-6 grid gap-4 md:grid-cols-2">
        <div className="rounded-2xl bg-white px-4 py-4">
          <div className="text-sm text-black/55">Total clicks</div>
          <div className="mt-2 text-3xl font-bold">{data?.total_clicks ?? 0}</div>
        </div>
        <div className="rounded-2xl bg-white px-4 py-4">
          <div className="text-sm text-black/55">Short links</div>
          <div className="mt-2 text-3xl font-bold">{data?.total_short_links ?? 0}</div>
        </div>
      </div>
      <div className="mt-6 space-y-3">
        {data?.top_links?.map((item) => (
          <div key={item.id} className="rounded-2xl bg-white px-4 py-4">
            <div className="font-semibold">{item.title || item.code}</div>
            <div className="text-sm text-black/55">{item.clicks} clicks</div>
          </div>
        ))}
      </div>
    </div>
  );
}

