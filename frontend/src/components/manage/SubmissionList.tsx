"use client";

import { useState } from "react";
import type { SalarySubmission, SubmissionInput } from "@/types";
import {
  ARRANGEMENT_LABEL,
  BONUS_FREQ_LABEL,
  BONUS_TYPE_LABEL,
  SENIORITY_LABEL,
} from "@/lib/constants";
import SubmissionForm from "@/components/submission/SubmissionForm";

interface Props {
  submissions: SalarySubmission[];
  onUpdate: (id: string, data: SubmissionInput) => Promise<void>;
  onDelete: (id: string) => Promise<void>;
}

function formatMKD(amount: number): string {
  return new Intl.NumberFormat("mk-MK", {
    style: "currency",
    currency: "MKD",
    maximumFractionDigits: 0,
  }).format(amount);
}

export default function SubmissionList({ submissions, onUpdate, onDelete }: Props) {
  const [editingId, setEditingId] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [deleteLoading, setDeleteLoading] = useState(false);

  if (submissions.length === 0) {
    return null;
  }

  const handleDelete = async (id: string) => {
    setDeleteLoading(true);
    try {
      await onDelete(id);
      setDeletingId(null);
    } finally {
      setDeleteLoading(false);
    }
  };

  return (
    <div className="space-y-4">
      {submissions.map((s) => (
        <div key={s.id} className="card">
          {editingId === s.id ? (
            <div>
              <p className="label mb-3">Уреди запис</p>
              <SubmissionForm
                initial={s}
                onSubmit={async (data) => {
                  await onUpdate(s.id, data);
                  setEditingId(null);
                }}
                onCancel={() => setEditingId(null)}
              />
            </div>
          ) : (
            <div>
              {/* Header */}
              <div className="flex justify-between items-start gap-4">
                <div>
                  <div className="text-lg font-black text-ink">{s.company_name}</div>
                  <div className="font-semibold opacity-60">{s.job_title}</div>
                </div>
                <div className="text-right shrink-0">
                  <div className="text-2xl font-black text-ink">{formatMKD(s.base_salary)}</div>
                  <div className="label">месечно</div>
                </div>
              </div>

              {/* Tags */}
              <div className="mt-3 flex flex-wrap gap-2">
                <span className="tag tag-lime">{s.industry_name}</span>
                <span className="tag tag-yellow">{s.city_name}</span>
                <span className="tag tag-pink">{SENIORITY_LABEL[s.seniority] || s.seniority}</span>
                <span className="tag tag-green">{ARRANGEMENT_LABEL[s.work_arrangement] || s.work_arrangement}</span>
                <span className="tag tag-white">{s.years_experience} год. искуство</span>
                <span className="tag tag-white">{s.salary_year}</span>
              </div>

              {/* Bonuses */}
              {s.bonuses && s.bonuses.length > 0 && (
                <div className="mt-3 pt-3" style={{ borderTop: "1px solid rgba(40,40,37,0.15)" }}>
                  <p className="label mb-1">Бонуси</p>
                  <div className="space-y-0.5">
                    {s.bonuses.map((b) => (
                      <div key={b.id} className="text-sm font-medium">
                        {BONUS_TYPE_LABEL[b.bonus_type] || b.bonus_type} — {formatMKD(b.amount)} ({BONUS_FREQ_LABEL[b.frequency] || b.frequency})
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Actions */}
              <div className="mt-4 flex gap-2 justify-end">
                {deletingId === s.id ? (
                  <>
                    <span className="text-sm font-semibold self-center opacity-60">Сигурни сте?</span>
                    <button
                      onClick={() => handleDelete(s.id)}
                      disabled={deleteLoading}
                      className="btn-danger disabled:opacity-50"
                    >
                      {deleteLoading ? "..." : "Избриши"}
                    </button>
                    <button onClick={() => setDeletingId(null)} className="btn-white">
                      Откажи
                    </button>
                  </>
                ) : (
                  <>
                    <button onClick={() => setEditingId(s.id)} className="btn-white text-sm">
                      Уреди
                    </button>
                    <button onClick={() => setDeletingId(s.id)} className="btn-danger text-sm">
                      Избриши
                    </button>
                  </>
                )}
              </div>
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
