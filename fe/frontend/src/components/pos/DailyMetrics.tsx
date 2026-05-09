"use client";

import React from "react";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api/axios";
import { AnalyticsSummary } from "@/types";
import { TrendingUp, ShoppingBag, Loader2 } from "lucide-react";

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(amount);

const getTodayRange = () => {
  const now = new Date();
  const start = new Date(now.getFullYear(), now.getMonth(), now.getDate()).toISOString();
  const end = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59).toISOString();
  return { start, end };
};

export function DailyMetrics() {
  const { start, end } = getTodayRange();

  const { data, isLoading } = useQuery<AnalyticsSummary>({
    queryKey: ["daily-metrics", start],
    queryFn: async () => {
      const res = await api.get("/analytics/sales-summary", {
        params: { start_date: start, end_date: end },
      });
      return res.data.data as AnalyticsSummary;
    },
    staleTime: 60 * 1000, // 1 min
    retry: 1,
  });

  return (
    <div className="flex gap-3 mb-4">
      {/* Total Transaksi */}
      <div className="flex-1 bg-white border border-slate-100 rounded-2xl px-4 py-3 flex items-center gap-3 shadow-sm">
        <div className="w-9 h-9 bg-sky-50 rounded-xl flex items-center justify-center shrink-0">
          <ShoppingBag className="w-4.5 h-4.5 text-sky-500" />
        </div>
        <div>
          <p className="text-[10px] text-slate-400 font-semibold uppercase tracking-wide leading-none">
            Penjualan Hari Ini
          </p>
          {isLoading ? (
            <Loader2 className="w-4 h-4 text-slate-300 animate-spin mt-1" />
          ) : (
            <p className="text-lg font-black text-slate-800 leading-tight mt-0.5">
              {data?.total_transactions ?? data?.total_orders ?? "—"}
              <span className="text-xs font-semibold text-slate-400 ml-1">transaksi</span>
            </p>
          )}
        </div>
      </div>

      {/* Total Pendapatan */}
      <div className="flex-1 bg-white border border-slate-100 rounded-2xl px-4 py-3 flex items-center gap-3 shadow-sm">
        <div className="w-9 h-9 bg-emerald-50 rounded-xl flex items-center justify-center shrink-0">
          <TrendingUp className="w-4.5 h-4.5 text-emerald-500" />
        </div>
        <div>
          <p className="text-[10px] text-slate-400 font-semibold uppercase tracking-wide leading-none">
            Pendapatan Hari Ini
          </p>
          {isLoading ? (
            <Loader2 className="w-4 h-4 text-slate-300 animate-spin mt-1" />
          ) : (
            <p className="text-lg font-black text-emerald-600 leading-tight mt-0.5">
              {data?.total_revenue !== undefined ? formatCurrency(data.total_revenue) : "—"}
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
