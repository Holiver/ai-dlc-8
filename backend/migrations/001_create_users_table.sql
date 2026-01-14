-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    full_name VARCHAR(100) NOT NULL COMMENT '姓名',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱（唯一标识）',
    phone VARCHAR(20) NOT NULL COMMENT '手机号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role ENUM('employee', 'admin') DEFAULT 'employee' COMMENT '角色',
    points_balance INT DEFAULT 0 COMMENT '积分余额',
    is_first_login BOOLEAN DEFAULT TRUE COMMENT '是否首次登录',
    is_active BOOLEAN DEFAULT TRUE COMMENT '账户是否有效',
    preferred_language VARCHAR(10) DEFAULT 'zh' COMMENT '首选语言',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
