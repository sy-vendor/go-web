# Go-Web

> 基于 GIN + Ent + GraphQL + WIRE (DI) 构建的现代化 Web 应用框架。

[English](README.md) | 中文

## 项目概述

Go-Web 是一个基于 Go 语言构建的强大且现代化的 Web 应用框架。它集成了最佳实践和流行的库，为构建可扩展的 Web 应用提供了坚实的基础。

## 特性

- 🚀 **高性能**：基于 Go 最快的 Web 框架 Gin 构建
- 🔄 **丰富的中间件**：内置 CORS 等必要的中间件支持
- 🗄️ **强大的 ORM**：使用 Ent 简化数据库操作，支持自动代码生成
- 🔍 **灵活的 API**：支持 GraphQL，提供灵活高效的数据查询
- 🎯 **依赖注入**：使用 Wire 实现编译时依赖注入，保持架构清晰
- 📦 **数据库迁移**：内置数据库版本控制系统

## 技术栈

- [Gin](https://gin-gonic.com/) - Go 语言最快的全功能 Web 框架
- [Wire](https://github.com/google/wire) - Go 语言的编译时依赖注入工具
- [GraphQL](https://github.com/graphql) - API 查询语言
- [Ent](https://github.com/ent) - Go 语言的实体框架

## 环境要求

- Go 1.16 或更高版本
- Wire CLI 工具

```bash
go get -u github.com/google/wire/cmd/wire
```

## 快速开始

1. 克隆仓库：
```bash
git clone https://github.com/sy-vendor/go-web
cd go-web
```

2. 安装依赖：
```bash
go mod download
```

3. 生成依赖注入代码：
```bash
wire ./cmd/apiserver
```

4. 运行应用：
```bash
go run cmd/apiserver/main.go
```

## 项目结构

```
.
├── cmd/            # 应用程序入口点
├── pkg/            # 核心业务逻辑
├── interface/      # 接口层
├── graph/          # GraphQL 定义
├── ent/            # 数据模型
├── migrations/     # 数据库迁移
└── server.go       # 服务器配置
```

## 开发指南

### 创建新实体

1. 初始化新实体：
```bash
ent init User
```

2. 修改 `./ent/schema/user.go` 中的模式定义

3. 生成实体代码：
```bash
go generate ./ent
```

### 添加 GraphQL 类型

1. 在 `graph/` 目录下创建新的 GraphQL 类型定义：

```graphql
type User {
    id: ID!
    name: String!
    sex: Boolean!
    age: Int!
    Account: String!
    Password: String!
}

extend type Query {
    userByAccount(account: String!): User!
}

extend type Mutation {
    updatePasswordByAccount(account: String!, password: String!): User!
}
```

2. 生成 GraphQL 代码：
```bash
go run github.com/99designs/gqlgen generate
```

### API 测试

你可以使用以下工具测试 GraphQL API：
- [Apollo Studio](https://studio.apollographql.com/)
- [Apollo Sandbox](https://studio.apollographql.com/sandbox/explorer)

## 贡献指南

欢迎贡献代码！请随时提交 Pull Request。

## 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件。
