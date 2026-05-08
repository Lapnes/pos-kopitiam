"use client";

import React, { useState } from "react";
import { useCartStore } from "@/store/useCartStore";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ShoppingCart, CheckCircle2 } from "lucide-react";

export function CheckoutDialog({ children }: { children: React.ReactNode }) {
  const { items, getGrandTotal, clearCart } = useCartStore();
  const [open, setOpen] = useState(false);
  const [cashReceived, setCashReceived] = useState<number | "">("");

  const grandTotal = getGrandTotal();
  const change = typeof cashReceived === "number" ? cashReceived - grandTotal : 0;
  const isSubmitDisabled = typeof cashReceived !== "number" || cashReceived < grandTotal;

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitDisabled) return;

    // Simulate API Call
    toast.success("Transaction Successful!", {
      description: `Change to return: ${formatCurrency(change)}`,
      icon: <CheckCircle2 className="w-5 h-5 text-emerald-500" />,
    });

    clearCart();
    setOpen(false);
    setCashReceived("");
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger render={React.isValidElement(children) ? children : undefined}>
        {!React.isValidElement(children) ? children : null}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl">
            <ShoppingCart className="w-5 h-5 text-emerald-600" />
            Checkout Validation
          </DialogTitle>
          <DialogDescription>
            Enter the amount of cash received from the customer.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6 py-4">
          <div className="space-y-2">
            <div className="flex justify-between items-center bg-slate-50 p-4 rounded-lg border border-slate-100">
              <span className="text-sm font-medium text-slate-500">Grand Total</span>
              <span className="text-xl font-bold text-slate-800">
                {formatCurrency(grandTotal)}
              </span>
            </div>
          </div>

          <div className="space-y-3">
            <Label htmlFor="cashReceived" className="text-slate-700">
              Cash Received (Nominal Pembayaran)
            </Label>
            <Input
              id="cashReceived"
              type="number"
              placeholder="e.g. 100000"
              value={cashReceived}
              onChange={(e) => setCashReceived(e.target.value ? Number(e.target.value) : "")}
              className="text-lg font-medium h-12"
              autoFocus
            />
          </div>

          <div className="space-y-2">
            <div className={`flex justify-between items-center p-4 rounded-lg border ${change >= 0 ? 'bg-emerald-50 border-emerald-100' : 'bg-red-50 border-red-100'}`}>
              <span className={`text-sm font-medium ${change >= 0 ? 'text-emerald-700' : 'text-red-600'}`}>
                Change (Kembalian)
              </span>
              <span className={`text-xl font-bold ${change >= 0 ? 'text-emerald-700' : 'text-red-600'}`}>
                {change >= 0 ? formatCurrency(change) : "-"}
              </span>
            </div>
          </div>

          <DialogFooter className="pt-2">
            <button
              type="button"
              onClick={() => setOpen(false)}
              className="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg font-medium transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitDisabled}
              className="px-6 py-2 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-300 disabled:cursor-not-allowed text-white rounded-lg font-bold transition-all shadow-sm active:scale-95"
            >
              Confirm Payment
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
