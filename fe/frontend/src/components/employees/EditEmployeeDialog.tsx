"use client";

import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import api from "@/lib/api/axios";
import { Employee } from "@/types";
import { AxiosError } from "axios";
import { toast } from "sonner";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Pencil, Loader2, Save, X } from "lucide-react";

interface Props {
  employee: Employee;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

const ROLES = [
  { value: "superadmin", label: "Super Admin" },
  { value: "manager", label: "Manajer" },
  { value: "cashier", label: "Kasir" },
  { value: "kitchen", label: "Dapur" },
  { value: "admin", label: "Admin" },
];

export function EditEmployeeDialog({ employee, open, onOpenChange, onSuccess }: Props) {
  const [name, setName] = useState(employee.name);
  const [email, setEmail] = useState(employee.email);
  const [pin, setPin] = useState("");
  const [role, setRole] = useState(employee.role);

  const mutation = useMutation({
    mutationFn: () =>
      api.put(`/api/v1/employees/${employee.id}`, {
        name,
        email,
        role,
        ...(pin ? { pin_code: pin } : {}),
      }),
    onSuccess: () => {
      toast.success("Data karyawan berhasil diperbarui");
      onSuccess();
    },
    onError: (err: AxiosError<{ message: string }>) => {
      toast.error("Gagal memperbarui karyawan", {
        description: err.response?.data?.message ?? "Terjadi error saat menyimpan.",
      });
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate();
  };

  return (
    <Dialog open={open} onOpenChange={(v) => { if (!mutation.isPending) onOpenChange(v); }}>
      <DialogContent className="sm:max-w-[420px] p-0 gap-0 rounded-3xl overflow-hidden">
        {/* Header */}
        <div className="bg-violet-600 px-6 py-5 text-white">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2.5 text-white text-lg font-bold">
              <div className="w-8 h-8 bg-white/20 rounded-xl flex items-center justify-center">
                <Pencil className="w-4 h-4" />
              </div>
              Edit Karyawan
            </DialogTitle>
          </DialogHeader>
          <p className="text-violet-100 text-sm mt-1">Perbarui data untuk <strong>{employee.name}</strong></p>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="emp-name" className="text-slate-700 font-semibold text-sm">Nama</Label>
            <Input
              id="emp-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="h-11 rounded-xl border-slate-200 focus:border-violet-400 focus:ring-violet-400/20"
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="emp-email" className="text-slate-700 font-semibold text-sm">Email</Label>
            <Input
              id="emp-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="h-11 rounded-xl border-slate-200 focus:border-violet-400 focus:ring-violet-400/20"
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="emp-pin" className="text-slate-700 font-semibold text-sm">
              PIN Baru <span className="font-normal text-slate-400">(kosongkan jika tidak ingin mengubah)</span>
            </Label>
            <Input
              id="emp-pin"
              type="password"
              placeholder="6 digit PIN..."
              value={pin}
              onChange={(e) => setPin(e.target.value)}
              maxLength={6}
              className="h-11 rounded-xl border-slate-200 font-mono tracking-widest focus:border-violet-400 focus:ring-violet-400/20"
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="emp-role" className="text-slate-700 font-semibold text-sm">Peran</Label>
            <select
              id="emp-role"
              value={role}
              onChange={(e) => setRole(e.target.value as Employee["role"])}
              className="w-full h-11 rounded-xl border border-slate-200 bg-white px-3 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-violet-400/30 focus:border-violet-400"
            >
              {ROLES.map((r) => (
                <option key={r.value} value={r.value}>{r.label}</option>
              ))}
            </select>
          </div>

          <DialogFooter className="gap-2 pt-2">
            <button
              type="button"
              onClick={() => onOpenChange(false)}
              disabled={mutation.isPending}
              className="flex-1 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl font-semibold text-sm transition-colors flex items-center justify-center gap-2"
            >
              <X className="w-4 h-4" />
              Batal
            </button>
            <button
              type="submit"
              disabled={mutation.isPending || !name.trim() || !email.trim()}
              className="flex-1 py-2.5 bg-violet-600 hover:bg-violet-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-xl font-bold text-sm transition-all shadow-md shadow-violet-200 active:scale-[0.98] flex items-center justify-center gap-2"
            >
              {mutation.isPending ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <>
                  <Save className="w-4 h-4" />
                  Simpan
                </>
              )}
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
