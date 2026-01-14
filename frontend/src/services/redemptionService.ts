import apiService from './apiService';
import { Product } from './productService';
import { User } from './authService';

export interface RedemptionOrder {
  id: number;
  order_number: string;
  user_id: number;
  product_id: number;
  product_name: string;
  points_cost: number;
  points_balance_after: number;
  status: string;
  created_at: string;
  updated_at: string;
  user?: User;
  product?: Product;
}

export interface RedeemProductRequest {
  product_id: number;
}

const redemptionService = {
  // Create redemption order
  redeemProduct: async (productId: number): Promise<RedemptionOrder> => {
    const response = await apiService.post<{ order: RedemptionOrder }>('/redemptions', {
      product_id: productId,
    });
    return response.data.order;
  },

  // Get redemption history
  getRedemptionHistory: async (): Promise<RedemptionOrder[]> => {
    const response = await apiService.get<{ orders: RedemptionOrder[] }>('/redemptions');
    return response.data.orders;
  },

  // Get redemption by ID
  getRedemptionById: async (id: number): Promise<RedemptionOrder> => {
    const response = await apiService.get<{ order: RedemptionOrder }>(`/redemptions/${id}`);
    return response.data.order;
  },
};

export default redemptionService;
