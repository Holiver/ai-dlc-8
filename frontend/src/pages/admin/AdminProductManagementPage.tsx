import React, { useState, useEffect } from 'react';
import { Card, Button, Modal, Form, Input, InputNumber, message, Tag, Space, Tabs } from 'antd';
import { PlusOutlined, EditOutlined, UploadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import DataTable from '../../components/DataTable';
import MarkdownTableInput from '../../components/MarkdownTableInput';
import { useConfirmDialog } from '../../components/ConfirmDialog';
import notificationService from '../../components/NotificationToast';
import adminService from '../../services/adminService';
import type { ColumnsType } from 'antd/es/table';

interface Product {
  id: number;
  name: string;
  image_url: string;
  points_required: number;
  stock_quantity: number;
  status: string;
  created_at: string;
}

interface ProductForm {
  name: string;
  image_url?: string;
  points_required: number;
  stock_quantity: number;
}

const AdminProductManagementPage: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isBatchModalVisible, setIsBatchModalVisible] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const { t } = useTranslation();
  const { showConfirm } = useConfirmDialog();
  const [form] = Form.useForm();

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const data = await adminService.getProducts();
      setProducts(data);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const showCreateModal = () => {
    setEditingProduct(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  const showEditModal = (product: Product) => {
    setEditingProduct(product);
    form.setFieldsValue(product);
    setIsModalVisible(true);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
    setEditingProduct(null);
    form.resetFields();
  };

  const handleSubmit = async (values: ProductForm) => {
    setSubmitting(true);
    try {
      if (editingProduct) {
        await adminService.updateProduct(editingProduct.id, values);
      } else {
        await adminService.createProduct(values);
      }
      
      notificationService.success(t('common.success'));
      setIsModalVisible(false);
      form.resetFields();
      setEditingProduct(null);
      fetchProducts();
    } catch (error: any) {
      notificationService.error(t('common.error'), error.message);
    } finally {
      setSubmitting(false);
    }
  };

  const handleToggleStatus = (product: Product) => {
    const newStatus = product.status === 'active' ? 'inactive' : 'active';
    const action = newStatus === 'active' ? 'Activate' : 'Deactivate';
    
    showConfirm({
      title: `${action} Product`,
      content: product.name,
      onOk: async () => {
        try {
          await adminService.updateProductStatus(product.id, newStatus);
          notificationService.success(t('common.success'));
          fetchProducts();
        } catch (error: any) {
          notificationService.error(t('common.error'), error.message);
        }
      },
    });
  };

  const handleBatchImport = async (markdownData: string) => {
    setSubmitting(true);
    try {
      await adminService.batchImportProducts(markdownData);
      notificationService.success(t('common.success'));
      setIsBatchModalVisible(false);
      fetchProducts();
    } catch (error: any) {
      notificationService.error(t('common.error'), error.message);
    } finally {
      setSubmitting(false);
    }
  };

  const columns: ColumnsType<Product> = [
    {
      title: t('product.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('product.points'),
      dataIndex: 'points_required',
      key: 'points_required',
      render: (points: number) => `${points} ${t('points.title')}`,
    },
    {
      title: t('product.stock'),
      dataIndex: 'stock_quantity',
      key: 'stock_quantity',
    },
    {
      title: t('product.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'success' : 'default'}>
          {status === 'active' ? t('product.active') : t('product.inactive')}
        </Tag>
      ),
    },
    {
      title: t('common.action'),
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => showEditModal(record)}
          >
            {t('common.edit')}
          </Button>
          <Button
            type="link"
            onClick={() => handleToggleStatus(record)}
          >
            {record.status === 'active' ? 'Deactivate' : 'Activate'}
          </Button>
        </Space>
      ),
    },
  ];

  const exampleMarkdown = `| Name | Points | Stock | Image URL |
|------|--------|-------|-----------|
| Product A | 100 | 50 | https://example.com/a.jpg |
| Product B | 200 | 30 | https://example.com/b.jpg |`;

  return (
    <div style={{ padding: '24px' }}>
      <Card
        title={t('nav.productManagement')}
        extra={
          <Space>
            <Button
              icon={<UploadOutlined />}
              onClick={() => setIsBatchModalVisible(true)}
            >
              {t('admin.batchImport')}
            </Button>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={showCreateModal}
            >
              {t('admin.createProduct')}
            </Button>
          </Space>
        }
      >
        <DataTable
          columns={columns}
          dataSource={products}
          loading={loading}
          rowKey="id"
        />
      </Card>

      <Modal
        title={editingProduct ? t('admin.updateProduct') : t('admin.createProduct')}
        open={isModalVisible}
        onCancel={handleCancel}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="name"
            label={t('product.name')}
            rules={[{ required: true, message: 'Product name is required' }]}
          >
            <Input placeholder={t('product.name')} />
          </Form.Item>

          <Form.Item
            name="image_url"
            label={t('product.image')}
          >
            <Input placeholder="https://example.com/image.jpg" />
          </Form.Item>

          <Form.Item
            name="points_required"
            label={t('product.points')}
            rules={[{ required: true, message: 'Points required is required' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            name="stock_quantity"
            label={t('product.stock')}
            rules={[{ required: true, message: 'Stock quantity is required' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              {t('common.submit')}
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={t('admin.batchImport')}
        open={isBatchModalVisible}
        onCancel={() => setIsBatchModalVisible(false)}
        footer={null}
        width={800}
      >
        <MarkdownTableInput
          onSubmit={handleBatchImport}
          placeholder="Enter product data in Markdown table format"
          exampleText={exampleMarkdown}
          loading={submitting}
        />
      </Modal>
    </div>
  );
};

export default AdminProductManagementPage;
