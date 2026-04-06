"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useCallback, useEffect, useState, Suspense } from "react";
import FilterBar from "@/components/salary/FilterBar";
import SalaryTable from "@/components/salary/SalaryTable";
import SalaryStatsComponent from "@/components/salary/SalaryStats";
import { getCities, getIndustries, getSalaries, getSalaryStats } from "@/lib/api";
import type { City, Industry, SalaryListResponse, SalaryStats } from "@/types";
import { PAGE_SIZE } from "@/lib/constants";
import Link from "next/link";

function HomeContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [industries, setIndustries] = useState<Industry[]>([]);
  const [cities, setCities] = useState<City[]>([]);
  const [salaries, setSalaries] = useState<SalaryListResponse>({
    data: [],
    total: 0,
    page: 1,
    page_size: PAGE_SIZE,
  });
  const [stats, setStats] = useState<SalaryStats[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Load meta on mount
  useEffect(() => {
    Promise.all([getIndustries(), getCities()]).then(([ind, cit]) => {
      setIndustries(ind);
      setCities(cit);
    });
  }, []);

  const filters = {
    industry: searchParams.get("industry") || undefined,
    city: searchParams.get("city") || undefined,
    seniority: searchParams.get("seniority") || undefined,
    arrangement: searchParams.get("arrangement") || undefined,
    min_salary: searchParams.get("min_salary") || undefined,
    max_salary: searchParams.get("max_salary") || undefined,
    page: searchParams.get("page") || "1",
    page_size: String(PAGE_SIZE),
  };

  useEffect(() => {
    setLoading(true);
    setError(null);

    Promise.all([
      getSalaries(filters),
      getSalaryStats(searchParams.get("industry") ? "city" : "industry", filters),
    ])
      .then(([sal, st]) => {
        setSalaries(sal);
        setStats(st);
      })
      .catch(() => setError("Грешка при вчитување на податоците."))
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams.toString()]);

  const handlePageChange = useCallback(
    (newPage: number) => {
      const params = new URLSearchParams(searchParams.toString());
      params.set("page", String(newPage));
      router.push(`/?${params.toString()}`);
    },
    [router, searchParams]
  );

  const currentPage = parseInt(searchParams.get("page") || "1");
  const groupBy = searchParams.get("industry") ? "city" : "industry";

  return (
    <div className="space-y-6">
      {/* Hero */}
      <div className="py-10 flex flex-col sm:flex-row sm:items-end sm:justify-between gap-6">
        <div>
          <p className="label mb-2">Плати во Македонија</p>
          <h1 className="text-4xl sm:text-7xl font-black tracking-tight leading-none text-ink">
            Чиста сметка,<br/> долга љубов
          </h1>
          <p className="mt-3 text-xl font-medium opacity-60 max-w-md">
            Анонимни плати споделени од вработени низ цела Македонија. Без регистрација, без лични податоци.
          </p>
        </div>
        <Link href="/submit" className="btn-primary whitespace-nowrap self-start sm:self-auto">
          Споделете плата →
        </Link>
      </div>

      {/* Stats */}
      {stats.length > 0 && <SalaryStatsComponent stats={stats} groupBy={groupBy} />}

      {/* Filters */}
      <FilterBar industries={industries} cities={cities} />

      {/* Table */}
      <div id="salary-table" />
      {loading ? (
        <div className="text-center py-16 text-gray-500">Се вчитува...</div>
      ) : error ? (
        <div className="text-center py-16 text-red-500">{error}</div>
      ) : (
        <SalaryTable
          submissions={salaries.data}
          total={salaries.total}
          page={currentPage}
          pageSize={PAGE_SIZE}
          onPageChange={handlePageChange}
        />
      )}
    </div>
  );
}

export default function HomePage() {
  return (
    <Suspense fallback={<div className="text-center py-16 text-gray-500">Се вчитува...</div>}>
      <HomeContent />
    </Suspense>
  );
}