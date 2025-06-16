# Go Web GraphQL Backend

A modern Go web backend featuring:
- **GraphQL API** powered by [gqlgen](https://github.com/99designs/gqlgen)
- **Ent ORM** for database modeling
- **Redis** for Automatic Persisted Query (APQ) caching
- **Custom error handling** and i18n
- **Developer experience** enhancements (Playground, hot reload, etc.)

## Features
- GraphQL API with schema-first development
- Ent ORM integration
- Redis APQ cache for persisted queries
- Unified error handling and internationalization
- GraphQL Playground for API exploration
- Transaction and query complexity middleware

## Getting Started

### Prerequisites
- Go 1.18+
- Redis
- (Optional) a SQL database supported by Ent

### Installation
```sh
git clone <your-repo-url>
cd go-web
make install # or go mod tidy
```

### Configuration
- Edit `.env` for database and Redis connection settings.

### Running
```sh
make run # or go run main.go
```

- GraphQL endpoint: `POST /query`
- Playground: `GET /playground`

### Development
- Hot reload: use [air](https://github.com/cosmtrek/air) or [fresh](https://github.com/gravityblast/fresh)
- Schema changes: edit `graph/*.graphql`, then run `make generate`
- Add resolvers in `interface/resolvers/`

### Testing
```sh
make test
```

## Project Structure
- `graph/` - GraphQL schema and generated code
- `interface/resolvers/` - Resolver implementations
- `pkg/redis/` - Redis APQ cache
- `pkg/errors/` - Custom error types
- `pkg/i18n/` - Internationalization

## License
MIT
