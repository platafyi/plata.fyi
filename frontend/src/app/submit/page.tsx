"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Turnstile } from "@marsidev/react-turnstile";
import SubmissionForm from "@/components/submission/SubmissionForm";
import SubmitConfirmModal from "@/components/submission/SubmitConfirmModal";
import { createSubmission, storeSession } from "@/lib/api";
import type { SubmissionInput } from "@/types";

const TURNSTILE_SITE_KEY =
  process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY || "1x00000000000000000000AA";

export default function SubmitPage() {
  const router = useRouter();
  const [turnstileToken, setTurnstileToken] = useState<string | null>(null);
  const [pendingData, setPendingData] = useState<SubmissionInput | null>(null);

  const handleFormSubmit = (data: SubmissionInput) => {
    setPendingData(data);
  };

  const handleConfirm = async () => {
    if (!pendingData) return;
    const result = await createSubmission(null, pendingData, turnstileToken ?? undefined);
    if (result.session_token) {
      storeSession(result.session_token);
    }
    router.replace("/manage");
  };

  const handleCancel = () => {
    setPendingData(null);
  };

  return (
    <>
    {pendingData && (
      <SubmitConfirmModal onConfirm={handleConfirm} onCancel={handleCancel} />
    )}
    <div className="max-w-2xl mx-auto py-4">
      <h1 className="text-4xl sm:text-6xl font-black tracking-tight text-ink mb-2">
        Споделете плата
      </h1>
      <p className="text-xl font-medium opacity-60 mb-8">
        Нема регистрација, нема лозинка. Сите записи се анонимни.
      </p>

      <div className="card-lg">
        <SubmissionForm onSubmit={handleFormSubmit} />
      </div>

      <div className="mt-4 flex justify-center">
        <Turnstile
          siteKey={TURNSTILE_SITE_KEY}
          onSuccess={setTurnstileToken}
          onError={() => setTurnstileToken(null)}
          onExpire={() => setTurnstileToken(null)}
        />
      </div>

      <div className="card mt-4" style={{ backgroundColor: "#e7fe05" }}>
        <p className="text-xs font-bold uppercase tracking-widest mb-2">Како е ова анонимно?</p>
        <ul className="space-y-1.5 text-sm font-medium mb-3">
          <li>→ Не е потребна е-маил адреса или регистрација</li>
          <li>→ Платата не може да се поврзе со вас</li>
          <li>→ Сите јавни записи се анонимни</li>
        </ul>
        <a href="/privacy" className="text-xs font-bold underline opacity-60 hover:opacity-100">
          Прочитај ја политиката за приватност →
        </a>
      </div>
    </div>
    </>
  );
}