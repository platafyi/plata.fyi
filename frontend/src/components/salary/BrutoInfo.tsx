"use client";

import { useState, useRef, useEffect } from "react";
import { netoToBruto, formatMKD } from "@/lib/salary";

export function BrutoInfo({ neto }: { neto: number }) {
  const bruto = netoToBruto(neto);
  const [pos, setPos] = useState<{ bottom: number; right: number } | null>(null);
  const btnRef = useRef<HTMLSpanElement>(null);

  useEffect(() => {
    if (!pos) return;
    function handleClick(e: MouseEvent) {
      if (btnRef.current && !btnRef.current.contains(e.target as Node)) {
        setPos(null);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, [pos]);

  function handleToggle(e: React.MouseEvent) {
    e.stopPropagation();
    if (pos) {
      setPos(null);
      return;
    }
    if (btnRef.current) {
      const rect = btnRef.current.getBoundingClientRect();
      setPos({
        bottom: window.innerHeight - rect.top + 8,
        right: window.innerWidth - rect.right,
      });
    }
  }

  return (
    <span ref={btnRef} className="relative inline-flex items-center ml-1.5 align-middle">
      <span
        onClick={handleToggle}
        className="inline-flex items-center justify-center w-5 h-5 rounded-full text-xs font-bold cursor-pointer select-none"
        style={{ backgroundColor: "rgba(40,40,37,0.12)", color: "rgba(40,40,37,0.55)" }}
      >
        i
      </span>
      {pos && (
        <span
          className="w-64 text-sm p-3.5 text-left"
          style={{
            position: "fixed",
            bottom: pos.bottom,
            right: pos.right,
            zIndex: 9999,
            border: "2px solid rgb(40,40,37)",
            borderRadius: "4px",
            backgroundColor: "#fff",
            boxShadow: "3px 3px 0 0 rgb(40,40,37)",
          }}
        >
          <span className="font-black text-base block">
            <span style={{ color: "rgba(40,40,37,0.45)" }}>Бруто: </span>
            <span style={{ color: "#16a34a" }}>~{formatMKD(bruto)}</span>
          </span>
          <span className="opacity-40 block mt-2 text-xs leading-tight">
            По МК даночен закон (придонеси 28% + ДЛД 10%).{" "}
            <a
              href="https://github.com/skopjehacklab/kalkulator.ot.mk"
              target="_blank"
              rel="noopener noreferrer"
              className="underline"
              onClick={(e) => e.stopPropagation()}
            >
              kalkulator.ot.mk
            </a>
          </span>
        </span>
      )}
    </span>
  );
}