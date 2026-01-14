import React, { useState } from 'react';
import { Card, Tabs, Form, Select, InputNumber, Input, Button, message } from 'antd';
import { useTranslation } from 'react-i18next';
import MarkdownTableInput from '../../components/MarkdownTableInput';
import notificationService from '../../components/NotificationToast';
import adminService from '../../services/adminService';

const { TabPane } = Tabs;
const { TextArea } = Input;

interface GrantPointsForm {
  user_email: string;
  amount: number;
  reason: string;
}

interface DeductPointsForm {
  user_email: string;
  amount: number;
  reason: string;
}

const AdminPointsManagementPage: React.FC = () => {
  const [grantLoading, setGrantLoading] = useState(false);
  const [deductLoading, setDeductLoading] = useState(false);
  const [batchLoading, setBatchLoading] = useState(false);
  const { t } = useTranslation();
  const [grantForm] = Form.useForm();
  const [deductForm] = Form.useForm();

  const handleGrantPoints = async (values: GrantPointsForm) => {
    setGrantLoading(true);
    try {
      await adminService.grantPoints(values);
      notificationService.success(t('common.success'));
      grantForm.resetFields();
    } catch (error: any) {
      notificationService.error(t('common.error'), error.message);
    } finally {
      setGrantLoading(false);
    }
  };

  const handleDeductPoints = async (values: DeductPointsForm) => {
    setDeductLoading(true);
    try {
      await adminService.deductPoints(values);
      notificationService.success(t('common.success'));
      deductForm.resetFields();
    } catch (error: any) {
      notificationService.error(t('common.error'), error.message);
    } finally {
      setDeductLoading(false);
    }
  };

  const handleBatchGrant = async (markdownData: string) => {
    setBatchLoading(true);
    try {
      await adminService.batchGrantPoints(markdownData);
      notificationService.success(t('common.success'));
    } catch (error: any) {
      notificationService.error(t('common.error'), error.message);
    } finally {
      setBatchLoading(false);
    }
  };

  const exampleMarkdown = `| Email | Amount | Reason |
|-------|--------|--------|
| user1@example.com | 100 | Monthly bonus |
| user2@example.com | 200 | Performance reward |
| user3@example.com | 150 | Project completion |`;

  return (
    <div style={{ padding: '24px' }}>
      <Card title={t('nav.pointsManagement')}>
        <Tabs defaultActiveKey="grant">
          <TabPane tab={t('admin.grantPoints')} key="grant">
            <Form
              form={grantForm}
              layout="vertical"
              onFinish={handleGrantPoints}
              style={{ maxWidth: 600 }}
            >
              <Form.Item
                name="user_email"
                label={t('user.email')}
                rules={[
                  { required: true, message: 'Email is required' },
                  { type: 'email', message: 'Please enter a valid email' },
                ]}
              >
                <Input placeholder={t('user.email')} />
              </Form.Item>

              <Form.Item
                name="amount"
                label={t('common.amount')}
                rules={[
                  { required: true, message: 'Amount is required' },
                  { type: 'number', min: 1, message: 'Amount must be greater than 0' },
                ]}
              >
                <InputNumber
                  min={1}
                  style={{ width: '100%' }}
                  placeholder="Enter points amount"
                />
              </Form.Item>

              <Form.Item
                name="reason"
                label={t('common.reason')}
                rules={[{ required: true, message: 'Reason is required' }]}
              >
                <TextArea
                  rows={3}
                  placeholder="Enter reason for granting points"
                />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={grantLoading} block>
                  {t('admin.grantPoints')}
                </Button>
              </Form.Item>
            </Form>
          </TabPane>

          <TabPane tab={t('admin.deductPoints')} key="deduct">
            <Form
              form={deductForm}
              layout="vertical"
              onFinish={handleDeductPoints}
              style={{ maxWidth: 600 }}
            >
              <Form.Item
                name="user_email"
                label={t('user.email')}
                rules={[
                  { required: true, message: 'Email is required' },
                  { type: 'email', message: 'Please enter a valid email' },
                ]}
              >
                <Input placeholder={t('user.email')} />
              </Form.Item>

              <Form.Item
                name="amount"
                label={t('common.amount')}
                rules={[
                  { required: true, message: 'Amount is required' },
                  { type: 'number', min: 1, message: 'Amount must be greater than 0' },
                ]}
              >
                <InputNumber
                  min={1}
                  style={{ width: '100%' }}
                  placeholder="Enter points amount"
                />
              </Form.Item>

              <Form.Item
                name="reason"
                label={t('common.reason')}
                rules={[{ required: true, message: 'Reason is required' }]}
              >
                <TextArea
                  rows={3}
                  placeholder="Enter reason for deducting points"
                />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={deductLoading} block danger>
                  {t('admin.deductPoints')}
                </Button>
              </Form.Item>
            </Form>
          </TabPane>

          <TabPane tab={t('admin.batchGrant')} key="batch">
            <MarkdownTableInput
              onSubmit={handleBatchGrant}
              placeholder="Enter batch grant data in Markdown table format"
              exampleText={exampleMarkdown}
              loading={batchLoading}
            />
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default AdminPointsManagementPage;
