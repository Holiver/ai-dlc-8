import React, { useState } from 'react';
import { Card, Tabs, Button, message } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import DataTable from '../../components/DataTable';
import adminService from '../../services/adminService';
import type { ColumnsType } from 'antd/es/table';

const { TabPane } = Tabs;

interface PointsGrantReport {
  user_email: string;
  user_name: string;
  total_granted: number;
  grant_count: number;
}

interface PointsBalanceReport {
  user_email: string;
  user_name: string;
  current_balance: number;
  total_earned: number;
  total_spent: number;
}

interface RedemptionReport {
  order_number: string;
  user_email: string;
  product_name: string;
  points_cost: number;
  status: string;
  created_at: string;
}

const AdminReportsPage: React.FC = () => {
  const [grantsData, setGrantsData] = useState<PointsGrantReport[]>([]);
  const [balancesData, setBalancesData] = useState<PointsBalanceReport[]>([]);
  const [redemptionsData, setRedemptionsData] = useState<RedemptionReport[]>([]);
  const [loading, setLoading] = useState({ grants: false, balances: false, redemptions: false });
  const { t } = useTranslation();

  const fetchGrantsReport = async () => {
    setLoading(prev => ({ ...prev, grants: true }));
    try {
      const data = await adminService.getPointsGrantsReport();
      setGrantsData(data);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(prev => ({ ...prev, grants: false }));
    }
  };

  const fetchBalancesReport = async () => {
    setLoading(prev => ({ ...prev, balances: true }));
    try {
      const data = await adminService.getPointsBalancesReport();
      setBalancesData(data);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(prev => ({ ...prev, balances: false }));
    }
  };

  const fetchRedemptionsReport = async () => {
    setLoading(prev => ({ ...prev, redemptions: true }));
    try {
      const data = await adminService.getRedemptionsReport();
      setRedemptionsData(data);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(prev => ({ ...prev, redemptions: false }));
    }
  };

  const handleExport = (data: any[], filename: string) => {
    const csv = convertToCSV(data);
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `${filename}_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
  };

  const convertToCSV = (data: any[]) => {
    if (data.length === 0) return '';
    
    const headers = Object.keys(data[0]).join(',');
    const rows = data.map(row => 
      Object.values(row).map(val => `"${val}"`).join(',')
    );
    
    return [headers, ...rows].join('\n');
  };

  const grantsColumns: ColumnsType<PointsGrantReport> = [
    { title: 'User Email', dataIndex: 'user_email', key: 'user_email' },
    { title: 'User Name', dataIndex: 'user_name', key: 'user_name' },
    { 
      title: 'Total Granted', 
      dataIndex: 'total_granted', 
      key: 'total_granted',
      render: (val: number) => `${val} ${t('points.title')}`,
    },
    { title: 'Grant Count', dataIndex: 'grant_count', key: 'grant_count' },
  ];

  const balancesColumns: ColumnsType<PointsBalanceReport> = [
    { title: 'User Email', dataIndex: 'user_email', key: 'user_email' },
    { title: 'User Name', dataIndex: 'user_name', key: 'user_name' },
    { 
      title: 'Current Balance', 
      dataIndex: 'current_balance', 
      key: 'current_balance',
      render: (val: number) => `${val} ${t('points.title')}`,
    },
    { 
      title: 'Total Earned', 
      dataIndex: 'total_earned', 
      key: 'total_earned',
      render: (val: number) => `${val} ${t('points.title')}`,
    },
    { 
      title: 'Total Spent', 
      dataIndex: 'total_spent', 
      key: 'total_spent',
      render: (val: number) => `${val} ${t('points.title')}`,
    },
  ];

  const redemptionsColumns: ColumnsType<RedemptionReport> = [
    { title: t('order.orderNumber'), dataIndex: 'order_number', key: 'order_number', width: 200 },
    { title: 'User Email', dataIndex: 'user_email', key: 'user_email' },
    { title: t('order.productName'), dataIndex: 'product_name', key: 'product_name' },
    { 
      title: t('order.pointsCost'), 
      dataIndex: 'points_cost', 
      key: 'points_cost',
      render: (val: number) => `${val} ${t('points.title')}`,
    },
    { title: 'Status', dataIndex: 'status', key: 'status' },
    { 
      title: t('order.createdAt'), 
      dataIndex: 'created_at', 
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString(),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card title={t('nav.reports')}>
        <Tabs defaultActiveKey="grants">
          <TabPane tab="Points Grants Report" key="grants">
            <div style={{ marginBottom: 16 }}>
              <Button onClick={fetchGrantsReport} loading={loading.grants}>
                Load Report
              </Button>
              <Button
                icon={<DownloadOutlined />}
                onClick={() => handleExport(grantsData, 'points_grants')}
                disabled={grantsData.length === 0}
                style={{ marginLeft: 8 }}
              >
                {t('admin.exportReport')}
              </Button>
            </div>
            <DataTable
              columns={grantsColumns}
              dataSource={grantsData}
              loading={loading.grants}
              rowKey="user_email"
            />
          </TabPane>

          <TabPane tab="Points Balances Report" key="balances">
            <div style={{ marginBottom: 16 }}>
              <Button onClick={fetchBalancesReport} loading={loading.balances}>
                Load Report
              </Button>
              <Button
                icon={<DownloadOutlined />}
                onClick={() => handleExport(balancesData, 'points_balances')}
                disabled={balancesData.length === 0}
                style={{ marginLeft: 8 }}
              >
                {t('admin.exportReport')}
              </Button>
            </div>
            <DataTable
              columns={balancesColumns}
              dataSource={balancesData}
              loading={loading.balances}
              rowKey="user_email"
            />
          </TabPane>

          <TabPane tab="Redemptions Report" key="redemptions">
            <div style={{ marginBottom: 16 }}>
              <Button onClick={fetchRedemptionsReport} loading={loading.redemptions}>
                Load Report
              </Button>
              <Button
                icon={<DownloadOutlined />}
                onClick={() => handleExport(redemptionsData, 'redemptions')}
                disabled={redemptionsData.length === 0}
                style={{ marginLeft: 8 }}
              >
                {t('admin.exportReport')}
              </Button>
            </div>
            <DataTable
              columns={redemptionsColumns}
              dataSource={redemptionsData}
              loading={loading.redemptions}
              rowKey="order_number"
            />
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default AdminReportsPage;
