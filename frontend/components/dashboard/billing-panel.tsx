"use client";

import { useEffect, useState } from "react";

import { fetchAPI } from "@/lib/api";
import { getAccessToken } from "@/lib/auth";

type Plan = {
  id: string;
  name: string;
  price: number;
  currency: string;
  features: string[];
};

type Payment = {
  id: string;
  status: string;
  amount: number;
  currency: string;
};

type Subscription = {
  plan_id: string;
  status: string;
  expired_at: string;
};

export function BillingPanel() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [history, setHistory] = useState<Payment[]>([]);
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [message, setMessage] = useState("");

  async function load() {
    const token = getAccessToken();
    const [plansData, historyData] = await Promise.all([
      fetchAPI<Plan[]>("/api/v1/plans"),
      fetchAPI<Payment[]>("/api/v1/billing/history", {
        headers: { Authorization: `Bearer ${token}` },
      }),
    ]);
    setPlans(plansData);
    setHistory(historyData);

    try {
      const current = await fetchAPI<Subscription>("/api/v1/subscription/current", {
        headers: { Authorization: `Bearer ${token}` },
      });
      setSubscription(current);
    } catch {
      setSubscription(null);
    }
  }

  useEffect(() => {
    load().catch(() => {});
  }, []);

  async function upgrade() {
    const token = getAccessToken();
    const payment = await fetchAPI<Payment & { checkout_url: string }>("/api/v1/payments/create-checkout", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({ plan_id: "pro" }),
    });
    await fetchAPI("/api/v1/payments/mock-success", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({ payment_id: payment.id }),
    });
    setMessage("Mock payment completed. Subscription upgraded to Pro.");
    await load();
  }

  return (
    <div className="card rounded-[28px] p-8 shadow-card">
      <h1 className="text-3xl font-bold">Billing roadmap</h1>
      <p className="mt-3 text-black/65">Trang này đang gọi plans, history và mock payment/subscription từ backend.</p>
      <div className="mt-6 grid gap-4 md:grid-cols-2">
        {plans.map((plan) => (
          <div key={plan.id} className="rounded-2xl bg-white px-4 py-4">
            <div className="text-2xl font-semibold">{plan.name}</div>
            <div className="mt-1 text-black/55">{plan.price.toLocaleString()} {plan.currency}</div>
            <div className="mt-3 text-sm text-black/60">{plan.features.join(" • ")}</div>
          </div>
        ))}
      </div>
      <div className="mt-6 flex flex-wrap items-center gap-4">
        <button className="rounded-full bg-[#17181f] px-6 py-3 text-sm font-semibold text-white" onClick={upgrade}>
          Mock upgrade to Pro
        </button>
        <div className="text-sm text-black/60">
          Current: {subscription ? `${subscription.plan_id} (${subscription.status})` : "free / none"}
        </div>
      </div>
      {message ? <p className="mt-4 text-sm text-emerald-700">{message}</p> : null}
      <div className="mt-6 space-y-3">
        {history.map((item) => (
          <div key={item.id} className="rounded-2xl bg-white px-4 py-4 text-sm">
            Payment {item.id} • {item.status} • {item.amount.toLocaleString()} {item.currency}
          </div>
        ))}
      </div>
    </div>
  );
}

