# Go-Web

> A modern web application framework based on GIN + Ent + GraphQL + WIRE (DI).

English | [中文](README_CN.md)

## Overview

Go-Web is a powerful and modern web application framework built with Go. It combines the best practices and popular libraries to provide a robust foundation for building scalable web applications.

## Features

- 🚀 **High Performance**: Built on top of Gin, one of the fastest web frameworks for Go
- 🔄 **Rich Middleware**: Built-in CORS and other essential middleware support
- 🗄️ **Powerful ORM**: Simplified database operations with Ent, featuring automatic code generation
- 🔍 **Flexible API**: GraphQL support for flexible and efficient data querying
- 🎯 **Dependency Injection**: Clean architecture with Wire for compile-time dependency injection
- 📦 **Database Migrations**: Built-in database version control system

## Tech Stack

- [Gin](https://gin-gonic.com/) - The fastest full-featured web framework for Go
- [Wire](https://github.com/google/wire) - Compile-time Dependency Injection for Go
- [GraphQL](https://github.com/graphql) - A query language for APIs
- [Ent](https://github.com/ent) - An entity framework for Go

## Prerequisites

- Go 1.16 or higher
- Wire CLI tool

```bash
go get -u github.com/google/wire/cmd/wire
```

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/sy-vendor/go-web
cd go-web
```

2. Install dependencies:
```bash
go mod download
```

3. Generate dependency injection code:
```bash
wire ./cmd/apiserver
```

4. Run the application:
```bash
go run cmd/apiserver/main.go
```

## Project Structure

```
.
├── cmd/            # Application entry points
├── pkg/            # Core business logic
├── interface/      # Interface layer
├── graph/          # GraphQL definitions
├── ent/            # Data models
├── migrations/     # Database migrations
└── server.go       # Server configuration
```

## Development Guide

### Creating a New Entity

1. Initialize a new entity:
```bash
ent init User
```

2. Modify the schema in `./ent/schema/user.go`

3. Generate the entity code:
```bash
go generate ./ent
```

### Adding GraphQL Types

1. Create a new GraphQL type definition in `graph/` directory:

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

2. Generate GraphQL code:
```bash
go run github.com/99designs/gqlgen generate
```

### Testing the API

You can test the GraphQL API using:
- [Apollo Studio](https://studio.apollographql.com/)
- [Apollo Sandbox](https://studio.apollographql.com/sandbox/explorer)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
