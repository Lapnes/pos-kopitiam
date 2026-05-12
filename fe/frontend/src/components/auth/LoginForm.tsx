'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import {
  Eye, EyeOff, Mail, Lock, ChevronDown, Check, Coffee, Loader2,
} from 'lucide-react';
import api from '@/lib/api/axios';
import { useAuthStore } from '@/store/useAuthStore';
import type { User as UserType } from '@/types';

/* ─────────────────────────── constants ──────────────────────────── */
const ROLE_REDIRECT: Record<string, string> = {
  cashier:    '/kasir',
  manager:    '/manager',
  superadmin: '/admin',
  kitchen:    '/kitchen',
};

/* ──────────────────── custom select ───────────────────────────── */
const ROLES = [
  { value: 'cashier',    label: 'Kasir',              emoji: '🧾' },
  { value: 'manager',    label: 'Manager',             emoji: '📊' },
  { value: 'superadmin', label: 'Admin / Super Admin', emoji: '🔐' },
  { value: 'kitchen',    label: 'Kitchen',             emoji: '🍳' },
];

function RoleSelect({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const selected = ROLES.find((r) => r.value === value) ?? ROLES[0];

  useEffect(() => {
    function onOut(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener('mousedown', onOut);
    return () => document.removeEventListener('mousedown', onOut);
  }, []);

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        id="role-select-trigger"
        onClick={() => setOpen((o) => !o)}
        aria-haspopup="listbox"
        aria-expanded={open}
        className="w-full flex items-center gap-3 px-4 py-3.5 rounded-2xl border text-sm font-medium outline-none transition-all duration-200 cursor-pointer"
        style={{
          borderColor: open ? '#C07A24' : '#E8D9C4',
          backgroundColor: '#FDFAF6',
          color: '#1E0E00',
          boxShadow: open
            ? '0 0 0 3px rgba(192,122,36,0.15), 0 2px 8px rgba(60,20,0,0.08)'
            : '0 1px 4px rgba(60,20,0,0.06)',
        }}
      >
        <span className="flex items-center justify-center w-7 h-7 rounded-xl text-base shrink-0" style={{ backgroundColor: '#FEF3E2' }}>
          {selected.emoji}
        </span>
        <span className="flex-1 text-left">{selected.label}</span>
        <ChevronDown size={16} strokeWidth={2.5} className="shrink-0 transition-transform duration-200"
          style={{ color: '#A0633A', transform: open ? 'rotate(180deg)' : 'rotate(0deg)' }} />
      </button>

      <div role="listbox" className="absolute z-50 left-0 right-0 py-2 rounded-2xl border overflow-hidden"
        style={{
          top: 'calc(100% + 8px)',
          backgroundColor: '#FFFDF8',
          borderColor: '#E8D9C4',
          boxShadow: '0 20px 60px rgba(30,8,0,0.12), 0 8px 24px rgba(30,8,0,0.08)',
          transformOrigin: 'top center',
          transform: open ? 'scaleY(1) translateY(0)' : 'scaleY(0.9) translateY(-8px)',
          opacity: open ? 1 : 0,
          pointerEvents: open ? 'all' : 'none',
          transition: 'transform 0.18s cubic-bezier(0.34,1.56,0.64,1), opacity 0.15s ease',
        }}
      >
        {ROLES.map((r) => {
          const active = r.value === value;
          return (
            <button key={r.value} type="button" role="option" aria-selected={active}
              onClick={() => { onChange(r.value); setOpen(false); }}
              className="w-full flex items-center gap-3 px-4 py-3 text-sm font-medium transition-colors duration-150 cursor-pointer"
              style={{ backgroundColor: active ? '#FEF3E2' : 'transparent', color: active ? '#A0633A' : '#3D1F00' }}
              onMouseEnter={(e) => { if (!active) e.currentTarget.style.backgroundColor = '#FDF6ED'; }}
              onMouseLeave={(e) => { if (!active) e.currentTarget.style.backgroundColor = 'transparent'; }}
            >
              <span className="flex items-center justify-center w-7 h-7 rounded-xl text-base shrink-0"
                style={{ backgroundColor: active ? '#FDDDA6' : '#F5EDD6' }}>
                {r.emoji}
              </span>
              <span className="flex-1 text-left">{r.label}</span>
              {active && <Check size={15} strokeWidth={2.5} style={{ color: '#A0633A' }} />}
            </button>
          );
        })}
      </div>
    </div>
  );
}

/* ──────────────────── floating input ──────────────────────────── */
function FloatingInput({ id, label, type = 'text', value, onChange, placeholder, leftIcon, rightSlot, autoComplete }: {
  id: string; label: string; type?: string; value: string;
  onChange: (v: string) => void; placeholder?: string;
  leftIcon: React.ReactNode; rightSlot?: React.ReactNode;
  autoComplete?: string;
}) {
  const [focused, setFocused] = useState(false);
  return (
    <div>
      <label htmlFor={id} className="block text-xs font-semibold uppercase tracking-wider mb-2 transition-colors duration-200"
        style={{ color: focused ? '#A0633A' : '#8B6347' }}>
        {label}
      </label>
      <div className="relative flex items-center rounded-2xl border transition-all duration-200"
        style={{
          borderColor: focused ? '#C07A24' : '#E8D9C4',
          backgroundColor: '#FDFAF6',
          boxShadow: focused
            ? '0 0 0 3px rgba(192,122,36,0.15), 0 2px 8px rgba(60,20,0,0.08)'
            : '0 1px 4px rgba(60,20,0,0.06)',
        }}
      >
        <span className="flex items-center justify-center pl-4 shrink-0 transition-colors duration-200"
          style={{ color: focused ? '#A0633A' : '#C4A07A' }}>
          {leftIcon}
        </span>
        <input id={id} type={type} autoComplete={autoComplete ?? 'off'} required value={value}
          onChange={(e) => onChange(e.target.value)} placeholder={placeholder}
          onFocus={() => setFocused(true)} onBlur={() => setFocused(false)}
          className="flex-1 px-3 py-3.5 bg-transparent text-sm outline-none"
          style={{ color: '#1E0E00', caretColor: '#A0633A' }} />
        {rightSlot && <span className="flex items-center pr-2 shrink-0">{rightSlot}</span>}
      </div>
    </div>
  );
}

/* ────────────────────────── main component ────────────────────── */
export default function LoginForm() {
  const router  = useRouter();
  const setAuth = useAuthStore((s) => s.setAuth);
  const { isAuthenticated, user } = useAuthStore();

  // Redirect if already logged in
  useEffect(() => {
    if (isAuthenticated && user) {
      router.replace(ROLE_REDIRECT[user.role] ?? '/');
    }
  }, [isAuthenticated, user, router]);

  const [email,        setEmail]        = useState('');
  const [password,     setPassword]     = useState('');
  const [role,         setRole]         = useState('cashier');  // UI hint only
  const [showPassword, setShowPassword] = useState(false);
  const [loading,      setLoading]      = useState(false);
  const [error,        setError]        = useState('');
  const [shake,        setShake]        = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!email.trim() || !password) { setError('Email dan password tidak boleh kosong.'); return; }
    setLoading(true); setError('');

    try {
      // Backend expects: { email, password }
      const res = await api.post('/api/v1/auth/login', { email: email.trim(), password });
      const data = res.data?.data ?? res.data;

      const token: string = data.access_token ?? data.token ?? '';
      const u = data.user ?? {};

      const userObj: UserType = {
        id:    u.id    ?? '',
        name:  u.name  ?? email,
        email: u.email ?? email,
        role:  u.role  ?? role,
      };

      setAuth(userObj, token);
      router.push(ROLE_REDIRECT[userObj.role] ?? '/');
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message ??
        'Login gagal. Periksa email dan password Anda.';
      setError(msg);
      setShake(true);
      setTimeout(() => setShake(false), 520);
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <style>{`
        @keyframes ag-fade-up {
          from { opacity:0; transform:translateY(28px) scale(0.97); }
          to   { opacity:1; transform:translateY(0) scale(1); }
        }
        @keyframes ag-shake {
          0%,100% { transform:translateX(0); }
          15% { transform:translateX(-9px); }
          45% { transform:translateX(9px); }
          75% { transform:translateX(-5px); }
        }
        @keyframes ring-pulse {
          0%,100% { box-shadow: 0 0 0 0   rgba(192,122,36,0.40); }
          50%     { box-shadow: 0 0 0 14px rgba(192,122,36,0); }
        }
        .ag-card  { animation: ag-fade-up 0.55s cubic-bezier(0.34,1.2,0.64,1) both; }
        .ag-shake { animation: ag-shake 0.50s ease; }
        .coffee-ring { animation: ring-pulse 2.6s ease infinite; }
      `}</style>

      <div className="min-h-screen flex flex-col items-center justify-center px-4 py-12 relative overflow-hidden"
        style={{
          background: [
            'radial-gradient(ellipse 80% 60% at 20% 10%, rgba(255,230,180,0.38) 0%, transparent 60%)',
            'radial-gradient(ellipse 60% 50% at 80% 80%, rgba(220,160,80,0.18) 0%, transparent 55%)',
            'radial-gradient(ellipse 70% 80% at 50% 50%, rgba(255,248,235,0.9) 0%, transparent 100%)',
            '#F8F2E8',
          ].join(', '),
        }}
      >
        {/* blobs */}
        <div aria-hidden="true" className="pointer-events-none absolute -top-24 -left-28 w-96 h-96 rounded-full opacity-20"
          style={{ background: 'radial-gradient(circle, #D4850A 0%, transparent 70%)' }} />
        <div aria-hidden="true" className="pointer-events-none absolute -bottom-16 -right-20 w-72 h-72 rounded-full opacity-15"
          style={{ background: 'radial-gradient(circle, #A0522D 0%, transparent 70%)' }} />

        <main className={`ag-card relative z-10 w-full max-w-[440px] ${shake ? 'ag-shake' : ''}`}>
          <div className="rounded-3xl border px-8 py-10 sm:px-10"
            style={{
              backgroundColor: 'rgba(255,253,248,0.85)',
              borderColor: 'rgba(232,217,196,0.70)',
              backdropFilter: 'blur(20px) saturate(1.4)',
              WebkitBackdropFilter: 'blur(20px) saturate(1.4)',
              boxShadow: [
                '0 4px 6px rgba(30,8,0,0.03)',
                '0 12px 24px rgba(60,20,0,0.07)',
                '0 40px 80px rgba(80,30,0,0.10)',
                'inset 0 1px 0 rgba(255,255,255,0.80)',
              ].join(', '),
            }}
          >
            {/* header */}
            <div className="flex flex-col items-center text-center mb-8">
              <div className="coffee-ring flex items-center justify-center w-16 h-16 rounded-2xl mb-5 shrink-0"
                style={{
                  background: 'linear-gradient(145deg, #D4850A 0%, #A0522D 100%)',
                  boxShadow: '0 8px 20px rgba(164,82,10,0.35), inset 0 1px 0 rgba(255,255,255,0.25)',
                }}>
                <Coffee size={30} color="#FFFBF4" strokeWidth={1.8} />
              </div>
              <h1 className="text-2xl sm:text-3xl font-bold tracking-tight leading-tight mb-1.5" style={{ color: '#1A0800' }}>
                Selamat Datang Kembali
              </h1>
              <p className="text-sm sm:text-base" style={{ color: '#8B7355' }}>
                Masuk sesuai peran Anda untuk melanjutkan
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-5" noValidate>
              {/* email — field yang benar sesuai backend */}
              <FloatingInput id="login-email" label="Email" type="email" value={email}
                onChange={setEmail} placeholder="admin@kopitiam.com"
                autoComplete="email"
                leftIcon={<Mail size={17} strokeWidth={2} />} />

              {/* password */}
              <FloatingInput id="login-password" label="Password" type={showPassword ? 'text' : 'password'}
                value={password} onChange={setPassword} placeholder="••••••••"
                autoComplete="current-password"
                leftIcon={<Lock size={17} strokeWidth={2} />}
                rightSlot={
                  <button type="button" id="btn-toggle-password"
                    onClick={() => setShowPassword((v) => !v)}
                    className="flex items-center justify-center w-9 h-9 rounded-xl transition-colors duration-150"
                    style={{ color: '#A0633A' }}
                    aria-label={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
                    onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FEF3E2')}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}>
                    {showPassword ? <EyeOff size={17} strokeWidth={2} /> : <Eye size={17} strokeWidth={2} />}
                  </button>
                } />

              {/* role selector (hint saja — backend tentukan dari token) */}
              <div>
                <span className="block text-xs font-semibold uppercase tracking-wider mb-2" style={{ color: '#8B6347' }}>
                  Masuk Sebagai
                </span>
                <RoleSelect value={role} onChange={setRole} />
              </div>

              {error && (
                <div id="login-error-banner" className="flex items-start gap-2.5 px-4 py-3 rounded-2xl text-sm" role="alert"
                  style={{ backgroundColor: 'rgba(180,30,30,0.07)', color: '#9B2226', border: '1px solid rgba(155,34,38,0.18)' }}>
                  <span className="mt-px text-base leading-none">⚠</span>
                  <span>{error}</span>
                </div>
              )}

              <button id="btn-login-submit" type="submit" disabled={loading}
                className="relative w-full flex items-center justify-center gap-2.5 py-3.5 rounded-2xl text-sm font-semibold tracking-wide transition-all duration-200 disabled:opacity-70 disabled:cursor-not-allowed overflow-hidden"
                style={{
                  background: loading ? '#A0633A' : 'linear-gradient(135deg, #C07A24 0%, #8B3A10 100%)',
                  color: '#FFFBF4',
                  boxShadow: loading ? 'none' : '0 4px 12px rgba(160,82,10,0.35), 0 2px 4px rgba(60,20,0,0.20)',
                }}
                onMouseEnter={(e) => {
                  if (loading) return;
                  e.currentTarget.style.transform = 'translateY(-2px)';
                  e.currentTarget.style.boxShadow = '0 8px 20px rgba(160,82,10,0.45), 0 4px 8px rgba(60,20,0,0.25)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = '';
                  e.currentTarget.style.boxShadow = '0 4px 12px rgba(160,82,10,0.35), 0 2px 4px rgba(60,20,0,0.20)';
                }}
              >
                <span aria-hidden="true" className="absolute inset-0 rounded-2xl pointer-events-none"
                  style={{ background: 'linear-gradient(105deg, transparent 40%, rgba(255,255,255,0.18) 50%, transparent 60%)' }} />
                {loading ? <><Loader2 size={17} className="animate-spin" />Memproses…</> : <><Coffee size={17} strokeWidth={2} />Masuk</>}
              </button>
            </form>
          </div>
        </main>

        <footer className="mt-8 text-center text-xs select-none" style={{ color: 'rgba(120,80,40,0.50)' }}>
          POS Kopitiam &copy; {new Date().getFullYear()} — Semua hak dilindungi
        </footer>
      </div>
    </>
  );
}
