import React from 'react';
import { Layout, Menu } from 'antd';
import { Routes, Route, Link, useNavigate } from 'react-router-dom';
import {
  ShoppingOutlined,
  HistoryOutlined,
  WalletOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useAuth } from '../contexts/AuthContext';
import Header from '../components/Header';

const { Content, Sider } = Layout;

// Placeholder page components
const ProductListPage = () => <div>Product List Page</div>;
const RedemptionHistoryPage = () => <div>Redemption History Page</div>;
const PointsHistoryPage = () => <div>Points History Page</div>;
const ProfilePage = () => <div>Profile Page</div>;

const EmployeeLayout: React.FC = () => {
  const navigate = useNavigate();
  const { logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const menuItems = [
    {
      key: '/products',
      icon: <ShoppingOutlined />,
      label: <Link to="/products">Products</Link>,
    },
    {
      key: '/redemptions',
      icon: <HistoryOutlined />,
      label: <Link to="/redemptions">Redemption History</Link>,
    },
    {
      key: '/points',
      icon: <WalletOutlined />,
      label: <Link to="/points">Points History</Link>,
    },
    {
      key: '/profile',
      icon: <UserOutlined />,
      label: <Link to="/profile">Profile</Link>,
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
            defaultSelectedKeys={['/products']}
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
              <Route path="/products" element={<ProductListPage />} />
              <Route path="/redemptions" element={<RedemptionHistoryPage />} />
              <Route path="/points" element={<PointsHistoryPage />} />
              <Route path="/profile" element={<ProfilePage />} />
              <Route path="/" element={<ProductListPage />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default EmployeeLayout;
