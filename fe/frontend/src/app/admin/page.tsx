'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api/axios';
import { Employee } from '@/types';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import {
  Plus, Pencil, Trash2, Users, X, Eye, EyeOff,
  ShieldCheck, User, Loader2, Search, AlertCircle,
} from 'lucide-react';

/* ── utils ── */
const ROLE_OPTS = [
  { value: 'cashier',    label: 'Kasir' },
  { value: 'manager',    label: 'Manager' },
  { value: 'kitchen',    label: 'Kitchen' },
  { value: 'superadmin', label: 'Super Admin' },
];

const ROLE_COLOR: Record<string, string> = {
  superadmin: '#7C3AED',
  manager:    '#0369A1',
  cashier:    '#047857',
  kitchen:    '#B45309',
};

const ROLE_BG: Record<string, string> = {
  superadmin: '#F5F3FF',
  manager:    '#E0F2FE',
  cashier:    '#D1FAE5',
  kitchen:    '#FEF3C7',
};

function RoleBadge({ role }: { role: string }) {
  return (
    <span
      className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-bold capitalize"
      style={{ backgroundColor: ROLE_BG[role] ?? '#F3F4F6', color: ROLE_COLOR[role] ?? '#374151' }}
    >
      <ShieldCheck size={11} strokeWidth={2.5} />
      {ROLE_OPTS.find((r) => r.value === role)?.label ?? role}
    </span>
  );
}

/* ── modal ── */
interface ModalProps {
  open: boolean;
  onClose: () => void;
  initial?: Employee | null;
  onSaved: () => void;
}

