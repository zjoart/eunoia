# Eunoia - AI Mental Wellbeing Assistant

A Go-based AI agent for mental wellbeing that performs emotional check-ins, analyzes reflections, and provides supportive, context-aware responses via A2A protocol integration.

## Features

- **Intelligent Mood Detection**: Automatically detects and tracks emotional expressions in conversations
- **Automatic Check-ins**: Creates emotional check-ins from mood expressions (e.g., "feeling great", "I'm stressed")
- **Smart Reflection Analysis**: Detects reflective messages and performs AI-powered sentiment analysis
- **Context-Aware Conversations**: Maintains conversation history with personalized, empathetic responses
- **Platform-Agnostic Architecture**: Extensible platform interface supporting multiple messaging platforms
- **A2A Protocol Compliant**: Full JSON-RPC 2.0 compliance with agent discovery endpoint
- **Gemini AI Integration**: Powered by Google's Gemini 2.5 Flash for natural, empathetic interactions

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

## How It Works

### Intelligent Intent Detection

Eunoia automatically detects user intent and creates appropriate records:

**Mood Detection**: Recognizes emotional expressions
- "I'm feeling great today" → Creates check-in with mood score 8/10 (happy)
- "Feeling stressed and anxious" → Creates check-in with mood score 3/10 (anxious)

**Reflection Detection**: Identifies deeper thoughts (15+ words with reflection indicators)
- "Today I realized..." → Creates reflection with sentiment analysis
- "I've been thinking about..." → Stores reflection with AI-generated insights

### Platform Integration

Eunoia uses a platform-agnostic architecture with flexible metadata handling:

**Primary Endpoint:** `POST /a2a/agent/eunoia`  
**Protocol:** JSON-RPC 2.0  
**Agent Discovery:** `GET /.well-known/agent.json`

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
      "parts": [{"kind": "text", "text": "I'm feeling anxious today"}],
      "metadata": {
        "platform_user_id": "user-123",
        "platform_channel_id": "channel-456"
      },
      "messageId": "msg-789"
    }
  }
}
```

**Supported Metadata Keys:**
- User ID: `platform_user_id`, `telex_user_id`, or `user_id`
- Channel ID: `platform_channel_id`, `telex_channel_id`, or `channel_id`

### A2A Protocol Compliance

- Full JSON-RPC 2.0 specification adherence
- Proper error codes and structured responses
- Agent discovery via `.well-known/agent.json`
- Support for conversation history and context

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/a2a/agent/eunoia` | POST | A2A protocol message endpoint (JSON-RPC 2.0) |
| `/agent/health` | GET | Health check endpoint |
| `/.well-known/agent.json` | GET | A2A agent discovery endpoint |

## Architecture

- **Language:** Go 1.24.2
- **Database:** MySQL 8.0+ with golang-migrate
- **AI Service:** Google Gemini 2.5 Flash
- **Protocol:** A2A (JSON-RPC 2.0)
- **Routing:** Gorilla Mux
- **Logging:** Structured logging with Zap
- **Design:** Clean architecture with service-repository pattern

### Key Components

- **Conversation Service**: Handles message processing and intent detection
- **Check-in Service**: Manages mood tracking and emotional check-ins
- **Reflection Service**: Processes reflections with AI analysis
- **Platform Interface**: Abstraction layer for multi-platform support
- **Gemini Service**: AI integration for empathetic responses

---

**Important:** This is a wellbeing support tool and should not replace professional mental health care. In crisis situations, please contact a mental health professional or emergency services.
