import apiService from './apiService';

export interface User {
  id: number;
  full_name: string;
  email: string;
  phone: string;
  role: string;
  points_balance: number;
  is_first_login: boolean;
  is_active: boolean;
  preferred_language: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

const authService = {
  // Login
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await apiService.post<LoginResponse>('/auth/login', credentials);
    
    // Store token and user in localStorage
    if (response.data.token) {
      localStorage.setItem('auth_token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    
    return response.data;
  },

  // Logout
  logout: async (): Promise<void> => {
    try {
      await apiService.post('/auth/logout');
    } finally {
      // Clear local storage regardless of API response
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
    }
  },

  // Get current user
  getCurrentUser: async (): Promise<User> => {
    const response = await apiService.get<{ user: User }>('/auth/me');
    return response.data.user;
  },

  // Get stored token
  getToken: (): string | null => {
    return localStorage.getItem('auth_token');
  },

  // Get stored user
  getStoredUser: (): User | null => {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch (e) {
        return null;
      }
    }
    return null;
  },

  // Check if user is authenticated
  isAuthenticated: (): boolean => {
    return !!authService.getToken();
  },

  // Check if user is admin
  isAdmin: (): boolean => {
    const user = authService.getStoredUser();
    return user?.role === 'admin';
  },

  // Update stored user
  updateStoredUser: (user: User): void => {
    localStorage.setItem('user', JSON.stringify(user));
  },
};

export default authService;
