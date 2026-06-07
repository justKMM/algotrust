import { useMissionStore } from "@/hooks/useMissionStore";
import { BUDGET_TIERS, type BudgetTier } from "@/lib/budget";
import { cn } from "@/lib/utils";

const TIER_ORDER: BudgetTier[] = ["cheapass", "mid", "luxury"];

export function OpsCorner() {
  const m = useMissionStore((s) => s.metrics);
  const running = useMissionStore((s) => s.running);
  const budgetTier = useMissionStore((s) => s.budgetTier);
  const setBudgetTier = useMissionStore((s) => s.setBudgetTier);

  const stats = [
    { label: "Trust", value: m.trustScore },
    { label: "OK", value: m.successful },
    { label: "Blocked", value: m.blocked },
    { label: "Caught", value: m.violationsPrevented },
  ];

  return (
    <aside className="absolute right-3 top-3 z-10 w-[152px] rounded-md border border-border bg-surface/95 p-2 shadow-sm backdrop-blur-sm">
      <div className="mono-meta mb-1">Budget</div>
      <div className="flex flex-col gap-0.5" role="group" aria-label="Research budget tier">
        {TIER_ORDER.map((tier) => {
          const meta = BUDGET_TIERS[tier];
          const active = budgetTier === tier;
          return (
            <button
              key={tier}
              type="button"
              disabled={running}
              onClick={() => setBudgetTier(tier)}
              className={cn(
                "flex w-full items-center justify-between rounded px-1.5 py-0.5 text-left text-[10px] font-medium transition disabled:opacity-50",
                active
                  ? "bg-foreground text-background"
                  : "text-muted-foreground hover:bg-surface-2 hover:text-foreground",
              )}
            >
              <span className="truncate">{meta.label}</span>
              <span className="ml-1 shrink-0 font-mono tabular-nums">€{meta.eur}</span>
            </button>
          );
        })}
      </div>

      <div className="mt-2 grid grid-cols-2 gap-x-2 gap-y-1 border-t border-border/60 pt-2">
        {stats.map((s) => (
          <div key={s.label}>
            <div className="font-mono text-[9px] uppercase tracking-wide text-muted-foreground">
              {s.label}
            </div>
            <div className="text-[13px] font-semibold tabular-nums leading-tight">{s.value}</div>
          </div>
        ))}
      </div>
    </aside>
  );
}
