# Go Clean Architecture Boilerplate

Production-ready Go REST API boilerplate with Clean Architecture and Module Pattern.

## Features

- ✅ Clean Architecture with Module Pattern
- ✅ Gin Web Framework
- ✅ GORM with PostgreSQL
- ✅ JWT Authentication
- ✅ Rate Limiting & CORS
- ✅ Docker & Docker Compose
- ✅ Graceful Shutdown
- ✅ Shared utilities and responses

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL
- Docker (optional)

### Setup

1. Clone and install dependencies:
```bash
git clone <your-repo>
cd go-boilerplate
go mod download
```

2. Create `.env` file (copy from .env.example)

3. Run with Docker:
```bash
docker-compose up
```

4. Or run locally:
```bash
make run
```

### API Endpoints

#### User Module