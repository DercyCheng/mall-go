# Mall-Go

![Mall-Go](https://img.shields.io/badge/Mall--Go-v0.1.0-blue)
![Go Version](https://img.shields.io/badge/Go-1.23.1-brightgreen)
![License](https://img.shields.io/badge/License-MIT-yellow)

Mall-Go 是一个基于微服务架构和领域驱动设计(DDD)原则实现的全面电子商务系统的 Go 语言版本。本项目旨在提供一个高性能、可扩展且易于维护的电子商务平台解决方案。

## 🚀 项目概述

Mall-Go 是对原 Java 版本 Mall 项目的重写，充分利用了 Go 语言在性能、资源利用率和开发效率方面的优势。该项目包括前台商城系统和后台管理系统两部分。

### 核心功能

- **前台商城系统**: 首页、商品推荐、商品搜索、商品展示、购物车、订单流程、会员中心、客户服务等
- **后台管理系统**: 商品管理、订单管理、会员管理、促销管理、运营管理、内容管理、统计报表、权限管理等

## 📋 项目进度

| 模块         | 状态     | 进度 |
| ------------ | -------- | ---- |
| 用户服务     | 进行中   | 70%  |
| 商品服务     | 进行中   | 85%  |
| 订单服务     | 进行中   | 60%  |
| 购物车服务   | 进行中   | 50%  |
| 库存服务     | 进行中   | 40%  |
| 支付服务     | 进行中   | 30%  |
| 搜索服务     | 初始阶段 | 20%  |
| 网关服务     | 进行中   | 75%  |
| 认证服务     | 进行中   | 65%  |
| 后台管理服务 | 进行中   | 60%  |
| 前台门户服务 | 进行中   | 55%  |
| 内容服务     | 初始阶段 | 15%  |
| 通知服务     | 初始阶段 | 10%  |
| 促销服务     | 初始阶段 | 20%  |
| 推荐服务     | 计划中   | 5%   |

### 部署进度

- 核心服务的 Kubernetes 配置已就绪
- 所有服务的 Docker 容器化已完成
- 基础设施组件（MySQL、Redis、Consul）设置已完成
- CI/CD 流水线正在配置中

## 🔨 技术栈

### 核心技术栈

| 技术     | 用途         | 仓库/网站                           |
| -------- | ------------ | ----------------------------------- |
| Go       | 编程语言     | https://golang.org/                 |
| Gin      | Web 框架     | https://github.com/gin-gonic/gin    |
| GORM     | ORM 框架     | https://gorm.io/                    |
| JWT-Go   | JWT 认证     | https://github.com/golang-jwt/jwt   |
| Go-Redis | Redis 客户端 | https://github.com/go-redis/redis   |
| Consul   | 服务注册中心 | https://github.com/hashicorp/consul |
| gRPC     | 微服务通信   | https://github.com/grpc/grpc-go     |
| Zap      | 日志记录     | https://github.com/uber-go/zap      |
| Viper    | 配置管理     | https://github.com/spf13/viper      |

### 微服务架构

```
                     ┌────────────────┐
                     │    API 网关    │
                     └────────────────┘
                             │
            ┌────────────────┼────────────────┐
            │                │                │
      ┌─────▼─────┐    ┌─────▼─────┐    ┌─────▼─────┐
      │ 后台门户  │    │ 前台门户  │    │ 移动应用  │
      └─────┬─────┘    └─────┬─────┘    └─────┬─────┘
            │                │                │
┌───────────┴────────────────┴────────────────┴───────────┐
│                                                         │
│                      服务网格                           │
│                                                         │
└───┬─────────┬─────────┬──────────┬────────────┬─────────┘
    │         │         │          │            │
┌───▼───┐ ┌───▼───┐ ┌───▼───┐  ┌───▼───┐    ┌───▼───┐
│用户   │ │商品   │ │订单   │  │支付   │    │搜索   │
│服务   │ │服务   │ │服务   │  │服务   │... │服务   │
└───┬───┘ └───┬───┘ └───┬───┘  └───┬───┘    └───┬───┘
    │         │         │          │            │
┌───▼─────────▼─────────▼──────────▼────────────▼───┐
│                                                   │
│                    消息队列                       │
│                                                   │
└───┬─────────────────────────────────┬─────────────┘
    │                                 │
┌───▼───┐                        ┌────▼───┐
│主数据 │                        │副本数据│
│  库   │◄───────同步───────────►│   库   │
└───────┘                        └────────┘
```

## 💻 开始使用

### 前提条件

- Go 1.23 或更高版本
- Docker 和 Docker Compose
- Kubernetes（用于生产部署）
- MySQL 8.0+
- Redis 6.0+

### 开发环境设置

1. 克隆仓库

```bash
git clone https://github.com/yourusername/mall-go.git
cd mall-go
```

2. 设置环境

```bash
# 设置 Go 环境
go env -w GOPROXY=https://goproxy.cn,direct

# 安装依赖
go mod tidy
```

3. 使用 Docker Compose 启动基础设施服务

```bash
cd deployments/docker
docker-compose -f docker-compose-env.yml up -d
```

4. 本地运行服务（开发模式）

```bash
# 示例：运行商品服务
cd services/product-service
go run cmd/server/main.go
```

### Kubernetes 部署

对于生产部署，使用 Kubernetes 配置：

```bash
cd deployments/kubernetes
./deploy.sh
```

## 🏗️ 项目结构

Mall-Go 遵循一个组织良好的结构：

```
mall-go/
├── services/                  # 所有微服务
│   ├── user-service/          # 用户服务
│   ├── product-service/       # 商品服务
│   ├── order-service/         # 订单服务
│   ├── ...
├── pkg/                       # 共享库
│   ├── auth/                  # 认证
│   ├── cache/                 # 缓存
│   ├── ...
├── api/                       # API 定义
├── deployments/               # 部署配置
│   ├── docker/                # Docker 相关
│   ├── kubernetes/            # Kubernetes 配置
├── docs/                      # 文档
├── script/                    # 脚本和资源
```

## 📝 许可证

本项目采用 MIT 许可证 - 有关详细信息，请查看 LICENSE 文件。

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

## 👥 作者

- [Dercyc](https://github.com/DercyCheng) - _初始工作_

## 🙏 致谢

- 感谢原始的 Mall 项目提供灵感
- 感谢 Go 社区提供的优秀库和工具
