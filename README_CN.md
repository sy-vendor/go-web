# go-web

> 基于 GIN + ENT + GRAPH + WIRE 实现的RBAC权限管理脚手架，目标是提供一套轻量的中后台开发框架，方便、快速的完成业务需求的开发。

[English](README.md) | 中文

## 特性

- 遵循 `RESTful API` 设计规范 & 基于接口的编程规范
- 基于 `GIN` 框架，提供了中间件支持（CORS）
- 基于 `ENT` 简化应用程序与数据库之间的交互 & 自动生成模型代码
- 基于 `GRAPH` 支持多种查询语言 & 可扩展性 & 多种API
- 基于 `WIRE` 的依赖注入 -- 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程
- 基于 `MIGRATE` 数据库版本控制
- 基于 `go mod` 的依赖管理(国内源可使用：<https://goproxy.cn/>)

## 依赖工具

```
go get -u github.com/google/wire/cmd/wire
```

- [wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go

## 依赖框架

- [Gin](https://gin-gonic.com/) -- The fastest full-featured web framework for Go.
- [Wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go.
- [Graph](https://github.com/graphql) -- A query language and runtime for APIs that simplifies data retrieval and manipulation.
- [Ent](https://github.com/ent) -- A framework for building scalable and maintainable software with Go.

## 快速开始

```bash
$ git clone https://github.com/sy-vendor/go-web

$ cd go-web

# 使用go命令运行
$ go run cmd/apiserver/main.go
```

## 生成依赖注入文件

```bash
#  使用wire命令
wire gen ./cmd/apiserver/
```

## 快速生成数据表
```bash
#  使用ent
ent init User
```
#### 修改模块 ./ent/schema/user.go

```bash
#  使用ent生成
go generate ./ent
```

### 快速生成业务模块

#### 在graph下创建模版 user.graphql

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
    "find user by account"
    userByAccount(account: String!): User!
}

extend type Mutation {
    "update user account password"
    updatePasswordByAccount(account: String!, password: String!): User!
}
```

#### 执行命令并运行

```bash
go run github.com/99designs/gqlgen generate

# 生成依赖项
make ./cmd/apiserver/wire

# 运行服务
go run ./cmd/apiserver/main.go
```

#### API请求测试
```
https://studio.apollographql.com/
```
