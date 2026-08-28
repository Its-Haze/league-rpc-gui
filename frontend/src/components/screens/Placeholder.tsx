export function Placeholder({ title, ticket }: { title: string; ticket: number }) {
  return (
    <div className="flex flex-col gap-2">
      <h1 className="text-xl font-semibold">{title}</h1>
      <p className="text-muted text-sm">This section lands in ticket {ticket}.</p>
    </div>
  );
}
