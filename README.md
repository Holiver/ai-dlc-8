# AWSomeShop - 员工福利电商系统

一个基于 Web 的员工福利电商平台，员工可以使用积分兑换产品。

**项目状态：✅ 100% 完成，可以交付使用**

## 📚 完整文档导航

### 🚀 快速开始
- **[快速启动指南](QUICK_START_GUIDE.md)** - 详细的启动步骤和配置说明
- **[用户手册](USER_MANUAL.md)** - 员工端和管理员端完整使用指南

### 🧪 测试和部署
- **[测试清单](TESTING_CHECKLIST.md)** - 19 个详细测试用例和测试指南
- **[部署指南](DEPLOYMENT_GUIDE.md)** - Docker 和手动部署的完整方案

### 📊 项目总结
- **[项目完成总结](PROJECT_COMPLETION_SUMMARY.md)** - 项目完成情况总览（推荐阅读）
- **[项目交付总结](PROJECT_DELIVERY_SUMMARY.md)** - 详细的交付清单和统计
- **[当前进度](CURRENT_PROGRESS.md)** - 各阶段完成情况

### 📖 开发文档
- **[需求文档](.kiro/specs/awsome-shop/requirements.md)** - 19 个功能需求详细说明
- **[设计文档](.kiro/specs/awsome-shop/design.md)** - 架构设计和 C4 模型
- **[任务列表](.kiro/specs/awsome-shop/tasks.md)** - 7 个阶段的实施计划

---

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- GORM (MySQL)
- JWT Authentication
- bcrypt Password Hashing

### 前端
- React 18+
- TypeScript
- Ant Design
- React Router
- Axios
- react-i18next

### 基础设施
- MySQL 8.0
- Docker & Docker Compose
- Nginx

## 项目结构

```
awsome-shop/
├── backend/              # Go 后端 API
│   ├── cmd/             # 应用入口
│   ├── internal/        # 内部包
│   ├── pkg/             # 公共包
│   ├── configs/         # 配置文件
│   └── migrations/      # 数据库迁移
├── frontend/            # React 前端应用
│   ├── public/          # 静态资源
│   └── src/             # 源代码
├── nginx/               # Nginx 配置
├── scripts/             # 辅助脚本
└── docker-compose.yml   # Docker Compose 配置
```

## 快速开始

### 前置要求

- Docker & Docker Compose
- Go 1.21+ (开发模式)
- Node.js 18+ (开发模式)

### 使用 Docker Compose（推荐）

1. 克隆仓库
```bash
git clone <repository-url>
cd awsome-shop
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，填入实际配置
```

3. 启动服务
```bash
# 开发环境
docker-compose up -d

# 生产环境（包含 Nginx）
docker-compose --profile production up -d
```

4. 访问应用
- 前端: http://localhost
- 后端 API: http://localhost:8080
- 生产环境: https://localhost

### 开发模式

#### 后端开发

```bash
cd backend
go mod download
go run cmd/api/main.go
```

#### 前端开发

```bash
cd frontend
npm install
npm start
```

### 使用辅助脚本

```bash
# 启动开发环境
./scripts/start-dev.sh

# 启动生产环境
./scripts/start-prod.sh

# 停止所有服务
./scripts/stop-all.sh
```

## 功能特性

### 员工端
- ✅ 用户登录认证
- ✅ 产品浏览（平铺展示）
- ✅ 积分兑换产品
- ✅ 查看积分余额
- ✅ 查看兑换历史
- ✅ 查看积分变动历史
- ✅ 修改个人信息
- ✅ 中英双语支持

### 管理员端
- ✅ 员工账户管理
- ✅ 产品管理（CRUD、上下架）
- ✅ 积分发放和扣除
- ✅ 批量操作（产品导入、积分发放）
- ✅ 订单管理
- ✅ 统计报表

## API 文档

API 文档请参考 [backend/docs/API.md](backend/docs/API.md)

## 数据库

### 数据库表

- `users` - 用户表
- `products` - 产品表
- `redemption_orders` - 兑换订单表
- `points_transactions` - 积分交易表
- `product_price_history` - 产品价格历史表

### 数据库迁移

```bash
# 使用 GORM AutoMigrate（开发）
go run cmd/api/main.go

# 使用 SQL 脚本（生产）
mysql -u username -p database_name < backend/migrations/001_create_users_table.sql
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose --profile production up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 环境变量

参考 `.env.example` 文件配置以下环境变量：

- 数据库配置
- JWT 密钥
- 邮件服务配置
- IP 地理位置服务配置

## 开发指南

### 后端开发

1. 添加新功能时，遵循分层架构：
   - Model → Repository → Service → Handler
2. 使用 GORM 进行数据库操作
3. 使用 JWT 进行身份认证
4. 编写单元测试

### 前端开发

1. 使用 TypeScript 编写类型安全的代码
2. 使用 Ant Design 组件库
3. 使用 React Hooks 管理状态
4. 使用 i18next 实现国际化

## 测试

### 后端测试

```bash
cd backend
go test ./...
```

### 前端测试

```bash
cd frontend
npm test
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请联系项目维护者。
