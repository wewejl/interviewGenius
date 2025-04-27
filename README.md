# Interview Genius

面试系统后端 API 服务

## 技术栈

- Go: 1.19+
- Web 框架: [Gin](https://github.com/gin-gonic/gin)
- ORM: [GORM](https://gorm.io/)
- 配置管理: [Viper](https://github.com/spf13/viper)
- 身份认证: [JWT](https://github.com/dgrijalva/jwt-go)
- 日志系统: [Zap](https://github.com/uber-go/zap)
- API 文档: [Swagger](https://github.com/swaggo/gin-swagger)
- 数据库: MySQL

## 项目结构

```
.
├── api/                 # API 处理函数
│   └── v1/              # v1版本API
├── config/              # 配置相关
├── docs/                # Swagger文档
├── internal/            # 内部应用程序包
│   ├── middleware/      # 中间件
│   ├── model/           # 数据模型
│   └── service/         # 业务逻辑
├── pkg/                 # 公共包
│   ├── setting/         # 设置相关
│   └── util/            # 工具函数
├── router/              # 路由设置
├── main.go              # 应用入口
├── go.mod               # 依赖管理
└── README.md            # 项目说明
```

## RBAC 权限管理

本项目采用基于角色的访问控制（RBAC）模型进行权限管理，实现细粒度的API访问控制。

### 数据库设计

**用户表 (users)**
```sql
CREATE TABLE users (
  id CHAR(36) PRIMARY KEY,
  username VARCHAR(64) UNIQUE NOT NULL,
  password VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE
);
```

**角色表 (roles)**
```sql
CREATE TABLE roles (
  id INT AUTO_INCREMENT PRIMARY KEY,
  role_name VARCHAR(32) UNIQUE NOT NULL,
  is_super BOOLEAN DEFAULT false
);
```

**权限表 (permissions)**
```sql
CREATE TABLE permissions (
  id CHAR(36) PRIMARY KEY,
  method VARCHAR(8) NOT NULL,      -- GET/POST等
  path_pattern VARCHAR(128) NOT NULL  -- 如 /api/v1/users/*
);
```

**角色-权限关联表 (role_permission)**
```sql
CREATE TABLE role_permission (
  role_id INT,
  permission_id CHAR(36),
  PRIMARY KEY (role_id, permission_id),
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (permission_id) REFERENCES permissions(id)
);
```

**用户-角色关联表 (user_role)**
```sql
CREATE TABLE user_role (
  user_id CHAR(36),
  role_id INT,
  PRIMARY KEY (user_id, role_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (role_id) REFERENCES roles(id)
);
```

### 用户管理 API

本项目目前已实现用户管理相关 API:

1. **用户注册**
   - 路径: `/api/v1/users/register`
   - 方法: POST
   - 权限: 无需权限

2. **用户登录**
   - 路径: `/api/v1/users/login`
   - 方法: POST
   - 权限: 无需权限

3. **获取用户列表**
   - 路径: `/api/v1/users`
   - 方法: GET
   - 权限: 需要认证

4. **获取用户详情**
   - 路径: `/api/v1/users/:id`
   - 方法: GET
   - 权限: 需要认证，且为自己或拥有相应权限

5. **更新用户信息**
   - 路径: `/api/v1/users/:id`
   - 方法: PUT
   - 权限: 需要认证，且为自己或拥有相应权限

6. **删除用户**
   - 路径: `/api/v1/users/:id`
   - 方法: DELETE
   - 权限: 需要认证，且拥有相应权限

7. **为用户分配角色**
   - 路径: `/api/v1/users/:id/roles`
   - 方法: POST
   - 权限: 需要认证，且拥有相应权限

8. **获取用户角色**
   - 路径: `/api/v1/users/:id/roles`
   - 方法: GET
   - 权限: 需要认证，且为自己或拥有相应权限

9. **移除用户角色**
   - 路径: `/api/v1/users/:userId/roles/:roleId`
   - 方法: DELETE
   - 权限: 需要认证，且拥有相应权限

## 开始使用

### 1. 配置

编辑 `config/app.yaml` 文件，设置应用程序参数：

```yaml
app:
  port: 8080
  jwtSecret: "your_jwt_secret"

server:
  runMode: "debug"  # debug or release
  readTimeout: 60
  writeTimeout: 60

database:
  type: "mysql"
  user: "root"
  password: "your_password"
  host: "127.0.0.1"
  port: 3306
  name: "interview_genius"
  tablePrefix: ""
```

### 2. 数据库

创建 MySQL 数据库：

```sql
CREATE DATABASE interview_genius CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

应用程序会自动创建所需表。

### 3. 运行

```bash
go run main.go
```

## API 文档

启动服务后访问 Swagger 文档：

```
http://localhost:8080/swagger/index.html
```

## 主要功能

- 用户注册与登录
- JWT 身份验证
- 用户信息管理
- 管理员角色权限 