import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Merges Tailwind classes with conditional support
 * @param {...ClassValue[]} inputs - Class names, conditionally applied
 * @returns {string} - Merged class names
 */
export function cn(...inputs: (string | undefined | null | false)[]): string {
  return twMerge(clsx(inputs));
}