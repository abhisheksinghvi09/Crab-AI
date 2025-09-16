# Crab-AI Backend Architecture

## Overview

The Crab-AI backend is built using a layered architecture pattern with clean separation of concerns, dependency injection, and modern Go best practices. The system provides APIs for database query generation using Large Language Models (LLMs) and supports multiple database types.

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Frontend Layer                             │
├─────────────────────────────────────────────────────────────────────────┤
│  • React/Next.js Client (Port 5173)                                    │
│  • Landing Page                                                        │
│  • Authentication UI                                                   │
│  • Chat Interface                                                      │
└─────────────────────┬───────────────────────────────────────────────────┘
                      │ HTTP/HTTPS
                      │ CORS Enabled
                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          API Gateway Layer                              │
├─────────────────────────────────────────────────────────────────────────┤
│                         Gin HTTP Router                                │
│  ┌─────────────────┬─────────────────┬─────────────────────────────────┐ │
│  │   Middleware    │   Middleware    │        Middleware              │ │
│  │     Stack       │     Stack       │         Stack                  │ │
│  ├─────────────────┼─────────────────┼─────────────────────────────────┤ │
│  │ • CORS          │ • Recovery      │ • Request Logging              │ │
│  │ • Rate Limiting │ • Error Handling│ • Security Headers             │ │
│  │ • JWT Auth      │ • Panic Recovery│ • Request/Response Validation  │ │
│  └─────────────────┴─────────────────┴─────────────────────────────────┘ │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │                        Route Handlers                              │ │
│  │  /health • /api/auth/* • /api/chat/* • /api/github/*              │ │
│  │  /api/waitlist/* • /api/upload/* • /api/oauth/*                   │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
└─────────────────────┬───────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Dependency Injection                            │
├─────────────────────────────────────────────────────────────────────────┤
│                          Uber Dig Container                            │
│  • Handler Registration  • Service Wiring  • Repository Injection     │
└─────────────────────┬───────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Handler Layer                                  │
├─────────────────────────────────────────────────────────────────────────┤
│ ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐ │
│ │    Auth     │    Chat     │   GitHub    │  Waitlist   │   Upload    │ │
│ │   Handler   │   Handler   │   Handler   │   Handler   │   Handler   │ │
│ │             │             │             │             │             │ │
│ │ • Login     │ • Stream    │ • Stats     │ • Join      │ • File      │ │
│ │ • Register  │ • Messages  │ • Analytics │ • Email     │ • CSV       │ │
│ │ • Refresh   │ • Query Gen │             │             │ • Excel     │ │
│ │ • Logout    │ • DB Exec   │             │             │             │ │
│ └─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘ │
└─────────────────────┬───────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Business Logic Layer                            │
├─────────────────────────────────────────────────────────────────────────┤
│ ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐ │
│ │    Auth     │    Chat     │   Email     │  Waitlist   │   GitHub    │ │
│ │   Service   │   Service   │   Service   │   Service   │   Service   │ │
│ │             │             │             │             │             │ │
│ │ • JWT Gen   │ • LLM Comm  │ • SMTP      │ • Queue     │ • Stats     │ │
│ │ • Password  │ • DB Query  │ • Templates │ • Notify    │ • Analytics │ │
│ │ • Session   │ • Stream    │ • Send      │ • Validate  │             │ │
│ │ • Refresh   │ • Execute   │             │             │             │ │
│ └─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘ │
└─────────────────────┬───────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    Integration Layer                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐ │
│ │                      LLM Manager                                    │ │
│ │  ┌─────────────────┬─────────────────┬─────────────────────────────┐ │ │
│ │  │    OpenAI       │     Gemini      │       Future LLMs           │ │ │
│ │  │    Client       │     Client      │       (Extensible)          │ │ │
│ │  │                 │                 │                             │ │ │
│ │  │ • GPT Models    │ • Gemini Pro    │ • Claude, etc.              │ │ │
│ │  │ • Function Call │ • Function Call │ • Pluggable Architecture    │ │ │
│ │  │ • JSON Response │ • JSON Response │                             │ │ │
│ │  └─────────────────┴─────────────────┴─────────────────────────────┘ │ │
│ └─────────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐ │
│ │                    Database Manager                                 │ │
│ │  ┌─────────┬─────────┬─────────┬─────────┬─────────┬─────────────┐  │ │
│ │  │PostgreSQL│ MySQL  │ClickHouse│MongoDB │YugabyteDB│Spreadsheet│  │ │
│ │  │ Driver   │ Driver │  Driver  │ Driver │ Driver   │  Driver   │  │ │
│ │  │          │        │          │        │          │           │  │ │
│ │  │ • GORM   │ • GORM │ • Native │ • Native│ • GORM  │ • Custom  │  │ │
│ │  │ • Schema │ • Schema│ • Schema │ • Schema│ • Schema │ • Schema  │  │ │
│ │  │ • Query  │ • Query │ • Query  │ • Query │ • Query  │ • Query   │  │ │
│ │  └─────────┴─────────┴─────────┴─────────┴─────────┴─────────────┘  │ │
│ └─────────────────────────────────────────────────────────────────────┘ │
└─────────────────────┬───────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Data Access Layer                                 │
├─────────────────────────────────────────────────────────────────────────┤
│ ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐ │
│ │    User     │    Chat     │    LLM      │   Token     │  Waitlist   │ │
│ │ Repository  │ Repository  │ Repository  │ Repository  │ Repository  │ │
│ │             │             │             │             │             │ │
│ │ • CRUD Ops  │ • CRUD Ops  │ • CRUD Ops  │ • CRUD Ops  │ • CRUD Ops  │ │
│ │ • MongoDB   │ • MongoDB   │ • MongoDB   │ • Redis     │ • MongoDB   │ │
│ │ • Caching   │ • Indexing  │ • Logging   │ • TTL       │ • Email Que │ │
│ └─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘ │
└─────────────────────┬───────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     Persistence Layer                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ ┌─────────────────────────────┬─────────────────────────────────────┐   │
│ │          MongoDB             │              Redis                  │   │
│ │      (Primary Database)      │           (Cache & Sessions)       │   │
│ │                              │                                     │   │
│ │ Collections:                 │ Data Types:                         │   │
│ │ • users                      │ • JWT Tokens (TTL)                  │   │
│ │ • chats                      │ • Session Data                      │   │
│ │ • llm_messages               │ • Database Connections              │   │
│ │ • waitlist                   │ • Rate Limiting Counters            │   │
│ │ • connections_metadata       │ • Temporary Query Results           │   │
│ │                              │                                     │   │
│ │ Features:                    │ Features:                           │   │
│ │ • Document Storage           │ • Key-Value Store                   │   │
│ │ • Indexing                   │ • Pub/Sub Messaging                 │   │
│ │ • Aggregation                │ • Atomic Operations                 │   │
│ │ • Transactions               │ • Expiration (TTL)                  │   │
│ └─────────────────────────────┴─────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                        External Services                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────┐ │
│ │   OpenAI    │   Gemini    │    SMTP     │   Google    │  Customer   │ │
│ │     API     │     API     │   Server    │    OAuth    │ Databases   │ │
│ │             │             │             │             │             │ │
│ │ • GPT-4     │ • Gemini    │ • Email     │ • SSO       │ • PostgreSQL│ │
│ │ • GPT-3.5   │   Pro       │   Delivery  │ • Profile   │ • MySQL     │ │
│ │ • Function  │ • Function  │ • Templates │   Data      │ • ClickHouse│ │
│ │   Calling   │   Calling   │ • HTML/Text │             │ • MongoDB   │ │
│ │ • Streaming │ • Streaming │             │             │ • YugabyteDB│ │
│ └─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                     Configuration & Deployment                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐ │
│ │                    Configuration Management                         │ │
│ │                                                                     │ │
│ │ Environment Variables (.env):                                       │ │
│ │ • Database Connections (MongoDB, Redis, PostgreSQL)                │ │
│ │ • LLM API Keys (OpenAI, Gemini)                                     │ │
│ │ • JWT Secrets & Expiration                                          │ │
│ │ • SMTP Configuration                                                │ │
│ │ • CORS Origins                                                      │ │
│ │ • Encryption Keys                                                   │ │
│ │ • OAuth Credentials                                                 │ │
│ └─────────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐ │
│ │                       Docker Deployment                            │ │
│ │                                                                     │ │
│ │ Multi-Stage Build:                                                  │ │
│ │ 1. Builder Stage (golang:1.23-alpine)                              │ │
│ │    • Dependency Download                                            │ │
│ │    • Binary Compilation                                             │ │
│ │    • Optimization (-ldflags="-w -s")                               │ │
│ │                                                                     │ │
│ │ 2. Runtime Stage (alpine:3.19)                                     │ │
│ │    • Minimal Base Image                                             │ │
│ │    • Non-root User                                                  │ │
│ │    • Health Checks                                                  │ │
│ │    • CA Certificates                                                │ │
│ │    • Email Templates                                                │ │
│ └─────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. API Gateway Layer

**Framework**: Gin HTTP Router (Go)
**Port**: 3000 (configurable)

**Middleware Stack**:
- **CORS**: Cross-origin resource sharing with configurable origins
- **Recovery**: Panic recovery with graceful error responses
- **Logging**: Request/response logging with Gin's built-in logger
- **Custom Recovery**: Structured error responses using DTO pattern
- **Authentication**: JWT token validation for protected routes

**Route Groups**:
- `/health` - Health check endpoint
- `/api/auth/*` - Authentication endpoints
- `/api/chat/*` - Chat and query generation endpoints
- `/api/github/*` - GitHub integration endpoints
- `/api/waitlist/*` - Waitlist management endpoints
- `/api/upload/*` - File upload endpoints
- `/api/oauth/*` - OAuth integration endpoints

### 2. Authentication & Authorization

**JWT Implementation**:
- **Access Tokens**: Configurable expiration (default: 10 days)
- **Refresh Tokens**: Extended expiration (default: 30 days)
- **Token Storage**: Redis with TTL for automatic cleanup
- **Security**: HMAC-SHA256 signing with configurable secret

**Password Security**:
- **Hashing**: bcrypt with salt
- **Validation**: Minimum complexity requirements
- **Reset Flow**: Email-based password reset with tokens

**Session Management**:
- **Storage**: Redis-based session store
- **Expiration**: Automatic cleanup with TTL
- **Multi-device**: Support for multiple active sessions

### 3. Business Logic Layer

#### Auth Service
- User registration and login
- JWT token generation and validation
- Password reset functionality
- Session management
- Integration with email service for notifications

#### Chat Service
- LLM integration for query generation
- Database schema analysis
- Query execution and streaming results
- Chat history management
- Support for multiple database types

#### Email Service
- SMTP integration
- HTML email templates
- Welcome emails, password resets
- Waitlist notifications
- Template rendering with dynamic content

#### Waitlist Service
- User waitlist management
- Email notifications
- Priority handling
- Analytics and reporting

#### GitHub Service
- Repository statistics
- Integration analytics
- Performance metrics

### 4. LLM Integration Layer

#### LLM Manager
- **Multi-provider Support**: OpenAI, Gemini (extensible architecture)
- **Configuration Management**: Per-provider settings
- **Database-specific Prompts**: Tailored prompts for each database type
- **Response Schemas**: Structured JSON responses
- **Function Calling**: Advanced LLM capabilities
- **Streaming Support**: Real-time response streaming

#### Supported Database Types
- PostgreSQL
- MySQL
- ClickHouse
- MongoDB
- YugabyteDB
- Spreadsheet (Custom implementation)

### 5. Database Management Layer

#### DB Manager
- **Multi-database Support**: Pluggable driver architecture
- **Schema Fetching**: Automatic schema discovery
- **Query Execution**: Safe query execution with validation
- **Connection Pooling**: Efficient connection management
- **Encryption**: Schema and data encryption

#### Driver Architecture
- **PostgreSQL**: GORM-based driver
- **MySQL**: GORM-based driver
- **ClickHouse**: Native driver implementation
- **MongoDB**: Native MongoDB driver
- **YugabyteDB**: PostgreSQL-compatible driver
- **Spreadsheet**: Custom driver for Excel/CSV files

### 6. Data Access Layer

#### Repository Pattern
- **User Repository**: MongoDB-based user management
- **Chat Repository**: Chat history and metadata
- **LLM Message Repository**: LLM interaction logging
- **Token Repository**: Redis-based token management
- **Waitlist Repository**: Waitlist queue management

#### Caching Strategy
- **Redis**: Primary cache for sessions, tokens, and temporary data
- **MongoDB**: Document storage for persistent data
- **TTL Policies**: Automatic cleanup of expired data

### 7. Configuration Management

#### Environment Configuration
```go
type Environment struct {
    // Server Configuration
    Port                         string
    Environment                  string
    CorsAllowedOrigin           string
    
    // Database Configuration
    MongoURI                    string
    MongoDatabaseName           string
    RedisHost                   string
    RedisPort                   string
    RedisPassword               string
    
    // LLM Configuration
    OpenAIAPIKey                string
    OpenAIModel                 string
    GeminiAPIKey                string
    GeminiModel                 string
    
    // Email Configuration
    SMTPHost                    string
    SMTPPort                    int
    SMTPUser                    string
    SMTPPassword                string
    
    // Security Configuration
    JWTSecret                   string
    SchemaEncryptionKey         string
}
```

#### Validation
- Required field validation
- Format validation for URIs and keys
- Security validation for default credentials
- Encryption key length validation

### 8. Deployment & Infrastructure

#### Docker Configuration
- **Multi-stage Build**: Optimized for production
- **Security**: Non-root user execution
- **Health Checks**: HTTP endpoint monitoring
- **Resource Optimization**: Minimal Alpine-based image
- **Build Optimization**: Compiled binary with size optimization

#### Health Monitoring
- **Endpoint**: `/health`
- **Docker Health Check**: 30-second intervals
- **Dependency Checks**: Database and Redis connectivity
- **Graceful Shutdown**: Proper signal handling

### 9. Email System

#### Templates
- **Welcome Email**: User onboarding
- **Password Reset**: Secure reset flow
- **Waitlist Notification**: Queue updates
- **Enterprise Waitlist**: Business tier notifications

#### Features
- **HTML Templates**: Rich email formatting
- **Template Variables**: Dynamic content injection
- **SMTP Integration**: Configurable email providers
- **Error Handling**: Delivery failure management

### 10. Utilities & Helpers

#### Security Utilities
- **AES-GCM Encryption**: Data encryption/decryption
- **JWT Utilities**: Token generation and validation
- **Password Hashing**: bcrypt implementation
- **Connection Encryption**: Database connection security

#### Data Utilities
- **Type Inference**: Automatic data type detection
- **Hash Utilities**: Data integrity and checksums
- **Secret Management**: Secure secret generation
- **Data Type Conversion**: Type-safe conversions

## Data Flow

### Authentication Flow
1. **User Login**: Client sends credentials to `/api/auth/login`
2. **Validation**: Auth service validates credentials against MongoDB
3. **Token Generation**: JWT service generates access/refresh tokens
4. **Token Storage**: Tokens stored in Redis with TTL
5. **Response**: Client receives tokens and user data

### Query Generation Flow
1. **Database Connection**: User provides database credentials
2. **Schema Fetching**: DB manager retrieves database schema
3. **Schema Encryption**: Schema encrypted and cached
4. **LLM Request**: Chat service sends schema + query to LLM
5. **Query Generation**: LLM generates SQL/query based on schema
6. **Query Execution**: DB manager executes query safely
7. **Response Streaming**: Results streamed back to client

### File Upload Flow
1. **File Upload**: Client uploads CSV/Excel file
2. **File Processing**: Upload service processes file
3. **Schema Generation**: Automatic schema inference
4. **Data Import**: Data imported into temporary storage
5. **Query Interface**: File treated as queryable database

## Security Considerations

### Data Protection
- **Encryption at Rest**: Schema and sensitive data encryption
- **Encryption in Transit**: HTTPS/TLS for all communications
- **API Key Security**: Secure storage of LLM API keys
- **Database Credentials**: Encrypted storage of connection strings

### Access Control
- **JWT Authentication**: Stateless token-based auth
- **Route Protection**: Middleware-based access control
- **Rate Limiting**: Protection against abuse
- **CORS Configuration**: Restricted origin access

### Input Validation
- **SQL Injection Prevention**: Parameterized queries
- **XSS Protection**: Input sanitization
- **Schema Validation**: Structured data validation
- **File Upload Security**: Type and size restrictions

## Scalability & Performance

### Horizontal Scaling
- **Stateless Design**: No server-side session storage
- **Load Balancing**: Multiple instance support
- **Database Sharding**: MongoDB replica sets
- **Cache Distribution**: Redis clustering

### Performance Optimization
- **Connection Pooling**: Efficient database connections
- **Query Caching**: Redis-based query result caching
- **Streaming Responses**: Real-time data delivery
- **CDN Integration**: Static asset optimization

### Monitoring & Observability
- **Health Checks**: Application and dependency monitoring
- **Logging**: Structured logging with Gin middleware
- **Metrics**: Performance and usage analytics
- **Error Tracking**: Comprehensive error handling

## Future Enhancements

### LLM Integration
- Additional LLM providers (Claude, Llama, etc.)
- Custom fine-tuned models
- Multi-modal capabilities
- Advanced function calling

### Database Support
- Additional database types
- NoSQL database enhancement
- Real-time database streaming
- Advanced analytics engines

### Security Enhancements
- OAuth 2.0 / OIDC integration
- Multi-factor authentication
- API rate limiting per user
- Advanced audit logging

### Performance Improvements
- GraphQL API layer
- Advanced caching strategies
- Database connection optimization
- Response compression

This architecture provides a solid foundation for a scalable, secure, and maintainable AI-powered database query system with room for future expansion and enhancement.