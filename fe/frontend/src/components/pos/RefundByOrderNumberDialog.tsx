"use client";

import React, { useState } from "react";
import api from "@/lib/api/axios";
import { AxiosError } from "axios";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Search, RotateCcw, X, Loader2 } from "lucide-react";
import { RefundDialog } from "./RefundDialog";
import { Order } from "@/types";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export function RefundByOrderNumberDialog({ open, onOpenChange, onSuccess }: Props) {
  const [orderNumber, setOrderNumber] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [foundOrder, setFoundOrder] = useState<Order | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orderNumber.trim()) return;

    setIsLoading(true);
    try {
      const res = await api.get(`/api/v1/orders/by-number/${orderNumber.trim()}`);
      const order = res.data.data as Order;

      if (order.status !== "confirmed") {
        toast.error("Tidak Bisa Direfund", {
          description: `Order ini berstatus ${order.status}, hanya order selesai yang bisa direfund.`,
        });
        return;
      }

      if (order.has_returns) {
        toast.error("Sudah Direfund", {
          description: "Order ini sudah pernah direfund sebelumnya.",
        });
        return;
      }

      setFoundOrder(order);
    } catch (error) {
      const err = error as AxiosError<{ message: string }>;
      toast.error("Order Tidak Ditemukan", {
        description: err.response?.data?.message || "Pastikan nomor order sudah benar.",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    setOrderNumber("");
    setFoundOrder(null);
    onOpenChange(false);
  };

  return (
    <>
      <Dialog open={open && !foundOrder} onOpenChange={handleClose}>
        <DialogContent className="sm:max-w-[400px] p-0 gap-0 rounded-3xl overflow-hidden">
          <div className="bg-slate-900 px-6 py-5 text-white">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2.5 text-white text-lg font-bold">
                <div className="w-8 h-8 bg-white/20 rounded-xl flex items-center justify-center">
                  <RotateCcw className="w-4 h-4" />
                </div>
                Cari Transaksi
              </DialogTitle>
            </DialogHeader>
            <p className="text-slate-400 text-sm mt-1">
              Masukkan nomor order untuk melakukan refund
            </p>
          </div>

          <form onSubmit={handleSearch} className="p-6 space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="orderNum" className="text-slate-700 font-semibold">
                Nomor Order
              </Label>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                <Input
                  id="orderNum"
                  autoFocus
                  placeholder="cth: ORD-20260509-1234"
                  value={orderNumber}
                  onChange={(e) => setOrderNumber(e.target.value)}
                  className="pl-9 h-11 rounded-xl border-slate-200 font-mono focus:border-slate-400 focus:ring-slate-400/20"
                />
              </div>
            </div>

            <DialogFooter className="gap-2 pt-2">
              <button
                type="button"
                onClick={handleClose}
                className="flex-1 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl font-semibold transition-colors flex items-center justify-center gap-2"
              >
                <X className="w-4 h-4" />
                Batal
              </button>
              <button
                type="submit"
                disabled={!orderNumber.trim() || isLoading}
                className="flex-1 py-2.5 bg-slate-900 hover:bg-slate-800 text-white disabled:bg-slate-300 disabled:cursor-not-allowed rounded-xl font-bold transition-colors flex items-center justify-center gap-2"
              >
                {isLoading ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    <Search className="w-4 h-4" />
                    Cari
                  </>
                )}
              </button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {foundOrder && (
        <RefundDialog
          open={!!foundOrder}
          onOpenChange={(isOpen) => {
            if (!isOpen) {
              setFoundOrder(null);
            }
          }}
          orderId={foundOrder.id}
          orderNumber={foundOrder.order_number}
          totalAmount={foundOrder.total}
          onSuccess={() => {
            setFoundOrder(null);
            handleClose();
            onSuccess();
          }}
        />
      )}
    </>
  );
}
