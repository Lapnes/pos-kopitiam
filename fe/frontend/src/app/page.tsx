"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { useCartStore } from "@/store/useCartStore";
import { useAuthStore } from "@/store/useAuthStore";
import { Menu } from "@/types";
import {
  Minus, Trash2, ShoppingCart, Coffee,
  UserCircle, LogOut, History, Loader2,
  ChefHat, Package, Search, Receipt,
  LayoutGrid, Users, BarChart2, ClipboardList,
} from "lucide-react";
import { CheckoutDialog } from "@/components/pos/CheckoutDialog";
import { HistoryDialog } from "@/components/pos/HistoryDialog";
import { MenuCard } from "@/components/pos/MenuCard";
import { DailyMetrics } from "@/components/pos/DailyMetrics";
import { EmployeePanel } from "@/components/employees/EmployeePanel";
import { AnalyticsDashboard } from "@/components/analytics/AnalyticsDashboard";
import api from "@/lib/api/axios";

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(amount);

const getMenuColor = (name: string): string => {
  const colors = [
    "bg-orange-100 text-orange-700",
    "bg-emerald-100 text-emerald-700",
    "bg-sky-100 text-sky-700",
    "bg-violet-100 text-violet-700",
    "bg-rose-100 text-rose-700",
    "bg-amber-100 text-amber-700",
  ];
  let hash = 0;
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash);
  return colors[Math.abs(hash) % colors.length];
};

const getRoleLabel = (role: string) => {
  const map: Record<string, string> = {
    cashier: "Kasir", manager: "Manajer", admin: "Admin",
    superadmin: "Super Admin", kitchen: "Dapur",
  };
  return map[role] || role;
};

type View = "pos" | "employees" | "analytics" | "history";

interface NavTab {
  id: View;
  label: string;
  Icon: React.ElementType;
  roles: string[];
}

const NAV_TABS: NavTab[] = [
  { id: "pos", label: "POS", Icon: LayoutGrid, roles: ["cashier", "superadmin", "admin", "manager", "kitchen"] },
  { id: "employees", label: "Karyawan", Icon: Users, roles: ["superadmin"] },
  { id: "analytics", label: "Analitik", Icon: BarChart2, roles: ["superadmin", "manager"] },
  { id: "history", label: "Riwayat Global", Icon: ClipboardList, roles: ["superadmin", "manager"] },
];

