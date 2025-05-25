package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/client"
	"github.com/suyashmohan/go-aiagent/cmd"
	"github.com/suyashmohan/go-aiagent/internal"
	"github.com/urfave/cli/v3"
)

func main() {
	// Load ENV from .env file
	godotenv.Load()

	mcpClients, err := internal.LoadMCPServers("mcp.json")
	if err != nil {
		log.Println("failed to load mcp.json - %w", err)
	}

	// Load System prompt
	fileBytes, err := os.ReadFile("system.md")
	if err != nil {
		log.Fatalln(err)
	}
	systemPrompt := string(fileBytes)

	// Create an Agent with Tools
	agent, err := internal.NewAgentWithTools(systemPrompt, getTools(mcpClients))
	if err != nil {
		log.Fatalln("failed to create ai agent", err)
	}

	// Setup CLI commands
	cmd := &cli.Command{
		Name:   "agent",
		Usage:  "Run AI agent with a query",
		Action: cmd.RootCMD(agent),
		// Separate command to start as Slack bot
		Commands: []*cli.Command{
			{
				Name:   "slack",
				Usage:  "Run AI Agent as slack bot",
				Action: cmd.SlackCMD(agent),
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err)
	}

	for mcpName, mcpClient := range mcpClients {
		log.Println("Closing", mcpName)
		mcpClient.Close()
		log.Println("Closed", mcpName)
	}
}

func getTools(mcpClients map[string]client.MCPClient) map[string]internal.AgentTool {
	tools := map[string]internal.AgentTool{}

	for mcpName, mcpClient := range mcpClients {
		toolsResult, err := internal.ListMCPTools(mcpClient)
		if err != nil {
			log.Printf("Failed to list tools: %v", err)
			return nil
		}

		log.Printf("%s has %d tools available\n", mcpName, len(toolsResult.Tools))
		for i, tool := range toolsResult.Tools {
			log.Printf("%d. %s - %s\n", i+1, tool.Name, tool.Description)
			toolName := fmt.Sprintf("%s__%s", mcpName, tool.Name)
			tools[toolName] = internal.MCPToolAsLLMTool(tool, mcpClient)
		}

	}

	return tools
}
