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
import { History, RotateCcw } from "lucide-react";
import { RefundDialog } from "./RefundDialog";
import { Order } from "@/types";

interface HistoryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function HistoryDialog({ open, onOpenChange }: HistoryDialogProps) {
  const [refundOrder, setRefundOrder] = useState<Order | null>(null);

  const { data: orders, isLoading, refetch } = useQuery({
    queryKey: ["orders"],
    queryFn: async () => {
      const response = await api.get("/orders");
      return response.data.data as Order[];
    },
    enabled: open, // Only fetch when dialog is open
  });

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const renderStatusBadge = (status: string) => {
    switch (status) {
      case "completed":
        return <span className="bg-emerald-100 text-emerald-700 px-2.5 py-1 rounded-full text-xs font-bold uppercase">Completed</span>;
      case "cancelled":
        return <span className="bg-red-100 text-red-700 px-2.5 py-1 rounded-full text-xs font-bold uppercase">Cancelled</span>;
      case "pending":
        return <span className="bg-amber-100 text-amber-700 px-2.5 py-1 rounded-full text-xs font-bold uppercase">Pending</span>;
      default:
        return <span className="bg-slate-100 text-slate-700 px-2.5 py-1 rounded-full text-xs font-bold uppercase">{status}</span>;
    }
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-[800px] max-h-[85vh] flex flex-col">
          <DialogHeader className="shrink-0">
            <DialogTitle className="flex items-center gap-2 text-2xl font-bold text-slate-800">
              <History className="w-6 h-6 text-emerald-600" />
              Transaction History
            </DialogTitle>
          </DialogHeader>

          <div className="flex-1 overflow-auto mt-4 border rounded-xl shadow-sm">
            <Table>
              <TableHeader className="bg-slate-50 sticky top-0 z-10 shadow-sm">
                <TableRow>
                  <TableHead>Order #</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Grand Total</TableHead>
                  <TableHead className="text-center">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center py-10 text-slate-500">
                      Loading orders...
                    </TableCell>
                  </TableRow>
                ) : orders?.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center py-10 text-slate-500">
                      No transaction history found.
                    </TableCell>
                  </TableRow>
                ) : (
                  orders?.map((order) => (
                    <TableRow key={order.id} className="hover:bg-slate-50/50 transition-colors">
                      <TableCell className="font-medium text-slate-800">{order.orderNumber}</TableCell>
                      <TableCell className="text-slate-500">{formatDate(order.createdAt)}</TableCell>
                      <TableCell>{renderStatusBadge(order.status)}</TableCell>
                      <TableCell className="text-right font-bold text-emerald-600">
                        {formatCurrency(order.grandTotal)}
                      </TableCell>
                      <TableCell className="text-center">
                        <button
                          onClick={() => setRefundOrder(order)}
                          disabled={order.status !== "completed"}
                          className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-rose-50 text-rose-600 hover:bg-rose-100 disabled:opacity-50 disabled:hover:bg-rose-50 rounded-lg text-xs font-bold transition-colors"
                        >
                          <RotateCcw className="w-3.5 h-3.5" />
                          Refund
                        </button>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </DialogContent>
      </Dialog>

      {/* Nested Refund Dialog */}
      {refundOrder && (
        <RefundDialog
          open={!!refundOrder}
          onOpenChange={(isOpen) => {
            if (!isOpen) setRefundOrder(null);
          }}
          orderId={refundOrder.id}
          orderNumber={refundOrder.orderNumber}
          totalAmount={refundOrder.grandTotal}
          onSuccess={() => refetch()}
        />
      )}
    </>
  );
}
