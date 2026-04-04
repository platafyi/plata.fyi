"use client";

import type { BonusInput } from "@/types";
import { BONUS_FREQUENCIES, BONUS_TYPES } from "@/lib/constants";
import Select from "@/components/ui/Select";

interface Props {
  bonuses: BonusInput[];
  onChange: (bonuses: BonusInput[]) => void;
}

const emptyBonus = (): BonusInput => ({
  bonus_type: "annual",
  amount: 0,
  frequency: "annual",
});

export default function BonusFields({ bonuses, onChange }: Props) {
  const add = () => onChange([...bonuses, emptyBonus()]);
  const remove = (i: number) => onChange(bonuses.filter((_, idx) => idx !== i));
  const update = (i: number, field: keyof BonusInput, value: string | number) => {
    const next = [...bonuses];
    next[i] = { ...next[i], [field]: value };
    onChange(next);
  };

  return (
    <div className="space-y-3">
      {bonuses.map((b, i) => (
        <div key={i} className="grid grid-cols-1 sm:grid-cols-3 gap-2">
          <Select
            value={b.bonus_type}
            onChange={(v) => update(i, "bonus_type", String(v))}
            options={BONUS_TYPES.map((t) => ({ value: t.value, label: t.label }))}
            accentColor="#fe91e6"
          />

          <input
            type="number"
            value={b.amount || ""}
            onChange={(e) => update(i, "amount", parseInt(e.target.value) || 0)}
            placeholder="Износ (МКД)"
            min="1"
            className="input"
          />

          <div className="flex gap-2">
            <div style={{ flex: 1 }}>
              <Select
                value={b.frequency}
                onChange={(v) => update(i, "frequency", String(v))}
                options={BONUS_FREQUENCIES.map((f) => ({ value: f.value, label: f.label }))}
                accentColor="#38ed81"
              />
            </div>
            <button
              type="button"
              onClick={() => remove(i)}
              className="btn-danger px-3"
              aria-label="Отстрани бонус"
            >
              ×
            </button>
          </div>
        </div>
      ))}

      <button type="button" onClick={add} className="btn-white text-sm">
        + Додади бонус
      </button>
    </div>
  );
}
