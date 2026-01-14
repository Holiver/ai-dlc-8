import React, { useState, useEffect } from 'react';
import { Card, Button, Select, Space, message, Input } from 'antd';
import { CheckOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import DataTable from '../../components/DataTable';
import OrderStatusBadge from '../../components/OrderStatusBadge';
import { useConfirmDialog } from '../../components/ConfirmDialog';
import notificationService from '../../components/NotificationToast';
import adminService from '../../services/adminService';
import type { ColumnsType } from 'antd/es/table';

const { TextArea } = Input;

interface RedemptionOrder {
  id: number;
  order_number: string;
  user_id: number;
  user_email?: string;
  product_name: string;
  points_cost: number;
  status: 'preparing' | 'delivered';
  created_at: string;
}

const AdminOrderManagementPage: React.FC = () => {
  const [orders, setOrders] = useState<RedemptionOrder[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const { t } = useTranslation();
  const { showConfirm } = useConfirmDialog();

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const data = await adminService.getOrders();
      setOrders(data);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const handleBatchUpdateStatus = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('Please select orders to update');
      return;
    }

    const orderNumbers = orders
      .filter(order => selectedRowKeys.includes(order.id))
      .map(order => order.order_number);

    showConfirm({
      title: t('admin.batchUpdate'),
      content: (
        <div>
          <p>Update {orderNumbers.length} orders to "Delivered" status?</p>
          <TextArea
            rows={4}
            value={orderNumbers.join('\n')}
            readOnly
            style={{ marginTop: 8 }}
          />
        </div>
      ),
      onOk: async () => {
        try {
          await adminService.batchUpdateOrderStatus(orderNumbers, 'delivered');
          notificationService.success(t('common.success'));
          setSelectedRowKeys([]);
          fetchOrders();
        } catch (error: any) {
          notificationService.error(t('common.error'), error.message);
        }
      },
    });
  };

  const filteredOrders = statusFilter === 'all'
    ? orders
    : orders.filter(order => order.status === statusFilter);

  const columns: ColumnsType<RedemptionOrder> = [
    {
      title: t('order.orderNumber'),
      dataIndex: 'order_number',
      key: 'order_number',
      width: 200,
    },
    {
      title: 'User Email',
      dataIndex: 'user_email',
      key: 'user_email',
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

  const rowSelection = {
    selectedRowKeys,
    onChange: (newSelectedRowKeys: React.Key[]) => {
      setSelectedRowKeys(newSelectedRowKeys);
    },
    getCheckboxProps: (record: RedemptionOrder) => ({
      disabled: record.status === 'delivered',
    }),
  };

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={t('nav.orderManagement')}
        extra={
          <Space>
            <Select
              value={statusFilter}
              onChange={setStatusFilter}
              style={{ width: 150 }}
            >
              <Select.Option value="all">All Status</Select.Option>
              <Select.Option value="preparing">{t('order.status.preparing')}</Select.Option>
              <Select.Option value="delivered">{t('order.status.delivered')}</Select.Option>
            </Select>
            <Button
              type="primary"
              icon={<CheckOutlined />}
              onClick={handleBatchUpdateStatus}
              disabled={selectedRowKeys.length === 0}
            >
              {t('admin.batchUpdate')} ({selectedRowKeys.length})
            </Button>
          </Space>
        }
      >
        <DataTable
          columns={columns}
          dataSource={filteredOrders}
          loading={loading}
          rowKey="id"
          rowSelection={rowSelection}
        />
      </Card>
    </div>
  );
};

export default AdminOrderManagementPage;
