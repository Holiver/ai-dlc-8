-- Create product_price_history table
CREATE TABLE IF NOT EXISTS product_price_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL COMMENT '产品ID',
    old_points INT COMMENT '旧积分价格',
    new_points INT NOT NULL COMMENT '新积分价格',
    operator_id BIGINT COMMENT '操作员ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_product (product_id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (operator_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品价格历史表';