function EmployeeModal({ open, onClose, initial, onSaved }: ModalProps) {
  const [form, setForm] = useState({
    name:     initial?.name     ?? '',
    email:    initial?.email    ?? '',
    role:     initial?.role     ?? 'cashier',
    pin_code: initial?.pin_code ?? '',
    password: '',
  });
  const [showPass, setShowPass] = useState(false);
  const [error,    setError]    = useState('');
  const [loading,  setLoading]  = useState(false);

  const isEdit = !!initial;

  function set(key: string, val: string) {
    setForm((p) => ({ ...p, [key]: val }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      if (isEdit) {
        await api.put(`/api/v1/employees/${initial!.id}`, form);
      } else {
        await api.post('/api/v1/employees', form);
      }
      onSaved();
      onClose();
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message ??
        'Terjadi kesalahan';
      setError(msg);
    } finally {
      setLoading(false);
    }
  }

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div
        className="relative w-full max-w-md rounded-3xl border p-6"
        style={{
          backgroundColor: '#FFFDF8',
          borderColor: '#E8D9C4',
          boxShadow: '0 40px 80px rgba(30,8,0,0.18)',
        }}
      >
        <div className="flex items-center justify-between mb-5">
          <h2 className="font-bold text-lg" style={{ color: '#1A0800' }}>
            {isEdit ? 'Edit Karyawan' : 'Tambah Karyawan Baru'}
          </h2>
          <button
            onClick={onClose}
            className="w-8 h-8 flex items-center justify-center rounded-xl transition-colors"
            style={{ color: '#A08060' }}
            onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FEF3E2')}
            onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
          >
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Name */}
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider mb-1.5" style={{ color: '#8B6347' }}>
              Nama Lengkap
            </label>
            <input
              className="w-full px-4 py-3 rounded-xl border text-sm outline-none transition-all"
              style={{ borderColor: '#E8D9C4', backgroundColor: '#FDFAF6', color: '#1A0800' }}
              value={form.name}
              onChange={(e) => set('name', e.target.value)}
              required
              placeholder="Nama karyawan"
              onFocus={(e) => (e.currentTarget.style.borderColor = '#C07A24')}
              onBlur={(e) => (e.currentTarget.style.borderColor = '#E8D9C4')}
            />
          </div>

          {/* Email */}
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider mb-1.5" style={{ color: '#8B6347' }}>
              Email
            </label>
            <input
              type="email"
              className="w-full px-4 py-3 rounded-xl border text-sm outline-none transition-all"
              style={{ borderColor: '#E8D9C4', backgroundColor: '#FDFAF6', color: '#1A0800' }}
              value={form.email}
              onChange={(e) => set('email', e.target.value)}
              required
              placeholder="email@kopitiam.com"
              onFocus={(e) => (e.currentTarget.style.borderColor = '#C07A24')}
              onBlur={(e) => (e.currentTarget.style.borderColor = '#E8D9C4')}
            />
          </div>

          {/* Role */}
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider mb-1.5" style={{ color: '#8B6347' }}>
              Peran
            </label>
            <select
              className="w-full px-4 py-3 rounded-xl border text-sm outline-none transition-all appearance-none cursor-pointer"
              style={{ borderColor: '#E8D9C4', backgroundColor: '#FDFAF6', color: '#1A0800' }}
              value={form.role}
              onChange={(e) => set('role', e.target.value)}
              onFocus={(e) => (e.currentTarget.style.borderColor = '#C07A24')}
              onBlur={(e) => (e.currentTarget.style.borderColor = '#E8D9C4')}
            >
              {ROLE_OPTS.map((r) => (
                <option key={r.value} value={r.value}>{r.label}</option>
              ))}
            </select>
          </div>

          {/* PIN */}
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider mb-1.5" style={{ color: '#8B6347' }}>
              PIN Code (6 Digit)
            </label>
            <input
              maxLength={6}
              pattern="\d{6}"
              className="w-full px-4 py-3 rounded-xl border text-sm outline-none transition-all"
              style={{ borderColor: '#E8D9C4', backgroundColor: '#FDFAF6', color: '#1A0800' }}
              value={form.pin_code}
              onChange={(e) => set('pin_code', e.target.value.replace(/\D/g, ''))}
              placeholder="123456"
              onFocus={(e) => (e.currentTarget.style.borderColor = '#C07A24')}
              onBlur={(e) => (e.currentTarget.style.borderColor = '#E8D9C4')}
            />
          </div>

          {/* Password */}
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider mb-1.5" style={{ color: '#8B6347' }}>
              {isEdit ? 'Password Baru (opsional)' : 'Password'}
            </label>
            <div className="relative">
              <input
                type={showPass ? 'text' : 'password'}
                className="w-full px-4 py-3 pr-12 rounded-xl border text-sm outline-none transition-all"
                style={{ borderColor: '#E8D9C4', backgroundColor: '#FDFAF6', color: '#1A0800' }}
                value={form.password}
                onChange={(e) => set('password', e.target.value)}
                required={!isEdit}
                placeholder={isEdit ? 'Kosongkan jika tidak diubah' : '••••••••'}
                onFocus={(e) => (e.currentTarget.style.borderColor = '#C07A24')}
                onBlur={(e) => (e.currentTarget.style.borderColor = '#E8D9C4')}
              />
              <button
                type="button"
                onClick={() => setShowPass((v) => !v)}
                className="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded-lg"
                style={{ color: '#A08060' }}
              >
                {showPass ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
          </div>

          {error && (
            <div
              className="flex items-center gap-2 px-4 py-3 rounded-xl text-sm"
              style={{ backgroundColor: 'rgba(155,34,38,0.07)', color: '#9B2226', border: '1px solid rgba(155,34,38,0.2)' }}
            >
              <AlertCircle size={15} />
              {error}
            </div>
          )}

          <div className="flex gap-3 pt-1">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 py-3 rounded-xl text-sm font-semibold transition-all"
              style={{ backgroundColor: '#F5EDD6', color: '#6B4C2A' }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#E8D9C4')}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = '#F5EDD6')}
            >
              Batal
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 py-3 rounded-xl text-sm font-bold text-white transition-all flex items-center justify-center gap-2 disabled:opacity-70"
              style={{ background: 'linear-gradient(135deg, #C07A24, #8B3A10)' }}
            >
              {loading ? <Loader2 size={15} className="animate-spin" /> : null}
              {isEdit ? 'Simpan' : 'Tambahkan'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

/* ── delete confirm ── */
function DeleteConfirmModal({
  open, onClose, employee, onDeleted,
}: { open: boolean; onClose: () => void; employee: Employee | null; onDeleted: () => void }) {
  const [loading, setLoading] = useState(false);

  async function handleDelete() {
    if (!employee) return;
    setLoading(true);
    try {
      await api.delete(`/api/v1/employees/${employee.id}`);
      onDeleted();
      onClose();
    } finally {
      setLoading(false);
    }
  }

  if (!open || !employee) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div
        className="relative w-full max-w-sm rounded-3xl border p-6"
        style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4', boxShadow: '0 40px 80px rgba(30,8,0,0.18)' }}
      >
        <div className="flex flex-col items-center text-center gap-3 mb-6">
          <div className="w-14 h-14 rounded-2xl flex items-center justify-center" style={{ backgroundColor: '#FEE2E2' }}>
            <Trash2 size={24} style={{ color: '#DC2626' }} />
          </div>
          <div>
            <h3 className="font-bold text-base" style={{ color: '#1A0800' }}>Hapus Karyawan?</h3>
            <p className="text-sm mt-1" style={{ color: '#6B4C2A' }}>
              <strong>{employee.name}</strong> akan dihapus dari sistem secara permanen.
            </p>
          </div>
        </div>
        <div className="flex gap-3">
          <button onClick={onClose} className="flex-1 py-3 rounded-xl text-sm font-semibold" style={{ backgroundColor: '#F5EDD6', color: '#6B4C2A' }}>Batal</button>
          <button
            onClick={handleDelete}
            disabled={loading}
            className="flex-1 py-3 rounded-xl text-sm font-bold text-white flex items-center justify-center gap-2"
            style={{ backgroundColor: '#DC2626' }}
          >
            {loading ? <Loader2 size={14} className="animate-spin" /> : <Trash2 size={14} />}
            Hapus
          </button>
        </div>
      </div>
    </div>
  );
}

/* ── main page ── */
export default function AdminPage() {
  const qc = useQueryClient();
  const [search, setSearch] = useState('');
  const [addOpen,    setAddOpen]    = useState(false);
  const [editTarget, setEditTarget] = useState<Employee | null>(null);
  const [delTarget,  setDelTarget]  = useState<Employee | null>(null);

  const { data: employees = [], isLoading, error } = useQuery<Employee[]>({
    queryKey: ['employees'],
    queryFn: async () => {
      const r = await api.get('/api/v1/employees');
      return r.data?.data ?? r.data;
    },
  });

  const filtered = employees.filter(
    (e) =>
      e.name.toLowerCase().includes(search.toLowerCase()) ||
      e.email.toLowerCase().includes(search.toLowerCase()) ||
      e.role.toLowerCase().includes(search.toLowerCase()),
  );

  function refetch() {
    qc.invalidateQueries({ queryKey: ['employees'] });
  }

  return (
    <DashboardLayout title="Manajemen Karyawan" subtitle="Kelola akun dan akses pengguna sistem">
      <EmployeeModal
        open={addOpen}
        onClose={() => setAddOpen(false)}
        onSaved={refetch}
      />
      <EmployeeModal
        open={!!editTarget}
        onClose={() => setEditTarget(null)}
        initial={editTarget}
        onSaved={refetch}
      />
      <DeleteConfirmModal
        open={!!delTarget}
        onClose={() => setDelTarget(null)}
        employee={delTarget}
        onDeleted={refetch}
      />

      {/* Header bar */}
      <div className="flex flex-col sm:flex-row gap-3 mb-5">
        {/* Search */}
        <div
          className="flex items-center gap-2 flex-1 px-4 py-2.5 rounded-2xl border"
          style={{ backgroundColor: '#FFFDF8', borderColor: '#E8D9C4' }}
        >
          <Search size={16} style={{ color: '#A08060' }} />
          <input
            className="flex-1 bg-transparent text-sm outline-none"
            style={{ color: '#1A0800' }}
            placeholder="Cari karyawan…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        {/* Add button */}
        <button
          id="btn-add-employee"
          onClick={() => setAddOpen(true)}
          className="flex items-center gap-2 px-5 py-2.5 rounded-2xl text-sm font-bold text-white transition-all"
          style={{
            background: 'linear-gradient(135deg, #C07A24, #8B3A10)',
            boxShadow: '0 4px 12px rgba(160,82,10,0.35)',
          }}
          onMouseEnter={(e) => (e.currentTarget.style.transform = 'translateY(-1px)')}
          onMouseLeave={(e) => (e.currentTarget.style.transform = '')}
        >
          <Plus size={16} />
          Tambah Karyawan
        </button>
      </div>

      {/* Stats mini-cards */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-5">
        {ROLE_OPTS.map((r) => {
          const count = employees.filter((e) => e.role === r.value).length;
          return (
            <div
              key={r.value}
              className="rounded-2xl border p-4"
              style={{
                backgroundColor: '#FFFDF8',
                borderColor: '#E8D9C4',
                boxShadow: '0 2px 12px rgba(80,30,0,0.06)',
              }}
            >
              <p className="text-xs font-semibold" style={{ color: '#A08060' }}>{r.label}</p>
              <p className="text-2xl font-black mt-1" style={{ color: ROLE_COLOR[r.value] }}>{count}</p>
            </div>
          );
        })}
      </div>

      {/* Table Card */}
      <div
        className="rounded-3xl border overflow-hidden"
        style={{
          backgroundColor: '#FFFDF8',
          borderColor: '#E8D9C4',
          boxShadow: '0 4px 24px rgba(80,30,0,0.07)',
        }}
      >
        {isLoading ? (
          <div className="flex items-center justify-center py-16 gap-3">
            <Loader2 size={24} className="animate-spin" style={{ color: '#C07A24' }} />
            <span className="text-sm" style={{ color: '#A08060' }}>Memuat data…</span>
          </div>
        ) : error ? (
          <div className="flex flex-col items-center py-16 gap-2">
            <AlertCircle size={32} style={{ color: '#DC2626' }} />
            <p className="text-sm font-semibold" style={{ color: '#9B2226' }}>Gagal memuat data karyawan</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr style={{ borderBottom: '1px solid #E8D9C4', backgroundColor: '#FEF9F0' }}>
                <th className="text-left px-5 py-3.5 text-xs font-bold uppercase tracking-wider" style={{ color: '#A08060' }}>Nama</th>
                <th className="text-left px-5 py-3.5 text-xs font-bold uppercase tracking-wider hidden sm:table-cell" style={{ color: '#A08060' }}>Email</th>
                <th className="text-left px-5 py-3.5 text-xs font-bold uppercase tracking-wider" style={{ color: '#A08060' }}>Peran</th>
                <th className="text-left px-5 py-3.5 text-xs font-bold uppercase tracking-wider hidden md:table-cell" style={{ color: '#A08060' }}>PIN</th>
                <th className="text-center px-5 py-3.5 text-xs font-bold uppercase tracking-wider" style={{ color: '#A08060' }}>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {filtered.length === 0 ? (
                <tr>
                  <td colSpan={5} className="text-center py-12">
                    <Users size={32} style={{ color: '#D4B896', margin: '0 auto 8px' }} />
                    <p className="text-sm" style={{ color: '#A08060' }}>Tidak ada karyawan ditemukan</p>
                  </td>
                </tr>
              ) : (
                filtered.map((emp, idx) => (
                  <tr
                    key={emp.id}
                    style={{ borderTop: idx > 0 ? '1px solid #F0E4D0' : 'none' }}
                    className="group transition-colors"
                    onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FFFBF5')}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
                  >
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-3">
                        <div
                          className="w-9 h-9 rounded-xl flex items-center justify-center text-sm font-black shrink-0"
                          style={{ backgroundColor: ROLE_BG[emp.role] ?? '#F5EDD6', color: ROLE_COLOR[emp.role] ?? '#6B4C2A' }}
                        >
                          {emp.name[0]?.toUpperCase()}
                        </div>
                        <span className="font-semibold" style={{ color: '#1A0800' }}>{emp.name}</span>
                      </div>
                    </td>
                    <td className="px-5 py-4 hidden sm:table-cell" style={{ color: '#6B4C2A' }}>{emp.email}</td>
                    <td className="px-5 py-4">
                      <RoleBadge role={emp.role} />
                    </td>
                    <td className="px-5 py-4 hidden md:table-cell">
                      <code className="text-xs px-2 py-1 rounded-lg font-mono" style={{ backgroundColor: '#F5EDD6', color: '#8B3A10' }}>
                        {emp.pin_code ?? '—'}
                      </code>
                    </td>
                    <td className="px-5 py-4">
                      <div className="flex items-center justify-center gap-2">
                        <button
                          onClick={() => setEditTarget(emp)}
                          className="w-8 h-8 rounded-xl flex items-center justify-center transition-colors"
                          title="Edit"
                          style={{ color: '#0369A1' }}
                          onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#E0F2FE')}
                          onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
                        >
                          <Pencil size={15} />
                        </button>
                        <button
                          onClick={() => setDelTarget(emp)}
                          className="w-8 h-8 rounded-xl flex items-center justify-center transition-colors"
                          title="Hapus"
                          style={{ color: '#DC2626' }}
                          onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FEE2E2')}
                          onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
                        >
                          <Trash2 size={15} />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        )}

        {/* Footer */}
        <div className="px-5 py-3 border-t flex items-center gap-2" style={{ borderColor: '#E8D9C4', backgroundColor: '#FEF9F0' }}>
          <Users size={14} style={{ color: '#A08060' }} />
          <span className="text-xs font-medium" style={{ color: '#A08060' }}>
            {filtered.length} dari {employees.length} karyawan
          </span>
        </div>
      </div>
    </DashboardLayout>
  );
}
