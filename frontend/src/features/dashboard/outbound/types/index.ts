export interface Order {
  id: string;
  order_sn: string;
  shop_id: string;
  marketplace_status: string;
  shipping_status: string;
  wms_status: string;
  tracking_number?: string;
  total_amount: number;
  created_at: string;
  updated_at: string;
}

export interface OrderDetail {
  id: string;
  order_sn: string;
  shop_id: string;
  marketplace_status: string;
  shipping_status: string;
  wms_status: string;
  tracking_number?: string;
  items: OrderItem[];
  total_amount: number;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  id: string;
  order_id: string;
  sku: string;
  quantity: number;
  price: number;
  created_at: string;
}

export type ShipOrderPayload = {
  orderSn: string;
  channelId: string;
};
