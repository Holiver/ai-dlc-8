import React, { useState, useEffect } from 'react';
import { Card, Tag, message } from 'antd';
import { ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import DataTable from '../components/DataTable';
import pointsService from '../services/pointsService';
import type { ColumnsType } from 'antd/es/table';

interface PointsTransaction {
  id: number;
  transaction_type: 'grant' | 'deduct' | 'redemption';
  amount: number;
  balance_after: number;
  reason: string;
  created_at: string;
}

const PointsHistoryPage: React.FC = () => {
  const [transactions, setTransactions] = useState<PointsTransaction[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const { t } = useTranslation();

  useEffect(() => {
    fetchPointsHistory(pagination.current, pagination.pageSize);
  }, []);

  const fetchPointsHistory = async (page: number, pageSize: number) => {
    setLoading(true);
    try {
      const data = await pointsService.getPointsHistory(page, pageSize);
      setTransactions(data.transactions || data);
      
      // Update pagination if total is provided
      if (data.total !== undefined) {
        setPagination({
          current: page,
          pageSize,
          total: data.total,
        });
      }
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const handleTableChange = (newPagination: any) => {
    fetchPointsHistory(newPagination.current, newPagination.pageSize);
  };

  const getTransactionTypeTag = (type: string) => {
    const typeConfig: Record<string, { color: string; text: string; icon: React.ReactNode }> = {
      grant: {
        color: 'success',
        text: t('points.grant'),
        icon: <ArrowUpOutlined />,
      },
      deduct: {
        color: 'error',
        text: t('points.deduct'),
        icon: <ArrowDownOutlined />,
      },
      redemption: {
        color: 'warning',
        text: t('points.redemption'),
        icon: <ArrowDownOutlined />,
      },
    };

    const config = typeConfig[type] || typeConfig.grant;
    return (
      <Tag color={config.color} icon={config.icon}>
        {config.text}
      </Tag>
    );
  };

  const columns: ColumnsType<PointsTransaction> = [
    {
      title: t('points.transaction'),
      dataIndex: 'transaction_type',
      key: 'transaction_type',
      width: 120,
      render: (type: string) => getTransactionTypeTag(type),
    },
    {
      title: t('common.amount'),
      dataIndex: 'amount',
      key: 'amount',
      width: 120,
      render: (amount: number) => (
        <span style={{ color: amount > 0 ? '#52c41a' : '#ff4d4f' }}>
          {amount > 0 ? '+' : ''}{amount}
        </span>
      ),
    },
    {
      title: t('points.balance'),
      dataIndex: 'balance_after',
      key: 'balance_after',
      width: 120,
      render: (balance: number) => `${balance} ${t('points.title')}`,
    },
    {
      title: t('common.reason'),
      dataIndex: 'reason',
      key: 'reason',
      ellipsis: true,
    },
    {
      title: t('common.time'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString(),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card title={t('points.history')}>
        <DataTable
          columns={columns}
          dataSource={transactions}
          loading={loading}
          rowKey="id"
          pagination={pagination}
          onChange={handleTableChange}
        />
      </Card>
    </div>
  );
};

export default PointsHistoryPage;
