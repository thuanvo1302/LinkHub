import { LoginForm } from "@/components/auth/login-form";
import { TopNav } from "@/components/top-nav";

export default function LoginPage() {
  return (
    <main className="min-h-screen">
      <TopNav />
      <section className="mx-auto max-w-2xl px-6 py-12">
        <div className="card rounded-[32px] p-8 shadow-card">
          <div className="text-sm uppercase tracking-[0.24em] text-black/55">Login</div>
          <h1 className="mt-3 text-4xl font-bold">Sign in to manage your links</h1>
          <LoginForm />
        </div>
      </section>
    </main>
  );
}
