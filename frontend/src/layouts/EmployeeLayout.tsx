import React from 'react';
import { Layout, Menu } from 'antd';
import { Link, useNavigate, useLocation, Outlet } from 'react-router-dom';
import {
  ShoppingOutlined,
  HistoryOutlined,
  WalletOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';
import Header from '../components/Header';

const { Content, Sider } = Layout;

const EmployeeLayout: React.FC = () => {
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
      key: '/products',
      icon: <ShoppingOutlined />,
      label: <Link to="/products">{t('nav.products')}</Link>,
    },
    {
      key: '/redemptions',
      icon: <HistoryOutlined />,
      label: <Link to="/redemptions">{t('nav.redemptions')}</Link>,
    },
    {
      key: '/points',
      icon: <WalletOutlined />,
      label: <Link to="/points">{t('nav.points')}</Link>,
    },
    {
      key: '/profile',
      icon: <UserOutlined />,
      label: <Link to="/profile">{t('nav.profile')}</Link>,
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

export default EmployeeLayout;
