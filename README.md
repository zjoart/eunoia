# Eunoia - AI Mental Wellbeing Assistant

An AI agent for mental wellbeing that performs emotional check-ins, analyzes reflections, and provides supportive, context-aware responses. Built in Go with A2A protocol support for integration with any compatible messaging platform.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Development](#development)

## Features

- **Emotional Check-ins**: Track mood scores (1-10 scale) and emotional states over time
- **Reflection Analysis**: AI-powered sentiment analysis and theme extraction from journal entries
- **Context-Aware Conversations**: Maintains conversation history with personalized, empathetic responses
- **Mood Trend Analysis**: Identifies patterns in emotional wellbeing (improving, declining, stable)
- **Supportive AI Responses**: Powered by Google's Gemini AI for empathetic interactions
- **A2A Protocol**: Standard protocol for integration with any compatible messaging platform

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24.2 or higher** - [Download Go](https://go.dev/dl/)
- **MySQL 8.0 or higher** - [Download MySQL](https://dev.mysql.com/downloads/)
- **Git** - For cloning the repository

### Required Accounts

- **Google AI Studio** - Get your Gemini API key from [Google AI Studio](https://makersuite.google.com/app/apikey)

## Quick Start

Follow these steps to get Eunoia running on your local machine:

### 1. Clone the Repository

```bash
git clone https://github.com/zjoart/eunoia.git
cd eunoia
```

### 2. Install Dependencies

```bash
make tidy
```

This will download all required Go packages including:
- `gorilla/mux` - HTTP routing
- `golang-migrate/migrate` - Database migrations
- `go-sql-driver/mysql` - MySQL driver
- `uber-go/zap` - Structured logging
- `joho/godotenv` - Environment variable loading

### 3. Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` with your configuration values. See `.env.example` for all required variables.

## Configuration

Update `.env` with your values:

- `PORT` - Server port (default: 8080)
- `APP_ENV` - Environment (development/production)
- `DB_*` - Database connection details
- `GEMINI_API_KEY` - Your Gemini API key from Google AI Studio

## Database Setup

### Migration Commands

The project includes a full migration system:

```bash
# Apply all migrations
make migrate-up

# Rollback all migrations
make migrate-down

# Check current migration version
make migrate-version

# Apply specific number of migrations
make migrate-steps STEPS=1      # Apply next migration
make migrate-steps STEPS=-1     # Rollback last migration

# Force migration version (use with caution)
make migrate-force VERSION=1
```


## Running the Application

### Development Mode

```bash
# Start the server
make run
```

The server will start on `http://localhost:8080`

### Production Build

```bash
# Build the binary
go build -o bin/eunoia cmd/app/main.go

# Run the binary
./bin/eunoia
```

### Verify Application is Running

```bash
# Test agent health endpoint
curl http://localhost:8080/agent/health

# Expected JSON output with status information
```

## API Endpoints

### 1. Agent Health Check

**Endpoint:** `GET /agent/health`

**Description:** Returns agent-specific health information.

**Example:**
```bash
curl http://localhost:8080/agent/health
```

**Response:**
```json
{
  "status": "healthy",
  "agent": "eunoia",
  "service": "mental wellbeing assistant"
}
```

### 3. A2A Agent Endpoint

**Endpoint:** `POST /a2a/agent/eunoia`

**Description:** Main endpoint for receiving messages via A2A protocol and responding with AI-generated, context-aware responses.

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "message": "How are you feeling today?",
  "userId": "user-12345",
  "channelId": "channel-uuid",
  "messageId": "message-uuid",
  "timestamp": "2025-11-03T10:30:00Z"
}
```

**Request Fields:**
- `message` (required, string): The user's message text
- `userId` (required, string): Unique user identifier
- `channelId` (optional, string): Channel identifier
- `messageId` (optional, string): Unique message identifier
- `timestamp` (optional, string): ISO 8601 timestamp

**Success Response (200 OK):**
```json
{
  "response": "I'm here to support you. How would you rate your mood today on a scale of 1-10?",
  "messageId": "1730628000123456789",
  "metadata": {
    "agent": "eunoia",
    "timestamp": "2025-11-03T10:30:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Bad Request",
  "message": "message field is required",
  "status": 400
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "error": "Internal Server Error",
  "message": "failed to process message",
  "status": 500
}
```

**Example with cURL:**
```bash
curl -X POST http://localhost:8080/a2a/agent/eunoia \
  -H "Content-Type: application/json" \
  -d '{
    "message": "I am feeling anxious today",
    "userId": "user-123",
    "channelId": "channel-456",
    "messageId": "msg-789",
    "timestamp": "2025-11-03T10:30:00Z"
  }'
```

### Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Missing or invalid required fields |
| 405 | Method Not Allowed - Must use POST |
| 500 | Internal Server Error - Processing failure |

## Database Schema

The application uses MySQL with 5 main tables:
- `users` - User accounts and identifiers
- `emotional_checkins` - Mood scores and emotional states
- `reflections` - Journal entries with AI analysis
- `conversation_history` - Chat history for context
- `user_preferences` - User settings and preferences

See `migrations/000001_create_tables.up.sql` for full schema details.

## A2A Protocol Integration

The application implements the A2A (Agent-to-Agent) protocol, allowing integration with any compatible messaging platform.

### Integration Steps

1. **Deploy** your application to a publicly accessible server
2. **Configure** your messaging platform to send POST requests to `/a2a/agent/eunoia`
3. **Test** the integration with sample messages

The A2A endpoint handles:
- User identification and management
- Context building from user history
- AI-powered response generation
- Conversation history storage

## Development

### Project Structure

```
eunoia/
├── cmd/
│   ├── app/main.go              # Application entry point
│   ├── migrate/main.go          # Migration CLI
│   └── routes/routes.go         # HTTP routes setup
├── internal/
│   ├── agent/                   # AI service
│   ├── user/                    # User feature
│   ├── checkin/                 # Check-in feature
│   ├── reflection/              # Reflection feature
│   ├── conversation/            # Conversation feature + handler
│   ├── config/                  # Configuration
│   ├── database/                # Database connection
│   └── middleware/              # HTTP middleware
├── migrations/                  # SQL migrations
├── pkg/
│   └── logger/                  # Logging utilities
├── .env                         # Environment variables (gitignored)
├── .env.example                 # Environment template
├── go.mod                       # Go module definition
├── go.sum                       # Go dependencies checksum
├── Makefile                     # Build and run commands
└── telex-workflow.json          # Telex integration config
```

### Available Make Commands

```bash
make help                   # Show all available commands
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with detailed output
make test-log

# Run tests without cache (fresh run)
make test-force

# Run specific test function
make test-function TEST=TestFunctionName

# Run with race condition detection
make test-race

# Run with coverage report
make test-ci
```


## License

MIT License - See LICENSE file for details.

**Note**: This is a personal wellbeing tool and should not replace professional mental health care. If you're experiencing a mental health crisis, please contact a mental health professional or crisis hotline immediately.
