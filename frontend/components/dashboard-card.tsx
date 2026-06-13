type Props = {
  title: string;
  value: string;
  hint: string;
};

export function DashboardCard({ title, value, hint }: Props) {
  return (
    <div className="card rounded-[28px] p-6 shadow-card">
      <div className="text-sm uppercase tracking-[0.24em] text-black/55">{title}</div>
      <div className="mt-3 text-4xl font-bold">{value}</div>
      <div className="mt-2 text-sm text-black/60">{hint}</div>
    </div>
  );
}

