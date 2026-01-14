import apiService from './apiService';
import { User } from './authService';

export interface UpdatePhoneRequest {
  phone: string;
}

const userService = {
  // Get user profile
  getProfile: async (): Promise<User> => {
    const response = await apiService.get<{ user: User }>('/users/profile');
    return response.data.user;
  },

  // Update phone number
  updatePhone: async (phone: string): Promise<void> => {
    await apiService.put('/users/phone', { phone });
  },
};

export default userService;
