import Decimal from 'decimal.js';

// Configure decimal.js for financial precision
Decimal.set({
  precision: 40, // High precision for very large numbers
  rounding: Decimal.ROUND_HALF_UP,
  toExpNeg: -20, // Use exponential notation for very small numbers
  toExpPos: 30,  // Use exponential notation for very large numbers
});

/**
 * Creates a new Decimal instance from a string or number
 */
export function createMoney(value: string | number): Decimal {
  return new Decimal(value);
}

/**
 * Creates a zero money value
 */
export function zeroMoney(): Decimal {
  return new Decimal(0);
}

/**
 * Adds two money values
 */
export function addMoney(a: Decimal, b: Decimal): Decimal {
  return a.plus(b);
}

/**
 * Subtracts money value b from a
 */
export function subtractMoney(a: Decimal, b: Decimal): Decimal {
  return a.minus(b);
}

/**
 * Multiplies two money values
 */
export function multiplyMoney(a: Decimal, b: Decimal): Decimal {
  return a.times(b);
}

/**
 * Divides money value a by b
 */
export function divideMoney(a: Decimal, b: Decimal): Decimal {
  return a.dividedBy(b);
}

/**
 * Compares two money values
 * Returns -1 if a < b, 0 if a == b, 1 if a > b
 */
export function compareMoney(a: Decimal, b: Decimal): number {
  return a.comparedTo(b);
}

/**
 * Checks if money value is zero
 */
export function isMoneyZero(money: Decimal): boolean {
  return money.isZero();
}

/**
 * Checks if money value is positive
 */
export function isMoneyPositive(money: Decimal): boolean {
  return money.isPositive();
}

/**
 * Checks if money value is negative
 */
export function isMoneyNegative(money: Decimal): boolean {
  return money.isNegative();
}

/**
 * Converts money to a string for display
 */
export function moneyToString(money: Decimal): string {
  return money.toFixed(2);
}

/**
 * Converts money to a string with full precision for API calls
 */
export function moneyToStringFull(money: Decimal): string {
  return money.toString();
}

/**
 * Parses a string to a money value
 */
export function parseMoney(value: string): Decimal {
  if (!value || value.trim() === '') {
    return new Decimal(0);
  }
  return new Decimal(value);
}

/**
 * Formats money for display purposes
 */
export function formatMoney(money: Decimal): string {
  // For very large numbers, use exponential notation
  if (money.gte(1e15)) {
    return money.toExponential(2);
  }
  
  // For regular numbers, use fixed notation with commas
  return new Intl.NumberFormat(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(parseFloat(money.toFixed(2)));
}

/**
 * Validates if a string is a valid money value
 */
export function isValidMoney(value: string): boolean {
  try {
    const decimal = new Decimal(value);
    return decimal.isFinite();
  } catch {
    return false;
  }
}

/**
 * Converts money to a number (use with caution for large values)
 */
export function moneyToNumber(money: Decimal): number {
  return money.toNumber();
}
