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
