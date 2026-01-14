#!/bin/bash

echo "========================================="
echo "阶段一完成验证脚本"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check function
check() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $1"
        return 0
    else
        echo -e "${RED}✗${NC} $1"
        return 1
    fi
}

# 1. Check backend structure
echo "1. 检查后端项目结构..."
[ -f "backend/go.mod" ] && check "go.mod 存在"
[ -f "backend/cmd/api/main.go" ] && check "main.go 存在"
[ -d "backend/internal/config" ] && check "config 目录存在"
[ -d "backend/internal/database" ] && check "database 目录存在"
[ -d "backend/internal/models" ] && check "models 目录存在"
[ -d "backend/migrations" ] && check "migrations 目录存在"
[ -f "backend/Dockerfile" ] && check "后端 Dockerfile 存在"
echo ""

# 2. Check frontend structure
echo "2. 检查前端项目结构..."
[ -f "frontend/package.json" ] && check "package.json 存在"
[ -f "frontend/tsconfig.json" ] && check "tsconfig.json 存在"
[ -f "frontend/src/App.tsx" ] && check "App.tsx 存在"
[ -d "frontend/src/layouts" ] && check "layouts 目录存在"
[ -d "frontend/src/i18n" ] && check "i18n 目录存在"
[ -f "frontend/Dockerfile" ] && check "前端 Dockerfile 存在"
echo ""

# 3. Check Docker configuration
echo "3. 检查 Docker 配置..."
[ -f "docker-compose.yml" ] && check "docker-compose.yml 存在"
[ -f ".env.example" ] && check ".env.example 存在"
[ -f "nginx/nginx.conf" ] && check "nginx.conf 存在"
echo ""

# 4. Check migration scripts
echo "4. 检查数据库迁移脚本..."
[ -f "backend/migrations/001_create_users_table.sql" ] && check "users 表迁移脚本存在"
[ -f "backend/migrations/002_create_products_table.sql" ] && check "products 表迁移脚本存在"
[ -f "backend/migrations/003_create_redemption_orders_table.sql" ] && check "redemption_orders 表迁移脚本存在"
[ -f "backend/migrations/004_create_points_transactions_table.sql" ] && check "points_transactions 表迁移脚本存在"
[ -f "backend/migrations/005_create_product_price_history_table.sql" ] && check "product_price_history 表迁移脚本存在"
echo ""

# 5. Check i18n configuration
echo "5. 检查国际化配置..."
[ -f "frontend/src/i18n/config.ts" ] && check "i18n config 存在"
[ -f "frontend/src/i18n/locales/zh.json" ] && check "中文语言文件存在"
[ -f "frontend/src/i18n/locales/en.json" ] && check "英文语言文件存在"
echo ""

# 6. Check scripts
echo "6. 检查辅助脚本..."
[ -f "scripts/start-dev.sh" ] && check "start-dev.sh 存在"
[ -f "scripts/start-prod.sh" ] && check "start-prod.sh 存在"
[ -f "scripts/stop-all.sh" ] && check "stop-all.sh 存在"
echo ""

# 7. Check documentation
echo "7. 检查文档..."
[ -f "README.md" ] && check "项目 README 存在"
[ -f "backend/README.md" ] && check "后端 README 存在"
[ -f "frontend/README.md" ] && check "前端 README 存在"
[ -f "PHASE1_COMPLETION.md" ] && check "阶段一完成文档存在"
echo ""

echo "========================================="
echo "验证完成！"
echo "========================================="
echo ""
echo "如果所有检查都通过，阶段一已成功完成！"
echo ""
echo "下一步："
echo "1. 配置环境变量: cp .env.example .env"
echo "2. 启动开发环境: ./scripts/start-dev.sh"
echo "3. 或使用 Docker: docker-compose up -d"
echo ""
