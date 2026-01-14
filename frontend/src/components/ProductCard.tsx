import React from 'react';
import { Card, Button, Tag } from 'antd';
import { ShoppingCartOutlined } from '@ant-design/icons';
import { Product } from '../services/productService';
import { useTranslation } from 'react-i18next';

interface ProductCardProps {
  product: Product;
  onRedeem?: (productId: number) => void;
  loading?: boolean;
}

const ProductCard: React.FC<ProductCardProps> = ({ product, onRedeem, loading }) => {
  const { t } = useTranslation();

  const isOutOfStock = product.stock_quantity === 0;

  return (
    <Card
      hoverable
      cover={
        product.image_url ? (
          <img
            alt={product.name}
            src={product.image_url}
            style={{ height: 200, objectFit: 'cover' }}
          />
        ) : (
          <div style={{ height: 200, background: '#f0f0f0', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <span style={{ color: '#999' }}>No Image</span>
          </div>
        )
      }
      actions={[
        <Button
          type="primary"
          icon={<ShoppingCartOutlined />}
          onClick={() => onRedeem && onRedeem(product.id)}
          disabled={isOutOfStock || loading}
          block
        >
          {isOutOfStock ? t('product.outOfStock') : t('product.redeem')}
        </Button>,
      ]}
    >
      <Card.Meta
        title={product.name}
        description={
          <div>
            <div style={{ fontSize: 18, fontWeight: 'bold', color: '#1890ff', marginBottom: 8 }}>
              {product.points_required} {t('common.points')}
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span>{t('product.stock')}: {product.stock_quantity}</span>
              {isOutOfStock && <Tag color="red">{t('product.soldOut')}</Tag>}
            </div>
          </div>
        }
      />
    </Card>
  );
};

export default ProductCard;
