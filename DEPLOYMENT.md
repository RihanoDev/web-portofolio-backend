# GitHub Environment Secrets Setup

This repository uses **shared credentials** for both development and production environments.
The only difference between dev and prod is the **port** where the service runs.

## Required GitHub Secrets (Environment: production)

Add these secrets in GitHub repository settings → Environments → production:

### Server/SSH Access

- `SERVER_HOST` - IP address or domain of your deployment server
- `SERVER_PORT` - SSH port (usually 22)
- `SERVER_USER` - SSH username
- `SERVER_SSH_KEY` - Private SSH key for authentication

### Docker Hub

- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password or access token

### Database (Shared for Dev & Prod)

- `DB_HOST` - Database host (e.g., `rihanodev.com`)
- `DB_PORT` - Database port (usually `5432` for PostgreSQL)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name (e.g., `rihanodev_db`)

### Application

- `JWT_SECRET` - JWT signing secret (min 32 characters, strong random string)

## Port Assignment

The workflows automatically assign different ports:

| Environment | Port | Container Name         | Trigger             |
| ----------- | ---- | ---------------------- | ------------------- |
| Production  | 1200 | web-porto-backend-prod | Tag push (v\*)      |
| Development | 2200 | web-porto-backend-dev  | Branch: development |

## Deployment Flow

### Production Deployment

```bash
git tag v1.0.0
git push origin v1.0.0
```

- Triggers: `.github/workflows/deploy-prod.yml`
- Builds: `rihanodev/web-porto-backend:latest` & `:v1.0.0`
- Runs on: Port 1200
- Mode: `GIN_MODE=release`, `APP_DEBUG=false`

### Development Deployment

```bash
git checkout development
git push origin development
```

- Triggers: `.github/workflows/deploy-dev.yml`
- Builds: `rihanodev/web-porto-backend:dev-latest` & `:dev-{sha}`
- Runs on: Port 2200
- Mode: `GIN_MODE=debug`, `APP_DEBUG=true`

## Important Notes

1. **Same Database**: Both dev and prod use the same database credentials
2. **Same Server**: Both containers run on the same Docker host
3. **Port Isolation**: Different ports prevent conflicts
4. **No Volume Mounts**: Configuration via environment variables only (no config.json needed in production)

## Environment Variables Passed to Containers

### Production Container

```bash
GIN_MODE=release
APP_DEBUG=false
SERVER_PORT=8080
DB_HOST=${{ secrets.DB_HOST }}
DB_PORT=${{ secrets.DB_PORT }}
DB_USER=${{ secrets.DB_USER }}
DB_PASSWORD=${{ secrets.DB_PASSWORD }}
DB_NAME=${{ secrets.DB_NAME }}
JWT_SECRET=${{ secrets.JWT_SECRET }}
```

### Development Container

```bash
GIN_MODE=debug
APP_DEBUG=true
SERVER_PORT=8080
DB_HOST=${{ secrets.DB_HOST }}
DB_PORT=${{ secrets.DB_PORT }}
DB_USER=${{ secrets.DB_USER }}
DB_PASSWORD=${{ secrets.DB_PASSWORD }}
DB_NAME=${{ secrets.DB_NAME }}
JWT_SECRET=${{ secrets.JWT_SECRET }}
```

## Verification

After deployment, the workflow automatically:

1. Checks port availability
2. Pulls latest image
3. Stops old container
4. Starts new container
5. Waits 10 seconds for startup
6. Checks container status
7. Shows container logs
8. Tests health endpoint

## Troubleshooting

### Check container logs

```bash
# Production
docker logs web-porto-backend-prod --tail 50

# Development
docker logs web-porto-backend-dev --tail 50
```

### Check container status

```bash
docker ps | grep web-porto-backend
```

### Test API manually

```bash
# Production
curl http://localhost:1200/health

# Development
curl http://localhost:2200/health
```

### Restart container

```bash
# Production
docker restart web-porto-backend-prod

# Development
docker restart web-porto-backend-dev
```
