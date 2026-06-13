"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";

import { fetchAPI } from "@/lib/api";
import { getAccessToken } from "@/lib/auth";

type Profile = {
  username: string;
  display_name: string;
  bio: string;
  avatar_url: string;
  theme: string;
  is_public: boolean;
};

const emptyProfile: Profile = {
  username: "",
  display_name: "",
  bio: "",
  avatar_url: "",
  theme: "sunset-grid",
  is_public: true,
};

export function ProfileManager() {
  const router = useRouter();
  const [profile, setProfile] = useState<Profile>(emptyProfile);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    const token = getAccessToken();
    if (!token) return;

    fetchAPI<Profile>("/api/v1/profiles/me", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(setProfile)
      .catch(() => {});
  }, []);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setMessage("");
    try {
      const token = getAccessToken();
      const data = await fetchAPI<Profile>("/api/v1/profiles/me", {
        method: "PUT",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify(profile),
      });
      setProfile(data);
      setMessage("Profile updated.");
      router.push("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not update profile");
    }
  }

  return (
    <div className="card rounded-[28px] p-8 shadow-card">
      <h1 className="text-3xl font-bold">Profile settings</h1>
      <p className="mt-3 text-black/65">Thiết lập public username, bio và theme cho trang link-in-bio.</p>
      <form className="mt-6 grid gap-4" onSubmit={onSubmit}>
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="username" value={profile.username} onChange={(e) => setProfile({ ...profile, username: e.target.value })} />
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="display name" value={profile.display_name} onChange={(e) => setProfile({ ...profile, display_name: e.target.value })} />
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="avatar url" value={profile.avatar_url} onChange={(e) => setProfile({ ...profile, avatar_url: e.target.value })} />
        <input className="rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="theme" value={profile.theme} onChange={(e) => setProfile({ ...profile, theme: e.target.value })} />
        <textarea className="min-h-28 rounded-2xl border border-black/10 bg-white px-4 py-3" placeholder="bio" value={profile.bio} onChange={(e) => setProfile({ ...profile, bio: e.target.value })} />
        <label className="flex items-center gap-3 text-sm">
          <input type="checkbox" checked={profile.is_public} onChange={(e) => setProfile({ ...profile, is_public: e.target.checked })} />
          Public profile
        </label>
        <button className="w-fit rounded-full bg-[#17181f] px-6 py-3 text-sm font-semibold text-white">Save profile</button>
      </form>
      {message ? <p className="mt-4 text-sm text-emerald-700">{message}</p> : null}
      {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}
    </div>
  );
}
