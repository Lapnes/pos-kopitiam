const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3400/api/v1';

export async function apiFetch<T>(path: string, options?: RequestInit): Promise<T | null> {
  try {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: { 'Content-Type': 'application/json' },
      ...options,
    });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export function formatRp(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency', currency: 'IDR',
    minimumFractionDigits: 0, maximumFractionDigits: 0,
  }).format(amount);
}

// Sage green + warm cream palette untuk category badges
export const CATEGORY_COLORS: Record<string, { bg: string; text: string }> = {
  'Kopi':     { bg: 'rgba(90,122,90,0.15)',   text: '#3d5c3d' },  // sage
  'Non-Kopi': { bg: 'rgba(138,110,58,0.15)',  text: '#6b530f' },  // olive-gold
  'Makanan':  { bg: 'rgba(45,122,74,0.12)',   text: '#2d7a4a' },  // forest green
  'Minuman':  { bg: 'rgba(90,140,110,0.15)',  text: '#3d7055' },  // teal-sage
  'Dessert':  { bg: 'rgba(160,100,60,0.12)',  text: '#9a5a30' },  // terra cotta warm
};

export function getCategoryStyle(name: string) {
  const key = Object.keys(CATEGORY_COLORS).find(k => name?.toLowerCase().includes(k.toLowerCase()));
  return key ? CATEGORY_COLORS[key] : { bg: 'rgba(90,122,90,0.10)', text: '#5a7a5a' };
}

export const MENU_EMOJIS: Record<string, string> = {
  'kopi': '☕', 'latte': '☕', 'cappuccino': '☕', 'espresso': '☕',
  'americano': '☕', 'teh': '🍵', 'matcha': '🍵', 'milo': '🧆',
  'cokelat': '🍫', 'chocolate': '🍫', 'jus': '🥤', 'juice': '🥤',
  'smoothie': '🥤', 'air': '💧', 'nasi': '🍚', 'mie': '🍜',
  'roti': '🥐', 'sandwich': '🥪', 'cake': '🎂', 'kue': '🍰',
  'pudding': '🍮', 'waffle': '🧇', 'pancake': '🥞',
};

export function getMenuEmoji(name: string): string {
  const lower = name.toLowerCase();
  for (const [key, emoji] of Object.entries(MENU_EMOJIS)) {
    if (lower.includes(key)) return emoji;
  }
  return '🍽️';
}
