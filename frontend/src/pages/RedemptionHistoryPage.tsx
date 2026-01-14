import React, { useState, useEffect } from 'react';
import { Card, message } from 'antd';
import { useTranslation } from 'react-i18next';
import DataTable from '../components/DataTable';
import OrderStatusBadge from '../components/OrderStatusBadge';
import redemptionService from '../services/redemptionService';
import type { ColumnsType } from 'antd/es/table';

interface RedemptionOrder {
  id: number;
  order_number: string;
  product_name: string;
  points_cost: number;
  points_balance_after: number;
  status: 'preparing' | 'delivered';
  created_at: string;
}

const RedemptionHistoryPage: React.FC = () => {
  const [orders, setOrders] = useState<RedemptionOrder[]>([]);
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation();

  useEffect(() => {
    fetchRedemptionHistory();
  }, []);

  const fetchRedemptionHistory = async () => {
    setLoading(true);
    try {
      const data = await redemptionService.getRedemptionHistory();
      // Sort by created_at descending (newest first)
      const sortedData = data.sort((a: RedemptionOrder, b: RedemptionOrder) => 
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );
      setOrders(sortedData);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const columns: ColumnsType<RedemptionOrder> = [
    {
      title: t('order.orderNumber'),
      dataIndex: 'order_number',
      key: 'order_number',
      width: 200,
    },
    {
      title: t('order.productName'),
      dataIndex: 'product_name',
      key: 'product_name',
    },
    {
      title: t('order.pointsCost'),
      dataIndex: 'points_cost',
      key: 'points_cost',
      width: 120,
      render: (points: number) => `${points} ${t('points.title')}`,
    },
    {
      title: t('order.status.label'),
      dataIndex: 'status',
      key: 'status',
      width: 120,
      render: (status: 'preparing' | 'delivered') => (
        <OrderStatusBadge status={status} />
      ),
    },
    {
      title: t('order.createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString(),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card title={t('redemptionHistory')}>
        <DataTable
          columns={columns}
          dataSource={orders}
          loading={loading}
          rowKey="id"
          pagination={{
            pageSize: 10,
            showSizeChanger: true,
            pageSizeOptions: ['10', '20', '50'],
          }}
        />
      </Card>
    </div>
  );
};

export default RedemptionHistoryPage;
