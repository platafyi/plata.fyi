"use client";

import { useRouter, useSearchParams, usePathname } from "next/navigation";
import { useCallback } from "react";
import type { City, Industry } from "@/types";
import { SENIORITY_LEVELS, WORK_ARRANGEMENTS } from "@/lib/constants";
import Select from "@/components/ui/Select";

interface Props {
  industries: Industry[];
  cities: City[];
  lockedCity?: string;
  lockedIndustry?: string;
}

export default function FilterBar({ industries, cities, lockedCity, lockedIndustry }: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const pathname = usePathname();

  const updateFilter = useCallback(
    (key: string, value: string) => {
      const params = new URLSearchParams(searchParams.toString());
      params.delete("page");

      if (key === "city") {
        params.delete("city");
        const qs = params.toString();
        if (value) {
          router.push(`/city/${value}${qs ? `?${qs}` : ""}`);
        } else {
          router.push(`/${qs ? `?${qs}` : ""}`);
        }
        return;
      }

      if (key === "industry") {
        params.delete("industry");
        const qs = params.toString();
        if (value) {
          router.push(`/industry/${value}${qs ? `?${qs}` : ""}`);
        } else {
          router.push(`/${qs ? `?${qs}` : ""}`);
        }
        return;
      }

      if (value) {
        params.set(key, value);
      } else {
        params.delete(key);
      }
      router.push(`${pathname}?${params.toString()}`);
    },
    [router, searchParams, pathname, lockedCity, lockedIndustry]
  );

  const clearAll = () => router.push(lockedCity || lockedIndustry ? "/" : pathname);

  const hasFilters =
    lockedCity ||
    lockedIndustry ||
    (!lockedIndustry && searchParams.get("industry")) ||
    (!lockedCity && searchParams.get("city")) ||
    searchParams.get("seniority") ||
    searchParams.get("arrangement") ||
    searchParams.get("min_salary") ||
    searchParams.get("max_salary");

  return (
    <div style={{ border: "2px solid rgb(40,40,37)", borderRadius: "4px", boxShadow: "4px 4px 0 0 rgb(40,40,37)", padding: "16px", backgroundColor: "#ffffff" }}>
      {/* Row 1 — category filters: 2 cols on mobile, 4 on desktop */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <Select
          value={lockedIndustry || searchParams.get("industry") || ""}
          onChange={(v) => updateFilter("industry", String(v))}
          options={industries.map((i) => ({ value: i.slug, label: i.name }))}
          placeholder="Сите индустрии"
          accentColor="#e7fe05"
        />
        <Select
          value={lockedCity || searchParams.get("city") || ""}
          onChange={(v) => updateFilter("city", String(v))}
          options={cities.map((c) => ({ value: c.slug, label: c.name }))}
          placeholder="Сите градови"
          accentColor="#f6cb44"
        />
        <Select
          value={searchParams.get("seniority") || ""}
          onChange={(v) => updateFilter("seniority", String(v))}
          options={SENIORITY_LEVELS.map((s) => ({ value: s.value, label: s.label }))}
          placeholder="Сите нивоа"
          accentColor="#fe91e6"
        />
        <Select
          value={searchParams.get("arrangement") || ""}
          onChange={(v) => updateFilter("arrangement", String(v))}
          options={WORK_ARRANGEMENTS.map((a) => ({ value: a.value, label: a.label }))}
          placeholder="Начин"
          accentColor="#38ed81"
        />
      </div>

      {/* Row 2 — salary range: 2 cols on mobile, auto on desktop */}
      <div className="grid grid-cols-3 gap-3 items-center mt-3 pt-3" style={{ borderTop: "1px solid rgba(40,40,37,0.12)" }}>
        <input
          type="number"
          placeholder="Мин. плата (МКД)"
          value={searchParams.get("min_salary") || ""}
          onChange={(e) => updateFilter("min_salary", e.target.value)}
          className="input"
          min="0"
        />
        <input
          type="number"
          placeholder="Макс. плата (МКД)"
          value={searchParams.get("max_salary") || ""}
          onChange={(e) => updateFilter("max_salary", e.target.value)}
          className="input"
          min="0"
        />
        <button onClick={clearAll} className="btn-white text-sm" style={{ visibility: hasFilters ? "visible" : "hidden", backgroundColor: "#e2e2df" }}>
          Исчисти ×
        </button>
      </div>
    </div>
  );
}
