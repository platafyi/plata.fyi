"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { clearSession, getStoredSession, logout } from "@/lib/api";
import LogoutModal from "./LogoutModal";

export default function Header() {
  const [hasSession, setHasSession] = useState(false);
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);

  useEffect(() => {
    setHasSession(!!getStoredSession());
  }, []);

  const handleLogout = async () => {
    const token = getStoredSession();
    if (token) {
      try { await logout(token); } catch { /* ignore */ }
    }
    clearSession();
    setHasSession(false);
    setShowLogoutConfirm(false);
    window.location.href = "/";
  };

  return (
    <>
      <header className="bg-white sticky top-0 z-10" style={{ borderBottom: "2px solid rgb(40,40,37)" }}>
        <div className="container mx-auto px-4 max-w-6xl">
          <div className="flex items-center justify-between h-14">
            <Link href="/" className="flex items-center gap-2 text-xl font-black tracking-tight text-ink">
              <span>плата<span className="font-medium" style={{ opacity: 0.5 }}>.fyi</span></span>
            </Link>

            <nav className="flex items-center gap-3 text-sm font-semibold">
              {hasSession ? (
                <>
                  <Link href="/zosto" className="btn-primary text-xs py-1 px-2 sm:py-1.5 sm:px-4">
                    Зошто да споделам?
                  </Link>
                  <Link href="/manage" className="btn-white text-xs py-1.5 px-3" title="Мои записи">
                    <span className="hidden sm:inline">Мои записи</span>
                    <svg className="sm:hidden w-4 h-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                      <circle cx="12" cy="7" r="4" />
                    </svg>
                  </Link>
                  <button onClick={() => setShowLogoutConfirm(true)} className="btn-white text-xs py-1.5 px-3" title="Одјави се">
                    <span className="hidden sm:inline">Одјави се</span>
                    <svg className="sm:hidden w-4 h-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                      <polyline points="16 17 21 12 16 7" />
                      <line x1="21" y1="12" x2="9" y2="12" />
                    </svg>
                  </button>
                </>
              ) : (
                <Link href="/zosto" className="btn-primary text-xs py-1 px-2 sm:py-1.5 sm:px-4">
                  Зошто да споделам?
                </Link>
              )}
            </nav>
          </div>
        </div>
      </header>

      {showLogoutConfirm && (
        <LogoutModal
          onConfirm={handleLogout}
          onCancel={() => setShowLogoutConfirm(false)}
        />
      )}
    </>
  );
}