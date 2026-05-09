"use client";

import React, { useState } from "react";
import { toast } from "sonner";
import api from "@/lib/api/axios";
import { AxiosError } from "axios";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AlertCircle, Undo2, X, Loader2, CheckCircle2 } from "lucide-react";

export interface RefundDialogProps {
  orderId: string | null;
  orderNumber: string;
  totalAmount: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(amount);

export function RefundDialog({
  orderId,
  orderNumber,
  totalAmount,
  open,
  onOpenChange,
  onSuccess,
}: RefundDialogProps) {
  const [returnAmount, setReturnAmount] = useState<number | "">("");
  const [reason, setReason] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const maxRefund = totalAmount * 0.8;
  const isOverLimit = typeof returnAmount === "number" && returnAmount > maxRefund;
  const isUnderMin = typeof returnAmount === "number" && returnAmount <= 0;
  const isSubmitDisabled =
    typeof returnAmount !== "number" || isOverLimit || isUnderMin || !reason.trim() || isSubmitting;

  const pct = typeof returnAmount === "number" && totalAmount > 0
    ? Math.min((returnAmount / totalAmount) * 100, 100)
    : 0;

  const handleRefundSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orderId || isSubmitDisabled) return;

    setIsSubmitting(true);
    try {
      await api.post(`/orders/${orderId}/returns`, {
        return_amount: returnAmount,
        reason: reason,
      });

      toast.success("Refund Berhasil Diproses", {
        description: `${formatCurrency(returnAmount as number)} telah dikembalikan.`,
        icon: <CheckCircle2 className="w-5 h-5 text-emerald-500" />,
      });
      setReturnAmount("");
      setReason("");
      onOpenChange(false);
      onSuccess();
    } catch (error) {
      const err = error as AxiosError<{ message: string }>;
      toast.error("Refund Gagal", {
        description: err.response?.data?.message || "Terjadi error saat memproses refund.",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    if (!isSubmitting) {
      setReturnAmount("");
      setReason("");
      onOpenChange(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[460px] p-0 gap-0 rounded-3xl overflow-hidden">
        {/* Header */}
        <div className="bg-rose-600 px-6 py-5 text-white">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2.5 text-white text-lg font-bold">
              <div className="w-8 h-8 bg-white/20 rounded-xl flex items-center justify-center">
                <Undo2 className="w-4 h-4" />
              </div>
              Proses Refund
            </DialogTitle>
          </DialogHeader>
          <p className="text-rose-100 text-sm mt-1">
            Order: <strong className="text-white">{orderNumber}</strong>
          </p>
        </div>

        <form onSubmit={handleRefundSubmit} className="p-6 space-y-5">
          {/* 80% rule warning */}
          <div className="bg-amber-50 border border-amber-200 rounded-2xl p-4 flex gap-3 items-start">
            <AlertCircle className="w-5 h-5 text-amber-500 shrink-0 mt-0.5" />
            <div className="text-sm text-amber-800">
              <p className="font-bold mb-1">Aturan Refund Maksimal 80%</p>
              <div className="space-y-0.5 text-xs">
                <div className="flex justify-between">
                  <span>Total Order:</span>
                  <span className="font-bold">{formatCurrency(totalAmount)}</span>
                </div>
                <div className="flex justify-between text-amber-700">
                  <span>Maks. Refund (80%):</span>
                  <span className="font-bold">{formatCurrency(maxRefund)}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Amount input */}
          <div className="space-y-2">
            <Label htmlFor="returnAmount" className="text-slate-700 font-semibold">
              Jumlah Refund
            </Label>
            <Input
              id="returnAmount"
              type="number"
              placeholder="Masukkan jumlah refund..."
              value={returnAmount}
              onChange={(e) => setReturnAmount(e.target.value ? Number(e.target.value) : "")}
              className={`text-lg font-bold h-12 rounded-xl ${
                isOverLimit
                  ? "border-red-400 focus:border-red-500 focus:ring-red-400/20"
                  : "border-slate-200 focus:border-rose-400 focus:ring-rose-400/20"
              }`}
              min={1}
              max={maxRefund}
              autoFocus
            />

            {/* Progress bar */}
            {typeof returnAmount === "number" && returnAmount > 0 && (
              <div className="space-y-1.5">
                <div className="w-full h-2 bg-slate-100 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full transition-all duration-300 ${isOverLimit ? "bg-red-400" : "bg-rose-400"}`}
                    style={{ width: `${Math.min(pct, 100)}%` }}
                  />
                </div>
                <p className={`text-xs font-medium ${isOverLimit ? "text-red-500" : "text-slate-500"}`}>
                  {isOverLimit
                    ? `⚠ Melebihi batas maksimal (${pct.toFixed(0)}% dari total)`
                    : `${pct.toFixed(0)}% dari total order`}
                </p>
              </div>
            )}
          </div>

          {/* Reason input */}
          <div className="space-y-2">
            <Label htmlFor="reason" className="text-slate-700 font-semibold">
              Alasan Refund
            </Label>
            <Input
              id="reason"
              placeholder="cth: Pelanggan membatalkan item..."
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="h-11 rounded-xl border-slate-200 focus:border-rose-400 focus:ring-rose-400/20"
            />
            {!reason.trim() && returnAmount !== "" && (
              <p className="text-xs text-rose-500">Alasan refund wajib diisi.</p>
            )}
          </div>

          {/* Footer */}
          <DialogFooter className="gap-2 pt-1">
            <button
              type="button"
              onClick={handleClose}
              disabled={isSubmitting}
              className="flex-1 py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl font-semibold transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
            >
              <X className="w-4 h-4" />
              Batal
            </button>
            <button
              type="submit"
              disabled={isSubmitDisabled}
              className="flex-1 py-3 bg-rose-600 hover:bg-rose-700 disabled:bg-slate-200 disabled:text-slate-400 disabled:cursor-not-allowed text-white rounded-xl font-bold transition-all shadow-md shadow-rose-200 active:scale-[0.98] flex items-center justify-center gap-2"
            >
              {isSubmitting ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <>
                  <Undo2 className="w-4 h-4" />
                  Konfirmasi Refund
                </>
              )}
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
