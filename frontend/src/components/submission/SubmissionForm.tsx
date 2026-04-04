"use client";

import { useEffect, useState } from "react";
import CompanySearch from "./CompanySearch";
import BonusFields from "./BonusFields";
import { getCities, getIndustries } from "@/lib/api";
import type { BonusInput, City, Industry, SalarySubmission, SubmissionInput } from "@/types";
import { SENIORITY_LEVELS, WORK_ARRANGEMENTS } from "@/lib/constants";
import Select from "@/components/ui/Select";
import AutocompleteInput from "@/components/ui/AutocompleteInput";
import { searchJobTitles } from "@/lib/api";

interface Props {
  initial?: SalarySubmission;
  onSubmit: (data: SubmissionInput) => Promise<void>;
  onCancel?: () => void;
}

const ARRANGEMENT_COLORS: Record<string, string> = {
  office:  "#e7fe05",
  hybrid:  "#fe91e6",
  remote:  "#38ed81",
};

const SENIORITY_COLORS: Record<string, string> = {
  junior:    "#38ed81",
  mid:       "#f6cb44",
  senior:    "#fe91e6",
  lead:      "#93c5fd",
  manager:   "#fda4af",
  director:  "#fb923c",
  executive: "#f87171",
};

export default function SubmissionForm({ initial, onSubmit, onCancel }: Props) {
  const [industries, setIndustries] = useState<Industry[]>([]);
  const [cities, setCities] = useState<City[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [companyName, setCompanyName] = useState(initial?.company_name || "");
  const [companyRegNo, setCompanyRegNo] = useState(initial?.company_reg_no || "");
  const [jobTitle, setJobTitle] = useState(initial?.job_title || "");
  const [industryId, setIndustryId] = useState(initial?.industry_id || 0);
  const [cityId, setCityId] = useState(initial?.city_id || 0);
  const [seniority, setSeniority] = useState(initial?.seniority || "");
  const [yearsAtCompany, setYearsAtCompany] = useState(initial?.years_at_company ?? 0);
  const [yearsExperience, setYearsExperience] = useState(initial?.years_experience ?? 0);
  const [workArrangement, setWorkArrangement] = useState(initial?.work_arrangement || "");
  const [employmentType, setEmploymentType] = useState(initial?.employment_type || "full_time");
  const [hoursPerWeek, setHoursPerWeek] = useState(initial?.hours_per_week || 30);
  const [baseSalary, setBaseSalary] = useState(initial?.base_salary || 0);
  const currentYear = new Date().getFullYear();
  const [salaryYear, setSalaryYear] = useState(initial?.salary_year || currentYear);
  const [bonuses, setBonuses] = useState<BonusInput[]>(
    initial?.bonuses?.map((b) => ({
      bonus_type: b.bonus_type,
      amount: b.amount,
      frequency: b.frequency,
    })) || []
  );

  useEffect(() => {
    Promise.all([getIndustries(), getCities()]).then(([ind, cit]) => {
      setIndustries(ind);
      setCities(cit);
    });
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!companyName.trim()) return setError("Внесете име на компанија");
    if (!jobTitle.trim()) return setError("Внесете работна позиција");
    if (!industryId) return setError("Изберете индустрија");
    if (!cityId) return setError("Изберете град");
    if (!seniority) return setError("Изберете ниво на искуство");
    if (!workArrangement) return setError("Изберете начин на работа");
    if (!baseSalary || baseSalary <= 0) return setError("Внесете плата");

    const invalidBonus = bonuses.find((b) => !b.amount || b.amount <= 0);
    if (invalidBonus) return setError("Сите бонуси мора да имаат позитивен износ");

    setLoading(true);
    try {
      await onSubmit({
        company_name: companyName.trim(),
        company_reg_no: companyRegNo.trim() || undefined,
        job_title: jobTitle.trim(),
        industry_id: industryId,
        city_id: cityId,
        seniority,
        years_at_company: yearsAtCompany,
        years_experience: yearsExperience,
        work_arrangement: workArrangement,
        employment_type: employmentType,
        hours_per_week: employmentType === "part_time" ? hoursPerWeek : undefined,
        base_salary: baseSalary,
        salary_year: salaryYear,
        bonuses,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Грешка при зачувување");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">

      {/* Company */}
      <div>
        <label className="label block mb-2">Компанија *</label>
        <CompanySearch
          value={companyName}
          regNo={companyRegNo}
          onChange={(name, reg) => { setCompanyName(name); setCompanyRegNo(reg); }}
        />
      </div>

      {/* Job title */}
      <div>
        <label className="label block mb-2">Работна позиција *</label>
        <AutocompleteInput
          value={jobTitle}
          onChange={setJobTitle}
          onSearch={searchJobTitles}
          placeholder="пр. Софтверски инженер"
          accentColor="#e7fe05"
          newItemNotice="Не најдовме позиција, ќе креираме нова. Ве молиме внесете го целосното име."
        />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {/* Industry */}
        <div>
          <label className="label block mb-2">Индустрија *</label>
          <Select
            value={industryId || ""}
            onChange={(v) => setIndustryId(parseInt(String(v)))}
            options={industries.map((i) => ({ value: i.id, label: i.name }))}
            placeholder="Изберете..."
            accentColor="#e7fe05"
          />
        </div>

        {/* City */}
        <div>
          <label className="label block mb-2">Град *</label>
          <Select
            value={cityId || ""}
            onChange={(v) => setCityId(parseInt(String(v)))}
            options={cities.map((c) => ({ value: c.id, label: c.name }))}
            placeholder="Изберете..."
            accentColor="#f6cb44"
          />
        </div>
      </div>

      {/* Seniority — color pill grid */}
      <div>
        <label className="label block mb-2">Ниво на искуство *</label>
        <div className="flex flex-wrap gap-2">
          {SENIORITY_LEVELS.map((s) => {
            const active = seniority === s.value;
            const color = SENIORITY_COLORS[s.value] || "#ffffff";
            return (
              <button
                key={s.value}
                type="button"
                onClick={() => setSeniority(s.value)}
                style={{
                  backgroundColor: active ? color : "#ffffff",
                  border: "2px solid rgb(40,40,37)",
                  borderRadius: "4px",
                  boxShadow: active ? "none" : "3px 3px 0 0 rgb(40,40,37)",
                  transform: active ? "translate(2px,2px)" : "none",
                  padding: "6px 14px",
                  fontFamily: "var(--font-jost), sans-serif",
                  fontWeight: active ? 700 : 500,
                  fontSize: "0.72rem",
                  color: "rgb(40,40,37)",
                  cursor: "pointer",
                  transition: "all 80ms ease",
                }}
              >
                {s.label}
              </button>
            );
          })}
        </div>
      </div>

      {/* Work arrangement — 3-button toggle */}
      <div>
        <label className="label block mb-2">Начин на работа *</label>
        <div className="flex gap-2">
          {WORK_ARRANGEMENTS.map((a) => {
            const active = workArrangement === a.value;
            const color = ARRANGEMENT_COLORS[a.value] || "#ffffff";
            return (
              <button
                key={a.value}
                type="button"
                onClick={() => setWorkArrangement(a.value)}
                style={{
                  flex: 1,
                  backgroundColor: active ? color : "#ffffff",
                  border: "2px solid rgb(40,40,37)",
                  borderRadius: "4px",
                  boxShadow: active ? "none" : "3px 3px 0 0 rgb(40,40,37)",
                  transform: active ? "translate(2px,2px)" : "none",
                  padding: "10px 8px",
                  fontFamily: "var(--font-jost), sans-serif",
                  fontWeight: active ? 700 : 500,
                  fontSize: "0.76rem",
                  color: "rgb(40,40,37)",
                  cursor: "pointer",
                  transition: "all 80ms ease",
                  textAlign: "center",
                }}
              >
                {a.label}
              </button>
            );
          })}
        </div>
      </div>

      {/* Employment type */}
      <div>
        <label className="label block mb-2">Тип на вработување *</label>
        <div className="flex gap-2">
          {[
            { value: "full_time", label: "Полно работно време" },
            { value: "part_time", label: "Скратено работно време" },
          ].map((et) => {
            const active = employmentType === et.value;
            return (
              <button
                key={et.value}
                type="button"
                onClick={() => setEmploymentType(et.value)}
                style={{
                  flex: 1,
                  backgroundColor: active ? "#38ed81" : "#ffffff",
                  border: "2px solid rgb(40,40,37)",
                  borderRadius: "4px",
                  boxShadow: active ? "none" : "3px 3px 0 0 rgb(40,40,37)",
                  transform: active ? "translate(2px,2px)" : "none",
                  padding: "10px 8px",
                  fontFamily: "var(--font-jost), sans-serif",
                  fontWeight: active ? 700 : 500,
                  fontSize: "0.76rem",
                  color: "rgb(40,40,37)",
                  cursor: "pointer",
                  transition: "all 80ms ease",
                  textAlign: "center",
                }}
              >
                {et.label}
              </button>
            );
          })}
        </div>
        {employmentType === "part_time" && (
          <div className="mt-3">
            <label className="label block mb-2">Часови неделно</label>
            <input
              type="number"
              value={hoursPerWeek}
              onChange={(e) => setHoursPerWeek(parseInt(e.target.value) || 30)}
              min="1"
              max="40"
              className="input"
            />
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {/* Years experience */}
        <div>
          <label className="label block mb-2">Вкупно искуство (години)</label>
          <input
            type="number"
            value={yearsExperience}
            onChange={(e) => setYearsExperience(parseInt(e.target.value) || 0)}
            min="0"
            max="60"
            className="input"
          />
        </div>

        {/* Years at company */}
        <div>
          <label className="label block mb-2">Години во оваа компанија</label>
          <input
            type="number"
            value={yearsAtCompany}
            onChange={(e) => setYearsAtCompany(parseInt(e.target.value) || 0)}
            min="0"
            max="60"
            className="input"
          />
        </div>
      </div>

      {/* Base salary */}
      <div>
        <label className="label block mb-2">Месечна плата (МКД) — Нето *</label>
        <input
          type="number"
          value={baseSalary || ""}
          onChange={(e) => setBaseSalary(parseInt(e.target.value) || 0)}
          placeholder="пр. 55000"
          min="1"
          className="input"
          required
        />
      </div>

      {/* Salary year */}
      <div>
        <label className="label block mb-2">Година на плата *</label>
        <div className="flex gap-2 items-start">
          <button
            type="button"
            onClick={() => setSalaryYear(currentYear)}
            style={{
              flexShrink: 0,
              backgroundColor: salaryYear === currentYear ? "#e7fe05" : "#ffffff",
              border: "2px solid rgb(40,40,37)",
              borderRadius: "4px",
              boxShadow: salaryYear === currentYear ? "none" : "3px 3px 0 0 rgb(40,40,37)",
              transform: salaryYear === currentYear ? "translate(2px,2px)" : "none",
              padding: "10px 20px",
              fontFamily: "var(--font-jost), sans-serif",
              fontWeight: salaryYear === currentYear ? 700 : 500,
              fontSize: "0.76rem",
              color: "rgb(40,40,37)",
              cursor: "pointer",
              transition: "all 80ms ease",
            }}
          >
            Сегашна
          </button>
          <div style={{ flex: 1 }}>
            <Select
              value={salaryYear === currentYear ? "" : salaryYear}
              onChange={(v) => setSalaryYear(parseInt(String(v)))}
              options={Array.from({ length: currentYear - 2000 }, (_, i) => {
                const y = currentYear - 1 - i;
                return { value: y, label: String(y) };
              })}
              placeholder="Поранешна..."
              accentColor="#e7fe05"
            />
          </div>
        </div>
      </div>

      {/* Bonuses */}
      <div>
        <label className="label block mb-2">Бонуси (опционално)</label>
        <BonusFields bonuses={bonuses} onChange={setBonuses} />
      </div>

      {error && (
        <p className="text-sm font-semibold text-red-600 border-2 border-ink px-3 py-2" style={{ boxShadow: "2px 2px 0 0 rgb(40,40,37)" }}>
          {error}
        </p>
      )}

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="btn-primary flex-1 py-3 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Се зачувува..." : initial ? "Ажурирај" : "Зачувај →"}
        </button>
        {onCancel && (
          <button type="button" onClick={onCancel} className="btn-white px-6 py-3">
            Откажи
          </button>
        )}
      </div>
    </form>
  );
}
