# Eunoia - AI Mental Wellbeing Assistant

A Go-based AI agent for mental wellbeing that performs emotional check-ins, analyzes reflections, and provides supportive, context-aware responses via A2A protocol integration.

## Features

- **Emotional Check-ins**: Track mood scores (1-10 scale) and emotional states
- **Reflection Analysis**: AI-powered sentiment analysis and theme extraction
- **Context-Aware Conversations**: Maintains conversation history with personalized responses
- **Multi-Platform Support**: Extensible platform architecture for A2A protocol integration
- **Gemini AI Integration**: Powered by Google's Gemini AI for empathetic interactions

## Prerequisites

- Go 1.24.2+
- MySQL 8.0+
- Google Gemini API key

## Quick Start

```bash
git clone https://github.com/zjoart/eunoia.git
cd eunoia

# Install dependencies
make tidy

# Setup environment
cp .env.example .env

# Setup database
make migrate-up

# Run application
make run
```

## Configuration

Configure your `.env` file with required values. See [.env.example](.env.example) for all required variables including database credentials and Gemini API key.

## Commands

Run `make help` to see all available commands with descriptions.

### Key Commands

```bash
make run              # Start the application
make migrate-up       # Apply database migrations
make migrate-down     # Rollback migrations
make migrate-version  # Check migration status
```

## Platform Integration

### Telex

Eunoia integrates with Telex.im via A2A protocol. See [Telex Documentation](https://docs.telex.im/docs) for integration details.

**Endpoint:** `POST /a2a/agent/eunoia`

**Format:** JSON-RPC 2.0

**Example Request:**
```json
{
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
    }
  }
}
```

### Adding New Platforms

The platform architecture supports easy addition of new messaging platforms:

1. Implement the `Platform` interface
2. Register with the platform registry
3. Handle platform-specific metadata and response formatting

## API Endpoints

- `GET /agent/health` - Health check
- `POST /a2a/agent/eunoia` - A2A protocol endpoint

## Architecture

- **Language:** Go 1.24.2
- **Database:** MySQL with golang-migrate
- **AI Service:** Google Gemini API
- **Protocol:** A2A (JSON-RPC 2.0)
- **Routing:** Gorilla Mux
- **Logging:** Zap (structured)

---

**Note:** This is a personal wellbeing tool and should not replace professional mental health care.
