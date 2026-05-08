import { create } from 'zustand';
import { CartItem, Menu } from '@/types';

interface CartState {
  items: CartItem[];
  addItem: (menu: Menu) => void;
  updateQuantity: (cartItemId: string, quantity: number) => void;
  removeItem: (cartItemId: string) => void;
  clearCart: () => void;

  getSubtotal: () => number;
  getTaxAmount: () => number;
  getServiceCharge: () => number;
  getGrandTotal: () => number;
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],

  addItem: (menu) => {
    set((state) => {
      const existingItem = state.items.find((item) => item.id === menu.id);
      if (existingItem) {
        return {
          items: state.items.map((item) =>
            item.id === menu.id ? { ...item, quantity: item.quantity + 1 } : item
          ),
        };
      }
      return {
        items: [...state.items, { ...menu, cartItemId: crypto.randomUUID(), quantity: 1 }],
      };
    });
  },

  updateQuantity: (cartItemId, quantity) => {
    set((state) => {
      if (quantity <= 0) {
        return { items: state.items.filter((item) => item.cartItemId !== cartItemId) };
      }
      return {
        items: state.items.map((item) =>
          item.cartItemId === cartItemId ? { ...item, quantity } : item
        ),
      };
    });
  },

  removeItem: (cartItemId) => {
    set((state) => ({
      items: state.items.filter((item) => item.cartItemId !== cartItemId),
    }));
  },

  clearCart: () => set({ items: [] }),

  getSubtotal: () => {
    const { items } = get();
    return items.reduce((total, item) => total + item.price * item.quantity, 0);
  },

  getTaxAmount: () => {
    return get().getSubtotal() * 0.11;
  },

  getServiceCharge: () => {
    return get().getSubtotal() * 0.05;
  },

  getGrandTotal: () => {
    return get().getSubtotal() + get().getTaxAmount() + get().getServiceCharge();
  },
}));
