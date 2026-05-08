export interface Menu {
  id: string;
  name: string;
  price: number;
  imageUrl?: string;
  color?: string; // fallback if no image
}

export interface CartItem extends Menu {
  cartItemId: string;
  quantity: number;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'cashier' | 'manager';
}

export interface Order {
  id: string;
  orderNumber: string;
  cashierId: string;
  subtotal: number;
  taxAmount: number;
  serviceCharge: number;
  grandTotal: number;
  status: 'pending' | 'completed' | 'cancelled';
  paymentMethod?: 'cash' | 'card' | 'qris';
  items: CartItem[];
  createdAt: string;
}
