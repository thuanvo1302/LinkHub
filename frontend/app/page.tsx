import Link from "next/link";

import { TopNav } from "@/components/top-nav";
import { dashboardStats, sampleLinks } from "@/lib/mock-data";

export default function HomePage() {
  return (
    <main className="min-h-screen pb-16">
      <TopNav />
      <section className="mx-auto grid max-w-6xl gap-10 px-6 pt-8 md:grid-cols-[1.1fr_0.9fr] md:items-center">
        <div>
          <div className="inline-flex rounded-full border border-black/10 bg-white/70 px-4 py-2 text-xs uppercase tracking-[0.28em] text-black/60">
            SaaS MVP Scaffold
          </div>
          <h1 className="mt-6 max-w-3xl text-5xl font-bold leading-[1.05] md:text-7xl">
            One page for your brand. One short link for every campaign.
          </h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-black/65">
            LinkHub bam theo ban plan trong file markdown: auth, public profile,
            profile links, short links, analytics va ha tang Docker de ban phat
            trien tiep thanh SaaS that.
          </p>
          <div className="mt-8 flex flex-wrap gap-4">
            <Link
              href="/dashboard"
              className="rounded-full bg-[#17181f] px-6 py-3 text-sm font-semibold text-white transition hover:-translate-y-0.5"
            >
              Open dashboard
            </Link>
            <Link
              href="/register"
              className="rounded-full border border-black/10 bg-white/70 px-6 py-3 text-sm font-semibold"
            >
              Start with auth flow
            </Link>
          </div>
        </div>

        <div className="card rounded-[32px] p-6 shadow-card">
          <div className="rounded-[24px] bg-[#17181f] p-6 text-white">
            <div className="text-xs uppercase tracking-[0.28em] text-white/60">
              Sample Public Profile
            </div>
            <div className="mt-5 flex items-center gap-4">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-[#f26a4b] text-2xl font-bold">
                T
              </div>
              <div>
                <div className="text-2xl font-semibold">Thuan Nguyen</div>
                <div className="text-sm text-white/70">@thuan-builds</div>
              </div>
            </div>
            <p className="mt-5 text-sm leading-7 text-white/75">
              Building useful tools with Go, Next.js, PostgreSQL and Redis.
            </p>
            <div className="mt-6 space-y-3">
              {sampleLinks.map((item) => (
                <div
                  key={item.title}
                  className="rounded-2xl border border-white/10 bg-white/10 px-4 py-3"
                >
                  <div className="font-semibold">{item.title}</div>
                  <div className="mt-1 text-xs text-white/65">{item.url}</div>
                </div>
              ))}
            </div>
          </div>
          <div className="mt-5 grid gap-3 md:grid-cols-3">
            {dashboardStats.map((item) => (
              <div key={item.title} className="rounded-2xl bg-white px-4 py-4">
                <div className="text-xs uppercase tracking-[0.18em] text-black/45">
                  {item.title}
                </div>
                <div className="mt-2 text-2xl font-bold">{item.value}</div>
              </div>
            ))}
          </div>
        </div>
      </section>
    </main>
  );
}
