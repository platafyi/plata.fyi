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

function CityContent({ slug }: { slug: string }) {
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
  const [cityName, setCityName] = useState<string>(slug);

  useEffect(() => {
    Promise.all([getIndustries(), getCities()]).then(([ind, cit]) => {
      setIndustries(ind);
      setCities(cit);
      const found = cit.find((c) => c.slug === slug);
      if (found) setCityName(found.name);
    });
  }, [slug]);

  // Strip redundant ?city= param — the slug already encodes it
  useEffect(() => {
    if (searchParams.get("city")) {
      const params = new URLSearchParams(searchParams.toString());
      params.delete("city");
      const qs = params.toString();
      router.replace(`/city/${slug}${qs ? `?${qs}` : ""}`);
    }
  }, [searchParams, router, slug]);

  const filters = {
    industry: searchParams.get("industry") || undefined,
    city: searchParams.get("city") || slug,
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
      getSalaryStats("industry", filters),
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
      router.push(`/city/${slug}?${params.toString()}`);
    },
    [router, searchParams, slug]
  );

  const currentPage = parseInt(searchParams.get("page") || "1");

  return (
    <div className="space-y-6">
      <div className="py-8">
        <Link href="/" className="label hover:underline">← Сите плати</Link>
        <h1 className="text-3xl sm:text-5xl font-black tracking-tight leading-none text-ink mt-3">
          Плати во {cityName}
        </h1>
        <p className="mt-2 text-base font-medium opacity-60">
          Анонимни плати споделени од вработени во овој град.
        </p>
      </div>

      {stats.length > 0 && <SalaryStatsComponent stats={stats} groupBy="industry" />}

      <FilterBar industries={industries} cities={cities} lockedCity={slug} />

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

export default function CityPage() {
  const { slug } = useParams<{ slug: string }>();
  return (
    <Suspense fallback={<div className="text-center py-16 text-gray-500">Се вчитува...</div>}>
      <CityContent slug={slug} />
    </Suspense>
  );
}