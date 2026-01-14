# AWSomeShop 快速启动指南

## 📋 前置要求

### 必需软件
- Go 1.21+ 
- Node.js 16+ 和 npm
- MySQL 8.0+
- Docker 和 Docker Compose（可选，用于容器化部署）

### 检查安装
```bash
go version
node --version
npm --version
mysql --version
docker --version
docker-compose --version
```

## 🚀 快速启动（开发环境）

### 方式 1：手动启动（推荐用于开发）

#### 步骤 1：启动 MySQL 数据库

**选项 A：使用 Docker**
```bash
docker run -d \
  --name awsomeshop-mysql \
  -e MYSQL_ROOT_PASSWORD=root123 \
  -e MYSQL_DATABASE=awsomeshop \
  -p 3306:3306 \
  mysql:8.0
```

**选项 B：使用本地 MySQL**
```sql
CREATE DATABASE awsomeshop CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 步骤 2：配置后端

1. 进入后端目录
```bash
cd backend
```

2. 复制配置文件
```bash
cp configs/config.yaml.example configs/config.yaml
```

3. 编辑配置文件 `configs/config.yaml`
```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  user: root
  password: root123
  dbname: awsomeshop
  
jwt:
  secret: your-secret-key-change-this-in-production
  expiration: 24h
```

4. 安装依赖
```bash
go mod download
```

5. 运行数据库迁移
```bash
# 方式 1：使用 migrate 工具
migrate -path migrations -database "mysql://root:root123@tcp(localhost:3306)/awsomeshop" up

# 方式 2：手动执行 SQL 文件
mysql -u root -p awsomeshop < migrations/001_create_tables.sql
```

6. 启动后端服务器
```bash
go run cmd/api/main.go
```

后端应该在 `http://localhost:8080` 运行

#### 步骤 3：配置前端

1. 打开新终端，进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 创建环境配置文件
```bash
cp .env.example .env
```

4. 编辑 `.env` 文件
```env
REACT_APP_API_URL=http://localhost:8080
```

5. 启动前端开发服务器
```bash
npm start
```

前端应该在 `http://localhost:3000` 运行，浏览器会自动打开

### 方式 2：使用 Docker Compose（一键启动）

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down
```

服务地址：
- 前端：http://localhost
- 后端 API：http://localhost/api
- MySQL：localhost:3306

## 👤 测试账户

### 创建管理员账户

启动后端后，你需要手动创建第一个管理员账户：

```sql
-- 连接到数据库
mysql -u root -p awsomeshop

-- 插入管理员账户（密码：admin123）
INSERT INTO users (full_name, email, phone, password_hash, role, is_first_login, is_active) 
VALUES (
  'Admin User',
  'admin@awsomeshop.com',
  '1234567890',
  '$2a$10$YourBcryptHashHere',  -- 需要使用 bcrypt 加密
  'admin',
  FALSE,
  TRUE
);
```

**生成 bcrypt 密码哈希：**

使用 Go 生成：
```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "admin123"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
```

或使用在线工具：https://bcrypt-generator.com/

### 创建测试员工账户

**方式 1：通过管理员界面**
1. 使用管理员账户登录
2. 进入"员工管理"页面
3. 点击"创建员工账户"
4. 填写信息并提交
5. 记录生成的初始密码

**方式 2：通过 API**
```bash
# 先登录获取 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@awsomeshop.com",
    "password": "admin123"
  }'

# 使用 token 创建员工
curl -X POST http://localhost:8080/api/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "full_name": "Test Employee",
    "email": "employee@awsomeshop.com",
    "phone": "9876543210"
  }'
```

## 🧪 测试功能

### 1. 测试员工端功能

1. **登录**
   - 访问 http://localhost:3000/login
   - 使用员工账户登录
   - 首次登录应该自动获得 1000 积分

2. **浏览产品**
   - 应该重定向到产品列表页面
   - 查看积分余额
   - 浏览可用产品

3. **兑换产品**
   - 点击产品的"兑换"按钮
   - 确认兑换
   - 查看积分余额更新

4. **查看历史**
   - 查看兑换历史
   - 查看积分历史
   - 验证数据正确显示

5. **个人信息**
   - 查看个人信息
   - 修改手机号
   - 验证更新成功

### 2. 测试管理员端功能

1. **登录**
   - 使用管理员账户登录
   - 应该重定向到管理员仪表板

2. **员工管理**
   - 创建新员工账户
   - 记录初始密码
   - 设置员工离职状态

3. **产品管理**
   - 创建新产品
   - 编辑产品信息
   - 上下架产品
   - 批量导入产品（Markdown 表格）

4. **积分管理**
   - 发放积分给单个员工
   - 扣除积分
   - 批量发放积分（Markdown 表格）

5. **订单管理**
   - 查看所有订单
   - 筛选订单状态
   - 批量更新订单状态为"已发放"

6. **统计报表**
   - 查看积分发放表
   - 查看积分存量表
   - 查看兑换记录表
   - 导出 CSV 报表

### 3. 测试权限控制

1. **员工访问管理员页面**
   - 使用员工账户登录
   - 尝试访问 `/admin/dashboard`
   - 应该被重定向到 `/products`

2. **未登录访问受保护页面**
   - 登出
   - 尝试访问 `/products` 或 `/admin/dashboard`
   - 应该被重定向到 `/login`

## 📝 API 测试

### 使用 Postman 测试

1. **导入 API 集合**（如果有）
2. **设置环境变量**
   - `base_url`: http://localhost:8080
   - `token`: 登录后获取的 JWT token

### 常用 API 端点

#### 认证
```bash
# 登录
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password"
}

