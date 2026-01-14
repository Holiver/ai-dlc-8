import React from 'react';
import { Card, Statistic } from 'antd';
import { TrophyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

interface PointsBalanceProps {
  balance: number;
  loading?: boolean;
  style?: React.CSSProperties;
}

const PointsBalance: React.FC<PointsBalanceProps> = ({ balance, loading = false, style }) => {
  const { t } = useTranslation();

  return (
    <Card style={style}>
      <Statistic
        title={t('points.balance')}
        value={balance}
        prefix={<TrophyOutlined />}
        loading={loading}
        valueStyle={{ color: '#3f8600' }}
      />
    </Card>
  );
};

export default PointsBalance;
