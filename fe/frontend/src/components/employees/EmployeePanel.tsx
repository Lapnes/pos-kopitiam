"use client";

import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api/axios";
import { Employee } from "@/types";
import { AxiosError } from "axios";
import { toast } from "sonner";
import {
  Users, Pencil, Trash2, Loader2, UserPlus,
  ShieldCheck, ChefHat, Wallet, UserCog, AlertCircle,
} from "lucide-react";
import { EditEmployeeDialog } from "./EditEmployeeDialog";

const ROLE_CONFIG: Record<string, { label: string; cls: string; Icon: React.ElementType }> = {
  superadmin: { label: "Super Admin", cls: "bg-violet-100 text-violet-700 border-violet-200", Icon: ShieldCheck },
  manager: { label: "Manajer", cls: "bg-sky-100 text-sky-700 border-sky-200", Icon: UserCog },
  cashier: { label: "Kasir", cls: "bg-emerald-100 text-emerald-700 border-emerald-200", Icon: Wallet },
  kitchen: { label: "Dapur", cls: "bg-amber-100 text-amber-700 border-amber-200", Icon: ChefHat },
  admin: { label: "Admin", cls: "bg-slate-100 text-slate-600 border-slate-200", Icon: Users },
};

function RoleBadge({ role }: { role: string }) {
  const cfg = ROLE_CONFIG[role] ?? ROLE_CONFIG.admin;
  const Icon = cfg.Icon;
  return (
    <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-bold border ${cfg.cls}`}>
      <Icon className="w-3 h-3" />
      {cfg.label}
    </span>
  );
}

export function EmployeePanel() {
  const qc = useQueryClient();
  const [editTarget, setEditTarget] = useState<Employee | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Employee | null>(null);

  const { data: employees, isLoading, error } = useQuery<Employee[]>({
    queryKey: ["employees"],
    queryFn: async () => {
      const res = await api.get("/employees");
      return (res.data.data ?? res.data) as Employee[];
    },
    retry: 1,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.delete(`/employees/${id}`),
    onSuccess: () => {
      toast.success("Karyawan berhasil dihapus");
      qc.invalidateQueries({ queryKey: ["employees"] });
      setDeleteTarget(null);
    },
    onError: (err: AxiosError<{ message: string }>) => {
      toast.error("Gagal menghapus karyawan", {
        description: err.response?.data?.message ?? "Terjadi error saat menghapus.",
      });
    },
  });

  return (
    <div className="flex-1 h-full overflow-y-auto p-5">
      {/* Header */}
      <div className="flex items-center justify-between mb-5">
        <div>
          <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
            <Users className="w-5 h-5 text-violet-500" />
            Panel Karyawan
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">Kelola akun karyawan dan hak akses</p>
        </div>
        <button
          disabled
          className="flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-700 disabled:opacity-50 text-white text-sm font-semibold rounded-xl shadow-sm transition-colors"
          title="Fitur tambah karyawan dalam pengembangan"
        >
          <UserPlus className="w-4 h-4" />
          Tambah Karyawan
        </button>
      </div>

      {/* Loading */}
      {isLoading && (
        <div className="flex flex-col items-center justify-center h-64 gap-3">
          <Loader2 className="w-8 h-8 text-violet-500 animate-spin" />
          <p className="text-sm text-slate-500">Memuat data karyawan...</p>
        </div>
      )}

      {/* Error / API not implemented */}
      {!isLoading && error && (
        <div className="flex flex-col items-center justify-center h-64 gap-4 text-center">
          <div className="w-16 h-16 bg-amber-50 rounded-2xl flex items-center justify-center">
            <AlertCircle className="w-8 h-8 text-amber-400" />
          </div>
          <div>
            <p className="text-sm font-bold text-slate-700">Endpoint karyawan belum tersedia</p>
            <p className="text-xs text-slate-400 mt-1 max-w-xs">
              Backend belum mengimplementasi route <code className="bg-slate-100 px-1 rounded">/api/v1/employees</code>.
              Data di bawah menggunakan data benih (seed) untuk tampilan demo.
            </p>
          </div>
          {/* Fallback: show seeded data */}
          <EmployeeFallbackTable
            onEdit={setEditTarget}
            onDelete={setDeleteTarget}
          />
        </div>
      )}

      {/* Real data table */}
      {!isLoading && !error && employees && (
        <EmployeeTable
          employees={employees}
          onEdit={setEditTarget}
          onDelete={setDeleteTarget}
        />
      )}

      {/* Edit Dialog */}
      {editTarget && (
        <EditEmployeeDialog
          employee={editTarget}
          open={!!editTarget}
          onOpenChange={(open) => { if (!open) setEditTarget(null); }}
          onSuccess={() => {
            setEditTarget(null);
            qc.invalidateQueries({ queryKey: ["employees"] });
          }}
        />
      )}

      {/* Delete Confirm Modal */}
      {deleteTarget && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-3xl shadow-2xl p-6 max-w-sm w-full mx-4">
            <div className="w-14 h-14 bg-red-50 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <Trash2 className="w-7 h-7 text-red-400" />
            </div>
            <h3 className="text-base font-bold text-slate-800 text-center">Hapus Karyawan?</h3>
            <p className="text-sm text-slate-500 text-center mt-2">
              Yakin ingin menghapus <strong>{deleteTarget.name}</strong>? Tindakan ini tidak bisa dibatalkan.
            </p>
            <div className="flex gap-3 mt-5">
              <button
                onClick={() => setDeleteTarget(null)}
                disabled={deleteMutation.isPending}
                className="flex-1 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl font-semibold text-sm transition-colors"
              >
                Batal
              </button>
              <button
                onClick={() => deleteMutation.mutate(deleteTarget.id)}
                disabled={deleteMutation.isPending}
                className="flex-1 py-2.5 bg-red-600 hover:bg-red-700 text-white rounded-xl font-bold text-sm transition-colors flex items-center justify-center gap-2"
              >
                {deleteMutation.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : <Trash2 className="w-4 h-4" />}
                Hapus
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// ---- Shared Table UI ----
function EmployeeTable({
  employees,
  onEdit,
  onDelete,
}: {
  employees: Employee[];
  onEdit: (e: Employee) => void;
  onDelete: (e: Employee) => void;
}) {
  return (
    <div className="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
      <table className="w-full text-sm">
        <thead className="bg-slate-50 border-b border-slate-100">
          <tr>
            {["Nama", "Email", "PIN", "Peran", "Aksi"].map((h) => (
              <th key={h} className="px-4 py-3 text-left text-[11px] font-bold text-slate-500 uppercase tracking-wider">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-50">
          {employees.map((emp) => (
            <tr key={emp.id} className="hover:bg-slate-50/60 transition-colors group">
              <td className="px-4 py-3">
                <div className="flex items-center gap-2.5">
                  <div className="w-8 h-8 rounded-xl bg-gradient-to-br from-violet-400 to-indigo-500 flex items-center justify-center shrink-0">
                    <span className="text-white text-xs font-black">{emp.name.substring(0, 2).toUpperCase()}</span>
                  </div>
                  <span className="font-semibold text-slate-800">{emp.name}</span>
                </div>
              </td>
              <td className="px-4 py-3 text-slate-500 text-xs">{emp.email}</td>
              <td className="px-4 py-3">
                <span className="font-mono text-xs bg-slate-100 px-2 py-1 rounded-lg tracking-widest text-slate-400">
                  ••••••
                </span>
              </td>
              <td className="px-4 py-3">
                <RoleBadge role={emp.role} />
              </td>
              <td className="px-4 py-3">
                <div className="flex items-center gap-1.5">
                  <button
                    onClick={() => onEdit(emp)}
                    className="p-1.5 text-slate-400 hover:text-violet-600 hover:bg-violet-50 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <Pencil className="w-3.5 h-3.5" />
                  </button>
                  <button
                    onClick={() => onDelete(emp)}
                    className="p-1.5 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                    title="Hapus"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// Seeded fallback data for demo when API is not available
const SEEDED_EMPLOYEES: Employee[] = [
  { id: "1", branch_id: "1", name: "Super Admin", email: "admin@kopitiam.com", role: "superadmin" },
  { id: "2", branch_id: "1", name: "Manager", email: "manager@kopitiam.com", role: "manager" },
  { id: "3", branch_id: "1", name: "Cashier 1", email: "cashier1@kopitiam.com", role: "cashier" },
  { id: "4", branch_id: "1", name: "Kitchen 1", email: "kitchen1@kopitiam.com", role: "kitchen" },
];

function EmployeeFallbackTable({
  onEdit, onDelete,
}: {
  onEdit: (e: Employee) => void;
  onDelete: (e: Employee) => void;
}) {
  return (
    <div className="w-full mt-4">
      <p className="text-[11px] text-amber-600 font-semibold mb-2 text-center">⚠ Data demo (seed) — bukan data live</p>
      <EmployeeTable employees={SEEDED_EMPLOYEES} onEdit={onEdit} onDelete={onDelete} />
    </div>
  );
}
