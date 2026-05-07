'use client';

import { Plus } from 'lucide-react';
import { Menu } from '@/types';
import { formatRp, getCategoryStyle, getMenuEmoji } from '@/lib/api';

interface MenuCardProps {
  menu: Menu;
  onAdd: (menuId: number) => void;
}

export default function MenuCard({ menu, onAdd }: MenuCardProps) {
  const isOutOfStock = menu.daily_stock <= 0;
  const isLowStock   = menu.daily_stock > 0 && menu.daily_stock <= 5;
  const catStyle     = getCategoryStyle(menu.category?.category_name || '');
  const emoji        = getMenuEmoji(menu.menu_name);

  return (
    <div
      onClick={() => !isOutOfStock && onAdd(menu.menu_id)}
      className="group relative flex flex-col rounded-2xl border transition-all duration-200 overflow-hidden"
      style={{
        background:   'var(--surface)',
        borderColor:  'var(--border)',
        cursor:       isOutOfStock ? 'not-allowed' : 'pointer',
        opacity:      isOutOfStock ? 0.5 : 1,
        boxShadow:    '0 1px 4px rgba(61,92,61,0.06)',
      }}
      onMouseEnter={(e) => {
        if (!isOutOfStock) {
          const el = e.currentTarget as HTMLElement;
          el.style.borderColor = 'var(--caramel)';
          el.style.transform   = 'translateY(-3px)';
          el.style.boxShadow   = '0 8px 28px rgba(61,92,61,0.13)';
        }
      }}
      onMouseLeave={(e) => {
        const el = e.currentTarget as HTMLElement;
        el.style.borderColor = 'var(--border)';
        el.style.transform   = 'translateY(0)';
        el.style.boxShadow   = '0 1px 4px rgba(61,92,61,0.06)';
      }}
    >
      {/* Caramel top stripe on hover */}
      <div
        className="absolute top-0 left-0 right-0 h-[3px] opacity-0 group-hover:opacity-100 transition-opacity duration-200"
        style={{ background: 'linear-gradient(90deg, var(--caramel), var(--caramel-glow))' }}
      />

      <div className="p-4 flex flex-col flex-1">
        {/* Emoji icon */}
        <div
          className="w-11 h-11 rounded-xl flex items-center justify-center mb-3 text-2xl"
          style={{ background: 'var(--surface3)' }}
        >
          {emoji}
        </div>

        {/* Category badge */}
        {menu.category && (
          <span
            className="inline-block self-start rounded-md px-2 py-0.5 mb-2"
            style={{
              background:    catStyle.bg,
              color:         catStyle.text,
              fontSize:      10,
              fontWeight:    700,
              letterSpacing: '0.6px',
              textTransform: 'uppercase',
            }}
          >
            {menu.category.category_name}
          </span>
        )}

        {/* Name */}
        <p className="font-semibold leading-snug flex-1 mb-4" style={{ fontSize: 13, color: 'var(--text)' }}>
          {menu.menu_name}
        </p>

        {/* Footer */}
        <div className="flex items-end justify-between">
          <div>
            <p style={{ fontSize: 16, fontWeight: 800, color: 'var(--sage-dark)', letterSpacing: '-0.3px' }}>
              {formatRp(menu.price)}
            </p>
            <p
              style={{
                fontSize:   11,
                marginTop:  2,
                color:      isOutOfStock ? 'var(--danger)' : isLowStock ? 'var(--warning)' : 'var(--text-4)',
                fontWeight: isLowStock || isOutOfStock ? 700 : 400,
              }}
            >
              {isOutOfStock ? '✕ Habis' : isLowStock ? `⚠ Sisa ${menu.daily_stock}` : `Stok: ${menu.daily_stock}`}
            </p>
          </div>

          {!isOutOfStock && (
            <button
              onClick={(e) => { e.stopPropagation(); onAdd(menu.menu_id); }}
              className="w-9 h-9 flex items-center justify-center rounded-xl transition-all duration-150"
              style={{ background: 'var(--sage-dark)', color: '#fff' }}
              onMouseEnter={(e) => {
                (e.currentTarget as HTMLElement).style.background = 'var(--caramel)';
                (e.currentTarget as HTMLElement).style.transform  = 'scale(1.08)';
              }}
              onMouseLeave={(e) => {
                (e.currentTarget as HTMLElement).style.background = 'var(--sage-dark)';
                (e.currentTarget as HTMLElement).style.transform  = 'scale(1)';
              }}
            >
              <Plus size={18} strokeWidth={2.5} />
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
