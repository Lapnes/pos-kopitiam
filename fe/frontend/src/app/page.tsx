'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/useAuthStore';

const ROLE_HOME: Record<string, string> = {
  cashier:    '/kasir',
  manager:    '/manager',
  superadmin: '/admin',
  kitchen:    '/kitchen',
};

export default function RootPage() {
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();

  useEffect(() => {
    if (!isAuthenticated || !user) {
      router.replace('/login');
    } else {
      router.replace(ROLE_HOME[user.role] ?? '/login');
    }
  }, [isAuthenticated, user, router]);

  return null;
}
