import React, { useState, useEffect } from 'react';
import { Card, Button, Modal, Form, Input, message, Tag, Space } from 'antd';
import { PlusOutlined, StopOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import DataTable from '../../components/DataTable';
import { useConfirmDialog } from '../../components/ConfirmDialog';
import notificationService from '../../components/NotificationToast';
import adminService from '../../services/adminService';
import type { ColumnsType } from 'antd/es/table';

interface User {
  id: number;
  full_name: string;
  email: string;
  phone: string;
  role: string;
  points_balance: number;
  is_active: boolean;
  created_at: string;
}

interface CreateUserForm {
  full_name: string;
  email: string;
  phone: string;
}

const AdminUserManagementPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const { t } = useTranslation();
  const { showConfirm } = useConfirmDialog();
  const [form] = Form.useForm();

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await adminService.getUsers();
      setUsers(data);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const showCreateModal = () => {
    form.resetFields();
    setIsModalVisible(true);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
    form.resetFields();
  };

  const handleCreateUser = async (values: CreateUserForm) => {
    setSubmitting(true);
    try {
      const result = await adminService.createUser(values);
      
      notificationService.success(
        t('common.success'),
        `Initial password: ${result.initial_password}`
      );
      
      setIsModalVisible(false);
      form.resetFields();
      fetchUsers();
    } catch (error: any) {
      notificationService.error(t('common.error'), error.message);
    } finally {
      setSubmitting(false);
    }
  };

  const handleSetInactive = (user: User) => {
    showConfirm({
      title: t('user.setInactive'),
      content: `${user.full_name} (${user.email})`,
      type: 'warning',
      onOk: async () => {
        try {
          await adminService.setUserInactive(user.id);
          notificationService.success(t('common.success'));
          fetchUsers();
        } catch (error: any) {
          notificationService.error(t('common.error'), error.message);
        }
      },
    });
  };

  const columns: ColumnsType<User> = [
    {
      title: t('user.fullName'),
      dataIndex: 'full_name',
      key: 'full_name',
    },
    {
      title: t('user.email'),
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: t('user.phone'),
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: t('user.role'),
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => (
        <Tag color={role === 'admin' ? 'red' : 'blue'}>
          {role === 'admin' ? t('user.admin') : t('user.employee')}
        </Tag>
      ),
    },
    {
      title: t('points.balance'),
      dataIndex: 'points_balance',
      key: 'points_balance',
      render: (points: number) => `${points} ${t('points.title')}`,
    },
    {
      title: 'Status',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'success' : 'default'}>
          {isActive ? 'Active' : 'Inactive'}
        </Tag>
      ),
    },
    {
      title: t('common.action'),
      key: 'action',
      render: (_, record) => (
        <Space>
          {record.is_active && record.role !== 'admin' && (
            <Button
              type="link"
              danger
              icon={<StopOutlined />}
              onClick={() => handleSetInactive(record)}
            >
              {t('user.setInactive')}
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={t('nav.userManagement')}
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={showCreateModal}
          >
            {t('user.createUser')}
          </Button>
        }
      >
        <DataTable
          columns={columns}
          dataSource={users}
          loading={loading}
          rowKey="id"
        />
      </Card>

      <Modal
        title={t('user.createUser')}
        open={isModalVisible}
        onCancel={handleCancel}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateUser}
        >
          <Form.Item
            name="full_name"
            label={t('user.fullName')}
            rules={[{ required: true, message: 'Full name is required' }]}
          >
            <Input placeholder={t('user.fullName')} />
          </Form.Item>

          <Form.Item
            name="email"
            label={t('user.email')}
            rules={[
              { required: true, message: 'Email is required' },
              { type: 'email', message: 'Please enter a valid email' },
            ]}
          >
            <Input placeholder={t('user.email')} />
          </Form.Item>

          <Form.Item
            name="phone"
            label={t('user.phone')}
            rules={[
              { required: true, message: 'Phone is required' },
              { pattern: /^[0-9]{10,15}$/, message: 'Please enter a valid phone number' },
            ]}
          >
            <Input placeholder={t('user.phone')} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              {t('common.submit')}
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AdminUserManagementPage;
