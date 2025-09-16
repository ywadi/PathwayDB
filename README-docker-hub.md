# PathwayDB - Docker Hub Deployment

This guide shows how to run PathwayDB using the pre-built images from Docker Hub.

## 🐳 Quick Start with Docker Hub Images

### Prerequisites
- Docker and Docker Compose installed
- Create a `data` directory: `mkdir -p data`

### Run with Docker Hub Images
```bash
# Pull and run the latest images from Docker Hub
docker compose -f docker-compose.hub.yml up -d

# View logs
docker compose -f docker-compose.hub.yml logs -f

# Stop services
docker compose -f docker-compose.hub.yml down
```

### Available Images
- `yousefwadi/pathwaydb-redis:latest` - Redis server with PathwayDB protocol
- `yousefwadi/pathwaydb-backend:latest` - WebSocket backend service
- `yousefwadi/pathwaydb-frontend:latest` - React frontend IDE

### Environment Variables
You can customize ports using environment variables:
```bash
export REDIS_PORT=6379
export WEBSOCKET_PORT=8081
export PORT=3000

docker compose -f docker-compose.hub.yml up -d
```

### Access Points
- **Frontend IDE**: http://localhost:3000
- **WebSocket Backend**: ws://localhost:8081
- **Redis Protocol**: localhost:6379

## 🔄 Development vs Production

- **Development**: Use `docker-compose.yml` (builds from source)
- **Production**: Use `docker-compose.hub.yml` (uses Docker Hub images)

## 📦 Image Tags
Each image is available with:
- `latest` - Latest stable version
- `<commit-sha>` - Specific commit versions for reproducible deployments
