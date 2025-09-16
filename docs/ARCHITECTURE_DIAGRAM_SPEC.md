# Crab-AI Architecture Visual Diagram Specification

This document provides a structured specification for creating the visual architecture diagram using tools like Draw.io, Lucidchart, or similar diagramming tools.

## Diagram Components

### Layer 1: Frontend (Top)
**Color**: Light Blue (#E3F2FD)
**Position**: Top of diagram
**Components**:
- React/Next.js Client (Port 5173)
- Landing Page
- Authentication UI  
- Chat Interface

### Layer 2: API Gateway (Below Frontend)
**Color**: Light Green (#E8F5E8)
**Position**: Below Frontend layer
**Components**:
- Gin HTTP Router (Port 3000)
- Middleware Stack:
  - CORS
  - Recovery  
  - Logging
  - JWT Auth
- Route Groups:
  - /health
  - /api/auth/*
  - /api/chat/*
  - /api/github/*
  - /api/waitlist/*
  - /api/upload/*

### Layer 3: Dependency Injection (Below API Gateway)
**Color**: Light Purple (#F3E5F5)
**Position**: Below API Gateway
**Components**:
- Uber Dig Container
- Handler Registration
- Service Wiring
- Repository Injection

### Layer 4: Handler Layer (Below DI)
**Color**: Light Orange (#FFF3E0)
**Position**: Below Dependency Injection
**Components**:
- Auth Handler (Login, Register, Refresh, Logout)
- Chat Handler (Stream, Messages, Query Gen, DB Exec)
- GitHub Handler (Stats, Analytics)
- Waitlist Handler (Join, Email)
- Upload Handler (File, CSV, Excel)

### Layer 5: Business Logic (Below Handlers)
**Color**: Light Yellow (#FFFDE7)
**Position**: Below Handler Layer
**Components**:
- Auth Service (JWT Gen, Password, Session, Refresh)
- Chat Service (LLM Comm, DB Query, Stream, Execute)
- Email Service (SMTP, Templates, Send)
- Waitlist Service (Queue, Notify, Validate)
- GitHub Service (Stats, Analytics)

### Layer 6: Integration Layer (Below Business Logic)
**Color**: Light Cyan (#E0F2F1)
**Position**: Below Business Logic
**Components**:
- LLM Manager:
  - OpenAI Client (GPT Models, Function Call, JSON Response)
  - Gemini Client (Gemini Pro, Function Call, JSON Response)
  - Future LLMs (Claude, etc., Pluggable Architecture)
- Database Manager:
  - PostgreSQL Driver (GORM, Schema, Query)
  - MySQL Driver (GORM, Schema, Query)
  - ClickHouse Driver (Native, Schema, Query)
  - MongoDB Driver (Native, Schema, Query)
  - YugabyteDB Driver (GORM, Schema, Query)
  - Spreadsheet Driver (Custom, Schema, Query)

### Layer 7: Data Access (Below Integration)
**Color**: Light Pink (#FCE4EC)
**Position**: Below Integration Layer
**Components**:
- User Repository (CRUD Ops, MongoDB, Caching)
- Chat Repository (CRUD Ops, MongoDB, Indexing)
- LLM Repository (CRUD Ops, MongoDB, Logging)
- Token Repository (CRUD Ops, Redis, TTL)
- Waitlist Repository (CRUD Ops, MongoDB, Email Queue)

### Layer 8: Persistence (Below Data Access)
**Color**: Light Gray (#F5F5F5)
**Position**: Below Data Access
**Components**:
- MongoDB (Primary Database):
  - Collections: users, chats, llm_messages, waitlist, connections_metadata
  - Features: Document Storage, Indexing, Aggregation, Transactions
- Redis (Cache & Sessions):
  - Data Types: JWT Tokens (TTL), Session Data, Database Connections, Rate Limiting Counters, Temporary Query Results
  - Features: Key-Value Store, Pub/Sub Messaging, Atomic Operations, Expiration (TTL)

### Layer 9: External Services (Side/Bottom)
**Color**: Light Red (#FFEBEE)
**Position**: Side or bottom of diagram
**Components**:
- OpenAI API (GPT-4, GPT-3.5, Function Calling, Streaming)
- Gemini API (Gemini Pro, Function Calling, Streaming)
- SMTP Server (Email Delivery, Templates, HTML/Text)
- Google OAuth (SSO, Profile Data)
- Customer Databases (PostgreSQL, MySQL, ClickHouse, MongoDB, YugabyteDB)

### Layer 10: Configuration & Deployment (Side)
**Color**: Light Brown (#EFEBE9)
**Position**: Side of diagram
**Components**:
- Configuration Management:
  - Environment Variables (.env)
  - Database Connections
  - LLM API Keys
  - JWT Secrets & Expiration
  - SMTP Configuration
  - CORS Origins
  - Encryption Keys
  - OAuth Credentials
- Docker Deployment:
  - Multi-Stage Build
  - Builder Stage (golang:1.23-alpine)
  - Runtime Stage (alpine:3.19)
  - Health Checks
  - Non-root User
  - CA Certificates

## Connection Types

### HTTP/HTTPS Connections
**Color**: Blue arrows
**Style**: Solid lines
**Connections**:
- Frontend ↔ API Gateway
- API Gateway ↔ External Services (OpenAI, Gemini, SMTP, OAuth)

### Internal API Calls
**Color**: Green arrows  
**Style**: Solid lines
**Connections**:
- API Gateway → Handlers
- Handlers → Services
- Services → Repositories
- Services → LLM Manager
- Services → DB Manager

### Database Connections
**Color**: Purple arrows
**Style**: Dashed lines
**Connections**:
- Repositories → MongoDB
- Repositories → Redis
- DB Manager → Customer Databases

### Configuration Flow
**Color**: Orange arrows
**Style**: Dotted lines
**Connections**:
- Configuration → All Components

## Grouping and Containers

### Main Application Container
**Border**: Thick black border
**Contents**: Layers 2-7 (API Gateway through Data Access)

### External Dependencies Container  
**Border**: Thick red border
**Contents**: Layer 9 (External Services)

### Infrastructure Container
**Border**: Thick brown border
**Contents**: Layer 8 (Persistence) + Layer 10 (Configuration & Deployment)

## Text Labels and Annotations

### Key Features Box (Top Right)
```
Key Features:
✓ Multi-LLM Support (OpenAI, Gemini)
✓ Multi-Database Support (6 types)
✓ JWT Authentication
✓ Real-time Streaming
✓ Docker Deployment
✓ Health Monitoring
✓ Email Notifications
✓ File Upload Support
```

### Technology Stack Box (Bottom Right)
```
Technology Stack:
• Go 1.23
• Gin HTTP Framework
• MongoDB (Primary DB)
• Redis (Cache/Sessions)
• Docker & Alpine Linux
• JWT for Authentication
• SMTP for Email
• AES-GCM Encryption
```

### Data Flow Arrows (Between Layers)
Add numbered arrows showing the main data flows:
1. User Request → API Gateway
2. API Gateway → Handler → Service
3. Service → LLM Manager → External LLM
4. Service → DB Manager → Customer DB
5. Service → Repository → MongoDB/Redis
6. Response ← All layers back to Frontend

## Color Scheme Summary
- Frontend: Light Blue (#E3F2FD)
- API Gateway: Light Green (#E8F5E8)  
- Dependency Injection: Light Purple (#F3E5F5)
- Handlers: Light Orange (#FFF3E0)
- Business Logic: Light Yellow (#FFFDE7)
- Integration: Light Cyan (#E0F2F1)
- Data Access: Light Pink (#FCE4EC)
- Persistence: Light Gray (#F5F5F5)
- External Services: Light Red (#FFEBEE)
- Config/Deployment: Light Brown (#EFEBE9)

## Diagram Dimensions
- **Recommended Size**: 1920x1080 (landscape orientation)
- **Layer Height**: Approximately 120-150 pixels each
- **Component Width**: 150-200 pixels each
- **Spacing**: 20-30 pixels between components
- **Margins**: 50 pixels on all sides

## Export Formats
When creating the visual diagram, export in the following formats:
- PNG (for documentation embedding)
- SVG (for web display)
- PDF (for presentations)
- Source format (.drawio, .lucid, etc. for future editing)

This specification can be used to create a professional architecture diagram in any visual diagramming tool while maintaining consistency with the documented architecture.