import React from 'react';
import { Tag } from 'antd';
import { ClockCircleOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

interface OrderStatusBadgeProps {
  status: 'preparing' | 'delivered';
}

const OrderStatusBadge: React.FC<OrderStatusBadgeProps> = ({ status }) => {
  const { t } = useTranslation();

  const statusConfig = {
    preparing: {
      color: 'processing',
      icon: <ClockCircleOutlined />,
      text: t('order.status.preparing'),
    },
    delivered: {
      color: 'success',
      icon: <CheckCircleOutlined />,
      text: t('order.status.delivered'),
    },
  };

  const config = statusConfig[status];

  return (
    <Tag color={config.color} icon={config.icon}>
      {config.text}
    </Tag>
  );
};

export default OrderStatusBadge;