# 获取当前用户信息
GET /api/v1/auth/me
Authorization: Bearer {token}
```

#### 产品
```bash
# 获取产品列表
GET /api/v1/products
Authorization: Bearer {token}
```

#### 兑换
```bash
# 创建兑换订单
POST /api/v1/redemptions
Authorization: Bearer {token}
{
  "product_id": 1
}

# 获取兑换历史
GET /api/v1/redemptions
Authorization: Bearer {token}
```

#### 积分
```bash
# 获取积分余额
GET /api/v1/points/balance
Authorization: Bearer {token}

# 获取积分历史
GET /api/v1/points/transactions?page=1&page_size=10
Authorization: Bearer {token}
```

## 🐛 常见问题

### 问题 1：后端无法连接数据库

**错误信息：**
```
Error connecting to database: dial tcp 127.0.0.1:3306: connect: connection refused
```

**解决方案：**
1. 确保 MySQL 正在运行
2. 检查配置文件中的数据库连接信息
3. 验证数据库用户名和密码

### 问题 2：前端无法连接后端

**错误信息：**
```
Network Error
```

**解决方案：**
1. 确保后端服务器正在运行
2. 检查 `.env` 文件中的 API URL
3. 检查 CORS 配置

### 问题 3：JWT token 无效

**错误信息：**
```
401 Unauthorized
```

**解决方案：**
1. 重新登录获取新 token
2. 检查 token 是否过期
3. 确保 Authorization header 格式正确：`Bearer {token}`

### 问题 4：首次登录没有获得积分

**解决方案：**
1. 检查用户的 `is_first_login` 字段
2. 查看后端日志
3. 验证 AuthService 的首次登录逻辑

### 问题 5：前端页面空白

**解决方案：**
1. 打开浏览器开发者工具查看错误
2. 检查控制台是否有 JavaScript 错误
3. 确保所有依赖都已安装：`npm install`
4. 清除缓存并重新启动：`npm start`

## 📊 监控和日志

### 后端日志

后端日志会输出到控制台，包括：
- HTTP 请求日志
- 数据库查询日志
- 错误日志

### 前端日志

打开浏览器开发者工具（F12）查看：
- Console：JavaScript 错误和日志
- Network：API 请求和响应
- Application：本地存储和 token

### 数据库日志

```bash
# 查看 MySQL 日志
docker logs awsomeshop-mysql

# 或查看本地 MySQL 日志
tail -f /var/log/mysql/error.log
```

## 🔧 开发工具

### 推荐的 VS Code 扩展

**Go 开发：**
- Go (golang.go)
- Go Test Explorer

**React 开发：**
- ES7+ React/Redux/React-Native snippets
- Prettier - Code formatter
- ESLint

**通用：**
- GitLens
- Docker
- REST Client

### 数据库管理工具

- MySQL Workbench
- DBeaver
- phpMyAdmin
- TablePlus

## 📚 下一步

1. **完成测试**
   - 测试所有功能
   - 记录发现的 bug
   - 修复问题

2. **性能优化**
   - 添加数据库索引
   - 优化查询
   - 添加缓存

3. **安全加固**
   - 更改默认密钥
   - 配置 HTTPS
   - 添加速率限制

4. **准备部署**
   - 配置生产环境
   - 设置 CI/CD
   - 准备部署文档

## 🎉 成功标志

如果你看到以下内容，说明启动成功：

✅ 后端服务器运行在 http://localhost:8080
✅ 前端应用运行在 http://localhost:3000
✅ 可以成功登录
✅ 可以浏览产品
✅ 可以进行兑换
✅ 管理员可以管理数据

**恭喜！你的 AWSomeShop 应用已经成功运行！** 🚀

## 📞 获取帮助

如果遇到问题：
1. 查看本文档的"常见问题"部分
2. 检查后端和前端的日志
3. 查看项目的其他文档（CURRENT_PROGRESS.md, PROJECT_STATUS.md）
4. 检查 GitHub Issues（如果有）
