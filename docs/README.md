# Crab-AI Documentation

This directory contains comprehensive documentation for the Crab-AI backend architecture and system design.

## Files

### `ARCHITECTURE.md`
Complete technical documentation of the Crab-AI backend architecture including:
- **System Overview**: High-level architecture description
- **ASCII Architecture Diagram**: Detailed visual representation using ASCII art
- **Component Details**: In-depth explanation of each system layer
- **Data Flows**: How data moves through the system
- **Security Considerations**: Security measures and best practices
- **Scalability & Performance**: Scaling strategies and optimizations
- **Future Enhancements**: Planned improvements and extensions

### `ARCHITECTURE_DIAGRAM_SPEC.md`
Visual diagram specification for creating professional architecture diagrams using tools like:
- **Draw.io / Diagrams.net**
- **Lucidchart**
- **Microsoft Visio**
- **Other diagramming tools**

Includes:
- Layer-by-layer component specifications
- Color schemes and styling guidelines
- Connection types and arrow specifications
- Grouping and container definitions
- Export format recommendations

## Architecture Overview

The Crab-AI backend is built using a **layered architecture pattern** with the following key layers:

1. **Frontend Layer** - React/Next.js client interfaces
2. **API Gateway Layer** - Gin HTTP router with middleware
3. **Dependency Injection** - Uber Dig container for clean architecture
4. **Handler Layer** - HTTP request handlers
5. **Business Logic Layer** - Core application services
6. **Integration Layer** - LLM and database managers
7. **Data Access Layer** - Repository pattern implementation
8. **Persistence Layer** - MongoDB and Redis storage
9. **External Services** - Third-party integrations
10. **Configuration & Deployment** - Environment management and Docker

## Key Features

- ✅ **Multi-LLM Support**: OpenAI and Gemini integration with extensible architecture
- ✅ **Multi-Database Support**: PostgreSQL, MySQL, ClickHouse, MongoDB, YugabyteDB, Spreadsheets
- ✅ **JWT Authentication**: Secure token-based authentication with refresh tokens
- ✅ **Real-time Streaming**: WebSocket-like streaming for LLM responses
- ✅ **Docker Deployment**: Multi-stage builds with security best practices
- ✅ **Health Monitoring**: Comprehensive health checks and graceful shutdown
- ✅ **Email System**: SMTP integration with HTML templates
- ✅ **File Upload**: CSV/Excel file processing and querying
- ✅ **Dependency Injection**: Clean architecture with Uber Dig
- ✅ **Security**: AES-GCM encryption, password hashing, input validation

## Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Gin HTTP Framework
- **Primary Database**: MongoDB
- **Cache/Sessions**: Redis
- **Authentication**: JWT with refresh tokens
- **LLM Integration**: OpenAI GPT & Google Gemini APIs
- **Email**: SMTP with HTML templates
- **Containerization**: Docker with Alpine Linux
- **Dependency Injection**: Uber Dig
- **Encryption**: AES-GCM for data protection

## Component Relationships

```
Frontend (React/Next.js)
    ↓ HTTP/HTTPS + CORS
API Gateway (Gin Router + Middleware)
    ↓ Dependency Injection
Handlers (Auth, Chat, GitHub, Waitlist, Upload)
    ↓ Service Layer
Business Logic (Auth, Chat, Email, Waitlist, GitHub Services)
    ↓ Integration
LLM Manager + Database Manager
    ↓ Repository Pattern
Data Access Layer (User, Chat, LLM, Token, Waitlist Repos)
    ↓ Persistence
MongoDB (Primary) + Redis (Cache/Sessions)

External Integrations:
    OpenAI API, Gemini API, SMTP Server, Google OAuth, Customer Databases
```

## Usage

### For Developers
1. Read `ARCHITECTURE.md` for detailed technical understanding
2. Use the component and layer descriptions to navigate the codebase
3. Follow the data flow diagrams to understand request processing
4. Reference security considerations for safe development practices

### For System Architects
1. Review the overall architecture diagram for system understanding
2. Use layer descriptions for planning modifications or extensions
3. Reference scalability section for growth planning
4. Consider future enhancements for roadmap planning

### For DevOps/Infrastructure
1. Focus on the Deployment & Infrastructure section
2. Review Docker configuration and health check setup
3. Understand external service dependencies
4. Reference configuration management for environment setup

### For Creating Visual Diagrams
1. Use `ARCHITECTURE_DIAGRAM_SPEC.md` as a guide
2. Follow the color schemes and layout specifications
3. Import the specification into your preferred diagramming tool
4. Export in multiple formats for different use cases

## Maintenance

This documentation should be updated when:
- New services or components are added
- Architecture patterns change
- External integrations are modified
- Security measures are updated
- Performance optimizations are implemented
- New database types are supported
- Additional LLM providers are integrated

## Contributing

When updating this documentation:
1. Maintain consistency with the existing format
2. Update both the ASCII diagram and component descriptions
3. Add new components to the diagram specification
4. Update the technology stack if new dependencies are added
5. Document any breaking changes or migration requirements

This documentation serves as the definitive reference for understanding and working with the Crab-AI backend architecture.