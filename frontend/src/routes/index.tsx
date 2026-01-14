import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import EmployeeLayout from '../layouts/EmployeeLayout';
import AdminLayout from '../layouts/AdminLayout';

// Employee Pages
import {
  LoginPage,
  ProductListPage,
  RedemptionHistoryPage,
  PointsHistoryPage,
  ProfilePage,
} from '../pages';

// Admin Pages
import {
  AdminDashboardPage,
  AdminUserManagementPage,
  AdminProductManagementPage,
  AdminPointsManagementPage,
  AdminOrderManagementPage,
  AdminReportsPage,
} from '../pages/admin';

// Protected Route Component
const ProtectedRoute: React.FC<{ children: React.ReactNode; requiredRole?: string }> = ({
  children,
  requiredRole,
}) => {
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && user?.role !== requiredRole) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

const AppRoutes: React.FC = () => {
  const { isAuthenticated, user } = useAuth();

  return (
    <Routes>
      {/* Public routes */}
      <Route
        path="/login"
        element={
          isAuthenticated ? (
            <Navigate to={user?.role === 'admin' ? '/admin/dashboard' : '/products'} replace />
          ) : (
            <LoginPage />
          )
        }
      />

      {/* Employee routes */}
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <EmployeeLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<Navigate to="/products" replace />} />
        <Route path="products" element={<ProductListPage />} />
        <Route path="redemptions" element={<RedemptionHistoryPage />} />
        <Route path="points" element={<PointsHistoryPage />} />
        <Route path="profile" element={<ProfilePage />} />
      </Route>

      {/* Admin routes */}
      <Route
        path="/admin"
        element={
          <ProtectedRoute requiredRole="admin">
            <AdminLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<Navigate to="/admin/dashboard" replace />} />
        <Route path="dashboard" element={<AdminDashboardPage />} />
        <Route path="users" element={<AdminUserManagementPage />} />
        <Route path="products" element={<AdminProductManagementPage />} />
        <Route path="points" element={<AdminPointsManagementPage />} />
        <Route path="orders" element={<AdminOrderManagementPage />} />
        <Route path="reports" element={<AdminReportsPage />} />
      </Route>

      {/* Fallback route */}
      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  );
};

export default AppRoutes;
