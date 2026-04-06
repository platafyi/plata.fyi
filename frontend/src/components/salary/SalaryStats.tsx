"use client";

import { useRouter } from "next/navigation";
import { useMemo, useRef, useState, useEffect } from "react";
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
  const scrollRef = useRef<HTMLDivElement>(null);
  const [canPrev, setCanPrev] = useState(false);
  const [canNext, setCanNext] = useState(false);

  const shuffled = useMemo(
    () => [...stats].sort(() => Math.random() - 0.5),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [stats.length]
  );

  const updateArrows = () => {
    const el = scrollRef.current;
    if (!el) return;
    setCanPrev(el.scrollLeft > 4);
    setCanNext(el.scrollLeft < el.scrollWidth - el.clientWidth - 4);
  };

  useEffect(() => {
    updateArrows();
  }, [shuffled]);

  // Autoslide on mobile, stops when user touches or reaches the end
  useEffect(() => {
    const el = scrollRef.current;
    if (!el) return;
    if (!window.matchMedia("(max-width: 639px)").matches) return;

    let stopped = false;
    const stop = () => { stopped = true; };
    el.addEventListener("pointerdown", stop, { once: true });

    const interval = setInterval(() => {
      if (stopped) { clearInterval(interval); return; }
      if (el.scrollLeft >= el.scrollWidth - el.clientWidth - 4) {
        clearInterval(interval);
        return;
      }
      el.scrollBy({ left: el.clientWidth + 12, behavior: "smooth" });
    }, 3000);

    return () => {
      clearInterval(interval);
      el.removeEventListener("pointerdown", stop);
    };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shuffled]);

  const scroll = (dir: 1 | -1) => {
    const el = scrollRef.current;
    if (!el) return;
    const isDesktop = window.innerWidth >= 640;
    const amount = isDesktop ? el.clientWidth / 3 + 4 : el.clientWidth + 12;
    el.scrollBy({ left: dir * amount, behavior: "smooth" });
  };

  const handleClick = (key: string) => {
    router.push(groupBy === "city" ? `/city/${key}` : `/industry/${key}`);
  };

  if (stats.length === 0) return null;

  const label = groupBy === "city" ? "градови" : "индустрии";

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <p className="label">{groupBy === "city" ? "Градови" : "Индустрии"}</p>
        <div className={`hidden gap-2 ${shuffled.length > 3 ? "sm:flex" : ""}`}>
          <button
            onClick={() => scroll(-1)}
            disabled={!canPrev}
            className="btn btn-white px-3 py-1.5 text-base leading-none disabled:opacity-25 disabled:pointer-events-none"
            aria-label={`Претходна ${label}`}
          >
            ←
          </button>
          <button
            onClick={() => scroll(1)}
            disabled={!canNext}
            className="btn btn-white px-3 py-1.5 text-base leading-none disabled:opacity-25 disabled:pointer-events-none"
            aria-label={`Следна ${label}`}
          >
            →
          </button>
        </div>
      </div>

      {/* Slider track */}
      <div
        ref={scrollRef}
        onScroll={updateArrows}
        className="flex gap-3 overflow-x-auto pb-1.5 pr-1.5 [&::-webkit-scrollbar]:hidden"
        style={{ scrollSnapType: "x mandatory", scrollbarWidth: "none" }}
      >
        {shuffled.map((s) => (
          <div
            key={s.group_key}
            className="card space-y-3 flex-none w-full sm:w-[calc(33.333%-0.5rem)]"
            style={{ scrollSnapAlign: "start" }}
          >
            <button
              onClick={() => handleClick(s.group_key)}
              className="font-black text-ink truncate text-left hover:underline cursor-pointer w-full"
              title={s.group_val}
            >
              {s.group_val} →
            </button>
            <div className="grid grid-cols-2 gap-y-1">
              {(
                [
                  ["Записи", s.count],
                  ["Просек", formatMKD(s.average)],
                  ["Медиан", formatMKD(s.median)],
                  ["Мин", formatMKD(s.min)],
                  ["Макс", formatMKD(s.max)],
                ] as [string, string | number][]
              ).map(([lbl, val]) => (
                <>
                  <span key={`${s.group_key}-${lbl}-l`} className="label">
                    {lbl}
                  </span>
                  <span
                    key={`${s.group_key}-${lbl}-v`}
                    className="text-sm font-bold text-right"
                  >
                    {val}
                  </span>
                </>
              ))}
            </div>
          </div>
        ))}
      </div>

    </div>
  );
}