"use client";

import { useEffect, useRef, useState } from "react";

interface Props {
  value: string;
  onChange: (value: string) => void;
  onSearch: (q: string) => Promise<string[]>;
  placeholder?: string;
  accentColor?: string;
  newItemNotice?: string;
}

export default function AutocompleteInput({
  value,
  onChange,
  onSearch,
  placeholder = "Пребарајте...",
  accentColor = "#e7fe05",
  newItemNotice,
}: Props) {
  const [results, setResults] = useState<string[]>([]);
  const [open, setOpen] = useState(false);
  const [showNotice, setShowNotice] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const pickedRef = useRef(false);

  useEffect(() => {
    if (timerRef.current) clearTimeout(timerRef.current);
    setShowNotice(false);
    if (value.length < 2) { setResults([]); setOpen(false); return; }
    timerRef.current = setTimeout(async () => {
      if (pickedRef.current) return;
      try {
        const res = await onSearch(value);
        setResults(res);
        setOpen(res.length > 0);
      } catch {
        setResults([]);
      }
    }, 300);
  }, [value, onSearch]);

  useEffect(() => {
    function handle(e: MouseEvent) {
      if (!containerRef.current?.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, []);

  const select = (r: string) => {
    pickedRef.current = true;
    if (timerRef.current) clearTimeout(timerRef.current);
    onChange(r);
    setShowNotice(false);
    setOpen(false);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    pickedRef.current = false;
    setShowNotice(false);
    onChange(e.target.value);
  };

  const handleBlur = () => {
    if (!newItemNotice) return;
    setTimeout(() => {
      setOpen(false);
      if (!pickedRef.current && value.length >= 2) {
        setShowNotice(true);
      }
      pickedRef.current = false;
    }, 150);
  };

  return (
    <div ref={containerRef} style={{ position: "relative" }}>
      <input
        type="text"
        value={value}
        onChange={handleChange}
        onFocus={() => { setShowNotice(false); results.length > 0 && setOpen(true); }}
        onBlur={handleBlur}
        placeholder={placeholder}
        className="input"
        autoComplete="off"
      />

      {open && (
        <div style={{
          position: "absolute",
          top: "calc(100% + 4px)",
          left: 0,
          right: 0,
          backgroundColor: "#ffffff",
          border: "2px solid rgb(40,40,37)",
          borderRadius: "4px",
          boxShadow: "4px 4px 0 0 rgb(40,40,37)",
          zIndex: 50,
          maxHeight: 220,
          overflowY: "auto",
        }}>
          {results.map((r, i) => (
            <button
              key={i}
              type="button"
              onMouseDown={() => select(r)}
              style={{
                display: "block",
                width: "100%",
                textAlign: "left",
                padding: "10px 14px",
                backgroundColor: "#ffffff",
                fontFamily: "var(--font-jost), sans-serif",
                fontWeight: 500,
                fontSize: "0.64rem",
                color: "rgb(40,40,37)",
                cursor: "pointer",
                border: "none",
                borderTop: i > 0 ? "1px solid rgba(40,40,37,0.12)" : "none",
              }}
              onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = accentColor; }}
              onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = "#ffffff"; }}
            >
              {r}
            </button>
          ))}
        </div>
      )}

      {showNotice && newItemNotice && (
        <p style={{
          marginTop: 6,
          padding: "8px 12px",
          backgroundColor: "#f6cb44",
          border: "2px solid rgb(40,40,37)",
          borderRadius: "4px",
          boxShadow: "2px 2px 0 0 rgb(40,40,37)",
          fontFamily: "var(--font-jost), sans-serif",
          fontSize: "0.875rem",
          fontWeight: 600,
          color: "rgb(40,40,37)",
        }}>
          {newItemNotice}
        </p>
      )}
    </div>
  );
}
