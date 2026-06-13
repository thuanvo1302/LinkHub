import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { TopNav } from "@/components/top-nav";

export default function DashboardPage() {
  return (
    <main className="min-h-screen pb-16">
      <TopNav />
      <DashboardShell />
    </main>
  );
}
