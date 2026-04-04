"use client";

import { useEffect, useRef, useState } from "react";

interface Option {
  value: string | number;
  label: string;
}

interface Props {
  value: string | number;
  onChange: (value: string | number) => void;
  options: Option[];
  placeholder?: string;
  accentColor?: string;
}

export default function Select({
  value,
  onChange,
  options,
  placeholder = "Изберете...",
  accentColor = "#ffffff",
}: Props) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const selected = options.find((o) => String(o.value) === String(value));

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  return (
    <div ref={ref} style={{ position: "relative", width: "100%" }}>
      {/* Trigger */}
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        style={{
          width: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          backgroundColor: selected ? accentColor : "#ffffff",
          border: "2px solid rgb(40,40,37)",
          borderRadius: "4px",
          boxShadow: open ? "none" : "3px 3px 0 0 rgb(40,40,37)",
          transform: open ? "translate(2px,2px)" : "none",
          padding: "10px 14px",
          fontFamily: "var(--font-jost), sans-serif",
          fontWeight: selected ? 700 : 500,
          fontSize: "0.8rem",
          color: "rgb(40,40,37)",
          cursor: "pointer",
          textAlign: "left",
          transition: "all 80ms ease",
        }}
      >
        <span style={{ opacity: selected ? 1 : 0.4 }}>
          {selected ? selected.label : placeholder}
        </span>
        <span style={{ fontWeight: 700, marginLeft: 8, transition: "transform 120ms ease", display: "inline-block", transform: open ? "rotate(180deg)" : "none" }}>
          ↓
        </span>
      </button>

      {/* Dropdown */}
      {open && (
        <div
          style={{
            position: "absolute",
            top: "calc(100% + 4px)",
            left: 0,
            right: 0,
            backgroundColor: "#ffffff",
            border: "2px solid rgb(40,40,37)",
            borderRadius: "4px",
            boxShadow: "4px 4px 0 0 rgb(40,40,37)",
            zIndex: 50,
            maxHeight: 260,
            overflowY: "auto",
          }}
        >
          {options.map((o, i) => {
            const isSelected = String(o.value) === String(value);
            return (
              <button
                key={o.value}
                type="button"
                onClick={() => { onChange(o.value); setOpen(false); }}
                style={{
                  display: "block",
                  width: "100%",
                  textAlign: "left",
                  padding: "10px 14px",
                  backgroundColor: isSelected ? accentColor : "#ffffff",
                  fontFamily: "var(--font-jost), sans-serif",
                  fontWeight: isSelected ? 700 : 500,
                  fontSize: "0.8rem",
                  color: "rgb(40,40,37)",
                  cursor: "pointer",
                  border: "none",
                  borderTop: i > 0 ? "1px solid rgba(40,40,37,0.15)" : "none",
                }}
                onMouseEnter={(e) => {
                  if (!isSelected) (e.currentTarget as HTMLButtonElement).style.backgroundColor = "#f5f5f5";
                }}
                onMouseLeave={(e) => {
                  if (!isSelected) (e.currentTarget as HTMLButtonElement).style.backgroundColor = "#ffffff";
                }}
              >
                {o.label}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
