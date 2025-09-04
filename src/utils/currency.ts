import Decimal from 'decimal.js';
import { formatMoney, parseMoney } from './money';

export const CURRENCY_SYMBOL = 'Ƹ'; // Default fictional currency symbol

/**
 * Formats a number as a currency string using the fictional currency symbol.
 * @deprecated Use formatCurrencyFromDecimal instead for better precision
 */
export function formatCurrency(amount: number): string {
  const formatted = new Intl.NumberFormat(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
  return `${CURRENCY_SYMBOL}${formatted}`;
}

/**
 * Formats a Decimal amount as a currency string using the fictional currency symbol.
 */
export function formatCurrencyFromDecimal(amount: Decimal): string {
  const formatted = formatMoney(amount);
  return `${CURRENCY_SYMBOL}${formatted}`;
}

/**
 * Formats a string amount as a currency string using the fictional currency symbol.
 */
export function formatCurrencyFromString(amount: string): string {
  const decimal = parseMoney(amount);
  return formatCurrencyFromDecimal(decimal);
}
