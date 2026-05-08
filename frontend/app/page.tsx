'use client';

import { useState, useEffect, useCallback } from 'react';
import { Search, Coffee, ChevronDown } from 'lucide-react';
import { Menu, Employee, CartItem, Category } from '@/types';
import { apiFetch, formatRp, getMenuEmoji } from '@/lib/api';
import MenuCard from '@/components/MenuCard';
import CartPanel from '@/components/CartPanel';
import OrdersTab from '@/components/OrdersTab';
import Modal from '@/components/Modal';
import Toast from '@/components/Toast';

type Tab = 'order' | 'orders';

const MOCK_MENUS: Menu[] = [
  { menu_id: 1, category_id: 1, menu_name: 'Kopi Hitam',   price: 12000, daily_stock: 20, category: { category_id: 1, category_name: 'Kopi' } },
  { menu_id: 2, category_id: 1, menu_name: 'Kopi Susu',    price: 18000, daily_stock: 15, category: { category_id: 1, category_name: 'Kopi' } },
  { menu_id: 3, category_id: 1, menu_name: 'Cappuccino',   price: 25000, daily_stock: 10, category: { category_id: 1, category_name: 'Kopi' } },
  { menu_id: 4, category_id: 2, menu_name: 'Teh Tarik',    price: 14000, daily_stock: 18, category: { category_id: 2, category_name: 'Non-Kopi' } },
  { menu_id: 5, category_id: 2, menu_name: 'Milo Panas',   price: 16000, daily_stock: 12, category: { category_id: 2, category_name: 'Non-Kopi' } },
  { menu_id: 6, category_id: 3, menu_name: 'Nasi Lemak',   price: 22000, daily_stock: 8,  category: { category_id: 3, category_name: 'Makanan' } },
  { menu_id: 7, category_id: 3, menu_name: 'Roti Bakar',   price: 15000, daily_stock: 5,  category: { category_id: 3, category_name: 'Makanan' } },
  { menu_id: 8, category_id: 3, menu_name: 'Mie Goreng',   price: 20000, daily_stock: 0,  category: { category_id: 3, category_name: 'Makanan' } },
];

