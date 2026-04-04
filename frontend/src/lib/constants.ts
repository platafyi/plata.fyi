export const SENIORITY_LEVELS = [
  { value: "junior", label: "Почетник" },
  { value: "mid", label: "Средно ниво" },
  { value: "senior", label: "Високо ниво" },
  { value: "lead", label: "Технички лидер" },
  { value: "manager", label: "Менаџер" },
  { value: "director", label: "Директор" },
  { value: "executive", label: "Извршен директор / ВП" },
] as const;

export const WORK_ARRANGEMENTS = [
  { value: "office", label: "Канцелариски" },
  { value: "hybrid", label: "Хибриден" },
  { value: "remote", label: "Remote" },
] as const;

export const BONUS_TYPES = [
  { value: "annual", label: "Годишен бонус" },
  { value: "performance", label: "Бонус за успешност" },
  { value: "signing", label: "Бонус при вработување" },
  { value: "project", label: "Проектен бонус" },
  { value: "other", label: "Друго" },
] as const;

export const BONUS_FREQUENCIES = [
  { value: "monthly", label: "Месечно" },
  { value: "quarterly", label: "Квартално" },
  { value: "annual", label: "Годишно" },
  { value: "one_time", label: "Еднократно" },
] as const;

export const SENIORITY_LABEL: Record<string, string> = Object.fromEntries(
  SENIORITY_LEVELS.map((s) => [s.value, s.label])
);

export const ARRANGEMENT_LABEL: Record<string, string> = Object.fromEntries(
  WORK_ARRANGEMENTS.map((a) => [a.value, a.label])
);

export const BONUS_TYPE_LABEL: Record<string, string> = Object.fromEntries(
  BONUS_TYPES.map((b) => [b.value, b.label])
);

export const BONUS_FREQ_LABEL: Record<string, string> = Object.fromEntries(
  BONUS_FREQUENCIES.map((f) => [f.value, f.label])
);

export const PAGE_SIZE = 20;
