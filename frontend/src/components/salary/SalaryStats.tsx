"use client";

import { useRouter } from "next/navigation";
import type { SalaryStats } from "@/types";

interface Props {
  stats: SalaryStats[];
  groupBy: string;
}

function formatMKD(amount: number): string {
  return new Intl.NumberFormat("mk-MK", {
    style: "currency",
    currency: "MKD",
    maximumFractionDigits: 0,
  }).format(amount);
}

export default function SalaryStats({ stats, groupBy }: Props) {
  const router = useRouter();

  if (stats.length === 0) return null;

  const top = stats.slice(0, 3);

  const handleClick = (key: string) => {
    const base = groupBy === "city" ? `/city/${key}` : `/industry/${key}`;
    router.push(base);
  };

  return (
    <div>
      <p className="label mb-3">Најчести {groupBy === "city" ? "градови" : "индустрии"}</p>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        {top.map((s) => (
          <div
            key={s.group_key}
            className="card space-y-3"
          >
            <button
              onClick={() => handleClick(s.group_key)}
              className="font-black text-ink truncate text-left hover:underline cursor-pointer"
              title={s.group_val}
            >
              {s.group_val} →
            </button>
            <div className="grid grid-cols-2 gap-y-1">
              {[
                ["Записи", s.count],
                ["Просек", formatMKD(s.average)],
                ["Медиан", formatMKD(s.median)],
                ["Мин", formatMKD(s.min)],
                ["Макс", formatMKD(s.max)],
              ].map(([label, val]) => (
                <>
                  <span key={`${s.group_key}-${label}-l`} className="label">{label}</span>
                  <span key={`${s.group_key}-${label}-v`} className="text-sm font-bold text-right">{val}</span>
                </>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
