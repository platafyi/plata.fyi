import type {
  City,
  Company,
  CompaniesResponse,
  Industry,
  SalaryListResponse,
  SalaryStats,
  SalarySubmission,
  SearchFilters,
  SubmissionInput,
} from "@/types";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  if (!res.ok) {
    let errMsg = `HTTP ${res.status}`;
    try {
      const body = await res.json();
      if (body.error) errMsg = body.error;
    } catch {
      // ignore parse error
    }
    throw new Error(errMsg);
  }

  if (res.status === 204 || res.headers.get("content-length") === "0") {
    return undefined as T;
  }

  return res.json() as Promise<T>;
}

function authHeaders(token: string): HeadersInit {
  return { Authorization: `Bearer ${token}` };
}

// ── Meta ─────────────────────────────────────────────────────────────────────

export function getIndustries(): Promise<Industry[]> {
  return apiFetch("/api/industries");
}

export function getCities(): Promise<City[]> {
  return apiFetch("/api/cities");
}

// ── Salaries ──────────────────────────────────────────────────────────────────

export function getSalaries(filters: SearchFilters): Promise<SalaryListResponse> {
  const params = new URLSearchParams();
  Object.entries(filters).forEach(([k, v]) => {
    if (v) params.set(k, v);
  });
  return apiFetch(`/api/salaries?${params.toString()}`);
}

export function getSubmission(id: string): Promise<SalarySubmission> {
  return apiFetch(`/api/salaries/${encodeURIComponent(id)}`);
}

export function getSalaryStats(
  groupBy?: string,
  filters?: SearchFilters
): Promise<SalaryStats[]> {
  const params = new URLSearchParams();
  if (groupBy) params.set("group_by", groupBy);
  if (filters) {
    Object.entries(filters).forEach(([k, v]) => {
      if (v) params.set(k, v);
    });
  }
  return apiFetch(`/api/salaries/stats?${params.toString()}`);
}

// ── Companies ─────────────────────────────────────────────────────────────────

export function searchCompanies(q: string): Promise<CompaniesResponse> {
  return apiFetch(`/api/companies?q=${encodeURIComponent(q)}`);
}

export function searchJobTitles(q: string): Promise<string[]> {
  return apiFetch(`/api/job-titles?q=${encodeURIComponent(q)}`);
}

// ── Auth ──────────────────────────────────────────────────────────────────────

export function logout(token: string): Promise<void> {
  return apiFetch("/api/auth/session", {
    method: "DELETE",
    headers: authHeaders(token),
  });
}

// ── Submissions ───────────────────────────────────────────────────────────────

export function getMySubmissions(token: string): Promise<SalarySubmission[]> {
  return apiFetch("/api/submissions", {
    headers: authHeaders(token),
  });
}

export function createSubmission(
  token: string | null,
  data: SubmissionInput,
  turnstileToken?: string
): Promise<SalarySubmission & { session_token?: string }> {
  const body = turnstileToken ? { ...data, turnstile_token: turnstileToken } : data;
  return apiFetch("/api/submissions", {
    method: "POST",
    headers: token ? authHeaders(token) : {},
    body: JSON.stringify(body),
  });
}

export function updateSubmission(
  token: string,
  id: string,
  data: SubmissionInput
): Promise<void> {
  return apiFetch(`/api/submissions/${id}`, {
    method: "PUT",
    headers: authHeaders(token),
    body: JSON.stringify(data),
  });
}

export function deleteSubmission(token: string, id: string): Promise<void> {
  return apiFetch(`/api/submissions/${id}`, {
    method: "DELETE",
    headers: authHeaders(token),
  });
}

// ── Session helpers ───────────────────────────────────────────────────────────

export const SESSION_KEY = "platafyi_session";

export function getStoredSession(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(SESSION_KEY);
}

export function storeSession(token: string): void {
  localStorage.setItem(SESSION_KEY, token);
}

export function clearSession(): void {
  localStorage.removeItem(SESSION_KEY);
}
