"use client";

import React, { useState } from "react";
import { toast } from "sonner";
import api from "@/lib/api/axios";
import { AxiosError } from "axios";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AlertCircle, Undo2 } from "lucide-react";

export interface RefundDialogProps {
  orderId: string | null;
  orderNumber: string;
  totalAmount: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export function RefundDialog({ 
  orderId, 
  orderNumber, 
  totalAmount, 
  open, 
  onOpenChange,
  onSuccess 
}: RefundDialogProps) {
  const [returnAmount, setReturnAmount] = useState<number | "">("");
  const [reason, setReason] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const maxRefund = totalAmount * 0.8; // Max 80% rule
  const isOverLimit = typeof returnAmount === "number" && returnAmount > maxRefund;
  const isSubmitDisabled = typeof returnAmount !== "number" || returnAmount <= 0 || !reason.trim() || isOverLimit || isSubmitting;

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const handleRefundSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orderId || isSubmitDisabled) return;

    setIsSubmitting(true);
    try {
      await api.post(`/orders/${orderId}/returns`, {
        return_amount: returnAmount,
        reason: reason,
      });

      toast.success("Refund Processed Successfully");
      setReturnAmount("");
      setReason("");
      onOpenChange(false);
      onSuccess();
    } catch (error) {
      const err = error as AxiosError<{ message: string }>;
      toast.error("Refund Failed", { 
        description: err.response?.data?.message || "An error occurred while processing the refund." 
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[450px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl text-rose-600">
            <Undo2 className="w-5 h-5" />
            Process Refund
          </DialogTitle>
          <DialogDescription>
            Order Number: <strong className="text-slate-800">{orderNumber}</strong>
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleRefundSubmit} className="space-y-6 py-2">
          
          <div className="bg-rose-50 border border-rose-100 p-4 rounded-lg flex gap-3 items-start">
            <AlertCircle className="w-5 h-5 text-rose-600 shrink-0 mt-0.5" />
            <div className="text-sm text-rose-800">
              <p className="font-semibold mb-1">Max 80% Refund Rule</p>
              <p>Refund amount cannot exceed 80% of the original order total ({formatCurrency(totalAmount)}).</p>
              <p className="font-bold mt-1">Maximum allowed: {formatCurrency(maxRefund)}</p>
            </div>
          </div>

          <div className="space-y-3">
            <Label htmlFor="returnAmount" className="text-slate-700">
              Refund Amount (Nominal Pengembalian)
            </Label>
            <Input
              id="returnAmount"
              type="number"
              placeholder="e.g. 15000"
              value={returnAmount}
              onChange={(e) => setReturnAmount(e.target.value ? Number(e.target.value) : "")}
              className={`text-lg font-medium h-12 ${isOverLimit ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
            />
            {isOverLimit && (
              <p className="text-sm text-red-500 font-medium">Amount exceeds the 80% maximum limit.</p>
            )}
          </div>

          <div className="space-y-3">
            <Label htmlFor="reason" className="text-slate-700">
              Reason for Refund
            </Label>
            <Input
              id="reason"
              placeholder="e.g. Customer cancelled item"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="h-11"
            />
          </div>

          <DialogFooter className="pt-2">
            <button
              type="button"
              onClick={() => onOpenChange(false)}
              className="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg font-medium transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitDisabled}
              className="px-6 py-2 bg-rose-600 hover:bg-rose-700 disabled:bg-slate-300 disabled:cursor-not-allowed text-white rounded-lg font-bold transition-all shadow-sm active:scale-95"
            >
              {isSubmitting ? "Processing..." : "Confirm Refund"}
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
