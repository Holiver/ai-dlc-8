import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import EmployeeLayout from '../layouts/EmployeeLayout';
import AdminLayout from '../layouts/AdminLayout';

// Placeholder components - will be implemented later
const LoginPage = () => <div>Login Page</div>;

const AppRoutes: React.FC = () => {
  const { isAuthenticated, user } = useAuth();

  return (
    <Routes>
      {/* Public routes */}
      <Route path="/login" element={<LoginPage />} />

      {/* Protected routes */}
      {isAuthenticated ? (
        <>
          {user?.role === 'admin' ? (
            <Route path="/admin/*" element={<AdminLayout />} />
          ) : (
            <Route path="/*" element={<EmployeeLayout />} />
          )}
        </>
      ) : (
        <Route path="*" element={<Navigate to="/login" replace />} />
      )}
    </Routes>
  );
};

export default AppRoutes;
