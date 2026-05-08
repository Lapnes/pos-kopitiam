export interface Menu {
  id: string;
  category_id: string;
  name: string;
  description: string;
  price: number;
  cost_price: number;
  daily_stock: number;
  is_active: boolean;
  station: string;
  image_url?: string;
  bigcapital_item_id?: string;
  color?: string; // fallback if no image (frontend-only logic)
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
