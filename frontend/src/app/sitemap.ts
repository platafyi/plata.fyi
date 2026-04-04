import type { MetadataRoute } from "next";
import { getIndustries, getCities, getSalaries } from "@/lib/api";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const base = "https://plata.fyi";

  const [industries, cities, salaries] = await Promise.all([
    getIndustries().catch(() => []),
    getCities().catch(() => []),
    getSalaries({ page_size: "500" }).catch(() => ({ data: [] })),
  ]);

  const staticRoutes: MetadataRoute.Sitemap = [
    { url: base, changeFrequency: "daily", priority: 1 },
    { url: `${base}/submit`, changeFrequency: "monthly", priority: 0.8 },
    { url: `${base}/privacy`, changeFrequency: "yearly", priority: 0.3 },
  ];

  const industryRoutes: MetadataRoute.Sitemap = industries.map((i) => ({
    url: `${base}/industry/${i.slug}`,
    changeFrequency: "weekly",
    priority: 0.7,
  }));

  const cityRoutes: MetadataRoute.Sitemap = cities.map((c) => ({
    url: `${base}/city/${c.slug}`,
    changeFrequency: "weekly",
    priority: 0.7,
  }));

  const submissionRoutes: MetadataRoute.Sitemap = salaries.data.map((s) => ({
    url: `${base}/s/${s.id}`,
    changeFrequency: "monthly",
    priority: 0.5,
    lastModified: new Date(s.created_at),
  }));

  return [...staticRoutes, ...industryRoutes, ...cityRoutes, ...submissionRoutes];
}
