## Analytics (Page Views)

Endpoints under `/api/v1/analytics`:

- POST `/track` body: `{ page: string, visitorId: string, userAgent?: string, referrer?: string }` -> records a view and returns stats
- GET `/views?page=/optional-path` -> returns aggregated view stats

Database table: `page_views` with indexes for performance.

# Web Portfolio Backend CMS

A comprehensive, mature, and maintainable Content Management System (CMS) backend built with Go, following clean architecture principles with authentication and robust API design.

## 🚀 Features

### Core CMS Features

- **User Management**: Registration, authentication, and role-based access control
- **Content Management**: Posts, Pages, Categories, and Tags
- **Authentication**: JWT-based authentication with bcrypt password hashing
- **Clean Architecture**: Repository/Service/Handler pattern with dependency injection
- **Database Migration**: Automated database schema management with Goose
- **CORS Support**: Cross-origin resource sharing for frontend integration
- **Logging**: Structured logging with Logrus
- **Configuration**: JSON-based configuration with Viper
- **Input Validation**: Request validation with go-playground/validator

### API Endpoints

#### Authentication

- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `GET /auth/me` - Get current user profile

#### Categories

- `GET /categories` - List all categories
- `GET /categories/:id` - Get category by ID
- `POST /categories` - Create new category
- `PUT /categories/:id` - Update category
- `DELETE /categories/:id` - Delete category

#### Posts

- `GET /posts` - List all posts (with pagination)
- `GET /posts/published` - List published posts
- `GET /posts/:id` - Get post by ID
- `GET /posts/slug/:slug` - Get post by slug
- `GET /posts/author/:authorId` - Get posts by author
- `POST /posts` - Create new post
- `PUT /posts/:id` - Update post
- `DELETE /posts/:id` - Delete post

#### Pages

- `GET /pages` - List all pages (with pagination)
- `GET /pages/published` - List published pages
- `GET /pages/:id` - Get page by ID
- `GET /pages/slug/:slug` - Get page by slug
- `POST /pages` - Create new page
- `PUT /pages/:id` - Update page
- `DELETE /pages/:id` - Delete page

## 🏗️ Architecture

### Clean Architecture Layers

```
├── main.go                 # Application entry point
├── config/                 # Configuration management
├── internal/
│   ├── domain/
│   │   └── models/         # Domain entities
│   ├── repositories/       # Data access layer
│   │   ├── category/
│   │   ├── user/
│   │   ├── post/
│   │   ├── page/
│   │   └── registry.go
│   ├── services/           # Business logic layer
│   │   ├── category/
│   │   ├── user/
│   │   ├── post/
│   │   ├── page/
│   │   └── registry.go
│   ├── handlers/           # HTTP handlers (controllers)
│   │   ├── category/
│   │   ├── auth/
│   │   ├── post/
│   │   ├── page/
│   │   └── registry.go
│   └── auth/              # Authentication service
├── middleware/            # HTTP middlewares
├── common/               # Shared utilities
├── database_schema/      # Database migrations
└── log/                 # Application logs
```

### Technology Stack

- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT with golang-jwt/jwt/v5
- **Password Hashing**: bcrypt
- **Configuration**: Viper
- **Logging**: Logrus
- **Validation**: go-playground/validator
- **Migrations**: Goose
- **CORS**: gin-contrib/cors

## 🔧 Installation & Setup

### Prerequisites

- Go 1.21+
- PostgreSQL database
- Git

### Step 1: Clone Repository

```bash
git clone <repository-url>
cd web-portofolio-backend
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Configure Database

Update `config.json` with your database credentials:

```json
{
  "database": {
    "host": "your-db-host",
    "port": 5432,
    "user": "your-db-user",
    "password": "your-db-password",
    "name": "web_porto_cms",
    "sslmode": "disable"
  }
}
```

### Step 4: Install Goose (Migration Tool)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Step 5: Run Database Migrations

```bash
goose -dir database_schema postgres "host=your-host port=5432 user=your-user password=your-password dbname=web_porto_cms sslmode=disable" up
```

### Step 6: Start the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📁 Database Schema

### Tables Created by Migrations:

1. **users** - User accounts and authentication
2. **categories** - Content categories
3. **posts** - Blog posts and articles
4. **pages** - Static pages
5. **tags** - Content tags
6. **post_categories** - Many-to-many relationship
7. **post_tags** - Many-to-many relationship

## 🔐 Authentication

The system uses JWT (JSON Web Tokens) for authentication:

1. **Register**: Create a new user account
2. **Login**: Authenticate and receive JWT token
3. **Protect Routes**: Include `Authorization: Bearer <token>` header

### Example Registration:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123",
    "role": "admin"
  }'
```

### Example Login:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

## 📊 API Response Format

### Success Response:

```json
{
  "success": true,
  "data": {
    // response data
  }
}
```

### Error Response:

```json
{
  "success": false,
  "error": "error message"
}
```

### Paginated Response:

```json
{
  "success": true,
  "data": {
    "posts": [...],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_records": 50,
      "has_next_page": true,
      "has_prev_page": false,
      "limit": 10
    }
  }
}
```

## 🧪 Testing

### Test User Registration:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### Test Category Creation:

```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"Technology","description":"Tech related posts"}'
```

### Test Post Creation:

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title":"My First Post",
    "content":"This is the content of my first post",
    "status":"published",
    "author_id":1
  }'
```

## 📝 Configuration

### config.json Structure:

```json
{
  "server": {
    "port": 8080,
    "host": "localhost"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "password",
    "name": "web_porto_cms",
    "sslmode": "disable"
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "jwt": {
    "secret": "your-super-secret-jwt-key-change-this-in-production"
  },
  "app": {
    "name": "Web Porto CMS",
    "version": "1.0.0",
    "debug": true
  }
}
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: Request payload validation
- **CORS Configuration**: Controlled cross-origin access
- **SQL Injection Protection**: GORM ORM with prepared statements

## 🚀 Production Deployment

> **📘 Complete deployment guide:** See [DEPLOYMENT.md](./DEPLOYMENT.md) for detailed CI/CD setup and GitHub secrets configuration.

### Quick Overview

The application uses **dual workflows** for automated deployment:

- **Production**: Triggered by version tags (`v*`), deploys to port **1200**
- **Development**: Triggered by `development` branch push, deploys to port **2200**

Both environments share the same database and server credentials (via GitHub environment secrets).

### Required GitHub Secrets

Set these in repository Settings → Environments → production:

```bash
# Server/SSH
SERVER_HOST, SERVER_PORT, SERVER_USER, SERVER_SSH_KEY

# Docker Hub
DOCKER_USERNAME, DOCKER_PASSWORD

# Database (shared for dev & prod)
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

# Application
JWT_SECRET
```

### Deploy to Production

```bash
git tag v1.0.0
git push origin v1.0.0
```

### Deploy to Development

```bash
git checkout development
git push origin development
```

### Local Development Build

```bash
go build -o web-porto-backend main.go
./web-porto-backend
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For support, please create an issue in the repository or contact the development team.

---

**Built with ❤️ by the Web Porto Team**

---

## Quick Start (Local Dev)

1. Copy config.example.json to config.json and adjust values.
2. Start Postgres via Docker:

docker compose up -d

3. Build and run backend:

go build ./...
./web-porto-backend.exe

On startup, the app auto-migrates the PageView model and applies SQL files in database_schema once. Adminer UI: http://localhost:8081
