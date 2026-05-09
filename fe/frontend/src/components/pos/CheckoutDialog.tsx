"use client";

import React, { useState } from "react";
import { useCartStore } from "@/store/useCartStore";
import { toast } from "sonner";
import api from "@/lib/api/axios";
import { AxiosError } from "axios";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  ShoppingCart, CheckCircle2, Banknote, Calculator,
  Receipt, X, Loader2, ArrowRight,
} from "lucide-react";

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(amount);

// Seeded Table 1 UUID
const DEFAULT_TABLE_ID = "91610a61-5f2f-444a-a97f-1f0d6ad185a5";

export function CheckoutDialog({ children }: { children: React.ReactNode }) {
  const { items, getSubtotal, getTaxAmount, getServiceCharge, getGrandTotal, clearCart } = useCartStore();
  const [open, setOpen] = useState(false);
  const [cashReceived, setCashReceived] = useState<number | "">("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const subtotal = getSubtotal();
  const taxAmount = getTaxAmount();
  const serviceCharge = getServiceCharge();
  const grandTotal = getGrandTotal();
  const change = typeof cashReceived === "number" ? cashReceived - grandTotal : 0;
  const isEnoughCash = typeof cashReceived === "number" && cashReceived >= grandTotal;
  const isSubmitDisabled = !isEnoughCash || isSubmitting || items.length === 0;

  const quickAmounts = [
    Math.ceil(grandTotal / 10000) * 10000,
    Math.ceil(grandTotal / 50000) * 50000,
    Math.ceil(grandTotal / 100000) * 100000,
  ].filter((v, i, arr) => arr.indexOf(v) === i && v >= grandTotal);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitDisabled) return;

    setIsSubmitting(true);
    try {
      const payload = {
        table_id: DEFAULT_TABLE_ID,
        customer_name: "Walk-in Customer",
        customer_phone: "",
        order_type: "dine_in",
        notes: "",
        discount_amount: 0,
        items: items.map((item) => ({
          menu_id: item.id,
          menu_name: item.name,
          quantity: item.quantity,
          price: item.price,
          notes: "",
        })),
      };

      await api.post("/orders", payload);

      toast.success("Transaksi Berhasil!", {
        description: `Kembalian: ${formatCurrency(change)}`,
        icon: <CheckCircle2 className="w-5 h-5 text-emerald-500" />,
        duration: 5000,
      });

      clearCart();
      setOpen(false);
      setCashReceived("");
    } catch (error) {
      const err = error as AxiosError<{ message: string }>;
      console.error("CHECKOUT ERROR:", err.response?.data || err.message);
      toast.error("Checkout Gagal", {
        description: err.response?.data?.message || "Terjadi error saat memproses pesanan.",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleOpenChange = (val: boolean) => {
    if (!isSubmitting) {
      setOpen(val);
      if (!val) setCashReceived("");
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger render={<div onClick={() => items.length > 0 && setOpen(true)}>{children}</div>} />
      <DialogContent className="sm:max-w-[480px] p-0 gap-0 overflow-hidden rounded-3xl">

        {/* Header */}
        <div className="bg-emerald-600 px-6 py-5 text-white">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2.5 text-white text-lg font-bold">
              <div className="w-8 h-8 bg-white/20 rounded-xl flex items-center justify-center">
                <Receipt className="w-4 h-4" />
              </div>
              Konfirmasi Pembayaran
            </DialogTitle>
          </DialogHeader>
          <p className="text-emerald-100 text-sm mt-1">Masukkan jumlah uang yang diterima dari pelanggan.</p>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          {/* Order summary */}
          <div className="bg-slate-50 border border-slate-100 rounded-2xl p-4 space-y-2">
            <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">Ringkasan Pesanan</p>

            <div className="max-h-28 overflow-y-auto space-y-1.5 pr-1">
              {items.map((item) => (
                <div key={item.cartItemId} className="flex justify-between items-center text-xs">
                  <span className="text-slate-600">
                    {item.name} <span className="text-slate-400">×{item.quantity}</span>
                  </span>
                  <span className="font-semibold text-slate-700">
                    {formatCurrency(item.price * item.quantity)}
                  </span>
                </div>
              ))}
            </div>

            <div className="border-t border-slate-200 pt-2 mt-2 space-y-1">
              <div className="flex justify-between text-xs text-slate-500">
                <span>Subtotal</span><span>{formatCurrency(subtotal)}</span>
              </div>
              <div className="flex justify-between text-xs text-slate-500">
                <span>Pajak (11%)</span><span>{formatCurrency(taxAmount)}</span>
              </div>
              <div className="flex justify-between text-xs text-slate-500">
                <span>Service (5%)</span><span>{formatCurrency(serviceCharge)}</span>
              </div>
            </div>

            <div className="border-t-2 border-dashed border-slate-300 pt-3 flex justify-between items-center">
              <span className="font-bold text-slate-800">Total Bayar</span>
              <span className="text-2xl font-black text-emerald-600">{formatCurrency(grandTotal)}</span>
            </div>
          </div>

          {/* Cash input */}
          <div className="space-y-3">
            <Label htmlFor="cashReceived" className="text-slate-700 font-semibold flex items-center gap-2">
              <Banknote className="w-4 h-4 text-emerald-600" />
              Nominal Pembayaran
            </Label>
            <Input
              id="cashReceived"
              type="number"
              placeholder="Masukkan jumlah uang..."
              value={cashReceived}
              onChange={(e) => setCashReceived(e.target.value ? Number(e.target.value) : "")}
              className="text-lg font-bold h-13 rounded-xl border-slate-200 focus:border-emerald-400 focus:ring-emerald-400/20"
              autoFocus
              min={0}
            />

            {/* Quick amounts */}
            {quickAmounts.length > 0 && (
              <div className="flex gap-2 flex-wrap">
                {quickAmounts.slice(0, 3).map((amount) => (
                  <button
                    key={amount}
                    type="button"
                    onClick={() => setCashReceived(amount)}
                    className={`flex-1 py-2 px-3 text-xs font-bold rounded-xl border transition-all ${
                      cashReceived === amount
                        ? "bg-emerald-500 text-white border-emerald-500"
                        : "bg-white text-slate-700 border-slate-200 hover:border-emerald-300 hover:text-emerald-700"
                    }`}
                  >
                    {formatCurrency(amount)}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Change display */}
          <div className={`flex justify-between items-center p-4 rounded-2xl border-2 transition-all ${
            typeof cashReceived !== "number"
              ? "bg-slate-50 border-slate-200"
              : isEnoughCash
              ? "bg-emerald-50 border-emerald-200"
              : "bg-red-50 border-red-200"
          }`}>
            <div className="flex items-center gap-2">
              <Calculator className={`w-5 h-5 ${
                typeof cashReceived !== "number" ? "text-slate-400" : isEnoughCash ? "text-emerald-600" : "text-red-500"
              }`} />
              <span className={`font-semibold text-sm ${
                typeof cashReceived !== "number" ? "text-slate-500" : isEnoughCash ? "text-emerald-700" : "text-red-600"
              }`}>
                Kembalian
              </span>
            </div>
            <span className={`text-2xl font-black ${
              typeof cashReceived !== "number" ? "text-slate-400" : isEnoughCash ? "text-emerald-700" : "text-red-600"
            }`}>
              {typeof cashReceived !== "number"
                ? "—"
                : isEnoughCash
                ? formatCurrency(change)
                : "Kurang " + formatCurrency(Math.abs(change))}
            </span>
          </div>

          {/* Actions */}
          <DialogFooter className="gap-2 pt-1">
            <button
              type="button"
              onClick={() => handleOpenChange(false)}
              disabled={isSubmitting}
              className="flex-1 py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl font-semibold transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
            >
              <X className="w-4 h-4" />
              Batal
            </button>
            <button
              type="submit"
              disabled={isSubmitDisabled}
              className="flex-1 py-3 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-200 disabled:text-slate-400 disabled:cursor-not-allowed text-white rounded-xl font-bold transition-all shadow-md shadow-emerald-200 active:scale-[0.98] flex items-center justify-center gap-2"
            >
              {isSubmitting ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <>
                  Konfirmasi
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
