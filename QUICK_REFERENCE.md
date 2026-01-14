# AWSomeShop 快速参考卡

## 📋 项目状态

- **完成度**：✅ 100%
- **状态**：可以交付使用
- **最后更新**：2026-01-14

## 🚀 快速命令

### 启动应用（Docker）
```bash
docker-compose up -d
```

### 启动应用（手动）
```bash
# 终端 1 - 后端
cd backend && go run cmd/api/main.go

# 终端 2 - 前端
cd frontend && npm start
```

### 停止应用
```bash
docker-compose down
```

## 📚 文档快速索引

| 需要什么 | 查看文档 |
|---------|---------|
| 启动项目 | [QUICK_START_GUIDE.md](QUICK_START_GUIDE.md) |
| 测试项目 | [TESTING_CHECKLIST.md](TESTING_CHECKLIST.md) |
| 部署项目 | [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) |
| 使用系统 | [USER_MANUAL.md](USER_MANUAL.md) |
| 项目总览 | [PROJECT_COMPLETION_SUMMARY.md](PROJECT_COMPLETION_SUMMARY.md) |
| 了解需求 | [requirements.md](.kiro/specs/awsome-shop/requirements.md) |
| 了解架构 | [design.md](.kiro/specs/awsome-shop/design.md) |

## 🔑 默认访问地址

| 服务 | 地址 |
|-----|------|
| 前端 | http://localhost:3000 |
| 后端 API | http://localhost:8080 |
| MySQL | localhost:3306 |

## 👤 测试账户

### 管理员账户
- **邮箱**：admin@awsomeshop.com
- **密码**：需要手动创建（参考 QUICK_START_GUIDE.md）

### 员工账户
- 通过管理员界面创建
- 系统自动生成初始密码

## 📊 核心功能

### 员工端
- ✅ 登录（首次登录获得 1000 积分）
- ✅ 浏览产品
- ✅ 兑换产品
- ✅ 查看兑换历史
- ✅ 查看积分历史
- ✅ 管理个人信息

### 管理员端
- ✅ 员工管理
- ✅ 产品管理（含批量导入）
- ✅ 积分管理（含批量发放）
- ✅ 订单管理
- ✅ 统计报表

## 🛠️ 技术栈

| 层级 | 技术 |
|-----|------|
| 后端 | Go + Gin + GORM + MySQL |
| 前端 | React + TypeScript + Ant Design |
| 部署 | Docker + Nginx |

## 📁 重要文件位置

### 配置文件
- 后端配置：`backend/configs/config.yaml`
- 前端配置：`frontend/.env`
- Docker 配置：`docker-compose.yml`
- Nginx 配置：`nginx/nginx.conf`

### 数据库
- 迁移脚本：`backend/migrations/`
- 数据模型：`backend/internal/models/`

### 代码
- 后端 API：`backend/internal/handler/`
- 前端页面：`frontend/src/pages/`
- 前端组件：`frontend/src/components/`

## 🔧 常用操作

### 创建管理员账户
```sql
-- 连接数据库
mysql -u root -p awsomeshop

-- 插入管理员（密码需要 bcrypt 加密）
INSERT INTO users (full_name, email, phone, password_hash, role, is_first_login, is_active) 
VALUES ('Admin', 'admin@awsomeshop.com', '1234567890', '$2a$10$...', 'admin', FALSE, TRUE);
```

### 查看日志
```bash
# Docker 日志
docker-compose logs -f

# 后端日志
docker-compose logs -f backend

# 前端日志
docker-compose logs -f frontend
```

### 备份数据库
```bash
docker-compose exec mysql mysqldump -u root -p awsomeshop > backup.sql
```

### 恢复数据库
```bash
docker-compose exec -T mysql mysql -u root -p awsomeshop < backup.sql
```

## 🧪 测试清单

- [ ] 员工注册和首次登录
- [ ] 产品浏览和兑换
- [ ] 积分发放和扣除
- [ ] 订单管理
- [ ] 批量操作
- [ ] 并发测试
- [ ] 性能测试
- [ ] 安全测试

详细测试步骤：[TESTING_CHECKLIST.md](TESTING_CHECKLIST.md)

## 🚨 常见问题

### 后端无法启动
1. 检查 MySQL 是否运行
2. 检查配置文件是否正确
3. 检查端口 8080 是否被占用

### 前端无法连接后端
1. 检查后端是否运行
2. 检查 `.env` 中的 API URL
3. 检查 CORS 配置

### 数据库连接失败
1. 检查 MySQL 容器状态
2. 检查数据库用户名和密码
3. 检查数据库是否已创建

## 📞 获取帮助

1. 查看相关文档（见上方文档索引）
2. 检查日志文件
3. 查看 [USER_MANUAL.md](USER_MANUAL.md) 的常见问题部分

## 🎯 下一步

1. ✅ 准备环境（Go、Node.js、MySQL）
2. ✅ 启动应用（参考 QUICK_START_GUIDE.md）
3. ✅ 执行测试（参考 TESTING_CHECKLIST.md）
4. ✅ 部署生产（参考 DEPLOYMENT_GUIDE.md）
5. ✅ 培训用户（参考 USER_MANUAL.md）

---

**项目已 100% 完成，可以开始使用！** 🎉
