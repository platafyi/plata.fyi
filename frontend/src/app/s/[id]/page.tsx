import { notFound } from "next/navigation";
import type { Metadata } from "next";
import Link from "next/link";
import { getSubmission } from "@/lib/api";
import { SENIORITY_LABEL, ARRANGEMENT_LABEL, BONUS_TYPE_LABEL, BONUS_FREQ_LABEL } from "@/lib/constants";
import { netoToBruto } from "@/lib/salary";

function formatMKD(amount: number): string {
  return new Intl.NumberFormat("mk-MK", {
    style: "currency",
    currency: "MKD",
    maximumFractionDigits: 0,
  }).format(amount);
}

interface Props {
  params: Promise<{ id: string }>;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { id } = await params;
  try {
    const s = await getSubmission(id);
    return {
      title: `${s.job_title} @ ${s.company_name} — ${formatMKD(s.base_salary)}`,
      description: `${s.job_title} кај ${s.company_name} во ${s.city_name} (${s.industry_name}). Месечна нето плата: ${formatMKD(s.base_salary)}.`,
    };
  } catch {
    return { title: "Плата" };
  }
}

export default async function SubmissionDetailPage({ params }: Props) {
  const { id } = await params;
  let s;
  try {
    s = await getSubmission(id);
  } catch {
    notFound();
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6 py-8">
      {/* Back */}
      <Link href="/" className="label hover:underline">← Назад</Link>

      {/* Hero */}
      <div style={{ border: "2px solid rgb(40,40,37)", borderRadius: "4px", boxShadow: "4px 4px 0 0 rgb(40,40,37)", padding: "24px", backgroundColor: "#ffffff" }}>
        <p className="label mb-1">{s.company_name}</p>
        <h1 className="text-3xl sm:text-4xl font-black tracking-tight text-ink leading-tight">{s.job_title}</h1>
        <div className="mt-4 space-y-2">
          <div>
            <span className="font-black text-ink text-3xl" style={{ backgroundColor: "#fe91e6", padding: "4px 12px", borderRadius: "4px" }}>
              {formatMKD(s.base_salary)}
            </span>
            <span className="label ml-2">нето / месец</span>
          </div>
          <div>
            <span className="font-bold text-ink text-xl" style={{ backgroundColor: "#bbf7d0", padding: "3px 10px", borderRadius: "4px" }}>
              ~{formatMKD(netoToBruto(s.base_salary))}
            </span>
            <span className="label ml-2">
              бруто / месец{" "}
              <a
                href="https://github.com/skopjehacklab/kalkulator.ot.mk"
                target="_blank"
                rel="noopener noreferrer"
                className="opacity-50 hover:opacity-100 hover:underline"
              >
                [Како е пресметано]
              </a>
            </span>
          </div>
        </div>
      </div>

      {/* Tags */}
      <div style={{ border: "2px solid rgb(40,40,37)", borderRadius: "4px", boxShadow: "4px 4px 0 0 rgb(40,40,37)", padding: "20px", backgroundColor: "#ffffff" }}>
        <p className="label mb-3">Детали</p>
        <div className="flex flex-wrap gap-2">
          <Link href={`/industry/${s.industry_slug}`} className="tag tag-white hover:underline">{s.industry_name}</Link>
          <Link href={`/city/${s.city_slug}`} className="tag tag-white hover:underline">{s.city_name}</Link>
          <span className="tag tag-white">{SENIORITY_LABEL[s.seniority] || s.seniority}</span>
          <span className="tag tag-white">{ARRANGEMENT_LABEL[s.work_arrangement] || s.work_arrangement}</span>
          <span className="tag tag-white">{s.employment_type === "full_time" ? "Полно работно време" : "Скратено работно време"}</span>
          <span className="tag tag-white">{s.salary_year} год.</span>
        </div>
        <div className="mt-4 grid grid-cols-2 gap-3 text-sm">
          <div>
            <span className="label">Искуство (вкупно)</span>
            <p className="font-bold text-ink mt-0.5">{s.years_experience} год.</p>
          </div>
          <div>
            <span className="label">Искуство (во компанија)</span>
            <p className="font-bold text-ink mt-0.5">{s.years_at_company} год.</p>
          </div>
        </div>
      </div>

      {/* Bonuses */}
      {s.bonuses && s.bonuses.length > 0 && (
        <div style={{ border: "2px solid rgb(40,40,37)", borderRadius: "4px", boxShadow: "4px 4px 0 0 rgb(40,40,37)", padding: "20px", backgroundColor: "#ffffff" }}>
          <p className="label mb-3">Бонуси</p>
          <table className="w-full text-sm">
            <thead style={{ borderBottom: "2px solid rgb(40,40,37)" }}>
              <tr>
                {["Тип", "Износ", "Фреквенција"].map((h) => (
                  <th key={h} className="pb-2 label text-left">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {s.bonuses.map((b) => (
                <tr key={b.id} style={{ borderTop: "1px solid rgba(40,40,37,0.1)" }}>
                  <td className="py-2 font-medium text-ink">{BONUS_TYPE_LABEL[b.bonus_type] || b.bonus_type}</td>
                  <td className="py-2 font-bold text-ink">{formatMKD(b.amount)}</td>
                  <td className="py-2 font-medium opacity-60">{BONUS_FREQ_LABEL[b.frequency] || b.frequency}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
