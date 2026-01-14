import apiService from './apiService';
import { User } from './authService';
import { Product } from './productService';
import { RedemptionOrder } from './redemptionService';

// Admin User Management
export interface CreateEmployeeRequest {
  full_name: string;
  email: string;
  phone: string;
}

export interface SetEmployeeStatusRequest {
  is_active: boolean;
}

// Admin Product Management
export interface CreateProductRequest {
  name: string;
  image_url?: string;
  points_required: number;
  stock_quantity: number;
}

export interface UpdateProductRequest {
  name?: string;
  image_url?: string;
  points_required?: number;
  stock_quantity?: number;
}

export interface SetProductStatusRequest {
  status: 'active' | 'inactive';
}

export interface BatchImportProductsRequest {
  markdown: string;
}

// Admin Points Management
export interface GrantPointsRequest {
  user_id: number;
  amount: number;
  reason: string;
}

export interface DeductPointsRequest {
  user_id: number;
  amount: number;
  reason: string;
}

export interface BatchGrantPointsRequest {
  markdown: string;
}

// Admin Order Management
export interface BatchUpdateOrderStatusRequest {
  order_numbers: string;
  status: 'preparing' | 'delivered';
}

// Admin Reports
export interface PointsGrantStats {
  user_name: string;
  user_email: string;
  amount: number;
  reason: string;
  operator_name: string;
  created_at: string;
}

export interface PointsBalanceStats {
  user_name: string;
  user_email: string;
  points_balance: number;
}

export interface RedemptionStats {
  product_name: string;
  product_id: number;
  user_name: string;
  user_email: string;
  points_cost: number;
  status: string;
  created_at: string;
}

const adminService = {
  // User Management
  users: {
    create: async (data: CreateEmployeeRequest): Promise<User> => {
      const response = await apiService.post<{ user: User; message: string }>('/admin/users', data);
      return response.data.user;
    },

    setStatus: async (userId: number, isActive: boolean): Promise<void> => {
      await apiService.put(`/admin/users/${userId}/status`, { is_active: isActive });
    },

    list: async (isActive?: boolean): Promise<User[]> => {
      const params = isActive !== undefined ? { is_active: isActive } : {};
      const response = await apiService.get<{ users: User[] }>('/admin/users', { params });
      return response.data.users;
    },
  },

  // Product Management
  products: {
    create: async (data: CreateProductRequest): Promise<Product> => {
      const response = await apiService.post<{ product: Product }>('/admin/products', data);
      return response.data.product;
    },

    update: async (productId: number, data: UpdateProductRequest): Promise<Product> => {
      const response = await apiService.put<{ product: Product }>(`/admin/products/${productId}`, data);
      return response.data.product;
    },

    setStatus: async (productId: number, status: 'active' | 'inactive'): Promise<void> => {
      await apiService.put(`/admin/products/${productId}/status`, { status });
    },

    batchImport: async (markdown: string): Promise<Product[]> => {
      const response = await apiService.post<{ products: Product[]; count: number; message: string }>(
        '/admin/products/batch',
        { markdown }
      );
      return response.data.products;
    },

    list: async (status?: string): Promise<Product[]> => {
      const params = status ? { status } : {};
      const response = await apiService.get<{ products: Product[] }>('/admin/products', { params });
      return response.data.products;
    },
  },

  // Points Management
  points: {
    grant: async (data: GrantPointsRequest): Promise<void> => {
      await apiService.post('/admin/points/grant', data);
    },

    deduct: async (data: DeductPointsRequest): Promise<void> => {
      await apiService.post('/admin/points/deduct', data);
    },

    batchGrant: async (markdown: string): Promise<void> => {
      await apiService.post('/admin/points/batch-grant', { markdown });
    },
  },

  // Order Management
  orders: {
    list: async (status?: string, userId?: number): Promise<RedemptionOrder[]> => {
      const params: any = {};
      if (status) params.status = status;
      if (userId) params.user_id = userId;
      const response = await apiService.get<{ orders: RedemptionOrder[] }>('/admin/orders', { params });
      return response.data.orders;
    },

    batchUpdateStatus: async (orderNumbers: string, status: 'preparing' | 'delivered'): Promise<void> => {
      await apiService.put('/admin/orders/batch-status', { order_numbers: orderNumbers, status });
    },
  },

  // Reports
  reports: {
    getPointsGrants: async (): Promise<PointsGrantStats[]> => {
      const response = await apiService.get<{ grants: PointsGrantStats[] }>('/admin/reports/points-grants');
      return response.data.grants;
    },

    getPointsBalances: async (): Promise<PointsBalanceStats[]> => {
      const response = await apiService.get<{ balances: PointsBalanceStats[] }>('/admin/reports/points-balances');
      return response.data.balances;
    },

    getRedemptions: async (): Promise<RedemptionStats[]> => {
      const response = await apiService.get<{ redemptions: RedemptionStats[] }>('/admin/reports/redemptions');
      return response.data.redemptions;
    },
  },
};

export default adminService;
