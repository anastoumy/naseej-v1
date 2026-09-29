import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatCurrency(amount: number, currency = "SAR"): string {
  return new Intl.NumberFormat("ar-SA", {
    style: "currency",
    currency,
    minimumFractionDigits: 2,
  }).format(amount);
}

export function formatSAR(amount: number): string {
  return `${amount.toFixed(2)} SAR`;
}

export function formatTime(isoString: string): string {
  return new Intl.DateTimeFormat("en-SA", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  }).format(new Date(isoString));
}

export function formatDate(isoString: string): string {
  return new Intl.DateTimeFormat("en-SA", {
    year: "numeric",
    month: "short",
    day: "numeric",
  }).format(new Date(isoString));
}

export function formatDateTime(isoString: string): string {
  return new Intl.DateTimeFormat("en-SA", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  }).format(new Date(isoString));
}

export function generateId(): string {
  return Math.random().toString(36).substring(2, 9);
}
