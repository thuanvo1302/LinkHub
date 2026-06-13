"use client";

import { FormEvent, useEffect, useState } from "react";

import { fetchAPI } from "@/lib/api";
import { getAccessToken } from "@/lib/auth";

type LinkItem = {
  id: string;
  title: string;
  url: string;
  icon: string;
  is_active: boolean;
};

export function LinksManager() {
  const [items, setItems] = useState<LinkItem[]>([]);
  const [title, setTitle] = useState("");
  const [url, setUrl] = useState("");
  const [error, setError] = useState("");

  async function load() {
    const token = getAccessToken();
    const data = await fetchAPI<LinkItem[]>("/api/v1/profile-links", {
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
      await fetchAPI("/api/v1/profile-links", {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify({ title, url, icon: "link" }),
      });
      setTitle("");
      setUrl("");
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not create link");
    }
  }

  return (
    <div className="card rounded-[28px] p-8 shadow-card">
      <h1 className="text-3xl font-bold">Profile links</h1>
      <form className="mt-6 grid gap-4 md:grid-cols-[1fr_1fr_auto]" onSubmit={onSubmit}>
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="Link title" value={title} onChange={(e) => setTitle(e.target.value)} />
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="https://..." value={url} onChange={(e) => setUrl(e.target.value)} />
        <button className="rounded-full bg-[#f26a4b] px-6 py-3 text-sm font-semibold text-white">Add</button>
      </form>
      {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}
      <div className="mt-6 space-y-3">
        {items.map((item) => (
          <div key={item.id} className="rounded-2xl bg-white px-4 py-4">
            <div className="font-semibold">{item.title}</div>
            <div className="text-sm text-black/55">{item.url}</div>
          </div>
        ))}
      </div>
    </div>
  );
}

