"use client";

import React, { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api/axios";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  History, RotateCcw, Loader2, ClipboardList,
  RefreshCw, ChevronDown, ChevronUp, Package,
} from "lucide-react";
import { RefundDialog } from "./RefundDialog";
import { Order } from "@/types";

interface HistoryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(amount);

const formatDate = (dateString: string) =>
  new Date(dateString).toLocaleString("id-ID", {
    year: "numeric", month: "short", day: "numeric",
    hour: "2-digit", minute: "2-digit",
  });

const orderTypeLabel: Record<string, string> = {
  dine_in: "Dine In",
  takeaway: "Takeaway",
  delivery: "Delivery",
  online: "Online",
};

const StatusBadge = ({ status }: { status: string }) => {
  const styles: Record<string, string> = {
    pending: "bg-amber-100 text-amber-700 border-amber-200",
    confirmed: "bg-emerald-100 text-emerald-700 border-emerald-200",
    cancelled: "bg-red-100 text-red-700 border-red-200",
    refunded: "bg-orange-100 text-orange-700 border-orange-200",
    void: "bg-slate-100 text-slate-500 border-slate-200",
  };
  const labels: Record<string, string> = {
    pending: "Pending",
    confirmed: "Selesai",
    cancelled: "Dibatalkan",
    refunded: "Refund",
    void: "Void",
  };
  const style = styles[status] || "bg-slate-100 text-slate-600 border-slate-200";
  return (
    <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-[11px] font-bold border uppercase tracking-wide ${style}`}>
      {labels[status] || status}
    </span>
  );
};

export function HistoryDialog({ open, onOpenChange }: HistoryDialogProps) {
  const [refundOrder, setRefundOrder] = useState<Order | null>(null);
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const { data: orders, isLoading, error, refetch } = useQuery({
    queryKey: ["orders"],
    queryFn: async () => {
      const response = await api.get("/orders");
      return response.data.data as Order[];
    },
    enabled: open,
    refetchOnWindowFocus: false,
  });

  const toggleExpand = (id: string) => {
    setExpandedId((prev) => (prev === id ? null : id));
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-[900px] max-h-[88vh] flex flex-col p-0 gap-0 rounded-3xl overflow-hidden">
          {/* Header */}
          <div className="bg-slate-900 px-6 py-5 text-white shrink-0">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-3 text-white text-xl font-bold">
                <div className="w-9 h-9 bg-white/15 rounded-xl flex items-center justify-center">
                  <History className="w-5 h-5" />
                </div>
                Riwayat Transaksi
              </DialogTitle>
            </DialogHeader>
            <p className="text-slate-400 text-sm mt-1">
              {orders ? `${orders.length} transaksi ditemukan` : "Memuat data transaksi..."}
            </p>
          </div>

          {/* Toolbar */}
          <div className="px-6 py-3 border-b border-slate-100 bg-slate-50 flex justify-between items-center shrink-0">
            <p className="text-xs text-slate-500">Klik baris untuk melihat detail item pesanan.</p>
            <button
              onClick={() => refetch()}
              className="flex items-center gap-1.5 text-xs text-slate-600 hover:text-emerald-600 font-medium transition-colors"
            >
              <RefreshCw className="w-3.5 h-3.5" />
              Refresh
            </button>
          </div>

          {/* Table content */}
          <div className="flex-1 overflow-auto">
            {isLoading && (
              <div className="flex flex-col items-center justify-center h-64 gap-3">
                <Loader2 className="w-8 h-8 text-emerald-600 animate-spin" />
                <p className="text-sm text-slate-500">Memuat riwayat transaksi...</p>
              </div>
            )}

            {error && !isLoading && (
              <div className="flex flex-col items-center justify-center h-64 gap-3 text-center p-6">
                <div className="w-14 h-14 bg-red-50 rounded-2xl flex items-center justify-center">
                  <Package className="w-7 h-7 text-red-400" />
                </div>
                <p className="text-sm font-semibold text-slate-700">Gagal memuat riwayat</p>
                <button
                  onClick={() => refetch()}
                  className="text-xs text-emerald-600 hover:underline font-medium"
                >
                  Coba lagi
                </button>
              </div>
            )}

            {!isLoading && !error && (!orders || orders.length === 0) && (
              <div className="flex flex-col items-center justify-center h-64 gap-3 text-center p-6">
                <div className="w-14 h-14 bg-slate-100 rounded-2xl flex items-center justify-center">
                  <ClipboardList className="w-7 h-7 text-slate-400" />
                </div>
                <p className="text-sm font-semibold text-slate-700">Belum ada transaksi</p>
                <p className="text-xs text-slate-500">Selesaikan checkout pertama untuk melihat riwayat di sini.</p>
              </div>
            )}

            {!isLoading && !error && orders && orders.length > 0 && (
              <Table>
                <TableHeader className="bg-slate-50 sticky top-0 z-10">
                  <TableRow>
                    <TableHead className="w-8 pl-6" />
                    <TableHead className="font-semibold text-slate-600 text-xs uppercase tracking-wider">No. Order</TableHead>
                    <TableHead className="font-semibold text-slate-600 text-xs uppercase tracking-wider">Waktu</TableHead>
                    <TableHead className="font-semibold text-slate-600 text-xs uppercase tracking-wider">Tipe</TableHead>
                    <TableHead className="font-semibold text-slate-600 text-xs uppercase tracking-wider">Pelanggan</TableHead>
                    <TableHead className="font-semibold text-slate-600 text-xs uppercase tracking-wider">Status</TableHead>
                    <TableHead className="text-right font-semibold text-slate-600 text-xs uppercase tracking-wider">Total</TableHead>
                    <TableHead className="text-center font-semibold text-slate-600 text-xs uppercase tracking-wider pr-6">Aksi</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {orders.map((order) => (
                    <React.Fragment key={order.id}>
                      {/* Main row */}
                      <TableRow
                        className="hover:bg-slate-50/80 transition-colors cursor-pointer group"
                        onClick={() => toggleExpand(order.id)}
                      >
                        <TableCell className="pl-6 pr-2">
                          {expandedId === order.id ? (
                            <ChevronUp className="w-3.5 h-3.5 text-slate-400" />
                          ) : (
                            <ChevronDown className="w-3.5 h-3.5 text-slate-300 group-hover:text-slate-400" />
                          )}
                        </TableCell>
                        <TableCell className="font-mono text-sm font-bold text-slate-800">
                          {order.order_number}
                        </TableCell>
                        <TableCell className="text-xs text-slate-500">{formatDate(order.created_at)}</TableCell>
                        <TableCell>
                          <span className="text-xs bg-slate-100 text-slate-600 px-2 py-0.5 rounded-lg font-medium">
                            {orderTypeLabel[order.order_type] || order.order_type}
                          </span>
                        </TableCell>
                        <TableCell className="text-sm text-slate-600">
                          {order.customer_name || "—"}
                        </TableCell>
                        <TableCell>
                          <StatusBadge status={order.status} />
                          {order.has_returns && (
                            <span className="ml-1.5 text-[10px] bg-amber-100 text-amber-600 px-1.5 py-0.5 rounded font-semibold border border-amber-200">
                              REFUND
                            </span>
                          )}
                        </TableCell>
                        <TableCell className="text-right font-bold text-emerald-600 text-sm">
                          {formatCurrency(order.total)}
                        </TableCell>
                        <TableCell className="text-center pr-6" onClick={(e) => e.stopPropagation()}>
                          <button
                            onClick={() => setRefundOrder(order)}
                            disabled={order.status !== "confirmed"}
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-rose-50 text-rose-600 hover:bg-rose-100 disabled:opacity-40 disabled:cursor-not-allowed rounded-lg text-xs font-bold transition-colors border border-rose-100"
                            title={
                              order.status !== "confirmed"
                                ? "Hanya order yang selesai yang bisa direfund"
                                : "Proses refund"
                            }
                          >
                            <RotateCcw className="w-3.5 h-3.5" />
                            Refund
                          </button>
                        </TableCell>
                      </TableRow>

                      {/* Expanded detail row */}
                      {expandedId === order.id && (
                        <TableRow className="bg-slate-50/50">
                          <TableCell colSpan={8} className="py-0">
                            <div className="px-6 py-4">
                              <p className="text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-3">
                                Detail Item Pesanan
                              </p>
                              <div className="grid gap-1.5">
                                {order.order_details && order.order_details.length > 0 ? (
                                  order.order_details.map((detail) => (
                                    <div
                                      key={detail.id}
                                      className={`flex items-center justify-between py-2 px-3 rounded-xl text-xs ${
                                        detail.is_voided
                                          ? "bg-red-50 text-red-400 line-through"
                                          : "bg-white border border-slate-100 text-slate-700"
                                      }`}
                                    >
                                      <div className="flex items-center gap-3">
                                        <div className="w-7 h-7 bg-slate-100 rounded-lg flex items-center justify-center font-bold text-[10px] text-slate-500">
                                          {detail.menu_name.substring(0, 2).toUpperCase()}
                                        </div>
                                        <span className="font-semibold">{detail.menu_name}</span>
                                        <span className="text-slate-400">×{detail.quantity}</span>
                                        {detail.is_voided && (
                                          <span className="text-[10px] bg-red-100 text-red-500 px-1.5 py-0.5 rounded font-bold">
                                            VOID
                                          </span>
                                        )}
                                      </div>
                                      <div className="flex items-center gap-4">
                                        <span className="text-slate-400">{formatCurrency(detail.price)} /pcs</span>
                                        <span className="font-bold text-slate-700">{formatCurrency(detail.subtotal)}</span>
                                      </div>
                                    </div>
                                  ))
                                ) : (
                                  <p className="text-xs text-slate-400 text-center py-2">Tidak ada detail item.</p>
                                )}
                              </div>
                              {/* Totals */}
                              <div className="mt-3 pt-3 border-t border-slate-200 flex justify-end">
                                <div className="space-y-1 text-xs text-right">
                                  <div className="flex gap-8 text-slate-500">
                                    <span>Subtotal</span><span className="w-24 text-right">{formatCurrency(order.subtotal)}</span>
                                  </div>
                                  <div className="flex gap-8 text-slate-500">
                                    <span>Pajak (11%)</span><span className="w-24 text-right">{formatCurrency(order.tax_amount)}</span>
                                  </div>
                                  <div className="flex gap-8 text-slate-500">
                                    <span>Service (5%)</span><span className="w-24 text-right">{formatCurrency(order.service_charge)}</span>
                                  </div>
                                  <div className="flex gap-8 font-bold text-slate-800 text-sm pt-1 border-t border-slate-200">
                                    <span>Total</span><span className="w-24 text-right text-emerald-600">{formatCurrency(order.total)}</span>
                                  </div>
                                </div>
                              </div>
                            </div>
                          </TableCell>
                        </TableRow>
                      )}
                    </React.Fragment>
                  ))}
                </TableBody>
              </Table>
            )}
          </div>
        </DialogContent>
      </Dialog>

      {/* Refund Dialog */}
      {refundOrder && (
        <RefundDialog
          open={!!refundOrder}
          onOpenChange={(isOpen) => {
            if (!isOpen) setRefundOrder(null);
          }}
          orderId={refundOrder.id}
          orderNumber={refundOrder.order_number}
          totalAmount={refundOrder.total}
          onSuccess={() => {
            setRefundOrder(null);
            refetch();
          }}
        />
      )}
    </>
  );
}
