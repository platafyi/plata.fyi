"use client";

import { useRouter, useSearchParams, useParams } from "next/navigation";
import { useCallback, useEffect, useState, Suspense } from "react";
import FilterBar from "@/components/salary/FilterBar";
import SalaryTable from "@/components/salary/SalaryTable";
import SalaryStatsComponent from "@/components/salary/SalaryStats";
import { getCities, getIndustries, getSalaries, getSalaryStats } from "@/lib/api";
import type { City, Industry, SalaryListResponse, SalaryStats } from "@/types";
import { PAGE_SIZE } from "@/lib/constants";
import Link from "next/link";

function IndustryContent({ slug }: { slug: string }) {
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
  const [industryName, setIndustryName] = useState<string>(slug);

  useEffect(() => {
    Promise.all([getIndustries(), getCities()]).then(([ind, cit]) => {
      setIndustries(ind);
      setCities(cit);
      const found = ind.find((i) => i.slug === slug);
      if (found) setIndustryName(found.name);
    });
  }, [slug]);

  // Strip redundant ?industry= param — the slug already encodes it
  useEffect(() => {
    if (searchParams.get("industry")) {
      const params = new URLSearchParams(searchParams.toString());
      params.delete("industry");
      const qs = params.toString();
      router.replace(`/industry/${slug}${qs ? `?${qs}` : ""}`);
    }
  }, [searchParams, router, slug]);

  const filters = {
    industry: searchParams.get("industry") || slug,
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
      getSalaryStats("city", filters),
    ])
      .then(([sal, st]) => {
        setSalaries(sal);
        setStats(st);
      })
      .catch(() => setError("Грешка при вчитување на податоците."))
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams.toString(), slug]);

  const handlePageChange = useCallback(
    (newPage: number) => {
      const params = new URLSearchParams(searchParams.toString());
      params.set("page", String(newPage));
      router.push(`/industry/${slug}?${params.toString()}`);
    },
    [router, searchParams, slug]
  );

  const currentPage = parseInt(searchParams.get("page") || "1");

  return (
    <div className="space-y-6">
      <div className="py-8">
        <Link href="/" className="label hover:underline">← Сите плати</Link>
        <h1 className="text-3xl sm:text-5xl font-black tracking-tight leading-none text-ink mt-3">
          Плати во {industryName}
        </h1>
        <p className="mt-2 text-base font-medium opacity-60">
          Анонимни плати споделени од вработени во оваа индустрија.
        </p>
      </div>

      {stats.length > 0 && <SalaryStatsComponent stats={stats} groupBy="city" />}

      <FilterBar industries={industries} cities={cities} lockedIndustry={slug} />

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

export default function IndustryPage() {
  const { slug } = useParams<{ slug: string }>();
  return (
    <Suspense fallback={<div className="text-center py-16 text-gray-500">Се вчитува...</div>}>
      <IndustryContent slug={slug} />
    </Suspense>
  );
}