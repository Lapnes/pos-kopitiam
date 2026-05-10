'use client';

import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api/axios';
import { Order } from '@/types';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import {
  ChefHat, Check, Clock, Loader2, AlertCircle,
  RefreshCw, Timer, Utensils,
} from 'lucide-react';

/* ── helpers ── */
function elapsed(createdAt: string): string {
  const diff = Math.floor((Date.now() - new Date(createdAt).getTime()) / 1000);
  if (diff < 60)  return `${diff}d`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m`;
  return `${Math.floor(diff / 3600)}j`;
}

function urgencyColor(createdAt: string): string {
  const mins = (Date.now() - new Date(createdAt).getTime()) / 60000;
  if (mins > 15) return '#DC2626'; // very late
  if (mins > 8)  return '#D97706'; // late
  return '#047857';                // ok
}

const ORDER_TYPE_LABEL: Record<string, string> = {
  dine_in:  '🪑 Makan Di Sini',
  takeaway: '🛍 Bawa Pulang',
  delivery: '🛵 Delivery',
  online:   '📱 Online',
};

/* ── ticket card ── */
function OrderTicket({
  order,
  onComplete,
  completing,
}: {
  order: Order;
  onComplete: (id: string) => void;
  completing: boolean;
}) {
  const [fadeOut, setFadeOut] = useState(false);
  const color = urgencyColor(order.created_at);
  const mins = Math.floor((Date.now() - new Date(order.created_at).getTime()) / 60000);

  function handleComplete() {
    setFadeOut(true);
    setTimeout(() => onComplete(order.id), 350);
  }

  return (
    <div
      className="rounded-3xl border flex flex-col overflow-hidden transition-all duration-300"
      style={{
        backgroundColor: '#FFFDF8',
        borderColor: '#E8D9C4',
        boxShadow: '0 4px 20px rgba(80,30,0,0.08)',
        borderTop: `4px solid ${color}`,
        opacity: fadeOut ? 0 : 1,
        transform: fadeOut ? 'scale(0.95) translateY(-8px)' : 'scale(1)',
      }}
    >
      {/* Header */}
      <div className="px-4 py-3.5 flex items-center justify-between" style={{ borderBottom: '1px solid #F0E4D0' }}>
        <div>
          <p className="font-black text-base" style={{ color: '#1A0800' }}>
            #{order.order_number}
          </p>
          <p className="text-xs mt-0.5" style={{ color: '#A08060' }}>
            {ORDER_TYPE_LABEL[order.order_type] ?? order.order_type}
          </p>
        </div>
        <div className="flex flex-col items-end gap-1">
          <div
            className="flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold"
            style={{ backgroundColor: `${color}15`, color }}
          >
            <Timer size={11} strokeWidth={2.5} />
            {elapsed(order.created_at)}
          </div>
          {order.customer_name && (
            <p className="text-[10px] font-semibold" style={{ color: '#A08060' }}>
              {order.customer_name}
            </p>
          )}
        </div>
      </div>

      {/* Items */}
      <div className="flex-1 p-4 space-y-2.5">
        {order.order_details
          .filter((d) => !d.is_voided)
          .map((item) => (
            <div key={item.id} className="flex items-start gap-3">
              <div
                className="w-7 h-7 rounded-xl flex items-center justify-center text-sm font-black shrink-0 mt-0.5"
                style={{ backgroundColor: '#FEF3E2', color: '#C07A24' }}
              >
                {item.quantity}
              </div>
              <div className="flex-1 min-w-0">
                <p className="font-semibold text-sm" style={{ color: '#1A0800' }}>{item.menu_name}</p>
                {item.notes && (
                  <p className="text-xs mt-0.5 italic" style={{ color: '#B45309' }}>
                    📝 {item.notes}
                  </p>
                )}
                <p className="text-[10px] mt-0.5 font-semibold uppercase tracking-wider" style={{ color: '#A08060' }}>
                  {item.station}
                </p>
              </div>
            </div>
          ))}
      </div>

      {/* Urgency bar */}
      {mins > 8 && (
        <div
          className="px-4 py-2 flex items-center gap-2 text-xs font-semibold"
          style={{ backgroundColor: `${color}10`, color }}
        >
          <Clock size={12} />
          {mins > 15 ? 'Sangat terlambat! Segera selesaikan.' : 'Waktu melebihi 8 menit'}
        </div>
      )}

      {/* Complete button */}
      <div className="p-4 pt-0">
        <button
          onClick={handleComplete}
          disabled={completing}
          className="w-full py-3 rounded-2xl text-sm font-bold text-white flex items-center justify-center gap-2 transition-all active:scale-[0.97] disabled:opacity-60"
          style={{
            background: 'linear-gradient(135deg, #047857, #065F46)',
            boxShadow: '0 4px 12px rgba(4,120,87,0.3)',
          }}
          onMouseEnter={(e) => {
            if (!completing) e.currentTarget.style.transform = 'translateY(-1px)';
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.transform = '';
          }}
        >
          {completing
            ? <Loader2 size={15} className="animate-spin" />
            : <Check size={15} strokeWidth={2.5} />}
          Tandai Selesai
        </button>
      </div>
    </div>
  );
}

/* ── page ── */
export default function KitchenPage() {
  const qc = useQueryClient();
  const [completingId, setCompletingId] = useState<string | null>(null);

  const { data: orders = [], isLoading, error, refetch } = useQuery<Order[]>({
    queryKey: ['kitchen-orders'],
    queryFn: async () => {
      const r = await api.get('/api/v1/orders', { params: { status: 'confirmed' } });
      return r.data?.data ?? r.data ?? [];
    },
    refetchInterval: 15000, // auto-refresh every 15s
  });

  const confirmMutation = useMutation({
    mutationFn: (id: string) => api.put(`/api/v1/orders/${id}/confirm`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['kitchen-orders'] }),
  });

  async function handleComplete(id: string) {
    setCompletingId(id);
    try {
      await confirmMutation.mutateAsync(id);
    } finally {
      setCompletingId(null);
    }
  }

  // Separate by station
  const barOrders     = orders.filter((o) => o.order_details.some((d) => d.station === 'bar' && !d.is_voided));
  const kitchenOrders = orders.filter((o) => o.order_details.some((d) => d.station === 'kitchen' && !d.is_voided));

  return (
    <DashboardLayout
      title="Dapur & Bar"
      subtitle={`${orders.length} pesanan aktif • Refresh otomatis setiap 15 detik`}
    >
      {/* Header actions */}
      <div className="flex items-center justify-between mb-5">
        <div className="flex gap-3">
          <div
            className="flex items-center gap-2 px-4 py-2 rounded-2xl border"
            style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4' }}
          >
            <div className="w-2.5 h-2.5 rounded-full bg-green-500 animate-pulse" />
            <span className="text-xs font-semibold" style={{ color: '#6B4C2A' }}>Live</span>
          </div>
          <div
            className="flex items-center gap-2 px-4 py-2 rounded-2xl border"
            style={{ backgroundColor: '#FEF3E2', borderColor: '#E8D9C4', color: '#C07A24' }}
          >
            <ChefHat size={14} />
            <span className="text-xs font-bold">{orders.length} Order</span>
          </div>
        </div>
        <button
          onClick={() => refetch()}
          className="flex items-center gap-2 px-4 py-2 rounded-2xl text-xs font-semibold transition-colors"
          style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4', color: '#6B4C2A', border: '1px solid #E8D9C4' }}
          onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FEF3E2')}
          onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = '#FFFDF8')}
        >
          <RefreshCw size={13} />
          Refresh
        </button>
      </div>

      {isLoading ? (
        <div className="flex flex-col items-center justify-center py-24 gap-4">
          <Loader2 size={36} className="animate-spin" style={{ color: '#C07A24' }} />
          <p className="text-sm" style={{ color: '#A08060' }}>Memuat pesanan…</p>
        </div>
      ) : error ? (
        <div className="flex flex-col items-center py-16 gap-3">
          <AlertCircle size={32} style={{ color: '#DC2626' }} />
          <p className="font-semibold text-sm" style={{ color: '#9B2226' }}>Gagal memuat pesanan</p>
        </div>
      ) : orders.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-24 gap-4">
          <div
            className="w-20 h-20 rounded-3xl flex items-center justify-center"
            style={{ backgroundColor: '#FEF3E2' }}
          >
            <Utensils size={36} style={{ color: '#D4B896' }} />
          </div>
          <div className="text-center">
            <p className="font-bold" style={{ color: '#6B4C2A' }}>Tidak ada pesanan aktif</p>
            <p className="text-sm mt-1" style={{ color: '#A08060' }}>Pesanan baru akan muncul di sini secara otomatis</p>
          </div>
        </div>
      ) : (
        <div className="space-y-6">
          {/* Bar station */}
          {barOrders.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-3">
                <div
                  className="flex items-center gap-2 px-3 py-1.5 rounded-xl"
                  style={{ backgroundColor: '#E0F2FE', color: '#0369A1' }}
                >
                  <span className="text-base">☕</span>
                  <span className="font-bold text-sm">Bar</span>
                </div>
                <span
                  className="text-xs font-bold px-2.5 py-1 rounded-full"
                  style={{ backgroundColor: '#0369A1', color: '#fff' }}
                >
                  {barOrders.length}
                </span>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                {barOrders.map((order) => (
                  <OrderTicket
                    key={order.id}
                    order={order}
                    onComplete={handleComplete}
                    completing={completingId === order.id}
                  />
                ))}
              </div>
            </div>
          )}

          {/* Kitchen station */}
          {kitchenOrders.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-3">
                <div
                  className="flex items-center gap-2 px-3 py-1.5 rounded-xl"
                  style={{ backgroundColor: '#FEF3C7', color: '#B45309' }}
                >
                  <span className="text-base">🍳</span>
                  <span className="font-bold text-sm">Dapur</span>
                </div>
                <span
                  className="text-xs font-bold px-2.5 py-1 rounded-full"
                  style={{ backgroundColor: '#B45309', color: '#fff' }}
                >
                  {kitchenOrders.length}
                </span>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                {kitchenOrders.map((order) => (
                  <OrderTicket
                    key={order.id}
                    order={order}
                    onComplete={handleComplete}
                    completing={completingId === order.id}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </DashboardLayout>
  );
}
