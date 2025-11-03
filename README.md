# Eunoia - AI Mental Wellbeing Assistant

An AI agent for mental wellbeing that performs daily emotional check-ins, analyzes reflections, and provides supportive, context-aware responses. Built in Go and integrated with Telex.im using the A2A protocol.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Database Setup](#database-setup)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Telex.im Integration](#telexim-integration)
- [Development](#development)
- [Project Structure](#project-structure)
- [Deployment](#deployment)

## Features

- **Daily Emotional Check-ins**: Track mood scores (1-10 scale) and emotional states over time
- **Reflection Analysis**: AI-powered sentiment analysis and theme extraction from journal entries
- **Context-Aware Conversations**: Maintains conversation history with personalized, empathetic responses
- **Mood Trend Analysis**: Identifies patterns in emotional wellbeing over time (improving, declining, stable)
- **Supportive AI Responses**: Powered by Google's Gemini AI for empathetic, human-like interactions
- **Telex.im Integration**: Full A2A protocol implementation for seamless messaging platform integration
- **Database Migrations**: Automated migration system using golang-migrate
- **Structured Logging**: Clean, professional logs using uber/zap without emojis

## Architecture

### Feature-Based Structure

The project follows a true feature-based architecture where each feature is self-contained:

```
internal/
├── agent/                         # Shared AI service
│   └── gemini_service.go          # Gemini API integration
├── user/                          # User feature (self-contained)
│   ├── model.go                   # User models
│   └── repository.go              # User database operations
├── checkin/                       # Emotional check-in feature (self-contained)
│   ├── model.go                   # Check-in models
│   ├── repository.go              # Check-in database operations
│   └── service.go                 # Check-in business logic
├── reflection/                    # Reflection analysis feature (self-contained)
│   ├── model.go                   # Reflection models
│   ├── repository.go              # Reflection database operations
│   └── service.go                 # Reflection business logic
├── conversation/                  # Conversation feature (self-contained)
│   ├── model.go                   # Conversation models
│   ├── repository.go              # Conversation database operations
│   ├── service.go                 # Conversation business logic
│   └── handler.go                 # HTTP handlers for A2A protocol
├── config/                        # Configuration management
│   └── config.go                  # Centralized config loading
├── database/                      # Database connection
│   └── db.go                      # Database initialization
└── middleware/                    # HTTP middleware
    └── cors.go                    # CORS middleware

cmd/
├── app/                           # Main application
│   └── main.go                    # Application entry point
├── migrate/                       # Database migrations CLI
│   └── main.go                    # Migration runner
└── routes/                        # Route definitions
    └── routes.go                  # HTTP route setup

migrations/                        # Database migrations
├── 000001_create_tables.up.sql   # Create tables migration
└── 000001_create_tables.down.sql # Drop tables migration

pkg/
└── logger/                        # Shared logging utilities
    └── logger.go                  # Zap logger configuration
```

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24.2 or higher** - [Download Go](https://go.dev/dl/)
- **MySQL 8.0 or higher** - [Download MySQL](https://dev.mysql.com/downloads/)
- **Git** - For cloning the repository
- **Make** - For running Makefile commands (optional but recommended)

### Required Accounts

- **Google Cloud Account** - For Gemini API access
  - Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
- **Telex.im Account** - For agent integration
  - Request organization access with `/telex-invite your-email@example.com`

## Quick Start

Follow these steps to get Eunoia running on your local machine:

### 1. Clone the Repository

```bash
git clone https://github.com/zjoart/eunoia.git
cd eunoia
```

### 2. Install Dependencies

```bash
go mod tidy
```

This will download all required Go packages including:
- `gorilla/mux` - HTTP routing
- `golang-migrate/migrate` - Database migrations
- `go-sql-driver/mysql` - MySQL driver
- `uber-go/zap` - Structured logging
- `joho/godotenv` - Environment variable loading

### 3. Configure Environment Variables

Copy the example environment file and update it with your credentials:

```bash
cp .env.example .env
```

Edit `.env` with your actual configuration:

```env
# Server Configuration
PORT=8080
APP_ENV=development

# Database Configuration
DB_USER=root
DB_PASS=your_mysql_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=eunoia_db

# API Configuration
API_BASE=localhost:8080
SWAGGER_SCHEMES=http

# AI Configuration
GEMINI_API_KEY=your_actual_gemini_api_key_here
```

**Important**: Never commit your `.env` file to version control.

## Database Setup

### Option 1: Using Make Commands (Recommended)

```bash
# Create the database
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS eunoia_db;"

# Run migrations
make migrate-up
```

### Option 2: Manual Setup

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS eunoia_db;"

# Apply migrations manually
mysql -u root -p eunoia_db < migrations/000001_create_tables.up.sql
```

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

### Verify Database Setup

```bash
mysql -u root -p eunoia_db -e "SHOW TABLES;"
```

You should see:
- `users`
- `emotional_checkins`
- `reflections`
- `conversation_history`
- `user_preferences`
- `schema_migrations` (created by golang-migrate)

## Configuration

### Environment Variables Explained

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `APP_ENV` | Environment (development/production) | `development` |
| `DB_USER` | MySQL username | `root` |
| `DB_PASS` | MySQL password | `your_password` |
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_NAME` | Database name | `eunoia_db` |
| `API_BASE` | API base URL for Swagger | `localhost:8080` |
| `SWAGGER_SCHEMES` | API schemes | `http` or `https` |
| `GEMINI_API_KEY` | Google Gemini API key | `AIza...` |

### Getting Gemini API Key

1. Go to [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with your Google account
3. Click "Create API Key"
4. Copy the key and add it to your `.env` file

**Note**: Keep your API key secure and never share it publicly.

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
# Test health endpoint
curl http://localhost:8080/health

# Expected output: "Service is up and running"

# Test agent health endpoint
curl http://localhost:8080/agent/health

# Expected JSON output with status information
```

## API Endpoints

### 1. Health Check

**Endpoint:** `GET /health`

**Description:** Returns the overall service health status.

**Example:**
```bash
curl http://localhost:8080/health
```

**Response:**
```
Service is up and running
```

### 2. Agent Health Check

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

### 3. A2A Agent Endpoint (Telex.im Integration)

**Endpoint:** `POST /a2a/agent/eunoia`

**Description:** Main endpoint for receiving messages from Telex.im and responding with AI-generated, context-aware responses.

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "message": "How are you feeling today?",
  "userId": "telex-user-id-12345",
  "channelId": "channel-uuid",
  "messageId": "message-uuid",
  "timestamp": "2025-11-03T10:30:00Z",
  "context": {
    "optional": "metadata"
  }
}
```

**Request Fields:**
- `message` (required, string): The user's message text
- `userId` (required, string): Unique Telex user identifier
- `channelId` (optional, string): Telex channel identifier
- `messageId` (optional, string): Unique message identifier
- `timestamp` (optional, string): ISO 8601 timestamp
- `context` (optional, object): Additional context metadata

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

The application uses MySQL with the following tables:

### 1. Users Table

Stores Telex user information and maps Telex user IDs to internal user IDs.

| Column | Type | Description |
|--------|------|-------------|
| `id` | VARCHAR(36) | Primary key, internal user ID |
| `telex_user_id` | VARCHAR(255) | Unique Telex user identifier |
| `username` | VARCHAR(255) | Optional username |
| `created_at` | TIMESTAMP | Account creation timestamp |
| `updated_at` | TIMESTAMP | Last update timestamp |

### 2. Emotional Check-ins Table

Records daily mood scores and emotional states.

| Column | Type | Description |
|--------|------|-------------|
| `id` | VARCHAR(36) | Primary key |
| `user_id` | VARCHAR(36) | Foreign key to users table |
| `mood_score` | INT | Mood rating (1-10 scale) |
| `mood_label` | VARCHAR(50) | Mood description (e.g., "happy", "anxious") |
| `description` | TEXT | Optional detailed description |
| `check_in_date` | DATE | Date of check-in |
| `created_at` | TIMESTAMP | Record creation timestamp |

### 3. Reflections Table

Stores journal entries with AI-generated sentiment analysis and themes.

| Column | Type | Description |
|--------|------|-------------|
| `id` | VARCHAR(36) | Primary key |
| `user_id` | VARCHAR(36) | Foreign key to users table |
| `content` | TEXT | Journal entry content |
| `sentiment` | VARCHAR(20) | AI-detected sentiment (positive/negative/neutral/mixed) |
| `key_themes` | TEXT | Comma-separated key themes |
| `ai_analysis` | TEXT | AI-generated supportive analysis |
| `created_at` | TIMESTAMP | Entry creation timestamp |
| `updated_at` | TIMESTAMP | Last update timestamp |

### 4. Conversation History Table

Maintains chat history for context-aware responses.

| Column | Type | Description |
|--------|------|-------------|
| `id` | VARCHAR(36) | Primary key |
| `user_id` | VARCHAR(36) | Foreign key to users table |
| `message_role` | VARCHAR(20) | Role: "user" or "assistant" |
| `message_content` | TEXT | Message text |
| `context_data` | TEXT | Optional context information |
| `created_at` | TIMESTAMP | Message timestamp |

### 5. User Preferences Table

Stores user-specific settings and preferences (for future features).

| Column | Type | Description |
|--------|------|-------------|
| `user_id` | VARCHAR(36) | Primary key, foreign key to users |
| `reminder_time` | TIME | Preferred check-in reminder time |
| `reminder_frequency` | VARCHAR(20) | Frequency (daily/weekly) |
| `preferred_tone` | VARCHAR(20) | Preferred AI tone |
| `created_at` | TIMESTAMP | Record creation timestamp |
| `updated_at` | TIMESTAMP | Last update timestamp |

## Telex.im Integration

### Step 1: Get Telex Access

Request access to the Telex organization by sending this command in the Telex chat:

```
/telex-invite your-email@example.com
```

Wait for the confirmation email and complete the organization setup.

### Step 2: Deploy Your Application

Deploy your application to a publicly accessible server. Recommended platforms:

- **Railway** - Easy deployment with MySQL addon
- **Render** - Free tier available
- **Fly.io** - Global deployment
- **DigitalOcean** - VPS option
- **AWS/GCP/Azure** - Enterprise options

### Step 3: Configure Workflow JSON

Update `telex-workflow.json` with your deployed URL:

```json
{
  "active": true,
  "category": "wellness",
  "description": "AI mental wellbeing assistant for emotional check-ins and support",
  "id": "eunoia_agent_v1",
  "long_description": "System prompt and agent behavior description...",
  "name": "eunoia_wellbeing_agent",
  "nodes": [
    {
      "id": "eunoia_agent_node",
      "name": "Eunoia Wellbeing Agent",
      "parameters": {},
      "position": [500, 200],
      "type": "a2a/mastra-a2a-node",
      "typeVersion": 1,
      "url": "https://your-deployed-app.com/a2a/agent/eunoia"
    }
  ],
  "pinData": {},
  "settings": {
    "executionOrder": "v1"
  },
  "short_description": "Compassionate AI for mental wellbeing support and emotional check-ins"
}
```

Replace `your-deployed-app.com` with your actual deployed URL.

### Step 4: Register Agent with Telex

1. Log into Telex.im platform
2. Navigate to Workflow Management
3. Import your `telex-workflow.json` file
4. Activate the workflow

### Step 5: Test Integration

1. Create or join a channel in Telex.im
2. Add the Eunoia agent to the channel
3. Send test messages:
   - "Hello!"
   - "How are you feeling today?"
   - "I want to do a check-in"

### Viewing Agent Logs

Monitor your agent's interactions in real-time:

```
https://api.telex.im/agent-logs/{channel-id}.txt
```

**Finding the Channel ID:**
1. Open a channel in Telex.im
2. Look at the URL: `https://telex.im/.../colleagues/{CHANNEL-ID}/{MESSAGE-ID}`
3. Copy the first UUID (CHANNEL-ID)

**Example:**
```
https://api.telex.im/agent-logs/01989dec-0d08-71ee-9017-00e4556e1942.txt
```

## How It Works

### Conversation Flow

```
User Message (Telex.im)
        ↓
A2A Endpoint (/a2a/agent/eunoia)
        ↓
User Identification (Get or Create User)
        ↓
Context Building (Recent Check-ins, Reflections, Mood Trends)
        ↓
AI Processing (Gemini API with Context)
        ↓
Response Generation (Empathetic, Context-Aware)
        ↓
History Storage (Save Conversation)
        ↓
Response Delivery (Back to Telex.im)
```

### Detailed Flow

1. **Message Receipt**: User sends a message via Telex.im to the agent
2. **Request Validation**: Handler validates required fields (message, userId)
3. **User Management**: System retrieves existing user or creates new user record
4. **Context Building**: Service gathers user's context:
   - Last 5 emotional check-ins
   - Last 3 reflections with sentiment
   - 7-day mood statistics and trends
5. **Conversation History**: Retrieves last 10 messages for context continuity
6. **AI Processing**: 
   - Constructs system prompt with user context
   - Sends message + history to Gemini AI
   - Generates empathetic, context-aware response
7. **Storage**: Saves both user message and AI response to database
8. **Response Delivery**: Returns formatted JSON response to Telex.im
9. **Display**: Telex.im shows response to user in chat

### AI Capabilities

#### 1. Sentiment Analysis
- Analyzes emotional tone in reflections
- Classifies as: positive, negative, neutral, or mixed
- Used to understand user's emotional state over time

#### 2. Theme Extraction
- Identifies 3-5 key themes from journal entries
- Helps recognize recurring topics or concerns
- Provides insights into user's focus areas

#### 3. Mood Tracking
- Records daily mood scores (1-10)
- Calculates averages over time periods
- Identifies trends: improving, declining, or stable

#### 4. Context-Aware Responses
- References past conversations naturally
- Acknowledges mood patterns when relevant
- Personalizes support based on user history
- Maintains conversational continuity

### Example Interactions

**First-time user:**
```
User: "Hello"
Agent: "Hello! I'm Eunoia, here to support your mental wellbeing. 
       How are you feeling today?"
```

**User with history:**
```
User: "I'm feeling better today"
Agent: "That's wonderful to hear! I notice your mood has been 
       trending upward over the past week. What's contributing 
       to this positive change?"
```

**Check-in request:**
```
User: "I want to do a check-in"
Agent: "I'd be happy to help you check in. On a scale of 1-10, 
       how would you rate your mood right now?"
```

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
# Running
make run                    # Run the application
make clean                  # Clean build artifacts

# Database
make migrate-up             # Apply all migrations
make migrate-down           # Rollback all migrations
make migrate-version        # Show current migration version
make migrate-steps STEPS=1  # Apply/rollback n migrations
make migrate-force VERSION=1 # Force specific version

# Testing
make test                   # Run all tests
make test-log               # Run tests with verbose logging
make test-force             # Run tests without cache
make test-race              # Run tests with race detection
make test-ci                # Run tests with coverage

# Utilities
make tidy                   # Tidy go.mod and go.sum
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

### Building

```bash
# Development build
go build -o bin/eunoia cmd/app/main.go

# Production build (optimized)
go build -ldflags="-s -w" -o bin/eunoia cmd/app/main.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o bin/eunoia-linux cmd/app/main.go
GOOS=darwin GOARCH=arm64 go build -o bin/eunoia-mac cmd/app/main.go
```

### Code Style and Conventions

- **No global variables** - Use dependency injection
- **Feature-based structure** - Each feature is self-contained
- **Centralized config** - All `os.Getenv()` calls in `internal/config/config.go`
- **Clean logging** - No emojis, professional structured logs
- **Error handling** - Always log errors with context
- **Database** - Use prepared statements to prevent SQL injection

### Logging

The application uses structured logging with `uber/zap`:

```go
// Info logging
logger.Info("processing message", logger.Fields{
    "user_id": userID,
    "message_length": len(message),
})

// Error logging
logger.Error("failed to process", logger.WithError(err))

// Multiple fields
logger.Info("request completed", logger.Merge(
    logger.Fields{"status": 200},
    logger.Fields{"duration": duration},
))
```

**Log Output Example:**
```json
{
  "level": "info",
  "timestamp": "2025-11-03T10:30:00Z",
  "caller": "conversation/handler.go:42",
  "message": "received A2A message request",
  "method": "POST",
  "path": "/a2a/agent/eunoia"
}
```

### Configuration Management

All configuration is centralized in `internal/config/config.go`:

```go
// ✅ Correct - Load in config package
func LoadConfig() *Config {
    return &Config{
        Port: getEnv("PORT"),
        // ...
    }
}

// ❌ Never do this elsewhere in the codebase
apiKey := os.Getenv("GEMINI_API_KEY")
```

## Deployment

### Deployment Checklist

Before deploying to production:

- [ ] Set `APP_ENV=production` in environment variables
- [ ] Use strong database password
- [ ] Configure HTTPS/TLS
- [ ] Set up database backups
- [ ] Configure monitoring and logging
- [ ] Test all endpoints
- [ ] Run migrations on production database
- [ ] Update `telex-workflow.json` with production URL
- [ ] Set appropriate CORS origins

### Platform-Specific Guides

#### Railway

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Initialize project
railway init

# Add environment variables
railway variables set PORT=8080
railway variables set APP_ENV=production
railway variables set GEMINI_API_KEY=your_key
# ... add all other variables

# Deploy
railway up
```

#### Render

1. Connect GitHub repository
2. Create new Web Service
3. Build Command: `go build -o eunoia cmd/app/main.go`
4. Start Command: `./eunoia`
5. Add environment variables in dashboard
6. Deploy

#### Fly.io

```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Launch app
fly launch

# Set secrets
fly secrets set GEMINI_API_KEY=your_key
fly secrets set DB_PASS=your_password

# Deploy
fly deploy
```

### Environment Variables for Production

```env
PORT=8080
APP_ENV=production
DB_USER=prod_user
DB_PASS=secure_random_password
DB_HOST=your-db-host.com
DB_PORT=3306
DB_NAME=eunoia_production
API_BASE=your-deployed-domain.com
SWAGGER_SCHEMES=https
GEMINI_API_KEY=your_gemini_api_key
```

### Database Migration in Production

```bash
# Using the built-in migration tool
go run cmd/migrate/main.go up

# Or with Make
make migrate-up
```

### Monitoring

**Application Logs:**
- Most platforms provide built-in log viewing
- Railway: `railway logs`
- Fly.io: `fly logs`
- Render: View in dashboard

**Health Checks:**
Set up monitoring services to ping:
- `https://your-app.com/health`
- `https://your-app.com/agent/health`

**Recommended Monitoring Services:**
- UptimeRobot (Free tier available)
- Pingdom
- Better Uptime
- StatusCake

## Security Considerations

### Best Practices

1. **API Keys and Secrets**
   - Never commit `.env` file to Git
   - Use platform secret management
   - Rotate keys regularly
   - Use different keys for dev/prod

2. **Database Security**
   - Use strong passwords
   - Enable SSL/TLS connections
   - Restrict database access by IP
   - Regular backups
   - Use prepared statements (already implemented)

3. **Application Security**
   - CORS configured properly
   - Input validation on all endpoints
   - Rate limiting (recommended to add)
   - HTTPS in production
   - Keep dependencies updated

4. **Logging Security**
   - Never log sensitive data (API keys, passwords)
   - Log authentication attempts
   - Monitor error patterns
   - Set up alerts for anomalies

### Security Headers (Recommended to Add)

```go
// Add to middleware
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
```

## Troubleshooting

### Common Issues

**Problem: `panic: DB_USER is required`**
```bash
# Solution: Ensure .env file exists and contains all required variables
cp .env.example .env
# Edit .env with your values
```

**Problem: `Error 1146: Table doesn't exist`**
```bash
# Solution: Run migrations
make migrate-up
```

**Problem: `Failed to connect to database`**
```bash
# Solution: Check MySQL is running
mysql -u root -p -e "SELECT 1;"

# Verify credentials in .env match MySQL
```

**Problem: `Gemini API error`**
```bash
# Solution: Verify API key is correct
# Check quota at https://makersuite.google.com/
# Ensure API key has no extra spaces
```

**Problem: Agent not responding in Telex**
```bash
# Solution: 
# 1. Verify deployment is running: curl https://your-app.com/health
# 2. Check Telex agent logs
# 3. Verify workflow JSON has correct URL
# 4. Check application logs for errors
```

### Debug Mode

Enable verbose logging in development:

```bash
# Add to .env
APP_ENV=development

# Run with verbose output
go run cmd/app/main.go
```

### Getting Help

If you encounter issues:

1. Check application logs
2. Verify environment variables
3. Test endpoints manually with cURL
4. Review database connection
5. Check Telex agent logs
6. Open GitHub issue with:
   - Error message
   - Steps to reproduce
   - Environment details
   - Relevant log snippets

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow the existing code style
4. Write tests for new features
5. Ensure all tests pass (`make test`)
6. Commit with clear messages
7. Push to your branch
8. Open a Pull Request

## Future Enhancements

Planned features and improvements:

- [ ] Push notifications for check-in reminders
- [ ] Data visualization dashboard for mood trends
- [ ] Weekly/monthly mood reports via email
- [ ] Multi-language support (i18n)
- [ ] Integration with additional AI models (Claude, GPT-4)
- [ ] Export functionality for user data (JSON, CSV)
- [ ] Advanced analytics and insights
- [ ] Group therapy session support
- [ ] Crisis detection and emergency resource suggestions
- [ ] Mobile app integration
- [ ] Voice note support for reflections
- [ ] Meditation and breathing exercise suggestions

## License

MIT License

Copyright (c) 2025 Eunoia Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Acknowledgments

- **Google Gemini AI** - For providing the AI capabilities
- **Telex.im** - For the messaging platform integration
- **golang-migrate** - For database migration management
- **uber-go/zap** - For structured logging
- **gorilla/mux** - For HTTP routing

## Support

- **Documentation**: This README and inline code comments
- **Issues**: [GitHub Issues](https://github.com/zjoart/eunoia/issues)
- **Email**: support@eunoia.example.com

---

**Built with care for mental wellbeing support** 💙

**Note**: This is a personal wellbeing tool and should not replace professional mental health care. If you're experiencing a mental health crisis, please contact a mental health professional or crisis hotline immediately.
