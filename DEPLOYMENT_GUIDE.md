# AWSomeShop 部署指南

## 📋 部署前检查清单

### 代码准备
- [ ] 所有功能测试通过
- [ ] 所有安全测试通过
- [ ] 性能测试达标
- [ ] 代码已提交到 Git 仓库
- [ ] 创建生产分支（production）

### 环境准备
- [ ] 生产服务器已准备
- [ ] 域名已配置
- [ ] SSL 证书已获取
- [ ] 数据库服务器已准备
- [ ] 备份策略已制定

---

## 🚀 部署方式

### 方式 1：Docker Compose 部署（推荐）

#### 优点
- 一键部署
- 环境隔离
- 易于维护
- 易于扩展

#### 部署步骤

##### 1. 准备服务器

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

##### 2. 克隆代码

```bash
# 创建应用目录
sudo mkdir -p /opt/awsomeshop
cd /opt/awsomeshop

# 克隆代码
git clone <your-repo-url> .
git checkout production
```

##### 3. 配置环境变量

```bash
# 复制环境配置文件
cp .env.example .env

# 编辑环境变量
nano .env
```

`.env` 文件内容：
```env
# 数据库配置
MYSQL_ROOT_PASSWORD=<strong-password>
MYSQL_DATABASE=awsomeshop
MYSQL_USER=awsomeshop
MYSQL_PASSWORD=<strong-password>

# 后端配置
JWT_SECRET=<generate-strong-secret>
JWT_EXPIRATION=24h
SERVER_PORT=8080
SERVER_MODE=release

# 前端配置
REACT_APP_API_URL=https://yourdomain.com/api
```

##### 4. 配置 SSL 证书

```bash
# 创建 SSL 目录
mkdir -p nginx/ssl

# 使用 Let's Encrypt 获取证书
sudo apt install certbot
sudo certbot certonly --standalone -d yourdomain.com

# 复制证书到项目目录
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem nginx/ssl/
```

##### 5. 更新 Nginx 配置

编辑 `nginx/nginx.conf`：

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 前端
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API
    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

##### 6. 构建和启动服务

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 检查服务状态
docker-compose ps
```

##### 7. 运行数据库迁移

```bash
# 进入后端容器
docker-compose exec backend sh

# 运行迁移（如果有迁移工具）
# 或者手动执行 SQL
mysql -h mysql -u awsomeshop -p awsomeshop < /app/migrations/001_create_tables.sql

# 退出容器
exit
```

##### 8. 创建管理员账户

```bash
# 连接到数据库
docker-compose exec mysql mysql -u awsomeshop -p awsomeshop

# 插入管理员账户
INSERT INTO users (full_name, email, phone, password_hash, role, is_first_login, is_active, created_at, updated_at) 
VALUES (
  'Admin User',
  'admin@yourdomain.com',
  '1234567890',
  '$2a$10$...',  -- 使用 bcrypt 生成的密码哈希
  'admin',
  FALSE,
  TRUE,
  NOW(),
  NOW()
);

# 退出
exit
```

##### 9. 验证部署

```bash
# 检查服务健康
curl https://yourdomain.com
curl https://yourdomain.com/api/health

# 测试登录
curl -X POST https://yourdomain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@yourdomain.com","password":"your-password"}'
```

---

### 方式 2：手动部署

#### 后端部署

##### 1. 编译 Go 应用

```bash
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o awsomeshop-api cmd/api/main.go
```

##### 2. 上传到服务器

```bash
scp awsomeshop-api user@server:/opt/awsomeshop/
scp -r configs user@server:/opt/awsomeshop/
scp -r migrations user@server:/opt/awsomeshop/
```

##### 3. 创建 systemd 服务

创建 `/etc/systemd/system/awsomeshop-api.service`：

```ini
[Unit]
Description=AWSomeShop API Service
After=network.target mysql.service

[Service]
Type=simple
User=awsomeshop
WorkingDirectory=/opt/awsomeshop
ExecStart=/opt/awsomeshop/awsomeshop-api
Restart=on-failure
RestartSec=5s

