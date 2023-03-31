# gin-web

> RBAC scaffolding based on GIN + Ent + Graph + WIRE (DI).

English | [中文](README_CN.md)

## Features

- Follow the `RESTful API` design specification
- Provides rich `Gin` middlewares (CORS)
- Simplifies the interaction between the application and database with `ENT`, and automatically generates model code.
- Supports multiple query languages, scalability, and multiple APIs with `GRAPH`.
- Use `Wire` to resolve dependencies between modules
- Uses `Migrate` for database version control
- Support `Swagger`

## Dependent Tools

```bash
go get -u github.com/google/wire/cmd/wire
```

- [wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go
## Dependent Library

- [Gin](https://gin-gonic.com/) -- The fastest full-featured web framework for Go.
- [Wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go.
- [Graph](https://github.com/graphql) -- A query language and runtime for APIs that simplifies data retrieval and manipulation.
- [Ent](https://github.com/ent) -- A framework for building scalable and maintainable software with Go.

## Getting Started

```bash
$ git clone https://github.com/sy-vendor/go-web

$ cd go-web

$ go run cmd/go-web/main.go
```

### Use `wire` to generate dependency injection

```bash
wire gen ./cmd/apiserver/
```

### Create Database: `user.go`
```bash
ent init User
```
#### Modify module ./ent/schema/user.go

```bash
go generate ./ent
```

### Rapidly generate business modules

#### Create a template under graph user.graphql

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

### Execute command and run

```bash
go run github.com/99designs/gqlgen generate

make ./cmd/apiserver/wire

go run ./cmd/apiserver/main.go
```

#### API request testing
```
https://studio.apollographql.com/
```
