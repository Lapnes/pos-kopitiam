'use client';

import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api/axios';
import { useCartStore } from '@/store/useCartStore';
import { Menu, Order, OrderDetail } from '@/types';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { MenuCard } from '@/components/pos/MenuCard';
import { CheckoutDialog } from '@/components/pos/CheckoutDialog';
import { HistoryDialog } from '@/components/pos/HistoryDialog';
import { RefundByOrderNumberDialog } from '@/components/pos/RefundByOrderNumberDialog';
import {
  ShoppingCart, Search, Trash2, Minus, Plus,
  History, RotateCcw, Pencil, X, Loader2,
  Receipt, ChefHat, Coffee,
} from 'lucide-react';

const fmt = (n: number) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);

const STATION_FILTERS = ['Semua', 'bar', 'kitchen'];

function getMenuColor(name: string): string {
  const colors = ['bg-amber-100 text-amber-700', 'bg-emerald-100 text-emerald-700', 'bg-sky-100 text-sky-700', 'bg-violet-100 text-violet-700', 'bg-rose-100 text-rose-700'];
  let h = 0;
  for (let i = 0; i < name.length; i++) h = name.charCodeAt(i) + ((h << 5) - h);
  return colors[Math.abs(h) % colors.length];
}

/* ── Menu Edit Modal ── */
function MenuFormModal({ open, menu, onClose, onSaved }: {
  open: boolean; menu?: Menu | null; onClose: () => void; onSaved: () => void;
}) {
  const [form, setForm] = useState({
    name: menu?.name ?? '', price: menu?.price ?? 0,
    daily_stock: menu?.daily_stock ?? 50, station: menu?.station ?? 'bar',
    category_id: menu?.category_id ?? '', description: menu?.description ?? '',
  });
  const [loading, setLoading] = useState(false);

  function set(k: string, v: string | number) { setForm((p) => ({ ...p, [k]: v })); }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault(); setLoading(true);
    try {
      if (menu?.id) await api.put(`/api/v1/menus/${menu.id}`, { ...form, price: Number(form.price), daily_stock: Number(form.daily_stock) });
      else          await api.post('/api/v1/menus', { ...form, price: Number(form.price), daily_stock: Number(form.daily_stock), is_active: true });
      onSaved(); onClose();
    } finally { setLoading(false); }
  }

  if (!open) return null;
  const inputCls = "w-full px-4 py-3 rounded-xl border text-sm outline-none transition-all";
  const inputStyle = { borderColor: '#E8D9C4', backgroundColor: '#FDFAF6', color: '#1A0800' };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-md rounded-3xl border p-6" style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4', boxShadow: '0 40px 80px rgba(30,8,0,0.18)' }}>
        <div className="flex items-center justify-between mb-5">
          <h2 className="font-bold text-lg" style={{ color: '#1A0800' }}>{menu ? 'Edit Menu' : 'Tambah Menu'}</h2>
          <button onClick={onClose} className="w-8 h-8 flex items-center justify-center rounded-xl" style={{ color: '#A08060' }}><X size={16} /></button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          {[
            { label: 'Nama Menu', key: 'name', type: 'text', placeholder: 'Kopi Susu' },
            { label: 'Harga (Rp)', key: 'price', type: 'number', placeholder: '18000' },
            { label: 'Stok Harian', key: 'daily_stock', type: 'number', placeholder: '50' },
          ].map(({ label, key, type, placeholder }) => (
            <div key={key}>
              <label className="block text-xs font-semibold uppercase tracking-wider mb-1.5" style={{ color: '#8B6347' }}>{label}</label>
              <input type={type} className={inputCls} style={inputStyle} value={(form as Record<string, string | number>)[key]} onChange={(e) => set(key, e.target.value)}
                placeholder={placeholder} required
                onFocus={(e) => (e.currentTarget.style.borderColor = '#C07A24')}
                onBlur={(e)  => (e.currentTarget.style.borderColor = '#E8D9C4')} />
            </div>
          ))}
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider mb-1.5" style={{ color: '#8B6347' }}>Stasiun</label>
            <select className={`${inputCls} appearance-none cursor-pointer`} style={inputStyle} value={form.station} onChange={(e) => set('station', e.target.value)}>
              <option value="bar">Bar</option>
              <option value="kitchen">Dapur</option>
            </select>
          </div>
          <div className="flex gap-3 pt-1">
            <button type="button" onClick={onClose} className="flex-1 py-3 rounded-xl text-sm font-semibold" style={{ backgroundColor: '#F5EDD6', color: '#6B4C2A' }}>Batal</button>
            <button type="submit" disabled={loading} className="flex-1 py-3 rounded-xl text-sm font-bold text-white flex items-center justify-center gap-2" style={{ background: 'linear-gradient(135deg,#C07A24,#8B3A10)' }}>
              {loading && <Loader2 size={14} className="animate-spin" />}
              {menu ? 'Simpan' : 'Tambahkan'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

/* ── Main Kasir Page ── */
export default function KasirPage() {
  const qc = useQueryClient();
  const { items, addItem, updateQuantity, removeItem, clearCart, getSubtotal, getTaxAmount, getServiceCharge, getGrandTotal } = useCartStore();
  const [search,    setSearch]    = useState('');
  const [station,   setStation]   = useState('Semua');
  const [editMode,  setEditMode]  = useState(false);
  const [menuModal, setMenuModal] = useState<{ open: boolean; menu?: Menu | null }>({ open: false });
  const [histOpen,  setHistOpen]  = useState(false);
  const [refundOpen, setRefundOpen] = useState(false);

  const { data: menus = [], isLoading } = useQuery<Menu[]>({
    queryKey: ['menus'],
    queryFn: async () => { const r = await api.get('/api/v1/menus'); return r.data?.data ?? r.data ?? []; },
  });

  const filtered = menus.filter((m) => {
    const matchSearch  = m.name.toLowerCase().includes(search.toLowerCase());
    const matchStation = station === 'Semua' || m.station === station;
    return matchSearch && matchStation;
  });

  const totalItems  = items.reduce((s, i) => s + i.quantity, 0);
  const subtotal    = getSubtotal();
  const tax         = getTaxAmount();
  const service     = getServiceCharge();
  const grandTotal  = getGrandTotal();

  function refetch() { qc.invalidateQueries({ queryKey: ['menus'] }); }

  return (
    <DashboardLayout title="Kasir POS" subtitle="Kelola pesanan & pembayaran">
      <MenuFormModal open={menuModal.open} menu={menuModal.menu} onClose={() => setMenuModal({ open: false })} onSaved={refetch} />
      <HistoryDialog open={histOpen} onOpenChange={setHistOpen} />
      <RefundByOrderNumberDialog open={refundOpen} onOpenChange={setRefundOpen} onSuccess={refetch} />

      <div className="flex gap-4 h-[calc(100vh-130px)]">
        {/* ── LEFT: Menu Grid ── */}
        <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
          {/* Toolbar */}
          <div className="flex flex-wrap items-center gap-2 mb-4">
            <div className="flex items-center gap-2 flex-1 min-w-[180px] px-4 py-2.5 rounded-2xl border" style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4' }}>
              <Search size={15} style={{ color: '#A08060' }} />
              <input className="flex-1 bg-transparent text-sm outline-none" style={{ color: '#1A0800' }} placeholder="Cari menu…" value={search} onChange={(e) => setSearch(e.target.value)} />
            </div>
            {/* Station filter */}
            <div className="flex gap-1 p-1 rounded-2xl" style={{ backgroundColor: '#E8D9C4' }}>
              {STATION_FILTERS.map((s) => (
                <button key={s} onClick={() => setStation(s)} className="px-3 py-1.5 rounded-xl text-xs font-bold transition-all capitalize"
                  style={{ backgroundColor: station === s ? '#FFFDF8' : 'transparent', color: station === s ? '#C07A24' : '#6B4C2A' }}>
                  {s === 'Semua' ? 'Semua' : s === 'bar' ? '☕ Bar' : '🍳 Dapur'}
                </button>
              ))}
            </div>
            {/* Action buttons */}
            <button onClick={() => setHistOpen(true)} className="flex items-center gap-2 px-4 py-2.5 rounded-2xl text-sm font-semibold border transition-colors" style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4', color: '#6B4C2A' }}>
              <History size={15} /> Riwayat
            </button>
            <button onClick={() => setRefundOpen(true)} className="flex items-center gap-2 px-4 py-2.5 rounded-2xl text-sm font-semibold border transition-colors" style={{ backgroundColor: '#FEF3E2', borderColor: '#E8D9C4', color: '#C07A24' }}>
              <RotateCcw size={15} /> Refund
            </button>
            <button onClick={() => setEditMode((v) => !v)} className="flex items-center gap-2 px-4 py-2.5 rounded-2xl text-sm font-semibold transition-all"
              style={{ backgroundColor: editMode ? '#1A0800' : '#FFFDF8', borderColor: '#E8D9C4', color: editMode ? '#FFFBF4' : '#6B4C2A', border: '1px solid #E8D9C4' }}>
              <Pencil size={15} /> {editMode ? 'Selesai Edit' : 'Edit Menu'}
            </button>
            {editMode && (
              <button onClick={() => setMenuModal({ open: true, menu: null })} className="flex items-center gap-2 px-4 py-2.5 rounded-2xl text-sm font-bold text-white"
                style={{ background: 'linear-gradient(135deg,#C07A24,#8B3A10)', boxShadow: '0 4px 12px rgba(160,82,10,0.3)' }}>
                + Tambah Menu
              </button>
            )}
          </div>

          {/* Grid */}
          <div className="flex-1 overflow-y-auto pr-1">
            {isLoading ? (
              <div className="flex items-center justify-center h-40 gap-3"><Loader2 size={24} className="animate-spin" style={{ color: '#C07A24' }} /></div>
            ) : filtered.length === 0 ? (
              <div className="flex flex-col items-center py-16 gap-3">
                <Coffee size={32} style={{ color: '#D4B896' }} />
                <p className="text-sm" style={{ color: '#A08060' }}>Menu tidak ditemukan</p>
              </div>
            ) : (
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
                {filtered.map((menu) => {
                  const cartItem = items.find((i) => i.id === menu.id);
                  return (
                    <MenuCard key={menu.id} menu={menu} cartQty={cartItem?.quantity ?? 0}
                      onAdd={() => addItem(menu)}
                      canEdit={editMode}
                      onEdit={() => setMenuModal({ open: true, menu })} />
                  );
                })}
              </div>
            )}
          </div>
        </div>

        {/* ── RIGHT: Cart ── */}
        <div className="w-80 xl:w-96 shrink-0 flex flex-col rounded-3xl border overflow-hidden"
          style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4', boxShadow: '0 8px 32px rgba(80,30,0,0.10)' }}>
          {/* Cart header */}
          <div className="px-5 py-4 border-b flex items-center gap-3" style={{ borderColor: '#F0E4D0', backgroundColor: '#FEF9F0' }}>
            <div className="w-9 h-9 rounded-xl flex items-center justify-center" style={{ backgroundColor: '#FEF3E2' }}>
              <ShoppingCart size={18} style={{ color: '#C07A24' }} />
            </div>
            <div>
              <h2 className="font-bold text-sm" style={{ color: '#1A0800' }}>Pesanan Saat Ini</h2>
              <p className="text-xs" style={{ color: '#A08060' }}>{totalItems} item</p>
            </div>
            {totalItems > 0 && (
              <button onClick={clearCart} className="ml-auto w-8 h-8 flex items-center justify-center rounded-xl transition-colors" style={{ color: '#DC2626' }}
                onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FEE2E2')}
                onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}>
                <Trash2 size={15} />
              </button>
            )}
          </div>

          {/* Items */}
          <div className="flex-1 overflow-y-auto p-4 space-y-2">
            {items.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full gap-4 py-12">
                <div className="w-14 h-14 rounded-2xl flex items-center justify-center" style={{ backgroundColor: '#FEF3E2' }}>
                  <Receipt size={24} style={{ color: '#D4B896' }} />
                </div>
                <p className="text-sm font-semibold" style={{ color: '#A08060' }}>Keranjang kosong</p>
                <p className="text-xs text-center" style={{ color: '#C4A07A' }}>Klik menu di sebelah kiri untuk menambahkan item</p>
              </div>
            ) : (
              items.map((item) => (
                <div key={item.cartItemId} className="flex items-center gap-3 p-3 rounded-2xl border group transition-all"
                  style={{ backgroundColor: '#FDFAF6', borderColor: '#F0E4D0' }}
                  onMouseEnter={(e) => (e.currentTarget.style.borderColor = '#E8D9C4')}
                  onMouseLeave={(e) => (e.currentTarget.style.borderColor = '#F0E4D0')}>
                  <div className={`w-9 h-9 rounded-xl flex items-center justify-center shrink-0 text-xs font-black ${getMenuColor(item.name)}`}>
                    {item.name.substring(0, 2).toUpperCase()}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold text-xs truncate" style={{ color: '#1A0800' }}>{item.name}</p>
                    <p className="text-xs font-bold" style={{ color: '#C07A24' }}>{fmt(item.price)}</p>
                  </div>
                  <div className="flex items-center gap-1 shrink-0">
                    <button onClick={() => updateQuantity(item.cartItemId, item.quantity - 1)} className="w-7 h-7 flex items-center justify-center rounded-lg transition-colors" style={{ backgroundColor: '#F5EDD6' }}>
                      <Minus size={12} style={{ color: '#6B4C2A' }} />
                    </button>
                    <span className="w-7 text-center text-sm font-bold" style={{ color: '#1A0800' }}>{item.quantity}</span>
                    <button onClick={() => updateQuantity(item.cartItemId, item.quantity + 1)} className="w-7 h-7 flex items-center justify-center rounded-lg transition-colors" style={{ backgroundColor: '#FEF3E2' }}>
                      <Plus size={12} style={{ color: '#C07A24' }} />
                    </button>
                    <button onClick={() => removeItem(item.cartItemId)} className="w-7 h-7 flex items-center justify-center rounded-lg transition-all opacity-0 group-hover:opacity-100" style={{ color: '#DC2626' }}
                      onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FEE2E2')}
                      onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}>
                      <X size={12} />
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>

          {/* Summary + Checkout */}
          {items.length > 0 && (
            <div className="p-4 border-t space-y-3" style={{ borderColor: '#F0E4D0', backgroundColor: '#FEF9F0' }}>
              <div className="space-y-1.5">
                {[
                  { label: 'Subtotal', value: fmt(subtotal) },
                  { label: 'Pajak (11%)', value: fmt(tax) },
                  { label: 'Service (5%)', value: fmt(service) },
                ].map(({ label, value }) => (
                  <div key={label} className="flex justify-between text-xs" style={{ color: '#A08060' }}>
                    <span>{label}</span><span className="font-semibold">{value}</span>
                  </div>
                ))}
                <div className="flex justify-between items-center pt-2 border-t border-dashed" style={{ borderColor: '#E8D9C4' }}>
                  <span className="font-bold text-sm" style={{ color: '#1A0800' }}>Total</span>
                  <span className="text-xl font-black" style={{ color: '#C07A24' }}>{fmt(grandTotal)}</span>
                </div>
              </div>
              <CheckoutDialog>
                <button id="checkout-btn" className="w-full py-3.5 rounded-2xl font-bold text-sm text-white flex items-center justify-center gap-2 transition-all"
                  style={{ background: 'linear-gradient(135deg,#C07A24,#8B3A10)', boxShadow: '0 4px 12px rgba(160,82,10,0.35)' }}
                  onMouseEnter={(e) => (e.currentTarget.style.transform = 'translateY(-1px)')}
                  onMouseLeave={(e) => (e.currentTarget.style.transform = '')}>
                  <ShoppingCart size={16} /> Bayar — {fmt(grandTotal)}
                </button>
              </CheckoutDialog>
            </div>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}