Environment="JWT_SECRET=your-secret"
Environment="DB_HOST=localhost"
Environment="DB_PORT=3306"
Environment="DB_USER=awsomeshop"
Environment="DB_PASSWORD=your-password"
Environment="DB_NAME=awsomeshop"

[Install]
WantedBy=multi-user.target
```

##### 4. 启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl enable awsomeshop-api
sudo systemctl start awsomeshop-api
sudo systemctl status awsomeshop-api
```

#### 前端部署

##### 1. 构建生产版本

```bash
cd frontend
npm install
npm run build
```

##### 2. 上传到服务器

```bash
scp -r build/* user@server:/var/www/awsomeshop/
```

##### 3. 配置 Nginx

创建 `/etc/nginx/sites-available/awsomeshop`：

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    root /var/www/awsomeshop;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

##### 4. 启用站点

```bash
sudo ln -s /etc/nginx/sites-available/awsomeshop /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 🔒 安全加固

### 1. 防火墙配置

```bash
# 安装 UFW
sudo apt install ufw

# 配置规则
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 启用防火墙
sudo ufw enable
sudo ufw status
```

### 2. 数据库安全

```bash
# 运行 MySQL 安全脚本
sudo mysql_secure_installation

# 配置建议：
# - 设置强密码
# - 删除匿名用户
# - 禁止 root 远程登录
# - 删除测试数据库
```

### 3. 应用安全

- [ ] 更改所有默认密码
- [ ] 使用强 JWT 密钥
- [ ] 启用 HTTPS
- [ ] 配置 CORS 白名单
- [ ] 启用速率限制
- [ ] 配置日志记录
- [ ] 定期更新依赖

### 4. SSL/TLS 配置

```bash
# 自动续期 Let's Encrypt 证书
sudo crontab -e

# 添加以下行（每月 1 号凌晨 2 点续期）
0 2 1 * * certbot renew --quiet && systemctl reload nginx
```

---

## 💾 数据备份策略

### 1. 数据库备份

#### 自动备份脚本

创建 `/opt/awsomeshop/backup.sh`：

```bash
#!/bin/bash

# 配置
BACKUP_DIR="/opt/awsomeshop/backups"
DB_NAME="awsomeshop"
DB_USER="awsomeshop"
DB_PASSWORD="your-password"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/awsomeshop_$DATE.sql.gz"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -u $DB_USER -p$DB_PASSWORD $DB_NAME | gzip > $BACKUP_FILE

# 删除 30 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

# 输出结果
if [ -f $BACKUP_FILE ]; then
    echo "Backup successful: $BACKUP_FILE"
    ls -lh $BACKUP_FILE
else
    echo "Backup failed!"
    exit 1
fi
```

#### 设置定时任务

```bash
# 添加执行权限
chmod +x /opt/awsomeshop/backup.sh

# 配置 cron（每天凌晨 3 点备份）
sudo crontab -e

# 添加以下行
0 3 * * * /opt/awsomeshop/backup.sh >> /var/log/awsomeshop-backup.log 2>&1
```

### 2. 数据恢复

```bash
# 解压备份文件
gunzip awsomeshop_20260114_030000.sql.gz

# 恢复数据库
mysql -u awsomeshop -p awsomeshop < awsomeshop_20260114_030000.sql
```

### 3. 远程备份

```bash
# 使用 rsync 同步到远程服务器
rsync -avz /opt/awsomeshop/backups/ user@backup-server:/backups/awsomeshop/

# 或使用云存储（AWS S3 示例）
aws s3 sync /opt/awsomeshop/backups/ s3://your-bucket/awsomeshop-backups/
```

---

## 📊 监控和日志

### 1. 应用日志

#### 后端日志

```bash
# 查看实时日志
docker-compose logs -f backend

# 或使用 systemd
sudo journalctl -u awsomeshop-api -f

# 配置日志轮转
sudo nano /etc/logrotate.d/awsomeshop
```

`/etc/logrotate.d/awsomeshop`：
```
/var/log/awsomeshop/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 awsomeshop awsomeshop
    sharedscripts
    postrotate
        systemctl reload awsomeshop-api > /dev/null 2>&1 || true
    endscript
}
```

#### Nginx 日志

```bash
# 访问日志
tail -f /var/log/nginx/access.log

