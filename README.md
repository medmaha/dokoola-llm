# Dokoola LLM Service (Go)

A high-performance Go implementation of the Dokoola LLM service, providing AI-powered content generation for job categorization, text completion, and prompt-based content creation.

## Features

- **Job Categorization**: Automatically categorize job postings using LLM analysis
- **Text Completion**: Generate text completions with context-aware prompts
- **Prompt Templates**: Pre-built templates for:
  - Talent bio generation
  - Client "About Us" sections
  - Job descriptions
  - Proposal cover letters
- **Authentication**: Header-based service authentication via config.ini
- **Cerebras Cloud Integration**: Leverages Cerebras LLM API for fast inference

## Architecture

```
llm-go/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── clients/         # Backend API clients
│   ├── config/          # Configuration management
│   ├── constants/       # System constants and messages
│   ├── handlers/        # HTTP request handlers
│   ├── llm/            # LLM client implementation
│   ├── middleware/      # Authentication & logging middleware
│   ├── models/         # Data models
│   └── prompts/        # Prompt template builders
├── pkg/                # Public packages (if any)
├── Dockerfile          # Container build configuration
├── Makefile           # Build automation
├── go.mod             # Go module definition
└── config.ini         # Service authentication registry
```

## API Endpoints

All endpoints (except health check) require authentication via custom headers.

### Health Check
- `GET /health` - Service health status

### Jobs
- `POST /api/v1/llm/chat/jobs/categorize` - Categorize job postings

### Text Completion
- `POST /api/v1/llm/chat/completion` - Generate text completions

### Prompt Generation
- `POST /api/v1/llm/chat/actions/generate-prompt` - Generate content from templates

## Setup

### Prerequisites

- Go 1.23 or higher
- Docker (optional, for containerized deployment)
- Cerebras API key
- Access to Dokoola backend API

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `Dokoola LLM Service` |
| `APP_VERSION` | Application version | `0.1.0` |
| `DEBUG` | Enable debug mode | `true` |
| `API_PREFIX` | API route prefix | `/api/v1` |
| `HOST` | Server host | `0.0.0.0` |
| `PORT` | Server port | `8000` |
| `LOG_LEVEL` | Logging level | `info` |
| `LLM_API_KEY` | Cerebras API key (required) | - |
| `BACKEND_SERVER_API` | Backend API URL (required) | - |

### Service Authentication (config.ini)

```ini
[SERVICE_DKL001]
host = https://frontend.dokoola.com
client_name = FRONTEND_CLIENT
secret_hash = your_secret_hash

[SERVICE_DKL002]
host = https://backend.dokoola.com
client_name = BACKEND_CLIENT
secret_hash = another_secret_hash
```

Services authenticate by providing three headers:
- Service key (e.g., "DKL001")
- Client name (must match config)
- Secret hash (must match config)

## License

Proprietary - Dokoola Platform

## Support

For issues or questions, contact the Dokoola development team.
