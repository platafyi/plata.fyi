"use client";

import { useEffect, useRef, useState } from "react";
import { searchCompanies } from "@/lib/api";
import type { Company } from "@/types";

interface Props {
  value: string;
  regNo: string;
  onChange: (name: string, regNo: string) => void;
}

export default function CompanySearch({ value, regNo, onChange }: Props) {
  const [query, setQuery] = useState(value);
  const [results, setResults] = useState<Company[]>([]);
  const [open, setOpen] = useState(false);
  const [newCompany, setNewCompany] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const pickedRef = useRef(value.length > 0);

  useEffect(() => {
    if (timerRef.current) clearTimeout(timerRef.current);
    setNewCompany(false);
    if (query.length < 2) { setResults([]); setOpen(false); return; }
    timerRef.current = setTimeout(async () => {
      if (pickedRef.current) return;
      try {
        const res = await searchCompanies(query);
        setResults(res.results);
        setOpen(res.results.length > 0);
      } catch {
        setResults([]);
      }
    }, 300);
  }, [query]);

  useEffect(() => {
    function handle(e: MouseEvent) {
      if (!containerRef.current?.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, []);

  const select = (name: string, reg = "") => {
    pickedRef.current = true;
    if (timerRef.current) clearTimeout(timerRef.current);
    setQuery(name);
    onChange(name, reg);
    setNewCompany(false);
    setOpen(false);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    pickedRef.current = false;
    setNewCompany(false);
    setQuery(e.target.value);
    onChange(e.target.value, regNo);
  };

  const handleBlur = () => {
    setTimeout(() => {
      setOpen(false);
      if (!pickedRef.current && query.length >= 2) {
        setNewCompany(true);
      }
      pickedRef.current = false;
    }, 150);
  };

  return (
    <div ref={containerRef} style={{ position: "relative" }}>
      <input
        type="text"
        value={query}
        onChange={handleChange}
        onFocus={() => { setNewCompany(false); results.length > 0 && setOpen(true); }}
        onBlur={handleBlur}
        placeholder="Пребарајте компанија..."
        className="input"
        autoComplete="off"
      />

      {open && results.length > 0 && (
        <div style={{
          position: "absolute",
          top: "calc(100% + 4px)",
          left: 0, right: 0,
          backgroundColor: "#ffffff",
          border: "2px solid rgb(40,40,37)",
          borderRadius: "4px",
          boxShadow: "4px 4px 0 0 rgb(40,40,37)",
          zIndex: 50,
          maxHeight: 260,
          overflowY: "auto",
        }}>
          {results.map((c, i) => (
            <button
              key={i}
              type="button"
              onMouseDown={() => select(c.name, c.reg_no || "")}
              style={{
                display: "block", width: "100%", textAlign: "left",
                padding: "10px 14px",
                backgroundColor: "#ffffff",
                fontFamily: "var(--font-jost), sans-serif",
                fontWeight: 500, fontSize: "0.64rem",
                color: "rgb(40,40,37)", cursor: "pointer",
                border: "none",
                borderTop: i > 0 ? "1px solid rgba(40,40,37,0.12)" : "none",
              }}
              onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = "#fe91e6"; }}
              onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = "#ffffff"; }}
            >
              {c.name}
              {c.reg_no && <span style={{ marginLeft: 8, fontSize: "0.64rem", opacity: 0.4 }}>#{c.reg_no}</span>}
            </button>
          ))}
        </div>
      )}

      {newCompany && (
        <p style={{
          marginTop: 6,
          padding: "8px 12px",
          backgroundColor: "#f6cb44",
          border: "2px solid rgb(40,40,37)",
          borderRadius: "4px",
          boxShadow: "2px 2px 0 0 rgb(40,40,37)",
          fontFamily: "var(--font-jost), sans-serif",
          fontSize: "0.7rem",
          fontWeight: 600,
          color: "rgb(40,40,37)",
        }}>
          Не најдовме компанија, ќе креираме нова. Ве молиме внесете го целосното име.
        </p>
      )}
    </div>
  );
}
