"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { fetchAPI } from "@/lib/api";
import { setSession } from "@/lib/auth";

type AuthPayload = {
  access_token: string;
  refresh_token: string;
};

export function RegisterForm() {
  const router = useRouter();
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError("");

    try {
      const data = await fetchAPI<AuthPayload>("/api/v1/auth/register", {
        method: "POST",
        body: JSON.stringify({
          full_name: fullName,
          email,
          password,
        }),
      });
      setSession(data.access_token, data.refresh_token);
      router.push("/dashboard/profile");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Register failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <form className="mt-8 grid gap-4" onSubmit={onSubmit}>
        <input
          className="w-full rounded-2xl border border-black/10 bg-white px-4 py-3"
          placeholder="Full name"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
        />
        <input
          className="w-full rounded-2xl border border-black/10 bg-white px-4 py-3"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          className="w-full rounded-2xl border border-black/10 bg-white px-4 py-3"
          placeholder="Password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          className="w-fit rounded-full bg-[#f26a4b] px-6 py-3 text-sm font-semibold text-white disabled:opacity-50"
          disabled={loading}
        >
          {loading ? "Creating..." : "Register"}
        </button>
      </form>
      {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}
    </>
  );
}

