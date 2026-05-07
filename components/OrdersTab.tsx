'use client';

import { useEffect, useState, useCallback } from 'react';
import { RefreshCw, Trash2, TrendingUp, ShoppingBag, BarChart2, Receipt } from 'lucide-react';
import { OrderReportRow } from '@/types';
import { apiFetch, formatRp } from '@/lib/api';
import Modal from './Modal';
import Toast from './Toast';

export default function OrdersTab() {
  const [orders,  setOrders]  = useState<OrderReportRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [toast,   setToast]   = useState<{ msg: string; err?: boolean } | null>(null);

  const showToast = (msg: string, err = false) => {
    setToast({ msg, err });
    setTimeout(() => setToast(null), 3000);
  };

  const load = useCallback(async () => {
    setLoading(true);
    const res = await apiFetch<{ data: OrderReportRow[] }>('/orders/report');
    setOrders(res?.data || []);
    setLoading(false);
  }, []);

  useEffect(() => { load(); }, [load]);

  const uniqueOrders = [...new Set(orders.map((o) => o.order_id))];
  const totalSales   = uniqueOrders.reduce((s, oid) => {
    const row = orders.find((o) => o.order_id === oid);
    return s + (row?.total_price || 0);
  }, 0);
  const avg = uniqueOrders.length ? Math.round(totalSales / uniqueOrders.length) : 0;

  const handleDelete = async () => {
    if (!deleteId) return;
    const res = await apiFetch<{ status: string }>(`/orders/${deleteId}`, { method: 'DELETE' });
    setDeleteId(null);
    if (res?.status === 'success') { showToast('Transaksi berhasil dihapus'); load(); }
    else showToast('Gagal menghapus transaksi', true);
  };

  const stats = [
    { label: 'Total Pendapatan', value: formatRp(totalSales), icon: TrendingUp, color: 'var(--caramel)' },
    { label: 'Total Transaksi',  value: uniqueOrders.length.toString(), icon: ShoppingBag, color: 'var(--sage-dark)' },
    { label: 'Rata-rata / Order', value: formatRp(avg),     icon: BarChart2,  color: 'var(--brown-light)' },
  ];

  let lastOrderId = '';

  return (
    <div className="flex-1 overflow-y-auto" style={{ background: 'var(--bg)' }}>
      <div className="max-w-5xl mx-auto px-6 py-6">

        {/* Page header */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 style={{ fontFamily: "'DM Serif Display', serif", fontSize: 26, color: 'var(--text)', lineHeight: 1.1 }}>
              Riwayat Transaksi
            </h1>
            <p style={{ fontSize: 13, color: 'var(--text-3)', marginTop: 4 }}>
              Semua order hari ini
            </p>
          </div>
          <button
            onClick={load}
            className="flex items-center gap-2 px-4 py-2.5 rounded-xl border text-sm font-semibold transition-all"
            style={{ background: 'var(--surface)', borderColor: 'var(--border)', color: 'var(--text-3)' }}
            onMouseEnter={(e) => {
              (e.currentTarget as HTMLElement).style.borderColor = 'var(--caramel)';
              (e.currentTarget as HTMLElement).style.color       = 'var(--sage-dark)';
              (e.currentTarget as HTMLElement).style.background  = 'var(--caramel-pale)';
            }}
            onMouseLeave={(e) => {
              (e.currentTarget as HTMLElement).style.borderColor = 'var(--border)';
              (e.currentTarget as HTMLElement).style.color       = 'var(--text-3)';
              (e.currentTarget as HTMLElement).style.background  = 'var(--surface)';
            }}
          >
            <RefreshCw size={14} /> Refresh
          </button>
        </div>

        {/* Stat cards */}
        <div className="grid grid-cols-3 gap-4 mb-6">
          {stats.map(({ label, value, icon: Icon, color }) => (
            <div key={label}
              className="rounded-2xl border p-5 transition-all"
              style={{ background: 'var(--surface)', borderColor: 'var(--border)', boxShadow: '0 1px 4px rgba(90,60,20,0.05)' }}
              onMouseEnter={(e) => {
                (e.currentTarget as HTMLElement).style.borderColor = 'var(--border-med)';
                (e.currentTarget as HTMLElement).style.transform   = 'translateY(-2px)';
                (e.currentTarget as HTMLElement).style.boxShadow   = '0 6px 20px rgba(90,60,20,0.1)';
              }}
              onMouseLeave={(e) => {
                (e.currentTarget as HTMLElement).style.borderColor = 'var(--border)';
                (e.currentTarget as HTMLElement).style.transform   = 'translateY(0)';
                (e.currentTarget as HTMLElement).style.boxShadow   = '0 1px 4px rgba(90,60,20,0.05)';
              }}
            >
              <div className="flex items-center gap-2 mb-3">
                <div className="w-7 h-7 rounded-lg flex items-center justify-center"
                  style={{ background: 'var(--surface3)' }}>
                  <Icon size={14} style={{ color }} />
                </div>
                <p className="uppercase tracking-wider"
                  style={{ fontSize: 10, color: 'var(--text-4)', fontWeight: 700 }}>
                  {label}
                </p>
              </div>
              <p style={{ fontFamily: "'DM Serif Display', serif", fontSize: 22, color, letterSpacing: '-0.3px' }}>
                {loading ? <span className="skeleton inline-block w-24 h-6 rounded" /> : value}
              </p>
            </div>
          ))}
        </div>

        {/* Table */}
        <div className="rounded-2xl border overflow-hidden"
          style={{ background: 'var(--surface)', borderColor: 'var(--border)', boxShadow: '0 1px 4px rgba(90,60,20,0.05)' }}>

          {/* Table header */}
          <div className="px-5 py-3.5 flex items-center gap-2"
            style={{ borderBottom: '1px solid var(--border)', background: 'var(--surface2)' }}>
            <Receipt size={14} style={{ color: 'var(--text-3)' }} />
            <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-2)' }}>Daftar Order</span>
            {!loading && (
              <span className="ml-auto text-xs px-2 py-0.5 rounded-full font-semibold"
                style={{ background: 'var(--caramel-pale)', color: 'var(--sage-dark)' }}>
                {uniqueOrders.length} transaksi
              </span>
            )}
          </div>

          <div className="overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr style={{ background: 'var(--surface2)' }}>
                  {['Order ID', 'Kasir', 'Menu', 'Qty', 'Harga Satuan', 'Subtotal', 'Total', ''].map((h) => (
                    <th key={h} className="text-left px-4 py-3 border-b uppercase tracking-wider"
                      style={{ fontSize: 10, fontWeight: 700, color: 'var(--text-4)', borderColor: 'var(--border)' }}>
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr>
                    <td colSpan={8} className="text-center py-16">
                      <div className="flex flex-col items-center gap-3">
                        <div className="w-8 h-8 rounded-full border-2 border-t-transparent"
                          style={{ borderColor: 'var(--border)', borderTopColor: 'var(--caramel)', animation: 'spin 0.8s linear infinite' }} />
                        <span style={{ fontSize: 13, color: 'var(--text-4)' }}>Memuat data…</span>
                      </div>
                    </td>
                  </tr>
                ) : orders.length === 0 ? (
                  <tr>
                    <td colSpan={8} className="text-center py-16">
                      <div className="flex flex-col items-center gap-3">
                        <div className="text-4xl">🧾</div>
                        <span style={{ fontSize: 14, color: 'var(--text-3)' }}>Belum ada transaksi</span>
                      </div>
                    </td>
                  </tr>
                ) : (
                  orders.map((row, i) => {
                    const isNew = row.order_id !== lastOrderId;
                    lastOrderId = row.order_id;
                    return (
                      <tr key={i}
                        style={{ borderBottom: '1px solid var(--border)', background: isNew ? 'transparent' : 'rgba(245,240,232,0.3)' }}
                        onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--surface2)'; }}
                        onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = isNew ? 'transparent' : 'rgba(245,240,232,0.3)'; }}
                      >
                        {/* Order ID */}
                        <td className="px-4 py-3">
                          {isNew && (
                            <span className="font-mono text-xs px-2 py-1 rounded-md"
                              style={{ background: 'var(--surface3)', color: 'var(--text-3)', fontWeight: 600 }}>
                              #{row.order_id.substring(0, 8).toUpperCase()}
                            </span>
                          )}
                        </td>
                        {/* Kasir */}
                        <td className="px-4 py-3" style={{ fontSize: 13, color: 'var(--text-2)' }}>
                          {isNew && (
                            <div className="flex items-center gap-2">
                              <div className="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold"
                                style={{ background: 'var(--caramel-pale)', color: 'var(--sage-dark)' }}>
                                {(row.employee_name || 'W').charAt(0).toUpperCase()}
                              </div>
                              {row.employee_name || 'Walk-in'}
                            </div>
                          )}
                        </td>
                        {/* Menu */}
                        <td className="px-4 py-3 font-medium" style={{ fontSize: 13, color: 'var(--text)' }}>
                          {row.menu_name}
                        </td>
                        {/* Qty */}
                        <td className="px-4 py-3 text-center">
                          <span className="font-bold px-2 py-0.5 rounded-md"
                            style={{ fontSize: 12, background: 'var(--surface3)', color: 'var(--text-2)' }}>
                            {row.quantity}
                          </span>
                        </td>
                        {/* Harga */}
                        <td className="px-4 py-3" style={{ fontSize: 13, color: 'var(--text-3)' }}>
                          {formatRp(row.unit_price)}
                        </td>
                        {/* Subtotal */}
                        <td className="px-4 py-3" style={{ fontSize: 13, color: 'var(--text-2)', fontWeight: 600 }}>
                          {formatRp(row.subtotal)}
                        </td>
                        {/* Total */}
                        <td className="px-4 py-3">
                          {isNew && (
                            <span className="font-bold" style={{ fontSize: 14, color: 'var(--sage-dark)' }}>
                              {formatRp(row.total_price)}
                            </span>
                          )}
                        </td>
                        {/* Action */}
                        <td className="px-4 py-3">
                          {isNew && (
                            <button
                              onClick={() => setDeleteId(row.order_id)}
                              className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all"
                              style={{ background: 'var(--danger-bg)', color: 'var(--danger)', border: '1px solid rgba(192,56,56,0.2)' }}
                              onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'rgba(192,56,56,0.18)'; }}
                              onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--danger-bg)'; }}
                            >
                              <Trash2 size={12} /> Hapus
                            </button>
                          )}
                        </td>
                      </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Delete Modal */}
      <Modal isOpen={!!deleteId} onClose={() => setDeleteId(null)}>
        <div className="text-center p-1">
          <div className="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4"
            style={{ background: 'var(--danger-bg)' }}>
            <Trash2 size={28} style={{ color: 'var(--danger)' }} />
          </div>
          <h2 style={{ fontFamily: "'DM Serif Display', serif", fontSize: 22, marginBottom: 8, color: 'var(--text)' }}>
            Hapus Transaksi?
          </h2>
          <p style={{ fontSize: 14, color: 'var(--text-3)', marginBottom: 24 }}>
            Order{' '}
            <span style={{ color: 'var(--sage-dark)', fontWeight: 700 }}>
              #{deleteId?.substring(0, 8).toUpperCase()}
            </span>{' '}
            akan dihapus permanen.
          </p>
          <div className="flex gap-3">
            <button onClick={() => setDeleteId(null)}
              className="flex-1 py-3 rounded-xl border font-semibold text-sm transition-all"
              style={{ background: 'var(--surface2)', borderColor: 'var(--border)', color: 'var(--text-2)' }}
              onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.borderColor = 'var(--border-med)'; }}
              onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.borderColor = 'var(--border)'; }}>
              Batal
            </button>
            <button onClick={handleDelete}
              className="flex-1 py-3 rounded-xl font-bold text-sm transition-all"
              style={{ background: 'var(--danger)', color: '#fff', boxShadow: '0 4px 12px rgba(192,56,56,0.25)' }}
              onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.transform = 'translateY(-1px)'; }}
              onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.transform = 'translateY(0)'; }}>
              Hapus
            </button>
          </div>
        </div>
      </Modal>

      {toast && <Toast message={toast.msg} isError={toast.err} onClose={() => setToast(null)} />}
    </div>
  );
}
