# AWSomeShop Frontend

员工福利电商系统前端应用

## 技术栈

- React 18+
- TypeScript
- Ant Design
- React Router
- Axios
- react-i18next

## 项目结构

```
frontend/
├── public/              # 静态资源
├── src/
│   ├── components/      # 可复用组件
│   ├── pages/           # 页面组件
│   ├── services/        # API 服务
│   ├── contexts/        # React Context
│   ├── hooks/           # 自定义 Hooks
│   ├── utils/           # 工具函数
│   ├── types/           # TypeScript 类型定义
│   ├── locales/         # 国际化资源
│   ├── routes/          # 路由配置
│   ├── App.tsx          # 应用入口组件
│   └── index.tsx        # 应用入口文件
├── package.json
└── tsconfig.json
```

## 快速开始

### 安装依赖

```bash
cd frontend
npm install
```

### 开发模式

```bash
npm start
```

应用将在 [http://localhost:3000](http://localhost:3000) 启动

### 构建生产版本

```bash
npm run build
```

### 运行测试

```bash
npm test
```

## 环境变量

创建 `.env.local` 文件配置环境变量：

```
REACT_APP_API_BASE_URL=http://localhost:8080/api/v1
REACT_APP_IPGEO_API_KEY=your-api-key
```

## 功能模块

### 员工端
- 登录
- 产品浏览和兑换
- 积分查询
- 兑换历史
- 个人信息管理

### 管理员端
- 员工账户管理
- 产品管理
- 积分管理
- 订单管理
- 统计报表

## 开发指南

### 添加新页面

1. 在 `src/pages/` 创建页面组件
2. 在 `src/routes/` 注册路由
3. 在导航菜单中添加链接

### 调用 API

使用 `src/services/` 中的服务模块：

```typescript
import { productService } from 'services/productService';

const products = await productService.getProducts();
```

### 国际化

使用 `useTranslation` hook：

```typescript
import { useTranslation } from 'react-i18next';

const { t } = useTranslation();
return <h1>{t('welcome')}</h1>;
```

## License

MIT
