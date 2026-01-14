import apiService from './apiService';
import { User } from './authService';
import { RedemptionOrder } from './redemptionService';

export interface PointsTransaction {
  id: number;
  user_id: number;
  transaction_type: string;
  amount: number;
  balance_after: number;
  reason: string;
  operator_id?: number;
  related_order_id?: number;
  created_at: string;
  user?: User;
  operator?: User;
  related_order?: RedemptionOrder;
}

export interface PointsTransactionsResponse {
  transactions: PointsTransaction[];
  total: number;
  page: number;
  page_size: number;
}

const pointsService = {
  // Get points balance
  getPointsBalance: async (): Promise<number> => {
    const response = await apiService.get<{ balance: number }>('/points/balance');
    return response.data.balance;
  },

  // Get points transactions with pagination
  getPointsTransactions: async (page: number = 1, pageSize: number = 20): Promise<PointsTransactionsResponse> => {
    const response = await apiService.get<PointsTransactionsResponse>('/points/transactions', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  },
};

export default pointsService;
