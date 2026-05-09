"use client";

import React from "react";
import { Menu } from "@/types";
import { Plus, Pencil, Package } from "lucide-react";

// Deterministic food emoji based on menu name
const FOOD_EMOJIS = ["🍜", "🍛", "🍚", "🥘", "🍲", "🥗", "🍝", "🍱", "🍤", "🥩", "🍗", "🧆", "🥚", "🫕", "🍵", "☕", "🧋", "🥤", "🍰", "🧁"];

const getEmoji = (name: string): string => {
  let hash = 0;
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash);
  return FOOD_EMOJIS[Math.abs(hash) % FOOD_EMOJIS.length];
};

// Gradient palette based on name hash
const GRADIENTS = [
  "from-orange-400 to-amber-500",
  "from-emerald-400 to-teal-500",
  "from-sky-400 to-blue-500",
  "from-violet-400 to-purple-500",
  "from-rose-400 to-pink-500",
  "from-amber-400 to-yellow-500",
  "from-teal-400 to-cyan-500",
  "from-indigo-400 to-violet-500",
  "from-lime-400 to-green-500",
  "from-fuchsia-400 to-pink-500",
];

const getGradient = (name: string): string => {
  let hash = 0;
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash);
  return GRADIENTS[Math.abs(hash) % GRADIENTS.length];
};

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(amount);

const stationLabel: Record<string, string> = {
  kitchen: "Dapur",
  bar: "Bar",
  cashier: "Kasir",
};

interface MenuCardProps {
  menu: Menu;
  cartQty: number;
  onAdd: () => void;
  onEdit?: () => void;
  canEdit?: boolean;
}

export function MenuCard({ menu, cartQty, onAdd, onEdit, canEdit = false }: MenuCardProps) {
  const emoji = getEmoji(menu.name);
  const gradient = getGradient(menu.name);
  const isOutOfStock = menu.daily_stock === 0;
  const isLowStock = menu.daily_stock > 0 && menu.daily_stock <= 10;

  return (
    <div
      className={`relative flex flex-col bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-lg transition-all duration-200 border border-slate-100 hover:border-emerald-200 text-left group ${
        isOutOfStock ? "opacity-60" : ""
      }`}
    >
      {/* Stock badge top-left */}
      <div className="absolute top-2 left-2 z-10">
        <span
          className={`text-[10px] font-bold px-1.5 py-0.5 rounded-md backdrop-blur-sm ${
            isOutOfStock
              ? "bg-red-500/90 text-white"
              : isLowStock
              ? "bg-amber-400/90 text-white"
              : "bg-black/30 text-white"
          }`}
        >
          {isOutOfStock ? "Habis" : `S: ${menu.daily_stock}`}
        </span>
      </div>

      {/* Cart qty badge top-right */}
      {cartQty > 0 && (
        <div className="absolute top-2 right-2 z-10 w-6 h-6 bg-emerald-500 rounded-full flex items-center justify-center shadow-md ring-2 ring-white">
          <span className="text-white text-[10px] font-black">{cartQty}</span>
        </div>
      )}

      {/* Edit icon — superadmin only */}
      {canEdit && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            onEdit?.();
          }}
          className="absolute top-2 right-2 z-20 w-6 h-6 bg-white/80 hover:bg-white rounded-full flex items-center justify-center shadow transition-colors opacity-0 group-hover:opacity-100"
          title="Edit menu"
        >
          <Pencil className="w-3 h-3 text-slate-600" />
        </button>
      )}

      {/* Image / Emoji top area */}
      <button
        onClick={() => !isOutOfStock && onAdd()}
        disabled={isOutOfStock}
        className="w-full focus:outline-none"
        aria-label={`Tambah ${menu.name} ke keranjang`}
      >
        <div
          className={`h-[88px] w-full bg-gradient-to-br ${gradient} flex flex-col items-center justify-center gap-0.5 group-hover:brightness-105 transition-all`}
        >
          <span className="text-4xl drop-shadow-sm select-none">{emoji}</span>
          <span className="text-[9px] font-semibold text-white/80 uppercase tracking-widest">
            {stationLabel[menu.station] ?? menu.station}
          </span>
        </div>
      </button>

      {/* Info block */}
      <div className="p-2.5 flex flex-col gap-1">
        <h3 className="font-bold text-slate-800 text-xs leading-tight line-clamp-2 group-hover:text-emerald-700 transition-colors">
          {menu.name}
        </h3>
        <div className="flex items-center justify-between mt-0.5">
          <span className="font-black text-emerald-600 text-sm">{formatCurrency(menu.price)}</span>
          <button
            onClick={() => !isOutOfStock && onAdd()}
            disabled={isOutOfStock}
            className="w-7 h-7 rounded-lg bg-emerald-500 hover:bg-emerald-600 disabled:bg-slate-200 flex items-center justify-center shadow-sm transition-colors active:scale-90"
            aria-label="Tambah ke keranjang"
          >
            <Plus className="w-4 h-4 text-white disabled:text-slate-400" />
          </button>
        </div>
        {isLowStock && (
          <p className="text-[10px] text-amber-600 font-semibold flex items-center gap-1">
            <Package className="w-2.5 h-2.5" />
            Sisa {menu.daily_stock} porsi
          </p>
        )}
      </div>
    </div>
  );
}
