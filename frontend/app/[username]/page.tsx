import { API_URL } from "@/lib/api";

type Props = {
  params: Promise<{ username: string }>;
};

type PublicPayload = {
  profile: {
    username: string;
    display_name: string;
    bio: string;
  };
  links: {
    id: string;
    title: string;
    url: string;
  }[];
};

export default async function PublicProfilePage({ params }: Props) {
  const { username } = await params;
  let payload: PublicPayload | null = null;

  try {
    const res = await fetch(`${API_URL}/api/v1/public/profiles/${username}`, {
      cache: "no-store",
    });
    const body = await res.json();
    if (body.success) {
      payload = body.data;
    }
  } catch {}

  return (
    <main className="mx-auto flex min-h-screen max-w-2xl items-center px-6 py-16">
      <section className="card w-full rounded-[32px] p-8 shadow-card">
        <div className="text-xs uppercase tracking-[0.3em] text-black/55">Public profile</div>
        <h1 className="mt-4 text-4xl font-bold">
          {payload?.profile.display_name || `@${username}`}
        </h1>
        <p className="mt-4 text-black/65">
          {payload?.profile.bio || `Route này tương ứng với public page trong bản plan: yourdomain.com/${username}.`}
        </p>
        <div className="mt-6 space-y-3">
          {payload?.links?.map((link) => (
            <a key={link.id} href={link.url} className="block rounded-2xl bg-white px-4 py-4">
              <div className="font-semibold">{link.title}</div>
              <div className="text-sm text-black/55">{link.url}</div>
            </a>
          ))}
        </div>
      </section>
    </main>
  );
}
