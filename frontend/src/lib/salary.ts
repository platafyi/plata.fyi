// Neto Bruto estimation based on Macedonian tax law.
// Formula derived from kalkulator.ot.mk by Skopje Hacklab (AGPL-3.0).

const PERSONAL_EXEMPTION = 10_932;
const CONTRIBUTION_RATE = 0.28;
const TAX_RATE = 0.10;
const TAX_CREDIT = TAX_RATE * PERSONAL_EXEMPTION; // 1,093.2

const NET_RATE = 1 - CONTRIBUTION_RATE - TAX_RATE * (1 - CONTRIBUTION_RATE);
// = 1 - 0.28 - 0.10 * 0.72 = 0.648

export function netoToBruto(neto: number): number {
  return Math.round((neto + TAX_CREDIT) / NET_RATE);
}

export function formatMKD(amount: number): string {
  return new Intl.NumberFormat("mk-MK", {
    style: "currency",
    currency: "MKD",
    maximumFractionDigits: 0,
  }).format(amount);
}