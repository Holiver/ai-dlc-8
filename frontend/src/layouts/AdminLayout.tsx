import React from 'react';
import { Layout, Menu } from 'antd';
import { Routes, Route, Link, useNavigate } from 'react-router-dom';
import {
  DashboardOutlined,
  UserOutlined,
  ShoppingOutlined,
  WalletOutlined,
  FileTextOutlined,
  BarChartOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useAuth } from '../contexts/AuthContext';
import Header from '../components/Header';

const { Content, Sider } = Layout;

// Placeholder page components
const AdminDashboardPage = () => <div>Admin Dashboard Page</div>;
const AdminUserManagementPage = () => <div>User Management Page</div>;
const AdminProductManagementPage = () => <div>Product Management Page</div>;
const AdminPointsManagementPage = () => <div>Points Management Page</div>;
const AdminOrderManagementPage = () => <div>Order Management Page</div>;
const AdminReportsPage = () => <div>Reports Page</div>;

const AdminLayout: React.FC = () => {
  const navigate = useNavigate();
  const { logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const menuItems = [
    {
      key: '/admin/dashboard',
      icon: <DashboardOutlined />,
      label: <Link to="/admin/dashboard">Dashboard</Link>,
    },
    {
      key: '/admin/users',
      icon: <UserOutlined />,
      label: <Link to="/admin/users">User Management</Link>,
    },
    {
      key: '/admin/products',
      icon: <ShoppingOutlined />,
      label: <Link to="/admin/products">Product Management</Link>,
    },
    {
      key: '/admin/points',
      icon: <WalletOutlined />,
      label: <Link to="/admin/points">Points Management</Link>,
    },
    {
      key: '/admin/orders',
      icon: <FileTextOutlined />,
      label: <Link to="/admin/orders">Order Management</Link>,
    },
    {
      key: '/admin/reports',
      icon: <BarChartOutlined />,
      label: <Link to="/admin/reports">Reports</Link>,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'Logout',
      onClick: handleLogout,
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header />
      <Layout>
        <Sider width={200} theme="light">
          <Menu
            mode="inline"
            defaultSelectedKeys={['/admin/dashboard']}
            style={{ height: '100%', borderRight: 0 }}
            items={menuItems}
          />
        </Sider>
        <Layout style={{ padding: '24px' }}>
          <Content
            style={{
              padding: 24,
              margin: 0,
              minHeight: 280,
              background: '#fff',
            }}
          >
            <Routes>
              <Route path="/dashboard" element={<AdminDashboardPage />} />
              <Route path="/users" element={<AdminUserManagementPage />} />
              <Route path="/products" element={<AdminProductManagementPage />} />
              <Route path="/points" element={<AdminPointsManagementPage />} />
              <Route path="/orders" element={<AdminOrderManagementPage />} />
              <Route path="/reports" element={<AdminReportsPage />} />
              <Route path="/" element={<AdminDashboardPage />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AdminLayout;
