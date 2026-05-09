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
  color?: string; // frontend-only for UI color card
}

export interface CartItem extends Menu {
  cartItemId: string;
  quantity: number;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'cashier' | 'manager' | 'superadmin' | 'kitchen';
}

export interface Employee {
  id: string;
  branch_id: string;
  name: string;
  email: string;
  pin_code?: string;
  role: 'admin' | 'cashier' | 'manager' | 'superadmin' | 'kitchen';
  created_at?: string;
}

export interface OrderDetail {
  id: string;
  order_id: string;
  menu_id: string;
  menu_name: string;
  quantity: number;
  price: number;
  subtotal: number;
  notes: string;
  is_voided: boolean;
  station: string;
  kds_status: string;
}

export interface Order {
  id: string;
  order_number: string;
  branch_id: string;
  employee_id: string;
  customer_name: string;
  customer_phone: string;
  order_type: 'dine_in' | 'takeaway' | 'delivery' | 'online';
  status: 'pending' | 'confirmed' | 'cancelled' | 'void';
  subtotal: number;
  tax_amount: number;
  service_charge: number;
  discount_amount: number;
  total: number;
  notes: string;
  has_returns: boolean;
  order_details: OrderDetail[];
  created_at: string;
  updated_at: string;
}

export interface AnalyticsSummary {
  total_revenue: number;
  total_transactions: number;
  total_orders: number;
  average_order_value: number;
  period_start?: string;
  period_end?: string;
}

export interface BestSeller {
  menu_id: string;
  menu_name: string;
  total_quantity: number;
  total_revenue: number;
}

export interface ReturnImpact {
  total_returns: number;
  total_return_amount: number;
  return_rate: number;
}

