"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Store, KeyRound, ArrowRight, Mail, Lock, Coffee } from "lucide-react";
import api from "@/lib/api/axios";
import { useAuthStore } from "@/store/useAuthStore";
import { AxiosError } from "axios";

export default function LoginPage() {
  const router = useRouter();
  const { setAuth, isAuthenticated, user } = useAuthStore();

  const [pin, setPin] = useState("");
  const [isLoadingCashier, setIsLoadingCashier] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoadingAdmin, setIsLoadingAdmin] = useState(false);

  useEffect(() => {
    if (isAuthenticated && user) {
      if (user.role === "cashier") {
        router.replace("/");
      } else {
        router.replace("/");
      }
    }
  }, [isAuthenticated, user, router]);

  const handleCashierLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (pin.length < 4) {
      toast.error("PIN tidak valid", { description: "Masukkan PIN minimal 4 digit." });
      return;
    }
    setIsLoadingCashier(true);
    try {
      const response = await api.post("/auth/login-pin", { pin });
      const { user: userData, access_token } = response.data.data;
      setAuth(userData, access_token);
      toast.success("Login berhasil!", { description: `Selamat datang, ${userData.name}!` });
      router.push("/");
    } catch (error) {
      const err = error as AxiosError<{ message: string }>;
      console.error("CASHIER LOGIN ERROR:", err.response?.data || err.message);
      toast.error("Login gagal", {
        description: err.response?.data?.message || err.message || "PIN salah atau server error.",
      });
    } finally {
      setIsLoadingCashier(false);
    }
  };

  const handleAdminLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      toast.error("Data tidak lengkap", { description: "Masukkan email dan password." });
      return;
    }
    setIsLoadingAdmin(true);
    try {
      const response = await api.post("/auth/login", { email, password });
      const { user: userData, access_token } = response.data.data;
      setAuth(userData, access_token);
      toast.success("Login Admin berhasil!", { description: "Mengarahkan ke dashboard..." });
      router.push("/");
    } catch (error) {
      const err = error as AxiosError<{ message: string }>;
      console.error("ADMIN LOGIN ERROR:", err.response?.data || err.message);
      toast.error("Login gagal", {
        description: err.response?.data?.message || err.message || "Kredensial salah.",
      });
    } finally {
      setIsLoadingAdmin(false);
    }
  };

  return (
    <div className="min-h-screen flex bg-slate-950">
      {/* Left decorative panel */}
      <div className="hidden lg:flex lg:w-1/2 relative overflow-hidden bg-gradient-to-br from-emerald-600 via-emerald-700 to-teal-900">
        <div className="absolute inset-0">
          <div className="absolute top-20 left-20 w-72 h-72 bg-emerald-400/20 rounded-full blur-3xl" />
          <div className="absolute bottom-20 right-10 w-96 h-96 bg-teal-400/15 rounded-full blur-3xl" />
          <div className="absolute top-1/2 left-1/3 w-48 h-48 bg-white/5 rounded-full blur-2xl" />
          {/* Grid pattern */}
          <div
            className="absolute inset-0 opacity-10"
            style={{
              backgroundImage: `radial-gradient(circle at 1px 1px, white 1px, transparent 0)`,
              backgroundSize: "40px 40px",
            }}
          />
        </div>
        <div className="relative z-10 flex flex-col justify-center items-center w-full p-12 text-white">
          <div className="w-24 h-24 bg-white/15 backdrop-blur-sm rounded-3xl flex items-center justify-center border border-white/25 shadow-2xl mb-8">
            <Coffee className="w-12 h-12 text-white" />
          </div>
          <h1 className="text-4xl font-bold tracking-tight mb-3">POS KopiTiam</h1>
          <p className="text-emerald-100 text-lg text-center leading-relaxed max-w-sm">
            Sistem manajemen kafe terintegrasi untuk efisiensi operasional maksimal.
          </p>

          <div className="mt-16 grid grid-cols-3 gap-6 w-full max-w-xs">
            {[
              { label: "Menu Aktif", value: "24+" },
              { label: "Transaksi/Hari", value: "150+" },
              { label: "Kasir", value: "5" },
            ].map((stat) => (
              <div key={stat.label} className="text-center">
                <div className="text-2xl font-bold text-white">{stat.value}</div>
                <div className="text-xs text-emerald-200 mt-1">{stat.label}</div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Right login panel */}
      <div className="flex-1 flex items-center justify-center p-6 bg-slate-50">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-3 mb-8 justify-center">
            <div className="w-12 h-12 bg-emerald-600 rounded-2xl flex items-center justify-center">
              <Store className="w-6 h-6 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-slate-900">POS KopiTiam</h1>
          </div>

          <div className="bg-white rounded-3xl shadow-xl border border-slate-100 overflow-hidden">
            <div className="p-8 border-b border-slate-100 bg-gradient-to-r from-slate-50 to-white">
              <h2 className="text-2xl font-bold text-slate-900">Masuk ke Akun</h2>
              <p className="text-slate-500 text-sm mt-1">
                Pilih tipe login sesuai role Anda.
              </p>
            </div>

            <div className="p-8">
              <Tabs defaultValue="cashier" className="w-full">
                <TabsList className="grid w-full grid-cols-2 mb-8 bg-slate-100 p-1 rounded-2xl h-12">
                  <TabsTrigger
                    value="cashier"
                    className="rounded-xl data-[state=active]:bg-white data-[state=active]:shadow-sm data-[state=active]:text-emerald-700 text-sm font-semibold transition-all"
                  >
                    <KeyRound className="w-4 h-4 mr-2" />
                    Kasir
                  </TabsTrigger>
                  <TabsTrigger
                    value="admin"
                    className="rounded-xl data-[state=active]:bg-white data-[state=active]:shadow-sm data-[state=active]:text-emerald-700 text-sm font-semibold transition-all"
                  >
                    <Lock className="w-4 h-4 mr-2" />
                    Admin
                  </TabsTrigger>
                </TabsList>

                {/* Cashier Tab */}
                <TabsContent value="cashier" className="space-y-5 animate-in fade-in-0 zoom-in-95 duration-200">
                  <div className="text-center mb-2">
                    <div className="w-14 h-14 bg-emerald-50 rounded-2xl flex items-center justify-center mx-auto mb-3">
                      <KeyRound className="w-7 h-7 text-emerald-600" />
                    </div>
                    <h3 className="font-bold text-slate-800">Login Kasir</h3>
                    <p className="text-xs text-slate-500 mt-1">Masukkan PIN 6 digit Anda untuk mulai shift.</p>
                  </div>

                  <form onSubmit={handleCashierLogin} className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="pin" className="text-slate-700 font-medium">
                        PIN Kasir
                      </Label>
                      <div className="relative">
                        <KeyRound className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4.5 w-4.5 text-slate-400" />
                        <Input
                          id="pin"
                          type="password"
                          placeholder="• • • • • •"
                          className="pl-11 h-12 text-xl tracking-[0.5em] font-mono border-slate-200 focus:border-emerald-400 focus:ring-emerald-400/20 rounded-xl"
                          value={pin}
                          onChange={(e) => setPin(e.target.value.replace(/\D/g, ""))}
                          maxLength={6}
                          inputMode="numeric"
                          disabled={isLoadingCashier}
                          autoFocus
                        />
                      </div>
                    </div>

                    {/* PIN dots indicator */}
                    <div className="flex justify-center gap-3 py-2">
                      {Array.from({ length: 6 }).map((_, i) => (
                        <div
                          key={i}
                          className={`w-3 h-3 rounded-full transition-all duration-200 ${
                            i < pin.length
                              ? "bg-emerald-500 scale-110"
                              : "bg-slate-200"
                          }`}
                        />
                      ))}
                    </div>

                    <button
                      type="submit"
                      disabled={isLoadingCashier || pin.length < 4}
                      className="w-full h-12 flex items-center justify-center gap-2 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-200 disabled:text-slate-400 disabled:cursor-not-allowed text-white rounded-xl font-semibold transition-all active:scale-[0.98] shadow-md shadow-emerald-200"
                    >
                      {isLoadingCashier ? (
                        <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                      ) : (
                        <>
                          Masuk sebagai Kasir
                          <ArrowRight className="w-4 h-4" />
                        </>
                      )}
                    </button>
                  </form>
                </TabsContent>

                {/* Admin Tab */}
                <TabsContent value="admin" className="space-y-5 animate-in fade-in-0 zoom-in-95 duration-200">
                  <div className="text-center mb-2">
                    <div className="w-14 h-14 bg-slate-100 rounded-2xl flex items-center justify-center mx-auto mb-3">
                      <Lock className="w-7 h-7 text-slate-600" />
                    </div>
                    <h3 className="font-bold text-slate-800">Login Admin / Manager</h3>
                    <p className="text-xs text-slate-500 mt-1">Akses backoffice dan data master.</p>
                  </div>

                  <form onSubmit={handleAdminLogin} className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="email" className="text-slate-700 font-medium">
                        Email
                      </Label>
                      <div className="relative">
                        <Mail className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                        <Input
                          id="email"
                          type="email"
                          placeholder="admin@kopitiam.com"
                          className="pl-11 h-12 border-slate-200 focus:border-slate-400 rounded-xl"
                          value={email}
                          onChange={(e) => setEmail(e.target.value)}
                          disabled={isLoadingAdmin}
                        />
                      </div>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="password" className="text-slate-700 font-medium">
                        Password
                      </Label>
                      <div className="relative">
                        <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                        <Input
                          id="password"
                          type="password"
                          placeholder="••••••••"
                          className="pl-11 h-12 border-slate-200 focus:border-slate-400 rounded-xl"
                          value={password}
                          onChange={(e) => setPassword(e.target.value)}
                          disabled={isLoadingAdmin}
                        />
                      </div>
                    </div>

                    <button
                      type="submit"
                      disabled={isLoadingAdmin}
                      className="w-full h-12 flex items-center justify-center gap-2 bg-slate-900 hover:bg-slate-800 disabled:bg-slate-200 disabled:text-slate-400 disabled:cursor-not-allowed text-white rounded-xl font-semibold transition-all active:scale-[0.98] shadow-md mt-2"
                    >
                      {isLoadingAdmin ? (
                        <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                      ) : (
                        <>
                          Masuk sebagai Admin
                          <ArrowRight className="w-4 h-4" />
                        </>
                      )}
                    </button>
                  </form>
                </TabsContent>
              </Tabs>
            </div>

            <div className="px-8 py-4 bg-slate-50 border-t border-slate-100 text-center">
              <p className="text-xs text-slate-400">
                POS KopiTiam v1.0 &mdash; Enterprise Edition
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
