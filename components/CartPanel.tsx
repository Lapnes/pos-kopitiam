'use client';

import { Minus, Plus, ShoppingBag, Trash2 } from 'lucide-react';
import { CartItem } from '@/types';
import { formatRp } from '@/lib/api';

interface CartPanelProps {
  items: CartItem[];
  onUpdateQty: (menuId: number, delta: number) => void;
  onClear: () => void;
  onCheckout: () => void;
  payment: string;
  onPaymentChange: (val: string) => void;
}

export default function CartPanel({
  items, onUpdateQty, onClear, onCheckout, payment, onPaymentChange,
}: CartPanelProps) {
  const total  = items.reduce((s, c) => s + c.price * c.qty, 0);
  const count  = items.reduce((s, c) => s + c.qty, 0);
  const paid   = parseInt(payment) || 0;
  const change = paid - total;

  // Quick-pay denominasi
  const quickAmounts = [10000, 20000, 50000, 100000];

  return (
    <div className="flex flex-col h-full" style={{ background: 'var(--surface)', borderLeft: '1px solid var(--border)' }}>

      {/* Header */}
      <div className="flex items-center justify-between px-5 pt-5 pb-4 flex-shrink-0"
        style={{ borderBottom: '1px solid var(--border)' }}>
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg flex items-center justify-center"
            style={{ background: 'var(--sage-dark)', color: '#fff' }}>
            <ShoppingBag size={15} />
          </div>
          <span style={{ fontFamily: "'DM Serif Display', serif", fontSize: 18, color: 'var(--text)' }}>
            Pesanan
          </span>
          {count > 0 && (
            <span className="text-xs font-bold px-2 py-0.5 rounded-full"
              style={{ background: 'var(--caramel)', color: '#fff' }}>
              {count}
            </span>
          )}
        </div>
        {items.length > 0 && (
          <button onClick={onClear}
            className="flex items-center gap-1.5 text-xs px-2.5 py-1.5 rounded-lg font-semibold transition-all"
            style={{ color: 'var(--text-4)', background: 'transparent' }}
            onMouseEnter={(e) => {
              (e.currentTarget as HTMLElement).style.color = 'var(--danger)';
              (e.currentTarget as HTMLElement).style.background = 'var(--danger-bg)';
            }}
            onMouseLeave={(e) => {
              (e.currentTarget as HTMLElement).style.color = 'var(--text-4)';
              (e.currentTarget as HTMLElement).style.background = 'transparent';
            }}
          >
            <Trash2 size={13} /> Hapus Semua
          </button>
        )}
      </div>

      {/* Cart items */}
      <div className="flex-1 overflow-y-auto px-4 py-3 min-h-0">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full gap-3 min-h-[200px]">
            <div className="w-16 h-16 rounded-2xl flex items-center justify-center"
              style={{ background: 'var(--surface3)' }}>
              <ShoppingBag size={28} style={{ color: 'var(--text-4)' }} />
            </div>
            <div className="text-center">
              <p className="font-semibold" style={{ fontSize: 14, color: 'var(--text-3)' }}>Keranjang kosong</p>
              <p style={{ fontSize: 12, color: 'var(--text-4)', marginTop: 4 }}>Pilih menu untuk memulai</p>
            </div>
          </div>
        ) : (
          <div className="space-y-2">
            {items.map((item, i) => (
              <div key={item.menu_id}
                className="anim-slide-r flex items-center gap-3 p-3 rounded-xl border transition-all"
                style={{
                  background: 'var(--surface2)',
                  borderColor: 'var(--border)',
                  animationDelay: `${i * 0.04}s`,
                }}
                onMouseEnter={(e) => {
                  (e.currentTarget as HTMLElement).style.borderColor = 'var(--border-med)';
                }}
                onMouseLeave={(e) => {
                  (e.currentTarget as HTMLElement).style.borderColor = 'var(--border)';
                }}
              >
                <div className="flex-1 min-w-0">
                  <p className="font-semibold truncate" style={{ fontSize: 13, color: 'var(--text)' }}>
                    {item.name}
                  </p>
                  <p style={{ fontSize: 11, color: 'var(--text-4)', marginTop: 1 }}>
                    {formatRp(item.price)}
                  </p>
                </div>

                {/* Qty control */}
                <div className="flex items-center gap-1.5 rounded-lg border px-1.5 py-1"
                  style={{ background: 'var(--surface)', borderColor: 'var(--border)' }}>
                  <button
                    onClick={() => onUpdateQty(item.menu_id, -1)}
                    className="w-6 h-6 flex items-center justify-center rounded-md transition-all font-bold"
                    style={{ background: 'var(--surface3)', color: 'var(--text-3)', fontSize: 14 }}
                    onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--danger-bg)'; (e.currentTarget as HTMLElement).style.color = 'var(--danger)'; }}
                    onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--surface3)'; (e.currentTarget as HTMLElement).style.color = 'var(--text-3)'; }}
                  >
                    <Minus size={11} strokeWidth={3} />
                  </button>
                  <span className="font-bold text-center" style={{ minWidth: 18, fontSize: 13, color: 'var(--text)' }}>
                    {item.qty}
                  </span>
                  <button
                    onClick={() => onUpdateQty(item.menu_id, 1)}
                    className="w-6 h-6 flex items-center justify-center rounded-md transition-all"
                    style={{ background: 'var(--surface3)', color: 'var(--text-3)' }}
                    onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--caramel-pale)'; (e.currentTarget as HTMLElement).style.color = 'var(--sage-dark)'; }}
                    onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--surface3)'; (e.currentTarget as HTMLElement).style.color = 'var(--text-3)'; }}
                  >
                    <Plus size={11} strokeWidth={3} />
                  </button>
                </div>

                <span className="font-bold text-right flex-shrink-0"
                  style={{ fontSize: 13, color: 'var(--sage-dark)', minWidth: 70 }}>
                  {formatRp(item.price * item.qty)}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="flex-shrink-0 px-4 pb-5 pt-3 space-y-3"
        style={{ borderTop: '1px solid var(--border)' }}>

        {/* Subtotal row */}
        <div className="flex justify-between items-center">
          <span style={{ fontSize: 12, color: 'var(--text-4)' }}>{count} item</span>
          <span style={{ fontSize: 13, fontWeight: 700, color: 'var(--text-2)' }}>
            {count > 0 ? formatRp(total) : '—'}
          </span>
        </div>

        <div style={{ height: 1, background: 'var(--border)' }} />

        {/* Total */}
        <div className="flex justify-between items-baseline">
          <span style={{ fontFamily: "'DM Serif Display', serif", fontSize: 17, color: 'var(--text)' }}>Total</span>
          <span style={{ fontFamily: "'DM Serif Display', serif", fontSize: 24, color: 'var(--sage-dark)', letterSpacing: '-0.5px' }}>
            {formatRp(total)}
          </span>
        </div>

        {/* Quick pay buttons */}
        <div>
          <p className="mb-1.5 uppercase tracking-wider" style={{ fontSize: 10, color: 'var(--text-4)', fontWeight: 600 }}>
            Nominal Cepat
          </p>
          <div className="grid grid-cols-4 gap-1.5">
            {quickAmounts.map((amt) => (
              <button key={amt}
                onClick={() => onPaymentChange(String(amt))}
                className="py-1.5 rounded-lg text-xs font-bold border transition-all"
                style={{
                  background:   payment === String(amt) ? 'var(--sage-dark)' : 'var(--surface3)',
                  color:        payment === String(amt) ? '#fff' : 'var(--text-3)',
                  borderColor:  payment === String(amt) ? 'var(--sage-dark)' : 'var(--border)',
                  fontSize:     11,
                }}
                onMouseEnter={(e) => {
                  if (payment !== String(amt)) {
                    (e.currentTarget as HTMLElement).style.background = 'var(--caramel-pale)';
                    (e.currentTarget as HTMLElement).style.borderColor = 'var(--caramel)';
                    (e.currentTarget as HTMLElement).style.color = 'var(--sage-dark)';
                  }
                }}
                onMouseLeave={(e) => {
                  if (payment !== String(amt)) {
                    (e.currentTarget as HTMLElement).style.background = 'var(--surface3)';
                    (e.currentTarget as HTMLElement).style.borderColor = 'var(--border)';
                    (e.currentTarget as HTMLElement).style.color = 'var(--text-3)';
                  }
                }}
              >
                {amt >= 1000 ? `${amt / 1000}rb` : amt}
              </button>
            ))}
          </div>
        </div>

        {/* Payment input */}
        <div>
          <p className="mb-1.5 uppercase tracking-wider" style={{ fontSize: 10, color: 'var(--text-4)', fontWeight: 600 }}>
            Uang Bayar
          </p>
          <div className="relative">
            <span className="absolute left-3.5 top-1/2 -translate-y-1/2 font-bold"
              style={{ fontSize: 12, color: 'var(--text-3)' }}>Rp</span>
            <input
              type="number"
              value={payment}
              onChange={(e) => onPaymentChange(e.target.value)}
              placeholder="0"
              className="w-full pl-10 pr-4 py-2.5 rounded-xl border font-bold transition-all"
              style={{
                background:  'var(--surface2)',
                borderColor: 'var(--border)',
                color:       'var(--text)',
                fontSize:    15,
              }}
              onFocus={(e) => {
                (e.currentTarget as HTMLElement).style.borderColor = 'var(--caramel)';
                (e.currentTarget as HTMLElement).style.boxShadow   = '0 0 0 3px rgba(200,146,42,0.12)';
                (e.currentTarget as HTMLElement).style.background  = 'var(--surface)';
              }}
              onBlur={(e) => {
                (e.currentTarget as HTMLElement).style.borderColor = 'var(--border)';
                (e.currentTarget as HTMLElement).style.boxShadow   = 'none';
                (e.currentTarget as HTMLElement).style.background  = 'var(--surface2)';
              }}
            />
          </div>
        </div>

        {/* Change display */}
        {paid > 0 && (
          <div className="flex justify-between items-center rounded-xl px-4 py-2.5"
            style={{
              background:  change < 0 ? 'var(--danger-bg)'  : 'var(--success-bg)',
              border:      `1px solid ${change < 0 ? 'rgba(192,56,56,0.25)' : 'rgba(45,138,90,0.25)'}`,
            }}>
            <span style={{ fontSize: 12, fontWeight: 600, color: change < 0 ? 'var(--danger)' : 'var(--success)' }}>
              {change < 0 ? 'Kurang' : 'Kembalian'}
            </span>
            <span className="font-bold" style={{ fontSize: 16, color: change < 0 ? 'var(--danger)' : 'var(--success)' }}>
              {change < 0 ? '-' : ''}{formatRp(Math.abs(change))}
            </span>
          </div>
        )}

        {/* Checkout button */}
        <button
          onClick={onCheckout}
          disabled={items.length === 0}
          className="w-full py-3.5 rounded-xl font-bold transition-all"
          style={{
            background: items.length === 0
              ? 'var(--surface3)'
              : 'var(--sage-dark)',
            color:      items.length === 0 ? 'var(--text-4)' : '#fff',
            fontSize:   15,
            cursor:     items.length === 0 ? 'not-allowed' : 'pointer',
            boxShadow:  items.length > 0 ? '0 4px 16px rgba(61,92,61,0.25)' : 'none',
          }}
          onMouseEnter={(e) => {
            if (items.length > 0) {
              (e.currentTarget as HTMLElement).style.background  = 'var(--caramel)';
              (e.currentTarget as HTMLElement).style.boxShadow  = '0 6px 20px rgba(200,146,42,0.3)';
              (e.currentTarget as HTMLElement).style.transform  = 'translateY(-1px)';
            }
          }}
          onMouseLeave={(e) => {
            if (items.length > 0) {
              (e.currentTarget as HTMLElement).style.background = 'var(--sage-dark)';
              (e.currentTarget as HTMLElement).style.boxShadow = '0 4px 16px rgba(61,92,61,0.25)';
              (e.currentTarget as HTMLElement).style.transform  = 'translateY(0)';
            }
          }}
        >
          Proses Pembayaran →
        </button>
      </div>
    </div>
  );
}
