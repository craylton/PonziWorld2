export const CURRENCY_SYMBOL = 'Ƹ'; // Default fictional currency symbol

/**
 * Formats a number as a currency string using the fictional currency symbol.
 */
export function formatCurrency(amount: number): string {
  const formatted = new Intl.NumberFormat(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
  return `${CURRENCY_SYMBOL}${formatted}`;
}