export default function PosPage() {
  const [tab,              setTab]              = useState<Tab>('order');
  const [menus,            setMenus]            = useState<Menu[]>([]);
  const [categories,       setCategories]       = useState<Category[]>([]);
  const [employees,        setEmployees]        = useState<Employee[]>([]);
  const [loading,          setLoading]          = useState(true);
  const [search,           setSearch]           = useState('');
  const [activeCat,        setActiveCat]        = useState<number | null>(null);
  const [selectedEmployee, setSelectedEmployee] = useState<string>('');
  const [cart,             setCart]             = useState<CartItem[]>([]);
  const [payment,          setPayment]          = useState('');
  const [confirmOpen,      setConfirmOpen]      = useState(false);
  const [successOpen,      setSuccessOpen]      = useState(false);
  const [lastOrder,        setLastOrder]        = useState<{ id: string; items: CartItem[]; total: number; paid: number } | null>(null);
  const [toast,            setToast]            = useState<{ msg: string; err?: boolean } | null>(null);

  const showToast = (msg: string, err = false) => {
    setToast({ msg, err });
    setTimeout(() => setToast(null), 3000);
  };

  const loadData = useCallback(async () => {
    setLoading(true);
    const [menusRes, catsRes, empsRes] = await Promise.all([
      apiFetch<{ data: Menu[] }>('/menus'),
      apiFetch<{ data: Category[] }>('/categories'),
      apiFetch<{ data: Employee[] }>('/employees'),
    ]);
    setMenus(menusRes?.data || MOCK_MENUS);
    setCategories(catsRes?.data || []);
    setEmployees(empsRes?.data || []);
    setLoading(false);
  }, []);

  useEffect(() => { loadData(); }, [loadData]);

  const filtered = menus.filter((m) => {
    const matchCat    = activeCat === null || m.category_id === activeCat;
    const matchSearch = m.menu_name.toLowerCase().includes(search.toLowerCase());
    return matchCat && matchSearch;
  });

  const addToCart = (menuId: number) => {
    const menu = menus.find((m) => m.menu_id === menuId);
    if (!menu || menu.daily_stock <= 0) return;
    const existing = cart.find((c) => c.menu_id === menuId);
    if (existing && existing.qty >= menu.daily_stock) {
      showToast(`Stok ${menu.menu_name} tidak cukup`, true);
      return;
    }
    setCart((prev) => {
      const idx = prev.findIndex((c) => c.menu_id === menuId);
      if (idx >= 0) {
        const updated = [...prev];
        updated[idx] = { ...updated[idx], qty: updated[idx].qty + 1 };
        return updated;
      }
      return [...prev, { menu_id: menuId, name: menu.menu_name, price: menu.price, qty: 1 }];
    });
    showToast(`${menu.menu_name} ditambahkan ✓`);
  };

  const updateQty = (menuId: number, delta: number) => {
    setCart((prev) => prev.map((c) => c.menu_id === menuId ? { ...c, qty: c.qty + delta } : c).filter((c) => c.qty > 0));
  };

  const clearCart = () => { setCart([]); setPayment(''); };

  const total  = cart.reduce((s, c) => s + c.price * c.qty, 0);
  const paid   = parseInt(payment) || 0;
  const change = paid - total;

  const openConfirm = () => {
    if (!cart.length) return;
    if (paid > 0 && paid < total) { showToast('Uang bayar kurang!', true); return; }
    setConfirmOpen(true);
  };

  const processOrder = async () => {
    setConfirmOpen(false);
    const payload = {
      employee_id: selectedEmployee ? parseInt(selectedEmployee) : undefined,
      items: cart.map((c) => ({ menu_id: c.menu_id, quantity: c.qty })),
    };
    const res = await apiFetch<{ order_id: string }>('/orders', { method: 'POST', body: JSON.stringify(payload) });
    const orderId = res?.order_id || 'LOCAL-' + Date.now();
    setMenus((prev) => prev.map((m) => {
      const ci = cart.find((c) => c.menu_id === m.menu_id);
      return ci ? { ...m, daily_stock: m.daily_stock - ci.qty } : m;
    }));
    setLastOrder({ id: orderId, items: [...cart], total, paid });
    clearCart();
    setSuccessOpen(true);
    showToast('Transaksi berhasil! 🎉');
  };

  // Category counts for badges
  const catCounts = menus.reduce<Record<number, number>>((acc, m) => {
    acc[m.category_id] = (acc[m.category_id] || 0) + 1;
    return acc;
  }, {});

  const now = new Date();
  const greeting = now.getHours() < 11 ? 'Selamat Pagi' : now.getHours() < 15 ? 'Selamat Siang' : now.getHours() < 18 ? 'Selamat Sore' : 'Selamat Malam';

  return (
    <div className="flex flex-col h-screen" style={{ background: 'var(--bg)' }}>

      {/* ── Header ── */}
      <header
        className="flex items-center justify-between px-6 py-3 flex-shrink-0"
        style={{ background: 'var(--surface)', borderBottom: '1px solid var(--border)', boxShadow: '0 1px 0 var(--border)' }}
      >
        {/* Brand */}
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl flex items-center justify-center"
            style={{ background: 'var(--sage-dark)', boxShadow: '0 2px 8px rgba(61,92,61,0.25)' }}>
            <Coffee size={20} color="#fff" />
          </div>
          <div>
            <p style={{ fontFamily: "'DM Serif Display', serif", fontSize: 19, color: 'var(--sage-dark)', lineHeight: 1.1 }}>
              KopiTiam
            </p>
            <p style={{ fontSize: 10, color: 'var(--text-4)', letterSpacing: '2px', textTransform: 'uppercase', fontWeight: 600 }}>
              Point of Sale
            </p>
          </div>
        </div>

        {/* Tab switcher */}
        <div className="flex gap-1 rounded-xl p-1 border" style={{ background: 'var(--surface2)', borderColor: 'var(--border)' }}>
          {([['order', '☕  Order'], ['orders', '📋  Transaksi']] as [Tab, string][]).map(([key, label]) => (
            <button key={key} onClick={() => setTab(key)}
              className="px-5 py-2 rounded-lg text-sm font-semibold transition-all"
              style={tab === key
                ? { background: 'var(--sage-dark)', color: '#fff', boxShadow: '0 2px 8px rgba(61,92,61,0.2)' }
                : { color: 'var(--text-3)', background: 'transparent' }}>
              {label}
            </button>
          ))}
        </div>

        {/* Kasir selector */}
        <div className="flex items-center gap-3">
          <div className="text-right hidden sm:block">
            <p style={{ fontSize: 11, color: 'var(--text-4)', fontWeight: 600 }}>{greeting}</p>
            <p style={{ fontSize: 12, color: 'var(--text-2)', fontWeight: 700 }}>
              {selectedEmployee ? employees.find(e => String(e.employee_id) === selectedEmployee)?.employee_name : 'Walk-in'}
            </p>
          </div>
          <div className="relative">
            <select value={selectedEmployee} onChange={(e) => setSelectedEmployee(e.target.value)}
              className="rounded-xl border py-2.5 pl-3 text-sm font-medium transition-all"
              style={{
                background: 'var(--surface2)', borderColor: 'var(--border)',
                color: 'var(--text)', minWidth: 145,
              }}
              onFocus={(e) => { (e.currentTarget as HTMLElement).style.borderColor = 'var(--caramel)'; }}
              onBlur={(e)  => { (e.currentTarget as HTMLElement).style.borderColor = 'var(--border)'; }}
            >
              <option value="">— Walk-in —</option>
              {employees.map((emp) => (
                <option key={emp.employee_id} value={emp.employee_id}>{emp.employee_name}</option>
              ))}
            </select>
          </div>
        </div>
      </header>

      {/* ── Body ── */}
      <div className="flex flex-1 overflow-hidden">
        {tab === 'order' ? (
          <>
            {/* Menu area */}
            <div className="flex-1 flex flex-col overflow-hidden">

              {/* Search + filter toolbar */}
              <div className="px-5 pt-4 pb-3 flex-shrink-0"
                style={{ borderBottom: '1px solid var(--border)', background: 'var(--surface)' }}>
                {/* Search bar */}
                <div className="flex items-center gap-3 rounded-xl border px-3.5 py-2.5 mb-3 transition-all"
                  style={{ background: 'var(--surface2)', borderColor: 'var(--border)' }}
                  onFocusCapture={(e) => { e.currentTarget.style.borderColor = 'var(--caramel)'; e.currentTarget.style.background = 'var(--surface)'; }}
                  onBlurCapture={(e) => { e.currentTarget.style.borderColor = 'var(--border)'; e.currentTarget.style.background = 'var(--surface2)'; }}>
                  <Search size={15} style={{ color: 'var(--text-4)', flexShrink: 0 }} />
                  <input
                    type="text" placeholder="Cari menu…" value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    className="flex-1 bg-transparent text-sm"
                    style={{ color: 'var(--text)' }}
                  />
                  {search && (
                    <button onClick={() => setSearch('')}
                      className="text-xs px-2 py-0.5 rounded-md font-semibold"
                      style={{ background: 'var(--surface3)', color: 'var(--text-4)' }}>
                      ✕
                    </button>
                  )}
                </div>

                {/* Category filters */}
                <div className="flex gap-2 flex-wrap">
                  <button onClick={() => setActiveCat(null)}
                    className="px-3.5 py-1.5 rounded-full text-xs font-bold border transition-all"
                    style={activeCat === null
                      ? { background: 'var(--sage-dark)', color: '#fff', borderColor: 'var(--sage-dark)' }
                      : { background: 'var(--surface2)', color: 'var(--text-3)', borderColor: 'var(--border)' }}>
                    Semua ({menus.length})
                  </button>
                  {categories.map((cat) => (
                    <button key={cat.category_id} onClick={() => setActiveCat(cat.category_id)}
                      className="px-3.5 py-1.5 rounded-full text-xs font-bold border transition-all"
                      style={activeCat === cat.category_id
                        ? { background: 'var(--sage-dark)', color: '#fff', borderColor: 'var(--sage-dark)' }
                        : { background: 'var(--surface2)', color: 'var(--text-3)', borderColor: 'var(--border)' }}>
                      {cat.category_name} {catCounts[cat.category_id] ? `(${catCounts[cat.category_id]})` : ''}
                    </button>
                  ))}
                </div>
              </div>

              {/* Menu grid */}
              <div className="flex-1 overflow-y-auto p-5" style={{ background: 'var(--bg)' }}>
                {/* Section label */}
                <div className="flex items-center justify-between mb-4">
                  <p style={{ fontSize: 12, color: 'var(--text-4)', fontWeight: 600 }}>
                    {filtered.length} menu {search || activeCat ? 'ditemukan' : 'tersedia'}
                  </p>
                </div>

                {loading ? (
                  <div className="grid gap-3" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(170px, 1fr))' }}>
                    {Array.from({ length: 8 }).map((_, i) => (
                      <div key={i} className="skeleton rounded-2xl" style={{ height: 175 }} />
                    ))}
                  </div>
                ) : filtered.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-48 gap-3">
                    <div className="text-4xl">🔍</div>
                    <p style={{ fontSize: 14, color: 'var(--text-3)', fontWeight: 600 }}>Tidak ada menu</p>
                    <p style={{ fontSize: 12, color: 'var(--text-4)' }}>Coba kata kunci lain</p>
                    {(search || activeCat) && (
                      <button onClick={() => { setSearch(''); setActiveCat(null); }}
                        className="mt-1 px-4 py-2 rounded-xl text-sm font-semibold border transition-all"
                        style={{ background: 'var(--surface)', borderColor: 'var(--border)', color: 'var(--text-2)' }}>
                        Reset filter
                      </button>
                    )}
                  </div>
                ) : (
                  <div className="grid gap-3" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(170px, 1fr))' }}>
                    {filtered.map((menu, i) => (
                      <div key={menu.menu_id} className="anim-fade-up" style={{ animationDelay: `${i * 0.03}s` }}>
                        <MenuCard menu={menu} onAdd={addToCart} />
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* Cart */}
            <div className="w-[360px] flex-shrink-0">
              <CartPanel
                items={cart} onUpdateQty={updateQty} onClear={clearCart}
                onCheckout={openConfirm} payment={payment} onPaymentChange={setPayment}
              />
            </div>
          </>
        ) : (
          <OrdersTab />
        )}
      </div>

      {/* ── Confirm Modal ── */}
      <Modal isOpen={confirmOpen} onClose={() => setConfirmOpen(false)}>
        <div>
          <div className="flex items-center gap-3 mb-5">
            <div className="w-11 h-11 rounded-xl flex items-center justify-center"
              style={{ background: 'var(--caramel-pale)' }}>
              <span className="text-2xl">🧾</span>
            </div>
            <div>
              <h2 style={{ fontFamily: "'DM Serif Display', serif", fontSize: 20, color: 'var(--text)' }}>
                Konfirmasi Pesanan
              </h2>
              <p style={{ fontSize: 12, color: 'var(--text-4)', marginTop: 2 }}>Pastikan pesanan sudah benar</p>
            </div>
          </div>

          <div className="rounded-xl border mb-4 overflow-hidden"
            style={{ background: 'var(--surface2)', borderColor: 'var(--border)' }}>
            <div className="divide-y" style={{ borderColor: 'var(--border)' }}>
              {cart.map((item) => (
                <div key={item.menu_id} className="flex justify-between items-center px-4 py-2.5">
                  <span style={{ fontSize: 13, color: 'var(--text)' }}>
                    {getMenuEmoji(item.name)} {item.name} ×{item.qty}
                  </span>
                  <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-2)' }}>
                    {formatRp(item.price * item.qty)}
                  </span>
                </div>
              ))}
            </div>
            <div className="px-4 py-3 flex justify-between items-center"
              style={{ background: 'var(--surface3)', borderTop: '1px solid var(--border)' }}>
              <span style={{ fontFamily: "'DM Serif Display', serif", fontSize: 16, color: 'var(--text)' }}>Total</span>
              <span style={{ fontFamily: "'DM Serif Display', serif", fontSize: 20, color: 'var(--sage-dark)' }}>
                {formatRp(total)}
              </span>
            </div>
            {paid > 0 && (
              <div className="px-4 py-2.5 flex flex-col gap-1" style={{ borderTop: '1px solid var(--border)' }}>
                <div className="flex justify-between">
                  <span style={{ fontSize: 12, color: 'var(--text-4)' }}>Bayar</span>
                  <span style={{ fontSize: 12, color: 'var(--text-3)' }}>{formatRp(paid)}</span>
                </div>
                <div className="flex justify-between">
                  <span style={{ fontSize: 12, fontWeight: 700, color: change >= 0 ? 'var(--success)' : 'var(--danger)' }}>
                    Kembalian
                  </span>
                  <span style={{ fontSize: 13, fontWeight: 800, color: change >= 0 ? 'var(--success)' : 'var(--danger)' }}>
                    {formatRp(change)}
                  </span>
                </div>
              </div>
            )}
          </div>

          <div className="flex gap-3">
            <button onClick={() => setConfirmOpen(false)}
              className="flex-1 py-3 rounded-xl border font-semibold text-sm transition-all"
              style={{ background: 'var(--surface2)', borderColor: 'var(--border)', color: 'var(--text-2)' }}>
              Batal
            </button>
            <button onClick={processOrder}
              className="flex-1 py-3 rounded-xl font-bold text-sm transition-all"
              style={{ background: 'var(--sage-dark)', color: '#fff', boxShadow: '0 4px 14px rgba(61,92,61,0.25)' }}
              onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--caramel)'; }}
              onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--sage-dark)'; }}>
              Proses Pembayaran ✓
            </button>
          </div>
        </div>
      </Modal>

      {/* ── Success Modal ── */}
      <Modal isOpen={successOpen} onClose={() => setSuccessOpen(false)}>
        <div className="text-center">
          <div className="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4"
            style={{ background: 'var(--success-bg)' }}>
            <span className="text-3xl">🎉</span>
          </div>
          <h2 style={{ fontFamily: "'DM Serif Display', serif", fontSize: 22, marginBottom: 6, color: 'var(--text)' }}>
            Transaksi Berhasil!
          </h2>
          {lastOrder && (
            <p className="mb-4 inline-block font-mono text-xs px-3 py-1.5 rounded-lg"
              style={{ background: 'var(--surface3)', color: 'var(--text-3)' }}>
              Order #{lastOrder.id.substring(0, 8).toUpperCase()}
            </p>
          )}

          {lastOrder && (
            <div className="rounded-xl border mb-5 overflow-hidden text-left"
              style={{ background: 'var(--surface2)', borderColor: 'var(--border)' }}>
              <div className="divide-y" style={{ borderColor: 'var(--border)' }}>
                {lastOrder.items.map((item) => (
                  <div key={item.menu_id} className="flex justify-between px-4 py-2.5">
                    <span style={{ fontSize: 13, color: 'var(--text)' }}>{item.name} ×{item.qty}</span>
                    <span style={{ fontSize: 13, color: 'var(--text-2)' }}>{formatRp(item.price * item.qty)}</span>
                  </div>
                ))}
              </div>
              <div className="px-4 py-3 flex justify-between"
                style={{ background: 'var(--surface3)', borderTop: '1px solid var(--border)' }}>
                <span style={{ fontFamily: "'DM Serif Display', serif", fontSize: 16, color: 'var(--text)' }}>Total</span>
                <span style={{ fontFamily: "'DM Serif Display', serif", fontSize: 18, color: 'var(--sage-dark)' }}>
                  {formatRp(lastOrder.total)}
                </span>
              </div>
              {lastOrder.paid > 0 && (
                <div className="px-4 py-2.5 flex justify-between" style={{ borderTop: '1px solid var(--border)' }}>
                  <span style={{ fontSize: 12, color: 'var(--success)', fontWeight: 700 }}>Kembalian</span>
                  <span style={{ fontSize: 14, color: 'var(--success)', fontWeight: 800 }}>
                    {formatRp(lastOrder.paid - lastOrder.total)}
                  </span>
                </div>
              )}
            </div>
          )}

          <button onClick={() => setSuccessOpen(false)}
            className="w-full py-3.5 rounded-xl font-bold transition-all"
            style={{ background: 'var(--sage-dark)', color: '#fff', fontSize: 15, boxShadow: '0 4px 14px rgba(61,92,61,0.25)' }}
            onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--caramel)'; }}
            onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--sage-dark)'; }}>
            Order Berikutnya →
          </button>
        </div>
      </Modal>

      {toast && <Toast message={toast.msg} isError={toast.err} onClose={() => setToast(null)} />}
    </div>
  );
}
