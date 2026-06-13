"use client";

import { FormEvent, useEffect, useState } from "react";

import { API_URL, fetchAPI } from "@/lib/api";
import { getAccessToken } from "@/lib/auth";

type ShortLink = {
  id: string;
  code: string;
  title: string;
  original_url: string;
  is_active: boolean;
};

export function ShortLinksManager() {
  const [items, setItems] = useState<ShortLink[]>([]);
  const [title, setTitle] = useState("");
  const [originalURL, setOriginalURL] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState("");

  async function load() {
    const token = getAccessToken();
    const data = await fetchAPI<ShortLink[]>("/api/v1/short-links", {
      headers: { Authorization: `Bearer ${token}` },
    });
    setItems(data);
  }

  useEffect(() => {
    load().catch(() => {});
  }, []);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    try {
      const token = getAccessToken();
      await fetchAPI("/api/v1/short-links", {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify({
          title,
          original_url: originalURL,
          code,
        }),
      });
      setTitle("");
      setOriginalURL("");
      setCode("");
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not create short link");
    }
  }

  return (
    <div className="card rounded-[28px] p-8 shadow-card">
      <h1 className="text-3xl font-bold">Short links</h1>
      <form className="mt-6 grid gap-4" onSubmit={onSubmit}>
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} />
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="Original URL" value={originalURL} onChange={(e) => setOriginalURL(e.target.value)} />
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="Custom code (optional)" value={code} onChange={(e) => setCode(e.target.value)} />
        <button className="w-fit rounded-full bg-[#17181f] px-6 py-3 text-sm font-semibold text-white">Create short link</button>
      </form>
      {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}
      <div className="mt-6 space-y-3">
        {items.map((item) => (
          <div key={item.id} className="rounded-2xl bg-white px-4 py-4">
            <div className="font-semibold">{item.title || item.code}</div>
            <div className="text-sm text-black/55">{API_URL}/s/{item.code}</div>
            <div className="mt-1 text-xs text-black/45">{item.original_url}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
