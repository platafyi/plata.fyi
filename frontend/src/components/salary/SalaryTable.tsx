"use client";

import { useRouter, useSearchParams, usePathname } from "next/navigation";
import type { SalarySubmission } from "@/types";
import { ARRANGEMENT_LABEL, SENIORITY_LABEL } from "@/lib/constants";

interface Props {
  submissions: SalarySubmission[];
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

function formatMKD(amount: number): string {
  return new Intl.NumberFormat("mk-MK", {
    style: "currency",
    currency: "MKD",
    maximumFractionDigits: 0,
  }).format(amount);
}

export default function SalaryTable({ submissions, total, page, pageSize, onPageChange }: Props) {
  const totalPages = Math.ceil(total / pageSize);
  const router = useRouter();
  const searchParams = useSearchParams();
  const pathname = usePathname();

  const filterBy = (key: string, value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set(key, value);
    params.delete("page");
    router.push(`${pathname}?${params.toString()}`);
  };

  if (submissions.length === 0) {
    return (
      <div className="card p-12 text-center">
        <p className="text-lg font-bold">Нема резултати</p>
        <p className="text-sm mt-1 opacity-50">Обидете се со различни филтри</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <p className="label">
        Прикажани {(page - 1) * pageSize + 1}–{Math.min(page * pageSize, total)} од {total} записи
      </p>

      <p className="text-xs font-semibold sm:hidden" style={{ color: "#ec4899" }}>← Скрол за повеќе →</p>

      <div style={{ border: "2px solid rgb(40,40,37)", borderRadius: "4px", boxShadow: "4px 4px 0 0 rgb(40,40,37)", overflow: "hidden", backgroundColor: "#ffffff" }}>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead style={{ borderBottom: "2px solid rgb(40,40,37)", backgroundColor: "#f9f9f8" }}>
              <tr>
                {["Компанија / Позиција", "Индустрија", "Град", "Ниво", "Начин", "Искуство", "Месечна плата"].map((h, i) => (
                  <th key={h} className={`px-4 py-3 label text-left${i === 6 ? " text-right" : ""}`}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {submissions.map((s, idx) => (
                <tr
                  key={s.id}
                  onClick={() => router.push(`/s/${s.id}`)}
                  style={{
                    borderTop: idx > 0 ? "1px solid rgba(40,40,37,0.1)" : "none",
                    cursor: "pointer",
                    transition: "background 0.1s, box-shadow 0.1s",
                  }}
                  className="hover:bg-[rgba(56,237,129,0.5)]"
                >
                  <td className="px-4 py-3">
                    <div className="font-bold text-ink">{s.company_name}</div>
                    <div className="text-sm font-medium opacity-50 mt-0.5">{s.job_title}</div>
                  </td>
                  <td className="px-4 py-3">
                    <button onClick={(e) => {
                      e.stopPropagation();
                      if (pathname.startsWith("/city/")) {
                        filterBy("industry", s.industry_slug);
                      } else {
                        router.push(`/industry/${s.industry_slug}`);
                      }
                    }} className="tag tag-white" style={{ cursor: "pointer" }}>
                      {s.industry_name}
                    </button>
                  </td>
                  <td className="px-4 py-3">
                    <button onClick={(e) => {
                      e.stopPropagation();
                      if (pathname.startsWith("/industry/")) {
                        filterBy("city", s.city_slug);
                      } else {
                        router.push(`/city/${s.city_slug}`);
                      }
                    }} className="tag tag-white" style={{ cursor: "pointer" }}>
                      {s.city_name}
                    </button>
                  </td>
                  <td className="px-4 py-3">
                    <button
                      onClick={(e) => { e.stopPropagation(); filterBy("seniority", s.seniority); }}
                      className="tag tag-white"
                      style={{ cursor: "pointer" }}
                    >
                      {SENIORITY_LABEL[s.seniority] || s.seniority}
                    </button>
                  </td>
                  <td className="px-4 py-3 font-medium opacity-60 text-sm">
                    {ARRANGEMENT_LABEL[s.work_arrangement] || s.work_arrangement}
                  </td>
                  <td className="px-4 py-3 font-medium opacity-60 text-sm">{s.years_experience} год.</td>
                  <td className="px-4 py-3 text-right">
                    <span className="font-black text-ink text-lg" style={{ backgroundColor: "#fe91e6", padding: "2px 8px", borderRadius: "4px" }}>{formatMKD(s.base_salary)}</span>
                    {s.bonuses && s.bonuses.length > 0 && (
                      <span className="block text-xs font-semibold opacity-40 mt-1">+ бонус</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {totalPages > 1 && (
        <div className="flex justify-center items-center gap-3">
          <button
            onClick={() => onPageChange(page - 1)}
            disabled={page <= 1}
            className="btn-white disabled:opacity-30 disabled:cursor-not-allowed"
          >
            ← Претходна
          </button>
          <span className="label">{page} / {totalPages}</span>
          <button
            onClick={() => onPageChange(page + 1)}
            disabled={page >= totalPages}
            className="btn-white disabled:opacity-30 disabled:cursor-not-allowed"
          >
            Следна →
          </button>
        </div>
      )}
    </div>
  );
}
