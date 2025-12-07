作业要求概述
本次作业要求你使用 Go 语言结合 Gin 框架和 GORM 库开发一个个人博客系统的后端，实现博客文章的基本管理功能，
包括文章的创建、读取、更新和删除（CRUD）操作，同时支持用户认证和简单的评论功能。

1.项目初始化
创建一个新的 Go 项目，使用 go mod init 初始化项目依赖管理。
安装必要的库，如 Gin 框架、GORM 以及数据库驱动（ MySQL）。

2.数据库设计与模型定义
设计数据库表结构，至少包含以下几个表：
users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
使用 GORM 定义对应的 Go 模型结构体。

3.用户认证与授权
实现用户注册和登录功能，用户注册时需要对密码进行加密存储，登录时验证用户输入的用户名和密码。
使用 JWT（JSON Web Token）实现用户认证和授权，用户登录成功后返回一个 JWT，后续的需要认证的接口需要验证该 JWT 的有效性。

4.文章管理功能
实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
实现文章的更新功能，只有文章的作者才能更新自己的文章。
实现文章的删除功能，只有文章的作者才能删除自己的文章。

5.评论功能
实现评论的创建功能，已认证的用户可以对文章发表评论。
实现评论的读取功能，支持获取某篇文章的所有评论列表

6.错误处理与日志记录
对可能出现的错误进行统一处理，如数据库连接错误、用户认证失败、文章或评论不存在等，返回合适的 HTTP 状态码和错误信息。
使用日志库记录系统的运行信息和错误信息，方便后续的调试和维护。

项目结构：
blog-backend/
├── go.mod
├── go.sum
├── main.go
├── config/      # 配置管理（）
│   └── config.go
├── models/        # 模型
│   ├── user.go
│   ├── post.go
│   └── comment.go
├── controllers/       # 控制器
│   ├── auth.go
│   ├── post.go
│   └── comment.go
├── middleware/       # 中间件
│   ├── auth.go
│   └── logger.go
├── utils/          # 工具函数
│   ├── jwt.go
│   └── response.go
└── database/       # 数据库连接
    └── database.go


注册用户：
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

curl 命令的 -d 参数:
{"username":"testuser","email":"test@example.com","password":"password123"}
    ↓
HTTP 请求体 (Request Body)
    ↓
Gin 框架接收并解析
    ↓
c.ShouldBindJSON(&req) 绑定到局部变量 req

登录：
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

创建文章（使用返回的token）：
curl -X POST http://localhost:8080/api/posts \
  -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiaXNzIjoiYmxvZy1iYWNrZW5kIiwiZXhwIjoxNzY1MTcyMTY1LCJpYXQiOjE3NjUwODU3NjV9.erdyKu34rpz_ZMA-ST21s3E_yQotwZI1gnm7VM5Xi3E" \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Post","content":"This is the content of my first post."}'

获取所有文章：
curl -X GET http://localhost:8080/api/posts



安全的 JWT_SECRET
# 生成 32 字节的随机 Base64 字符串
openssl rand -base64 32
# 或者生成 64 个字符的随机字符串
head /dev/urandom | tr -dc A-Za-z0-9 | head -c 64 ; echo ''
