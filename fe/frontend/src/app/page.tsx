"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useCartStore } from "@/store/useCartStore";
import { useAuthStore } from "@/store/useAuthStore";
import { Menu } from "@/types";
import { Plus, Minus, Trash2, ShoppingCart, Store, UserCircle, LogOut, History } from "lucide-react";
import { CheckoutDialog } from "@/components/pos/CheckoutDialog";
import { HistoryDialog } from "@/components/pos/HistoryDialog";

import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api/axios";

const getMenuColor = (name: string) => {
  const colors = [
    "bg-orange-100 text-orange-700",
    "bg-yellow-100 text-yellow-700",
    "bg-amber-100 text-amber-800",
    "bg-red-100 text-red-800",
    "bg-stone-200 text-stone-800",
    "bg-rose-100 text-rose-800",
    "bg-emerald-100 text-emerald-800",
    "bg-blue-100 text-blue-800",
    "bg-indigo-100 text-indigo-800",
    "bg-violet-100 text-violet-800",
    "bg-fuchsia-100 text-fuchsia-800",
    "bg-pink-100 text-pink-800"
  ];
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  return colors[Math.abs(hash) % colors.length];
};

export default function CashierPage() {
  const router = useRouter();
  const { user, isAuthenticated, logout } = useAuthStore();
  
  const {
    items,
    addItem,
    updateQuantity,
    removeItem,
    getSubtotal,
    getTaxAmount,
    getServiceCharge,
    getGrandTotal,
  } = useCartStore();

  const [historyOpen, setHistoryOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  const { data: menus = [], isLoading: isLoadingMenus } = useQuery({
    queryKey: ["menus"],
    queryFn: async () => {
      const response = await api.get("/menus");
      return response.data.data as Menu[];
    },
    enabled: isAuthenticated,
  });

  useEffect(() => {
    setMounted(true);
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

  // Prevent hydration errors by not rendering until mounted
  if (!mounted || !isAuthenticated) return null;

  return (
    <div className="flex flex-col h-screen bg-slate-50 text-slate-900 overflow-hidden font-sans">
      {/* Top Header */}
      <header className="flex items-center justify-between px-6 py-4 bg-white shadow-sm shrink-0 z-20">
        <div className="flex items-center gap-2">
          <Store className="w-6 h-6 text-emerald-600" />
          <h1 className="text-xl font-bold tracking-tight text-slate-800">POS KopiTiam</h1>
        </div>
        
        <div className="flex items-center gap-4">
          <button 
            onClick={() => setHistoryOpen(true)}
            className="flex items-center gap-2 px-3 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg text-sm font-semibold transition-colors"
          >
            <History className="w-4 h-4" />
            Transaction History
          </button>
          
          <div className="w-px h-8 bg-slate-200 mx-1"></div>

          <div className="flex items-center gap-3">
            <div className="flex flex-col items-end">
              <span className="text-sm font-semibold text-slate-700">{user?.name || "Cashier"}</span>
              <span className="text-xs text-emerald-600 font-medium">Active Shift</span>
            </div>
            <div className="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center border border-slate-200">
              <UserCircle className="w-6 h-6 text-slate-500" />
            </div>
          </div>
          
          <button 
            onClick={handleLogout}
            className="ml-2 p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
            title="Logout"
          >
            <LogOut className="w-5 h-5" />
          </button>
        </div>
      </header>

      {/* Main Split-Screen */}
      <main className="flex-1 flex overflow-hidden">
        {/* Left Panel: Menu Grid (70%) */}
        <div className="w-[70%] h-full overflow-y-auto p-6 bg-slate-50/50">
          <h2 className="text-2xl font-bold text-slate-800 mb-6">Menu</h2>
          
          {isLoadingMenus ? (
            <div className="flex items-center justify-center h-40">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-600"></div>
            </div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {menus.map((menu) => (
                <button
                  key={menu.id}
                  onClick={() => addItem(menu)}
                  className="flex flex-col h-40 bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-all duration-200 border border-slate-100 hover:border-emerald-300 active:scale-95 text-left group"
                >
                  <div className={`flex-1 w-full flex items-center justify-center ${menu.color || getMenuColor(menu.name)} font-bold text-lg opacity-90 group-hover:opacity-100 transition-opacity`}>
                    {menu.name.substring(0, 2).toUpperCase()}
                  </div>
                  <div className="p-3 bg-white w-full">
                    <h3 className="font-semibold text-slate-800 line-clamp-1">{menu.name}</h3>
                    <div className="font-bold text-emerald-600 mt-1">{formatCurrency(menu.price)}</div>
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Right Panel: Cart Sidebar (30%) */}
        <div className="w-[30%] h-full bg-white border-l border-slate-200 flex flex-col shadow-xl z-10">
          <div className="p-5 border-b border-slate-100 flex items-center gap-2 shrink-0">
            <ShoppingCart className="w-5 h-5 text-slate-700" />
            <h2 className="text-lg font-bold text-slate-800">Current Order</h2>
            <span className="ml-auto bg-emerald-100 text-emerald-700 text-xs font-bold px-2.5 py-1 rounded-full">
              {items.reduce((acc, item) => acc + item.quantity, 0)} items
            </span>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {items.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-slate-400">
                <ShoppingCart className="w-12 h-12 text-slate-200 mb-2" />
                <p className="text-sm">Cart is empty</p>
              </div>
            ) : (
              items.map((item) => (
                <div key={item.cartItemId} className="flex flex-col bg-white p-3 rounded-xl border border-slate-100 shadow-sm">
                  <div className="flex justify-between items-start mb-2">
                    <h4 className="font-semibold text-sm text-slate-800 pr-2">{item.name}</h4>
                    <button
                      onClick={() => removeItem(item.cartItemId)}
                      className="text-slate-300 hover:text-red-500 transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                  <div className="text-emerald-600 font-medium text-sm mb-3">
                    {formatCurrency(item.price)}
                  </div>
                  <div className="flex justify-between items-center">
                    <div className="flex items-center bg-slate-50 rounded-lg border border-slate-200">
                      <button
                        onClick={() => updateQuantity(item.cartItemId, item.quantity - 1)}
                        className="w-8 h-8 flex items-center justify-center text-slate-600 hover:bg-slate-200 transition-colors"
                      >
                        <Minus className="w-3.5 h-3.5" />
                      </button>
                      <span className="w-8 text-center text-sm font-semibold text-slate-800">
                        {item.quantity}
                      </span>
                      <button
                        onClick={() => updateQuantity(item.cartItemId, item.quantity + 1)}
                        className="w-8 h-8 flex items-center justify-center text-slate-600 hover:bg-slate-200 transition-colors"
                      >
                        <Plus className="w-3.5 h-3.5" />
                      </button>
                    </div>
                    <span className="font-semibold text-slate-800">
                      {formatCurrency(item.price * item.quantity)}
                    </span>
                  </div>
                </div>
              ))
            )}
          </div>

          <div className="p-5 bg-slate-50 border-t border-slate-200 shrink-0">
            <div className="space-y-2 mb-4 text-sm">
              <div className="flex justify-between text-slate-500">
                <span>Subtotal</span>
                <span className="font-medium text-slate-700">{formatCurrency(getSubtotal())}</span>
              </div>
              <div className="flex justify-between text-slate-500">
                <span>Tax (11%)</span>
                <span className="font-medium text-slate-700">{formatCurrency(getTaxAmount())}</span>
              </div>
              <div className="flex justify-between text-slate-500">
                <span>Service Charge (5%)</span>
                <span className="font-medium text-slate-700">{formatCurrency(getServiceCharge())}</span>
              </div>
              <div className="pt-3 mt-3 border-t border-dashed border-slate-300 flex justify-between items-center">
                <span className="font-bold text-slate-800 text-base">Grand Total</span>
                <span className="font-bold text-2xl text-emerald-600">{formatCurrency(getGrandTotal())}</span>
              </div>
            </div>

            <CheckoutDialog>
              <button
                disabled={items.length === 0}
                className="w-full py-4 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-300 disabled:text-slate-500 disabled:cursor-not-allowed text-white font-bold rounded-xl transition-all shadow-md active:scale-[0.98]"
              >
                Checkout
              </button>
            </CheckoutDialog>
          </div>
        </div>
      </main>

      <HistoryDialog open={historyOpen} onOpenChange={setHistoryOpen} />
    </div>
  );
}
