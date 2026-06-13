"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { fetchAPI } from "@/lib/api";
import { setSession } from "@/lib/auth";

type AuthPayload = {
  access_token: string;
  refresh_token: string;
  user: {
    id: string;
    email: string;
    full_name: string;
  };
};

export function LoginForm() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError("");

    try {
      const data = await fetchAPI<AuthPayload>("/api/v1/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      setSession(data.access_token, data.refresh_token);
      router.push("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <form className="mt-8 space-y-4" onSubmit={onSubmit}>
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
          className="rounded-full bg-[#17181f] px-6 py-3 text-sm font-semibold text-white disabled:opacity-50"
          disabled={loading}
        >
          {loading ? "Logging in..." : "Login"}
        </button>
      </form>
      {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}
      <p className="mt-6 text-sm text-black/60">
        API target: <code>/api/v1/auth/login</code>.
      </p>
      <Link href="/register" className="mt-4 inline-block text-sm underline">
        Need an account?
      </Link>
    </>
  );
}

