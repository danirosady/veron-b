# TMS Docker Setup

## Quick Start (Production)

```bash
# Build and start all services
docker-compose up -d --build

# Check service health
docker-compose ps

# View logs
docker-compose logs -f backend

# Stop all services
docker-compose down
```

## Development

```bash
# Start development environment (with hot-reload)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Backend runs on http://localhost:8080
# Frontend Vite dev server on http://localhost:5173

# Stop
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
```

## First Run

On first startup, the database is initialized with migrations from `./migrations/`.
A superadmin user is created by the seed data:
- **Email**: superadmin@tms.local
- **Password**: password

## Environment Variables

### Backend

| Variable | Default | Description |
|----------|---------|-------------|
| APP_ENV | production | Environment mode |
| APP_PORT | 8080 | Server port |
| APP_HOST | 0.0.0.0 | Server host |
| APP_ALLOWED_ORIGINS | http://localhost:5173 | CORS origins (comma-separated) |
| DB_HOST | postgres | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_NAME | tms | Database name |
| DB_USER | tms_user | Database user |
| DB_PASSWORD | secret | Database password |
| DB_SSLMODE | disable | PostgreSQL SSL mode |
| JWT_SECRET | (required) | JWT signing secret — CHANGE THIS |
| JWT_ACCESS_EXPIRY | 15m | Access token lifetime |
| JWT_REFRESH_EXPIRY | 168h | Refresh token lifetime |
| REDIS_HOST | redis | Redis host |
| REDIS_PORT | 6379 | Redis port |

### Frontend (Vite)

| Variable | Default |
|----------|---------|
| VITE_API_URL | http://localhost:8080/api/v1 |

## Data Persistence

- `postgres_data` volume — PostgreSQL data
- `redis_data` volume — Redis data

To reset all data:
```bash
docker-compose down -v
docker-compose up -d
```

## Ports

| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache (optional) |
| Backend API | 8080 | Go server |
| Frontend | 80 | Nginx (production) |
| Frontend (dev) | 5173 | Vite dev server |
