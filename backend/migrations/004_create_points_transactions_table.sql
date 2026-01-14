-- Create points_transactions table
CREATE TABLE IF NOT EXISTS points_transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    transaction_type ENUM('grant', 'deduct', 'redemption') NOT NULL COMMENT '交易类型',
    amount INT NOT NULL COMMENT '变动数量（正数为增加，负数为减少）',
    balance_after INT NOT NULL COMMENT '交易后余额',
    reason VARCHAR(500) COMMENT '原因/备注',
    operator_id BIGINT COMMENT '操作员ID',
    related_order_id BIGINT COMMENT '关联订单ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_type (transaction_type),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (operator_id) REFERENCES users(id),
    FOREIGN KEY (related_order_id) REFERENCES redemption_orders(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分交易表';
