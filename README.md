# Webhook Example

A lightweight HTTP webhook server example designed for SIEM (Security Information and Event Management) systems. This service receives webhook payloads and logs them either to stdout or to rotating log files.

## Features

- Simple HTTP webhook endpoint (`/api/hook`)
- Configurable output (stdout or file logging)
- Log file rotation with configurable size, backup count, and age
- JSON or plain text log formatting
- Authorization via shared secret
- Health check endpoint
- Docker support
- Kubernetes deployment ready

## Project Structure

```
├── cmd/
│   └── main.go           # Main application entry point
├── vendor/               # Vendored dependencies
├── deployment.yaml       # Kubernetes deployment configuration
├── Dockerfile           # Docker image configuration
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Makefile            # Build automation
└── README.md           # This file
```

## Quick Start

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/adaptive-scale/webhook-example.git
cd webhook-example
```

2. Start the server:
```bash
# Using Go directly
go run cmd/main.go

# Or using Make (with default development settings)
make start
```

3. The server will start on port 8080 by default.

### Using Docker

```bash
docker run -p 8080:8080 -e SHARED_SECRET=your-secret <registry>/webhook-example:latest
```

## Configuration

The application is configured via environment variables:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `SHARED_SECRET` | Authorization token for webhook requests | (required) |
| `OUTPUT_TYPE` | Output destination: "stdout" or "file" | `stdout` |
| `FILE_LOCATION` | Log file path (when OUTPUT_TYPE=file) | `/tmp/adaptive.log` |
| `MAX_SIZE` | Max log file size in MB before rotation | `10` |
| `MAX_BACKUP` | Number of rotated log files to keep | `3` |
| `MAX_AGE` | Max days to retain old log files | `28` |
| `PORT` | HTTP server port | `8080` |
| `FORMATTER` | Log format: "json" or plain text | plain text |

## API Endpoints

### Webhook Endpoint
- **URL**: `/api/hook`
- **Method**: `POST`
- **Headers**: `Authorization: <SHARED_SECRET>`
- **Description**: Receives webhook payloads and logs them

### Health Check
- **URL**: `/healthz`
- **Method**: `GET`
- **Description**: Returns "ok" for health monitoring

## Usage Examples

### Basic webhook request:
```bash
curl -X POST \
  -H "Authorization: your-shared-secret" \
  -H "Content-Type: application/json" \
  -d '{"event": "alert", "message": "Security incident detected"}' \
  http://localhost:8080/api/hook
```

### Check server health:
```bash
curl http://localhost:8080/healthz
```

### File logging configuration:
```bash
export SHARED_SECRET=my-secret-key
export OUTPUT_TYPE=file
export FILE_LOCATION=/var/log/siem/webhook.log
export FORMATTER=json
export MAX_SIZE=100
export MAX_BACKUP=5
export MAX_AGE=30

./webhook-example
```

## Building

### Local Build
```bash
go build -o webhook-example ./cmd/main.go
```

### Using Makefile
```bash
make build-siem-webhook
```

This will build the binary for linux/amd64 as `webhook-example`.

## Deployment

### Docker Compose
```yaml
version: '3.8'
services:
  webhook-example:
    image: <registry>/siemwebhook:latest
    ports:
      - "8080:8080"
    environment:
      - SHARED_SECRET=your-shared-secret
      - OUTPUT_TYPE=file
      - FILE_LOCATION=/app/logs/webhook.log
      - FORMATTER=json
    volumes:
      - ./logs:/app/logs
```

### Kubernetes
Apply the provided Kubernetes deployment:
```bash
kubectl apply -f deployment.yaml
```

Make sure to update the `SHARED_SECRET` in the deployment file before applying.

## Log Rotation

When using file output (`OUTPUT_TYPE=file`), the application automatically rotates logs based on:
- **Size**: When log file exceeds `MAX_SIZE` MB
- **Time**: Removes files older than `MAX_AGE` days
- **Count**: Keeps only `MAX_BACKUP` number of rotated files
- **Compression**: Automatically compresses rotated files

## Security

- The webhook endpoint requires a shared secret via the `Authorization` header
- Only POST requests are accepted on the webhook endpoint
- Unauthorized requests return HTTP 401
- Invalid methods return HTTP 405

## Requirements

- Go 1.23.2 or later
- Docker (optional, for containerized deployment)
- Kubernetes (optional, for cluster deployment)

## Dependencies

- [logrus](https://github.com/sirupsen/logrus) - Structured logging
- [lumberjack](https://gopkg.in/natefinch/lumberjack.v2) - Log rotation

## Testing

You can test the webhook server using the provided curl examples or any HTTP client:

1. Start the server:
   ```bash
   # Using make (uses development-secret by default)
   make start
   
   # Or manually with custom secret
   SHARED_SECRET=test-secret go run cmd/main.go
   ```

2. Send a test webhook:
   ```bash
   curl -X POST \
     -H "Authorization: development-secret" \
     -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/api/hook
   ```

3. Check health:
   ```bash
   curl http://localhost:8080/healthz
   ```

## License

This project is part of the Adaptive Scale ecosystem.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
