export interface Category {
  category_id: number;
  category_name: string;
  menus?: Menu[];
}

export interface Menu {
  menu_id: number;
  category_id: number;
  menu_name: string;
  price: number;
  daily_stock: number;
  category?: Category;
}

export interface Employee {
  employee_id: number;
  employee_name: string;
  phone_number?: string;
}

export interface CartItem {
  menu_id: number;
  name: string;
  price: number;
  qty: number;
}

export interface OrderReportRow {
  order_id: string;
  employee_name: string;
  menu_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  total_price: number;
}

export interface CreateOrderPayload {
  employee_id?: number;
  items: { menu_id: number; quantity: number }[];
}
