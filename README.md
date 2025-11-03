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

### 2. A2A Agent Endpoint

**Endpoint:** `POST /a2a/agent/eunoia`

**Description:** Main endpoint for receiving messages via A2A protocol and responding with AI-generated, context-aware responses. Uses JSON-RPC 2.0 format.

**Headers:**
```
Content-Type: application/json
```

**Request Body (JSON-RPC 2.0):**
```json
{
  "jsonrpc": "2.0",
  "id": "unique-request-id",
  "method": "message/send",
  "params": {
    "message": {
      "kind": "message",
      "role": "user",
      "parts": [
        {
          "kind": "text",
          "text": "Hello, how are you feeling today?"
        }
      ],
      "metadata": {
        "telex_user_id": "user-uuid",
        "telex_channel_id": "channel-uuid",
        "org_id": "org-uuid"
      },
      "messageId": "message-uuid"
    },
    "configuration": {
      "acceptedOutputModes": ["text/plain"],
      "historyLength": 0,
      "blocking": true
    }
  }
}
```

**Request Fields:**
- `jsonrpc` (required, string): Must be "2.0"
- `id` (required, string): Unique request identifier
- `method` (required, string): Must be "message/send"
- `params.message.parts` (required, array): Message content parts
- `params.message.metadata.telex_user_id` (required, string): User identifier
- `params.configuration` (optional, object): Processing configuration

**Success Response (JSON-RPC 2.0):**
```json
{
  "jsonrpc": "2.0",
  "id": "unique-request-id",
  "result": {
    "message": {
      "kind": "message",
      "role": "assistant",
      "parts": [
        {
          "kind": "text",
          "text": "Hello! I'm Eunoia, here to support your mental wellbeing. How are you feeling today?"
        }
      ],
      "metadata": {
        "agent": "eunoia",
        "timestamp": "2025-11-03T14:43:11Z"
      },
      "messageId": "response-message-id"
    }
  }
}
```

**Error Response (JSON-RPC 2.0):**
```json
{
  "jsonrpc": "2.0",
  "id": "unique-request-id",
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": "telex_user_id is required"
  }
}
```

**Example with cURL:**
```bash
curl -X POST http://localhost:8080/a2a/agent/eunoia \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "req-123",
    "method": "message/send",
    "params": {
      "message": {
        "kind": "message",
        "role": "user",
        "parts": [{"kind": "text", "text": "Hello"}],
        "metadata": {
          "telex_user_id": "user-123",
          "telex_channel_id": "channel-456"
        },
        "messageId": "msg-789"
      },
      "configuration": {
        "acceptedOutputModes": ["text/plain"],
        "blocking": true
      }
    }
  }'
```

### Error Codes

| Code | Description |
|------|-------------|
| -32700 | Parse error - Invalid JSON |
| -32600 | Invalid Request - Wrong JSON-RPC format |
| -32601 | Method not found - Unsupported method |
| -32602 | Invalid params - Missing required fields |
| -32603 | Internal error - Server processing error |

## Database Schema

The application uses MySQL with 5 main tables:
- `users` - User accounts and identifiers
- `emotional_checkins` - Mood scores and emotional states
- `reflections` - Journal entries with AI analysis
- `conversation_history` - Chat history for context
- `user_preferences` - User settings and preferences

See `migrations/000001_create_tables.up.sql` for full schema details.

## A2A Protocol Integration

The application implements the A2A (Agent-to-Agent) protocol using JSON-RPC 2.0, allowing integration with any compatible messaging platform.

### Integration Steps

1. **Deploy** your application to a publicly accessible server
2. **Configure** your messaging platform to send JSON-RPC 2.0 requests to `/a2a/agent/eunoia`
3. **Test** the integration with sample messages

The A2A endpoint handles:
- JSON-RPC 2.0 request parsing
- User identification and management from `telex_user_id`
- Message content extraction from nested parts structure
- Context building from user history
- AI-powered response generation
- Conversation history storage
- JSON-RPC 2.0 formatted responses

### Message Format

Messages are sent in a nested structure with parts that can contain text, data, or other content types. The system automatically extracts text content from all parts and combines them into a single message for processing.

## Development

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



**Note**: This is a personal wellbeing tool and should not replace professional mental health care. If you're experiencing a mental health crisis, please contact a mental health professional or crisis hotline immediately.
