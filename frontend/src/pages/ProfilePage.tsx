import React, { useState } from 'react';
import { Card, Descriptions, Button, Modal, Form, Input, message } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';
import userService from '../services/userService';

const ProfilePage: React.FC = () => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation();
  const { user, updateUser } = useAuth();
  const [form] = Form.useForm();

  const showModal = () => {
    form.setFieldsValue({ phone: user?.phone });
    setIsModalVisible(true);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
    form.resetFields();
  };

  const handleUpdatePhone = async (values: { phone: string }) => {
    setLoading(true);
    try {
      await userService.updatePhone(values.phone);
      
      // Update user context
      if (user) {
        updateUser({
          ...user,
          phone: values.phone,
        });
      }
      
      message.success(t('common.success'));
      setIsModalVisible(false);
      form.resetFields();
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  if (!user) {
    return null;
  }

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={t('nav.profile')}
        extra={
          <Button
            type="primary"
            icon={<EditOutlined />}
            onClick={showModal}
          >
            {t('user.updatePhone')}
          </Button>
        }
      >
        <Descriptions bordered column={1}>
          <Descriptions.Item label={t('user.fullName')}>
            {user.full_name}
          </Descriptions.Item>
          <Descriptions.Item label={t('user.email')}>
            {user.email}
          </Descriptions.Item>
          <Descriptions.Item label={t('user.phone')}>
            {user.phone}
          </Descriptions.Item>
          <Descriptions.Item label={t('user.role')}>
            {user.role === 'admin' ? t('user.admin') : t('user.employee')}
          </Descriptions.Item>
          {user.role === 'employee' && (
            <Descriptions.Item label={t('points.balance')}>
              {user.points_balance} {t('points.title')}
            </Descriptions.Item>
          )}
        </Descriptions>
      </Card>

      <Modal
        title={t('user.updatePhone')}
        open={isModalVisible}
        onCancel={handleCancel}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleUpdatePhone}
        >
          <Form.Item
            name="phone"
            label={t('user.phone')}
            rules={[
              { required: true, message: 'Phone number is required' },
              { pattern: /^[0-9]{10,15}$/, message: 'Please enter a valid phone number' },
            ]}
          >
            <Input placeholder={t('user.phone')} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              {t('common.save')}
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ProfilePage;
