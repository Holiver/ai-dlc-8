import React from 'react';
import { Card, Row, Col, Statistic } from 'antd';
import {
  UserOutlined,
  ShoppingOutlined,
  TrophyOutlined,
  ShoppingCartOutlined,
  BarChartOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const AdminDashboardPage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const quickLinks = [
    {
      title: t('nav.userManagement'),
      icon: <UserOutlined style={{ fontSize: 48, color: '#1890ff' }} />,
      path: '/admin/users',
      description: 'Manage employee accounts',
    },
    {
      title: t('nav.productManagement'),
      icon: <ShoppingOutlined style={{ fontSize: 48, color: '#52c41a' }} />,
      path: '/admin/products',
      description: 'Manage products and inventory',
    },
    {
      title: t('nav.pointsManagement'),
      icon: <TrophyOutlined style={{ fontSize: 48, color: '#faad14' }} />,
      path: '/admin/points',
      description: 'Grant and manage points',
    },
    {
      title: t('nav.orderManagement'),
      icon: <ShoppingCartOutlined style={{ fontSize: 48, color: '#722ed1' }} />,
      path: '/admin/orders',
      description: 'Manage redemption orders',
    },
    {
      title: t('nav.reports'),
      icon: <BarChartOutlined style={{ fontSize: 48, color: '#eb2f96' }} />,
      path: '/admin/reports',
      description: 'View statistics and reports',
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={
          <span style={{ fontSize: 24 }}>
            {t('nav.dashboard')}
          </span>
        }
        style={{ marginBottom: 24 }}
      >
        <Row gutter={[16, 16]}>
          {quickLinks.map((link, index) => (
            <Col xs={24} sm={12} md={8} lg={8} key={index}>
              <Card
                hoverable
                onClick={() => navigate(link.path)}
                style={{
                  textAlign: 'center',
                  height: '100%',
                  cursor: 'pointer',
                }}
              >
                <div style={{ marginBottom: 16 }}>
                  {link.icon}
                </div>
                <Card.Meta
                  title={<span style={{ fontSize: 18 }}>{link.title}</span>}
                  description={link.description}
                />
              </Card>
            </Col>
          ))}
        </Row>
      </Card>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Total Users"
              value={0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Active Products"
              value={0}
              prefix={<ShoppingOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Total Points Granted"
              value={0}
              prefix={<TrophyOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Total Orders"
              value={0}
              prefix={<ShoppingCartOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default AdminDashboardPage;
