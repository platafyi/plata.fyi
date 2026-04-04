"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import SubmissionList from "@/components/manage/SubmissionList";
import SubmissionForm from "@/components/submission/SubmissionForm";
import {
  clearSession,
  createSubmission,
  deleteSubmission,
  getMySubmissions,
  getStoredSession,
  updateSubmission,
} from "@/lib/api";
import type { SalarySubmission, SubmissionInput } from "@/types";
import LogoutModal from "@/components/layout/LogoutModal";

function ManageContent() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [submissions, setSubmissions] = useState<SalarySubmission[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);

  useEffect(() => {
    const t = getStoredSession();
    if (!t) {
      router.replace("/submit");
      return;
    }
    setToken(t);

    getMySubmissions(t)
      .then(setSubmissions)
      .catch((err) => {
        if (err instanceof Error && err.message.includes("401")) {
          clearSession();
          router.replace("/submit");
        } else {
          setError("Грешка при вчитување на записите.");
        }
      })
      .finally(() => setLoading(false));
  }, [router]);

  const handleCreate = async (data: SubmissionInput) => {
    if (!token) return;
    await createSubmission(token, data);
    const updated = await getMySubmissions(token);
    setSubmissions(updated);
    setShowForm(false);
  };

  const handleUpdate = async (id: string, data: SubmissionInput) => {
    if (!token) return;
    await updateSubmission(token, id, data);
    const updated = await getMySubmissions(token);
    setSubmissions(updated);
  };

  const handleDelete = async (id: string) => {
    if (!token) return;
    await deleteSubmission(token, id);
    setSubmissions((prev) => prev.filter((s) => s.id !== id));
  };

  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);

  const handleLogout = () => {
    clearSession();
    router.push("/");
  };

  if (loading) return <div className="text-center py-20 font-semibold opacity-50">Се вчитува...</div>;

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <p className="label mb-1">Мои записи</p>
          <h1 className="text-4xl font-black tracking-tight text-ink">
            {submissions.length} {submissions.length === 1 ? "запис" : "записи"}
          </h1>
        </div>
        <button onClick={() => setShowLogoutConfirm(true)} className="btn-white text-sm">
          Одјави се
        </button>
      </div>

      {showLogoutConfirm && (
        <LogoutModal
          onConfirm={handleLogout}
          onCancel={() => setShowLogoutConfirm(false)}
        />
      )}

      {error && (
        <p className="text-sm font-semibold text-red-600 border-2 border-ink px-3 py-2" style={{ boxShadow: "2px 2px 0 0 rgb(40,40,37)" }}>
          {error}
        </p>
      )}

      {showForm ? (
        <div className="card-lg">
          <p className="label mb-3">Нов запис</p>
          <SubmissionForm
            onSubmit={handleCreate}
            onCancel={() => setShowForm(false)}
          />
        </div>
      ) : (
        <button
          onClick={() => setShowForm(true)}
          className="btn-black w-full py-3"
        >
          + Додади нов запис
        </button>
      )}

      <SubmissionList
        submissions={submissions}
        onUpdate={handleUpdate}
        onDelete={handleDelete}
      />
    </div>
  );
}

export default function ManagePage() {
  return (
    <Suspense fallback={<div className="text-center py-20 font-semibold opacity-50">Се вчитува...</div>}>
      <ManageContent />
    </Suspense>
  );
}
