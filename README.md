# Go AI Agent
This repository demonstrates how to build AI agents using Go, countering the common assumption that Python is the only viable language for AI development. This repo showcases how to integrate OpenAI's capabilities with Slack using popular Go libraries.

### Overview
This project creates a simple AI agent in Go that:

Connects to OpenAI's API for language model interaction
Integrates with Slack for messaging capabilities
Provides a convenient CLI interface

### Technologies Used

- Go (1.23+)
- OpenAI Go SDK
- Slack Go SDK
- urfave/cli for command-line functionality

### Installation
```bash
# Clone the repository
git clone https://github.com/suyashmohan/myaiagent
cd myaiagent
go mod download
```

### Configuration
Create a .env file in the project root with your API credentials:
```
OPENAI_API_KEY=
API_BASE_URL=
OPENAI_MODEL_NAME=
SLACK_BOT_TOKEN=
SLACK_APP_TOKEN=
```

### Usage
Basic Example

```bash
# Run the agent with slack integration
go run ./cli/main.go slack

# Test the agent on cli without slack
go run ./cli/main.go "What is 2+2?"
```

### Why Go for AI Agents?
While Python dominates AI development tutorials, Go offers several advantages:

- Strong typing and compilation catch errors earlier
- Excellent concurrency support via goroutines
- Lower memory footprint and faster execution
- Production-ready standard library
- Simple deployment with single binaries

### Limitations
This code is intended as a demonstration and learning resource, not for production use. It lacks:

- Comprehensive error handling
- Extensive testing
- Production security measures
- Advanced state management
