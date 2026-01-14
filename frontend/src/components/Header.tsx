import React from 'react';
import { Layout, Space, Typography, Button } from 'antd';
import { GlobalOutlined } from '@ant-design/icons';
import { useAuth } from '../contexts/AuthContext';
import { useTranslation } from 'react-i18next';

const { Header: AntHeader } = Layout;
const { Text } = Typography;

const Header: React.FC = () => {
  const { user } = useAuth();
  const { t, i18n } = useTranslation();

  const toggleLanguage = () => {
    const newLang = i18n.language === 'zh' ? 'en' : 'zh';
    i18n.changeLanguage(newLang);
  };

  return (
    <AntHeader
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        background: '#001529',
        padding: '0 24px',
      }}
    >
      <Typography.Title level={3} style={{ color: '#fff', margin: 0 }}>
        AWSomeShop
      </Typography.Title>
      
      <Space size="large">
        {user && user.role === 'employee' && (
          <Text style={{ color: '#fff' }}>
            {t('points')}: {user.points_balance}
          </Text>
        )}
        
        <Text style={{ color: '#fff' }}>
          {user?.full_name}
        </Text>
        
        <Button
          type="text"
          icon={<GlobalOutlined />}
          onClick={toggleLanguage}
          style={{ color: '#fff' }}
        >
          {i18n.language === 'zh' ? 'EN' : '中文'}
        </Button>
      </Space>
    </AntHeader>
  );
};

export default Header;