export default function CashierPage() {
  const router = useRouter();
  const { user, isAuthenticated, logout } = useAuthStore();
  const {
    items, addItem, updateQuantity, removeItem,
    getSubtotal, getTaxAmount, getServiceCharge, getGrandTotal,
  } = useCartStore();

  const [activeView, setActiveView] = useState<View>("pos");
  const [historyOpen, setHistoryOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [search, setSearch] = useState("");

  useEffect(() => {
    setMounted(true);
    if (!isAuthenticated) router.replace("/login");
  }, [isAuthenticated, router]);

  const { data: menus = [], isLoading: isLoadingMenus, error: menusError } = useQuery({
    queryKey: ["menus"],
    queryFn: async () => {
      const response = await api.get("/menus");
      return response.data.data as Menu[];
    },
    enabled: isAuthenticated && mounted,
    staleTime: 5 * 60 * 1000,
  });

  const filteredMenus = menus.filter((m) =>
    m.name.toLowerCase().includes(search.toLowerCase())
  );

  const totalItems = items.reduce((acc, item) => acc + item.quantity, 0);
  const grandTotal = getGrandTotal();

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

  const userRole = user?.role ?? "";
  const visibleTabs = NAV_TABS.filter((t) => t.roles.includes(userRole));
  const isSuperAdmin = userRole === "superadmin";

  if (!mounted || !isAuthenticated) return null;

  return (
    <div className="flex flex-col h-screen bg-slate-100 text-slate-900 overflow-hidden font-sans">
      {/* ===== TOP HEADER ===== */}
      <header className="flex items-center justify-between px-5 py-3 bg-white border-b border-slate-200 shrink-0 z-20 shadow-sm">
        {/* Brand */}
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 bg-emerald-600 rounded-xl flex items-center justify-center shadow">
            <Coffee className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="text-base font-bold text-slate-800 leading-none">POS KopiTiam</h1>
            <p className="text-[10px] text-slate-400 leading-none mt-0.5">Point of Sale System</p>
          </div>
        </div>

        {/* Nav Tabs (center) */}
        {visibleTabs.length > 1 && (
          <nav className="hidden md:flex items-center gap-1 bg-slate-100 rounded-xl p-1">
            {visibleTabs.map((tab) => {
              const Icon = tab.Icon;
              const isActive = activeView === tab.id ||
                (tab.id === "history" && historyOpen);
              return (
                <button
                  key={tab.id}
                  onClick={() => {
                    if (tab.id === "history") {
                      setHistoryOpen(true);
                    } else {
                      setActiveView(tab.id);
                      setHistoryOpen(false);
                    }
                  }}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
                    isActive
                      ? "bg-white text-emerald-700 shadow-sm"
                      : "text-slate-500 hover:text-slate-700"
                  }`}
                >
                  <Icon className="w-3.5 h-3.5" />
                  {tab.label}
                </button>
              );
            })}
          </nav>
        )}

        {/* Right: search (POS only) + user + logout */}
        <div className="flex items-center gap-3">
          {activeView === "pos" && (
            <div className="hidden md:flex relative w-52">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <input
                type="text"
                placeholder="Cari menu..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full pl-9 pr-4 py-2 text-sm bg-slate-100 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-400/30 focus:border-emerald-400 placeholder:text-slate-400"
              />
            </div>
          )}

          {/* History shortcut for cashier */}
          {!isSuperAdmin && (
            <button
              onClick={() => setHistoryOpen(true)}
              className="flex items-center gap-2 px-3 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-sm font-semibold transition-colors border border-slate-200"
            >
              <History className="w-4 h-4" />
              <span className="hidden sm:inline">Riwayat</span>
            </button>
          )}

          <div className="w-px h-7 bg-slate-200" />

          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-xl bg-gradient-to-br from-emerald-400 to-teal-600 flex items-center justify-center shadow-sm">
              <UserCircle className="w-5 h-5 text-white" />
            </div>
            <div className="hidden sm:block">
              <p className="text-sm font-semibold text-slate-800 leading-none">{user?.name || "Pengguna"}</p>
              <p className="text-[11px] text-emerald-600 font-medium leading-none mt-0.5">{getRoleLabel(userRole)}</p>
            </div>
          </div>

          <button
            onClick={handleLogout}
            className="p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-xl transition-colors"
            title="Keluar"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </header>

      {/* ===== MAIN CONTENT ===== */}
      <main className="flex-1 flex overflow-hidden">

        {/* ---- Non-POS views ---- */}
        {activeView === "employees" && <EmployeePanel />}
        {activeView === "analytics" && <AnalyticsDashboard />}

        {/* ---- POS View ---- */}
        {activeView === "pos" && (
          <>
            {/* LEFT PANEL: MENU GRID */}
            <div className="flex-1 h-full overflow-y-auto p-5">
              {/* Mobile search */}
              <div className="md:hidden mb-4 relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                <input
                  type="text"
                  placeholder="Cari menu..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="w-full pl-9 pr-4 py-2.5 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-400/30"
                />
              </div>

              {/* Daily metrics strip */}
              <DailyMetrics />

              <div className="flex items-center justify-between mb-4">
                <div>
                  <h2 className="text-lg font-bold text-slate-800">Daftar Menu</h2>
                  <p className="text-xs text-slate-500">{filteredMenus.length} item tersedia</p>
                </div>
                {search && (
                  <button onClick={() => setSearch("")} className="text-xs text-emerald-600 hover:underline font-medium">
                    Reset filter
                  </button>
                )}
              </div>

              {isLoadingMenus && (
                <div className="flex flex-col items-center justify-center h-64 gap-3">
                  <Loader2 className="w-8 h-8 text-emerald-600 animate-spin" />
                  <p className="text-sm text-slate-500">Memuat menu...</p>
                </div>
              )}

              {menusError && !isLoadingMenus && (
                <div className="flex flex-col items-center justify-center h-64 gap-3 text-center">
                  <div className="w-14 h-14 bg-red-50 rounded-2xl flex items-center justify-center">
                    <Package className="w-7 h-7 text-red-400" />
                  </div>
                  <p className="text-sm font-semibold text-slate-700">Gagal memuat menu</p>
                  <p className="text-xs text-slate-500">Pastikan server backend aktif dan coba lagi.</p>
                </div>
              )}

              {!isLoadingMenus && !menusError && filteredMenus.length === 0 && (
                <div className="flex flex-col items-center justify-center h-64 gap-3 text-center">
                  <div className="w-14 h-14 bg-slate-100 rounded-2xl flex items-center justify-center">
                    <ChefHat className="w-7 h-7 text-slate-400" />
                  </div>
                  <p className="text-sm font-semibold text-slate-700">Menu tidak ditemukan</p>
                  <p className="text-xs text-slate-500">Coba kata kunci lain atau tambah menu di backoffice.</p>
                </div>
              )}

              {/* Menu Grid — new MenuCard component */}
              {!isLoadingMenus && !menusError && filteredMenus.length > 0 && (
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
                  {filteredMenus.map((menu) => {
                    const cartItem = items.find((i) => i.id === menu.id);
                    return (
                      <MenuCard
                        key={menu.id}
                        menu={menu}
                        cartQty={cartItem?.quantity ?? 0}
                        onAdd={() => addItem(menu)}
                        canEdit={isSuperAdmin}
                      />
                    );
                  })}
                </div>
              )}
            </div>

            {/* RIGHT PANEL: CART */}
            <div className="w-80 xl:w-96 h-full bg-white border-l border-slate-200 flex flex-col shadow-xl z-10 shrink-0">
              {/* Cart header */}
              <div className="px-5 py-4 border-b border-slate-100 flex items-center gap-2 shrink-0 bg-slate-50">
                <div className="w-8 h-8 rounded-xl bg-emerald-100 flex items-center justify-center">
                  <ShoppingCart className="w-4 h-4 text-emerald-600" />
                </div>
                <div>
                  <h2 className="text-sm font-bold text-slate-800 leading-none">Pesanan Saat Ini</h2>
                  <p className="text-[11px] text-slate-500 mt-0.5">{totalItems} item dalam keranjang</p>
                </div>
                {totalItems > 0 && (
                  <span className="ml-auto bg-emerald-500 text-white text-xs font-bold px-2.5 py-1 rounded-full">
                    {totalItems}
                  </span>
                )}
              </div>

              {/* Cart items */}
              <div className="flex-1 overflow-y-auto p-4 space-y-2">
                {items.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-full gap-4 text-center py-10">
                    <div className="w-16 h-16 bg-slate-100 rounded-2xl flex items-center justify-center">
                      <Receipt className="w-8 h-8 text-slate-300" />
                    </div>
                    <div>
                      <p className="text-sm font-semibold text-slate-500">Keranjang kosong</p>
                      <p className="text-xs text-slate-400 mt-1">Klik menu untuk menambahkan item</p>
                    </div>
                  </div>
                ) : (
                  items.map((item) => (
                    <div
                      key={item.cartItemId}
                      className="flex items-center gap-3 bg-white p-3 rounded-xl border border-slate-100 shadow-sm hover:border-slate-200 transition-colors group"
                    >
                      <div className={`w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0 text-xs font-black ${getMenuColor(item.name)}`}>
                        {item.name.substring(0, 2).toUpperCase()}
                      </div>
                      <div className="flex-1 min-w-0">
                        <h4 className="font-semibold text-xs text-slate-800 line-clamp-1">{item.name}</h4>
                        <p className="text-emerald-600 text-xs font-bold mt-0.5">{formatCurrency(item.price)}</p>
                      </div>
                      <div className="flex items-center gap-1 shrink-0">
                        <button
                          onClick={() => updateQuantity(item.cartItemId, item.quantity - 1)}
                          className="w-7 h-7 flex items-center justify-center bg-slate-100 hover:bg-slate-200 rounded-lg transition-colors"
                        >
                          <Minus className="w-3 h-3 text-slate-600" />
                        </button>
                        <span className="w-7 text-center text-sm font-bold text-slate-800">{item.quantity}</span>
                        <button
                          onClick={() => updateQuantity(item.cartItemId, item.quantity + 1)}
                          className="w-7 h-7 flex items-center justify-center bg-emerald-50 hover:bg-emerald-100 rounded-lg transition-colors"
                        >
                          <span className="text-emerald-600 font-bold text-sm">+</span>
                        </button>
                        <button
                          onClick={() => removeItem(item.cartItemId)}
                          className="w-7 h-7 flex items-center justify-center text-slate-300 hover:text-red-400 hover:bg-red-50 rounded-lg transition-colors ml-1 opacity-0 group-hover:opacity-100"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </button>
                      </div>
                    </div>
                  ))
                )}
              </div>

              {/* Cart summary & checkout */}
              <div className="p-4 bg-slate-50 border-t border-slate-200 shrink-0 space-y-3">
                <div className="bg-white rounded-xl border border-slate-100 p-3 space-y-2">
                  <div className="flex justify-between items-center text-xs text-slate-500">
                    <span>Subtotal</span>
                    <span className="font-medium text-slate-700">{formatCurrency(getSubtotal())}</span>
                  </div>
                  <div className="flex justify-between items-center text-xs text-slate-500">
                    <span>Pajak (11%)</span>
                    <span className="font-medium text-slate-700">{formatCurrency(getTaxAmount())}</span>
                  </div>
                  <div className="flex justify-between items-center text-xs text-slate-500">
                    <span>Service (5%)</span>
                    <span className="font-medium text-slate-700">{formatCurrency(getServiceCharge())}</span>
                  </div>
                  <div className="pt-2 mt-1 border-t border-dashed border-slate-200 flex justify-between items-center">
                    <span className="text-sm font-bold text-slate-800">Total</span>
                    <span className="text-xl font-black text-emerald-600">{formatCurrency(grandTotal)}</span>
                  </div>

                  <CheckoutDialog>
                    <button
                      id="checkout-btn"
                      disabled={items.length === 0}
                      className="w-full py-3.5 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-200 disabled:text-slate-400 disabled:cursor-not-allowed text-white font-bold rounded-xl transition-all shadow-md shadow-emerald-200 active:scale-[0.98] flex items-center justify-center gap-2"
                    >
                      <ShoppingCart className="w-4 h-4" />
                      {items.length === 0 ? "Tambah item dulu" : `Checkout — ${formatCurrency(grandTotal)}`}
                    </button>
                  </CheckoutDialog>
                </div>
              </div>
            </div>
          </>
        )}
      </main>

      {/* History Dialog */}
      <HistoryDialog open={historyOpen} onOpenChange={setHistoryOpen} />
    </div>
  );
}
