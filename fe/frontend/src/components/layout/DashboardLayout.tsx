'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '@/store/useAuthStore';
import {
  Coffee, LogOut, User, LayoutDashboard, ChefHat,
  ShoppingCart, BarChart2, Shield, Menu as MenuIcon, X,
} from 'lucide-react';

const NAV_BY_ROLE: Record<string, { href: string; label: string; icon: React.ElementType }[]> = {
  superadmin: [
    { href: '/admin',   label: 'Karyawan',  icon: Shield },
    { href: '/manager', label: 'Analitik',  icon: BarChart2 },
    { href: '/kasir',   label: 'Kasir POS', icon: ShoppingCart },
    { href: '/kitchen', label: 'Dapur',     icon: ChefHat },
  ],
  manager: [
    { href: '/manager', label: 'Analitik', icon: BarChart2 },
  ],
  cashier: [
    { href: '/kasir', label: 'Kasir POS', icon: ShoppingCart },
  ],
  kitchen: [
    { href: '/kitchen', label: 'Dapur', icon: ChefHat },
  ],
};

interface DashboardLayoutProps {
  children: React.ReactNode;
  title: string;
  subtitle?: string;
}

export function DashboardLayout({ children, title, subtitle }: DashboardLayoutProps) {
  const router   = useRouter();
  const pathname = usePathname();
  const { user, isAuthenticated, logout } = useAuthStore();
  const [sideOpen, setSideOpen] = useState(false);

  useEffect(() => {
    if (!isAuthenticated) router.replace('/login');
  }, [isAuthenticated, router]);

  if (!isAuthenticated || !user) return null;

  const nav = NAV_BY_ROLE[user.role] ?? [];

  function handleLogout() {
    logout();
    router.replace('/login');
  }

  return (
    <>
      <style>{`
        @keyframes slide-in {
          from { transform: translateX(-100%); opacity: 0; }
          to   { transform: translateX(0);     opacity: 1; }
        }
        .sidebar-mobile { animation: slide-in 0.22s cubic-bezier(0.34,1.2,0.64,1); }
      `}</style>

      <div
        className="min-h-screen flex"
        style={{ backgroundColor: '#F7F3EE' }}
      >
        {/* ── SIDEBAR desktop ── */}
        <aside
          className="hidden lg:flex flex-col w-64 shrink-0 border-r"
          style={{
            backgroundColor: '#FFFDF8',
            borderColor: '#E8D9C4',
            boxShadow: '4px 0 20px rgba(80,30,0,0.06)',
          }}
        >
          {/* Brand */}
          <div className="flex items-center gap-3 px-6 py-5 border-b" style={{ borderColor: '#E8D9C4' }}>
            <div
              className="w-9 h-9 rounded-xl flex items-center justify-center shrink-0"
              style={{ background: 'linear-gradient(135deg, #C07A24 0%, #8B3A10 100%)' }}
            >
              <Coffee size={18} color="#FFFBF4" />
            </div>
            <div>
              <p className="font-bold text-sm leading-none" style={{ color: '#1A0800' }}>POS Kopitiam</p>
              <p className="text-xs mt-0.5" style={{ color: '#A08060' }}>v1.0</p>
            </div>
          </div>

          {/* Nav */}
          <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
            {nav.map(({ href, label, icon: Icon }) => {
              const active = pathname === href;
              return (
                <a
                  key={href}
                  href={href}
                  className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-150"
                  style={{
                    backgroundColor: active ? '#FEF3E2' : 'transparent',
                    color: active ? '#A0633A' : '#6B4C2A',
                    fontWeight: active ? 700 : 500,
                  }}
                  onMouseEnter={(e) => {
                    if (!active) e.currentTarget.style.backgroundColor = '#FDF6ED';
                  }}
                  onMouseLeave={(e) => {
                    if (!active) e.currentTarget.style.backgroundColor = 'transparent';
                  }}
                >
                  <Icon size={17} strokeWidth={active ? 2.5 : 2} />
                  {label}
                  {active && (
                    <span
                      className="ml-auto w-1.5 h-5 rounded-full"
                      style={{ backgroundColor: '#C07A24' }}
                    />
                  )}
                </a>
              );
            })}
          </nav>

          {/* User + Logout */}
          <div className="p-3 border-t" style={{ borderColor: '#E8D9C4' }}>
            <div
              className="flex items-center gap-3 px-3 py-2.5 rounded-xl mb-1"
              style={{ backgroundColor: '#FEF3E2' }}
            >
              <div
                className="w-8 h-8 rounded-xl flex items-center justify-center text-xs font-bold shrink-0"
                style={{ backgroundColor: '#D4850A', color: '#FFFBF4' }}
              >
                {user.name[0]?.toUpperCase()}
              </div>
              <div className="min-w-0">
                <p className="text-xs font-bold truncate" style={{ color: '#1A0800' }}>{user.name}</p>
                <p className="text-[10px] capitalize" style={{ color: '#A08060' }}>{user.role}</p>
              </div>
            </div>
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-2 px-3 py-2 rounded-xl text-xs font-semibold transition-colors"
              style={{ color: '#9B2226' }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FFF1F1')}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
            >
              <LogOut size={14} />
              Keluar
            </button>
          </div>
        </aside>

        {/* ── SIDEBAR mobile overlay ── */}
        {sideOpen && (
          <div
            className="lg:hidden fixed inset-0 z-50 flex"
            onClick={() => setSideOpen(false)}
          >
            <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" />
            <aside
              className="sidebar-mobile relative w-64 flex flex-col border-r"
              style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4' }}
              onClick={(e) => e.stopPropagation()}
            >
              <div className="flex items-center justify-between px-5 py-4 border-b" style={{ borderColor: '#E8D9C4' }}>
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 rounded-xl flex items-center justify-center" style={{ background: 'linear-gradient(135deg,#C07A24,#8B3A10)' }}>
                    <Coffee size={16} color="#FFFBF4" />
                  </div>
                  <span className="font-bold text-sm" style={{ color: '#1A0800' }}>POS Kopitiam</span>
                </div>
                <button onClick={() => setSideOpen(false)} style={{ color: '#A08060' }}>
                  <X size={18} />
                </button>
              </div>
              <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
                {nav.map(({ href, label, icon: Icon }) => {
                  const active = pathname === href;
                  return (
                    <a
                      key={href}
                      href={href}
                      onClick={() => setSideOpen(false)}
                      className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all"
                      style={{ backgroundColor: active ? '#FEF3E2' : 'transparent', color: active ? '#A0633A' : '#6B4C2A' }}
                    >
                      <Icon size={17} />
                      {label}
                    </a>
                  );
                })}
              </nav>
              <div className="p-3 border-t" style={{ borderColor: '#E8D9C4' }}>
                <button onClick={handleLogout} className="w-full flex items-center gap-2 px-3 py-2 rounded-xl text-xs font-semibold" style={{ color: '#9B2226' }}>
                  <LogOut size={14} />
                  Keluar
                </button>
              </div>
            </aside>
          </div>
        )}

        {/* ── MAIN ── */}
        <div className="flex-1 flex flex-col min-w-0 min-h-screen">
          {/* Topbar */}
          <header
            className="flex items-center gap-4 px-5 py-4 border-b sticky top-0 z-30"
            style={{
              backgroundColor: 'rgba(255,253,248,0.85)',
              borderColor: '#E8D9C4',
              backdropFilter: 'blur(16px)',
              boxShadow: '0 2px 16px rgba(80,30,0,0.05)',
            }}
          >
            <button
              className="lg:hidden w-9 h-9 flex items-center justify-center rounded-xl transition-colors"
              onClick={() => setSideOpen(true)}
              style={{ color: '#6B4C2A' }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FEF3E2')}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
            >
              <MenuIcon size={18} />
            </button>
            <div>
              <h1 className="font-bold text-lg leading-none" style={{ color: '#1A0800' }}>{title}</h1>
              {subtitle && <p className="text-xs mt-0.5" style={{ color: '#A08060' }}>{subtitle}</p>}
            </div>
            <div className="ml-auto flex items-center gap-2">
              <div
                className="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold"
                style={{ backgroundColor: '#FEF3E2', color: '#A0633A' }}
              >
                <User size={13} />
                {user.name}
              </div>
            </div>
          </header>

          {/* Content */}
          <main className="flex-1 p-5 lg:p-6">{children}</main>
        </div>
      </div>
    </>
  );
}
