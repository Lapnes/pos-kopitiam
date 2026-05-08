"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Store, UserCircle, KeyRound, ArrowRight } from "lucide-react";

export default function LoginPage() {
  const router = useRouter();
  
  // Cashier State
  const [pin, setPin] = useState("");
  
  // Admin State
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleCashierLogin = (e: React.FormEvent) => {
    e.preventDefault();
    if (pin.length < 4) {
      toast.error("Invalid PIN", { description: "Please enter a valid PIN." });
      return;
    }
    
    // Mock login success
    toast.success("Login Successful", { 
      description: "Welcome back, Cashier 01!",
      icon: <UserCircle className="w-5 h-5 text-emerald-500" />
    });
    router.push("/");
  };

  const handleAdminLogin = (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      toast.error("Missing Credentials", { description: "Please enter both email and password." });
      return;
    }
    
    // Mock login success
    toast.success("Admin Access Granted", { 
      description: "Redirecting to Dashboard..." 
    });
    router.push("/admin"); // Will be created in Phase 4
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50 p-4">
      <div className="w-full max-w-md bg-white rounded-2xl shadow-xl overflow-hidden border border-slate-100">
        
        {/* Header */}
        <div className="bg-emerald-600 p-8 text-center text-white relative overflow-hidden">
          <div className="absolute top-[-20%] left-[-10%] w-40 h-40 bg-emerald-500 rounded-full blur-3xl opacity-50"></div>
          <div className="absolute bottom-[-20%] right-[-10%] w-40 h-40 bg-emerald-700 rounded-full blur-3xl opacity-50"></div>
          
          <div className="relative z-10 flex flex-col items-center">
            <div className="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center backdrop-blur-sm mb-4 border border-white/30 shadow-inner">
              <Store className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-2xl font-bold tracking-tight">POS KopiTiam</h1>
            <p className="text-emerald-100 text-sm mt-1">Enterprise Management System</p>
          </div>
        </div>

        {/* Tabs Content */}
        <div className="p-6">
          <Tabs defaultValue="cashier" className="w-full">
            <TabsList className="grid w-full grid-cols-2 mb-6 bg-slate-100 p-1 rounded-xl h-12">
              <TabsTrigger value="cashier" className="rounded-lg data-[state=active]:bg-white data-[state=active]:shadow-sm text-sm font-semibold">
                Cashier
              </TabsTrigger>
              <TabsTrigger value="admin" className="rounded-lg data-[state=active]:bg-white data-[state=active]:shadow-sm text-sm font-semibold">
                Admin
              </TabsTrigger>
            </TabsList>

            <TabsContent value="cashier" className="space-y-6 animate-in fade-in zoom-in-95 duration-200">
              <div className="text-center mb-6">
                <h2 className="text-xl font-bold text-slate-800">Cashier Login</h2>
                <p className="text-sm text-slate-500 mt-1">Enter your assigned PIN to start shifting.</p>
              </div>
              
              <form onSubmit={handleCashierLogin} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="pin">PIN Code</Label>
                  <div className="relative">
                    <KeyRound className="absolute left-3 top-3 h-5 w-5 text-slate-400" />
                    <Input 
                      id="pin" 
                      type="password" 
                      placeholder="••••••" 
                      className="pl-10 h-12 text-lg tracking-widest font-mono"
                      value={pin}
                      onChange={(e) => setPin(e.target.value)}
                      maxLength={6}
                    />
                  </div>
                </div>
                
                <button 
                  type="submit"
                  className="w-full h-12 flex items-center justify-center gap-2 bg-slate-900 hover:bg-slate-800 text-white rounded-xl font-semibold transition-all active:scale-[0.98] shadow-md"
                >
                  Login as Cashier
                  <ArrowRight className="w-4 h-4" />
                </button>
              </form>
            </TabsContent>

            <TabsContent value="admin" className="space-y-6 animate-in fade-in zoom-in-95 duration-200">
              <div className="text-center mb-6">
                <h2 className="text-xl font-bold text-slate-800">Admin Login</h2>
                <p className="text-sm text-slate-500 mt-1">Access backoffice and master data.</p>
              </div>
              
              <form onSubmit={handleAdminLogin} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email Address</Label>
                  <Input 
                    id="email" 
                    type="email" 
                    placeholder="admin@kopitiam.com" 
                    className="h-11"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <Input 
                    id="password" 
                    type="password" 
                    placeholder="••••••••" 
                    className="h-11"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                </div>
                
                <button 
                  type="submit"
                  className="w-full h-12 flex items-center justify-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-semibold transition-all active:scale-[0.98] shadow-md mt-6"
                >
                  Login as Admin
                  <ArrowRight className="w-4 h-4" />
                </button>
              </form>
            </TabsContent>
          </Tabs>
        </div>
        
      </div>
    </div>
  );
}
