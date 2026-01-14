// Export all services
export { default as apiService } from './apiService';
export { default as authService } from './authService';
export { default as storageService } from './storageService';
export { default as productService } from './productService';
export { default as redemptionService } from './redemptionService';
export { default as pointsService } from './pointsService';
export { default as userService } from './userService';
export { default as adminService } from './adminService';

// Export types
export type { User, LoginRequest, LoginResponse } from './authService';
export type { Product } from './productService';
export type { RedemptionOrder, RedeemProductRequest } from './redemptionService';
export type { PointsTransaction, PointsTransactionsResponse } from './pointsService';
export type { UpdatePhoneRequest } from './userService';
export type {
  CreateEmployeeRequest,
  SetEmployeeStatusRequest,
  CreateProductRequest,
  UpdateProductRequest,
  SetProductStatusRequest,
  BatchImportProductsRequest,
  GrantPointsRequest,
  DeductPointsRequest,
  BatchGrantPointsRequest,
  BatchUpdateOrderStatusRequest,
  PointsGrantStats,
  PointsBalanceStats,
  RedemptionStats,
} from './adminService';
