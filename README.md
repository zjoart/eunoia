<div align="center">

# 🌱 Eunoia - AI Mental Wellbeing Assistant

**A Go-based AI agent for mental wellbeing that performs emotional check-ins, analyzes reflections, and provides supportive, context-aware responses via A2A protocol integration.**



📖 **[Read the full article: Building Eunoia - A Mental Wellbeing Companion](https://dev.to/oluwadahunsi_ifeoluwa_79e/building-eunoia-a-mental-wellbeing-companion-gei)**

</div>

---

## ✨ Features

- 🎯 **Intelligent Mood Detection**: Automatically detects and tracks emotional expressions in conversations
- 📊 **Automatic Check-ins**: Creates emotional check-ins from mood expressions (e.g., "feeling great", "I'm stressed")
- 🔍 **Smart Reflection Analysis**: Detects reflective messages and performs AI-powered sentiment analysis
- 💬 **Context-Aware Conversations**: Maintains conversation history with personalized, empathetic responses
- 🔌 **Platform-Agnostic Architecture**: Extensible platform interface supporting multiple messaging platforms
- ✅ **A2A Protocol Compliant**: Full JSON-RPC 2.0 compliance with agent discovery endpoint
- 🤖 **Gemini AI Integration**: Powered by Google's Gemini 2.5 Flash for natural, empathetic interactions

## 📋 Prerequisites

- Go 1.24.2+
- MySQL 8.0+
- Google Gemini API key

## 🚀 Quick Start

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

## ⚙️ Configuration

Configure your `.env` file with required values. See [.env.example](.env.example) for all required variables including database credentials and Gemini API key.

## 📦 Commands

Run `make help` to see all available commands with descriptions.

### Key Commands

```bash
make run              # Start the application
make migrate-up       # Apply database migrations
make migrate-down     # Rollback migrations
make migrate-version  # Check migration status
make test             # Run all tests
make test-ci          # Run tests with race detection and coverage
```

## 🧪 Testing

Eunoia includes comprehensive test coverage across all layers:

### Running Tests Locally

```bash
# Run all tests
make test

# Run tests with race condition detection and coverage report
make test-ci

# Run tests without cache
make test-force

# Run tests with verbose output
make test-log

# Run specific test function
make test-function TEST=TestHandleA2AMessage_ValidRequest
```

### 🔄 Continuous Integration

Tests run automatically via GitHub Actions on every push and pull request to the `main` branch. The CI pipeline:
- ✅ Runs the full test suite with race condition detection
- 📈 Generates coverage reports
- 📦 Uploads coverage artifacts for review

See the [GitHub Actions workflow configuration](.github/workflows/go-test.yml) for details.




## 🔧 How It Works

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

### Example Request

```json
{
  "jsonrpc": "2.0",
  "id": "req-12345",
  "method": "message/send",
  "params": {
    "message": {
      "kind": "message",
      "role": "user",
      "parts": [
        {
          "kind": "text",
          "text": "I'm feeling anxious today"
        }
      ],
      "metadata": {
        "telex_user_id": "user-123",
        "telex_channel_id": "channel-456"
      },
      "messageId": "msg-789"
    },
    "configuration": {
      "acceptedOutputModes": ["text/plain"],
      "historyLength": 0,
      "blocking": false
    }
  }
}
```

### Example Response

```json
{
  "jsonrpc": "2.0",
  "id": "req-12345",
  "result": {
    "id": "task-abc123",
    "contextId": "ctx-xyz789",
    "status": {
      "state": "completed",
      "timestamp": "2025-11-04T10:30:00.000Z",
      "message": {
        "messageId": "msg-789",
        "role": "agent",
        "parts": [
          {
            "kind": "text",
            "text": "I hear that you're feeling anxious today. That's completely valid, and I'm here to support you. Would you like to talk about what's contributing to your anxiety, or would you prefer some grounding techniques to help you feel more centered right now?"
          }
        ],
        "kind": "message",
        "taskId": "task-abc123",
        "metadata": {
          "agent": "eunoia"
        }
      }
    },
    "artifacts": [
      {
        "artifactId": "artifact-uuid-001",
        "name": "eunoia_response",
        "parts": [
          {
            "kind": "text",
            "text": "I hear that you're feeling anxious today. That's completely valid, and I'm here to support you. Would you like to talk about what's contributing to your anxiety, or would you prefer some grounding techniques to help you feel more centered right now?"
          }
        ]
      }
    ],
    "history": [
      {
        "messageId": "msg-789",
        "role": "agent",
        "parts": [
          {
            "kind": "text",
            "text": "I hear that you're feeling anxious today..."
          }
        ],
        "kind": "message",
        "taskId": "task-abc123",
        "metadata": {
          "agent": "eunoia"
        }
      }
    ],
    "kind": "task"
  }
}
```

**Supported Metadata Keys:**
- User ID: `platform_user_id`, `telex_user_id`, or `user_id`
- Channel ID: `platform_channel_id`, `telex_channel_id`, or `channel_id`

**Response Fields:**
- `status.message`: The agent's current response
- `history`: Full conversation history including the current response
- `artifacts`: Additional resources (empty for text-only conversations)

### A2A Protocol Compliance

- Full JSON-RPC 2.0 specification adherence
- Proper error codes and structured responses
- Agent discovery via `.well-known/agent.json`
- Support for conversation history and context

## 🌐 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/a2a/agent/eunoia` | POST | A2A protocol message endpoint (JSON-RPC 2.0) |
| `/agent/health` | GET | Health check endpoint |
| `/.well-known/agent.json` | GET | A2A agent discovery endpoint |

## 🏗️ Architecture

- **Language:** Go 1.24.2
- **Database:** MySQL 8.0+ with golang-migrate
- **AI Service:** Google Gemini 2.5 Flash
- **Protocol:** A2A (JSON-RPC 2.0)
- **Routing:** Gorilla Mux
- **Logging:** Structured logging with Zap

---

<div align="center">

### ⚠️ Important Notice

This is a wellbeing support tool and should not replace professional mental health care. 
In crisis situations, please contact a mental health professional or emergency services.

</div>
