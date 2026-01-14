import React from 'react';
import { Layout, Menu } from 'antd';
import { Link, useNavigate, useLocation, Outlet } from 'react-router-dom';
import {
  DashboardOutlined,
  UserOutlined,
  ShoppingOutlined,
  WalletOutlined,
  FileTextOutlined,
  BarChartOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';
import Header from '../components/Header';

const { Content, Sider } = Layout;

const AdminLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { logout } = useAuth();
  const { t } = useTranslation();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const menuItems = [
    {
      key: '/admin/dashboard',
      icon: <DashboardOutlined />,
      label: <Link to="/admin/dashboard">{t('nav.dashboard')}</Link>,
    },
    {
      key: '/admin/users',
      icon: <UserOutlined />,
      label: <Link to="/admin/users">{t('nav.userManagement')}</Link>,
    },
    {
      key: '/admin/products',
      icon: <ShoppingOutlined />,
      label: <Link to="/admin/products">{t('nav.productManagement')}</Link>,
    },
    {
      key: '/admin/points',
      icon: <WalletOutlined />,
      label: <Link to="/admin/points">{t('nav.pointsManagement')}</Link>,
    },
    {
      key: '/admin/orders',
      icon: <FileTextOutlined />,
      label: <Link to="/admin/orders">{t('nav.orderManagement')}</Link>,
    },
    {
      key: '/admin/reports',
      icon: <BarChartOutlined />,
      label: <Link to="/admin/reports">{t('nav.reports')}</Link>,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: t('auth.logout'),
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
            selectedKeys={[location.pathname]}
            style={{ height: '100%', borderRight: 0 }}
            items={menuItems}
          />
        </Sider>
        <Layout style={{ padding: '0' }}>
          <Content
            style={{
              minHeight: 280,
              background: '#f0f2f5',
            }}
          >
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AdminLayout;
