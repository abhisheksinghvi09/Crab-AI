# Crab-AI Backend Architecture

## Overview

The Crab-AI backend is a modern, scalable Go-based application that provides a robust foundation for AI-powered data analysis and chat functionality. This enhanced architecture implements industry best practices for security, scalability, and maintainability.

## Architecture Highlights

### 🏗️ Modular Architecture
- **Clean Architecture**: Separation of concerns with distinct layers (handlers, services, repositories, models)
- **Dependency Injection**: Using Uber's dig framework for loose coupling
- **Interface-driven design**: Easy testing and extensibility

### 🔒 Enhanced Security
- **Rate Limiting**: Configurable per-IP rate limiting
- **Request Validation**: Content type and size validation
- **Security Headers**: Comprehensive security headers middleware
- **Input Sanitization**: XSS and injection prevention
- **JWT Authentication**: Secure token-based authentication

### 🚀 Scalability Features
- **Database Connection Pooling**: Optimized connection management
- **Redis Caching**: Distributed caching for session management
- **LLM Provider Failover**: Automatic failover between AI providers
- **Health Checks**: Comprehensive monitoring endpoints

### 🤖 Enhanced LLM Integration
- **Multi-Provider Support**: OpenAI and Gemini with easy extensibility
- **Health Monitoring**: Automatic health checks for LLM providers
- **Request Metrics**: Performance tracking and monitoring
- **Failover Mechanism**: Automatic switching between providers

## Key Features

### 1. Enhanced Error Handling
- **Structured Errors**: Typed error system with proper HTTP status codes
- **Panic Recovery**: Graceful panic recovery with stack traces
- **Request Tracing**: Request ID tracking for debugging

### 2. Comprehensive Health Checks
- **Basic Health**: Simple service status check
- **Detailed Health**: Comprehensive dependency status
- **Readiness Probe**: Kubernetes-compatible readiness checks
- **Liveness Probe**: Kubernetes-compatible liveness checks

### 3. Security Enhancements
- **Rate Limiting**: Configurable request rate limiting
- **CORS Protection**: Enhanced CORS with origin validation
- **Security Headers**: XSS, clickjacking, and content-type protection
- **Input Validation**: Request size and content type validation

### 4. LLM Provider Management
- **Health Monitoring**: Continuous health checks for AI providers
- **Automatic Failover**: Switch between providers on failure
- **Performance Metrics**: Request tracking and performance monitoring
- **Configurable Timeouts**: Customizable request timeouts

## Development Setup

### Prerequisites
- Go 1.23+
- Docker and Docker Compose
- MongoDB 7.0+
- Redis 7.0+

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/abhisheksinghvi09/Crab-AI.git
   cd Crab-AI
   ```

2. **Set up environment variables**
   ```bash
   cd backend
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run with Docker Compose (Recommended)**
   ```bash
   # From project root
   docker-compose up -d
   ```

4. **Or run locally**
   ```bash
   cd backend
   go mod download
   go run cmd/main.go
   ```

### Development Commands

```bash
# Run tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Format code
go fmt ./...

# Build application
go build -o main cmd/main.go
```

## API Endpoints

### Health Check Endpoints
- `GET /health` - Basic health check
- `GET /health/detailed` - Detailed health with dependencies
- `GET /health/ready` - Kubernetes readiness probe
- `GET /health/live` - Kubernetes liveness probe

### Authentication Endpoints
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Token refresh
- `POST /api/auth/logout` - User logout

### Chat Endpoints
- `GET /api/chats` - List user chats
- `POST /api/chats` - Create new chat
- `GET /api/chats/:id` - Get chat details
- `POST /api/chats/:id/messages` - Send message

## Testing

The project includes comprehensive testing at multiple levels:

### Unit Tests
- **Service Layer**: Business logic testing
- **Middleware**: Security and validation testing
- **Utilities**: Helper function testing
- **LLM Manager**: AI provider management testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Deployment

### Docker Deployment

```bash
# Build production image
docker build -t crab-ai/backend:latest ./backend

# Run with environment variables
docker run -p 3000:3000 \
  -e CRAB_MONGODB_URI=mongodb://your-mongo-host:27017/crab-ai \
  -e CRAB_REDIS_HOST=your-redis-host \
  -e JWT_SECRET=your-jwt-secret \
  crab-ai/backend:latest
```

## Architecture Improvements

This enhanced backend architecture provides:

1. **Production-Ready Security**: Comprehensive security middleware stack
2. **Scalable LLM Integration**: Multi-provider support with health monitoring and failover
3. **Robust Error Handling**: Structured error responses with proper HTTP status codes
4. **Comprehensive Testing**: Unit and integration tests with CI/CD pipeline
5. **Container Optimization**: Multi-stage Docker builds for production deployment
6. **Monitoring & Observability**: Health checks and metrics for operations teams

The implementation follows Go best practices and provides a solid foundation for scaling the Crab-AI platform.