# 错误日志
tail -f /var/log/nginx/error.log
```

### 2. 系统监控

#### 安装监控工具

```bash
# 安装 htop
sudo apt install htop

# 安装 netdata（可选）
bash <(curl -Ss https://my-netdata.io/kickstart.sh)
```

#### 监控指标

- CPU 使用率
- 内存使用率
- 磁盘使用率
- 网络流量
- 数据库连接数
- API 响应时间

### 3. 告警配置

可以使用以下工具配置告警：
- Prometheus + Grafana
- Datadog
- New Relic
- CloudWatch（AWS）

---

## 🔄 更新和维护

### 1. 应用更新流程

```bash
# 1. 备份数据库
/opt/awsomeshop/backup.sh

# 2. 拉取最新代码
cd /opt/awsomeshop
git pull origin production

# 3. 重新构建（Docker）
docker-compose build

# 4. 停止服务
docker-compose down

# 5. 运行数据库迁移（如果有）
# ...

# 6. 启动服务
docker-compose up -d

# 7. 验证服务
docker-compose ps
curl https://yourdomain.com/api/health
```

### 2. 回滚流程

```bash
# 1. 停止服务
docker-compose down

# 2. 回滚代码
git checkout <previous-commit>

# 3. 重新构建
docker-compose build

# 4. 恢复数据库（如果需要）
mysql -u awsomeshop -p awsomeshop < backup.sql

# 5. 启动服务
docker-compose up -d
```

### 3. 维护窗口

建议设置定期维护窗口：
- 时间：每周日凌晨 2:00-4:00
- 内容：
  - 系统更新
  - 依赖更新
  - 数据库优化
  - 日志清理
  - 备份验证

---

## 📝 部署检查清单

### 部署前
- [ ] 代码已测试
- [ ] 配置文件已更新
- [ ] SSL 证书已配置
- [ ] 数据库已准备
- [ ] 备份策略已制定

### 部署中
- [ ] 服务器已准备
- [ ] 代码已部署
- [ ] 数据库已迁移
- [ ] 服务已启动
- [ ] 管理员账户已创建

### 部署后
- [ ] 功能验证通过
- [ ] 性能测试通过
- [ ] 安全检查通过
- [ ] 监控已配置
- [ ] 备份已验证
- [ ] 文档已更新

---

## 🆘 故障排查

### 问题 1：服务无法启动

```bash
# 检查日志
docker-compose logs backend
docker-compose logs mysql

# 检查端口占用
sudo netstat -tulpn | grep 8080
sudo netstat -tulpn | grep 3306

# 检查配置文件
cat configs/config.yaml
```

### 问题 2：数据库连接失败

```bash
# 检查 MySQL 状态
docker-compose exec mysql mysql -u root -p -e "SHOW DATABASES;"

# 检查网络连接
docker-compose exec backend ping mysql

# 检查用户权限
docker-compose exec mysql mysql -u root -p -e "SHOW GRANTS FOR 'awsomeshop'@'%';"
```

### 问题 3：前端无法访问后端

```bash
# 检查 Nginx 配置
sudo nginx -t

# 检查 Nginx 日志
tail -f /var/log/nginx/error.log

# 检查后端健康
curl http://localhost:8080/api/health
```

### 问题 4：SSL 证书问题

```bash
# 检查证书有效期
openssl x509 -in /etc/letsencrypt/live/yourdomain.com/fullchain.pem -noout -dates

# 手动续期
sudo certbot renew

# 测试 SSL 配置
curl -I https://yourdomain.com
```

---

## 📚 相关文档

- [快速启动指南](QUICK_START_GUIDE.md)
- [测试清单](TESTING_CHECKLIST.md)
- [项目状态](PROJECT_STATUS.md)
- [当前进度](CURRENT_PROGRESS.md)

---

## 🎉 部署成功

如果所有检查都通过，恭喜你成功部署了 AWSomeShop 系统！

访问 https://yourdomain.com 开始使用。

**记住：**
- 定期备份数据
- 监控系统状态
- 及时更新依赖
- 保持安全意识
