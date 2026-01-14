# AWSomeShop Backend API

员工福利电商系统后端 API 服务

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM (MySQL)
- JWT Authentication
- Viper Configuration

## 项目结构

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # 应用入口
├── internal/
│   ├── config/              # 配置管理
│   ├── database/            # 数据库连接和迁移
│   ├── models/              # 数据模型
│   ├── repository/          # 数据访问层
│   ├── service/             # 业务逻辑层
│   ├── handler/             # HTTP 处理器
│   ├── middleware/          # 中间件
│   └── router/              # 路由配置
├── pkg/                     # 公共包
├── configs/                 # 配置文件
├── migrations/              # 数据库迁移脚本
└── Makefile                 # 构建脚本
```

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0 或更高版本

### 安装依赖

```bash
cd backend
go mod download
```

### 配置

复制配置文件并修改：

```bash
cp configs/config.yaml configs/config.local.yaml
# 编辑 config.local.yaml 填入实际配置
```

### 运行

```bash
# 开发模式
make run

# 或直接运行
go run cmd/api/main.go
```

### 构建

```bash
make build
```

### 测试

```bash
make test
```

## API 文档

API 端点文档请参考 [API.md](docs/API.md)

## 开发指南

### 添加新功能

1. 在 `internal/models/` 添加数据模型
2. 在 `internal/repository/` 添加数据访问层
3. 在 `internal/service/` 添加业务逻辑
4. 在 `internal/handler/` 添加 HTTP 处理器
5. 在 `internal/router/` 注册路由

### 代码规范

- 使用 `go fmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查
- 编写单元测试

## License

MIT
