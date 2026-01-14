import React, { useState, useEffect } from 'react';
import { Row, Col, Spin, Empty, message } from 'antd';
import { useTranslation } from 'react-i18next';
import ProductCard from '../components/ProductCard';
import PointsBalance from '../components/PointsBalance';
import { useConfirmDialog } from '../components/ConfirmDialog';
import notificationService from '../components/NotificationToast';
import productService from '../services/productService';
import redemptionService from '../services/redemptionService';
import { useAuth } from '../contexts/AuthContext';

interface Product {
  id: number;
  name: string;
  image_url: string;
  points_required: number;
  stock_quantity: number;
  status: string;
}

const ProductListPage: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation();
  const { showConfirm } = useConfirmDialog();
  const { user, updateUser } = useAuth();

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const data = await productService.getProducts();
      setProducts(data);
    } catch (error: any) {
      message.error(error.message || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const handleRedeem = (product: Product) => {
    showConfirm({
      title: t('product.redeemConfirm'),
      content: `${product.name} - ${product.points_required} ${t('points.title')}`,
      onOk: async () => {
        try {
          const result = await redemptionService.createRedemption(product.id);
          
          notificationService.success(
            t('product.redeemSuccess'),
            `${t('order.orderNumber')}: ${result.order_number}`
          );
          
          // Update user points balance
          if (user) {
            updateUser({
              ...user,
              points_balance: result.points_balance_after,
            });
          }
          
          // Refresh product list to update stock
          fetchProducts();
        } catch (error: any) {
          notificationService.error(
            t('product.redeemFailed'),
            error.message
          );
        }
      },
    });
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div style={{ padding: '24px' }}>
      {/* Points Balance Card */}
      <PointsBalance
        balance={user?.points_balance || 0}
        style={{ marginBottom: 24, maxWidth: 300 }}
      />

      {/* Product Grid */}
      {products.length === 0 ? (
        <Empty description={t('common.noData')} />
      ) : (
        <Row gutter={[16, 16]}>
          {products.map((product) => (
            <Col xs={24} sm={12} md={8} lg={6} key={product.id}>
              <ProductCard
                product={product}
                onRedeem={handleRedeem}
                userPoints={user?.points_balance || 0}
              />
            </Col>
          ))}
        </Row>
      )}
    </div>
  );
};

export default ProductListPage;
