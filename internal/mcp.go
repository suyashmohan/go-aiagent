package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPConfig struct {
	McpServers map[string]MCPServerConfig `json:"mcpServers"`
}

type MCPServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

func LoadMCPServers(mcpConfigFilepath string) (map[string]client.MCPClient, error) {
	mcpClients := make(map[string]client.MCPClient)
	mcpServerConfig := MCPConfig{}

	file, err := os.ReadFile(mcpConfigFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mcp server config file - %w", err)
	}
	err = json.Unmarshal(file, &mcpServerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json mcp server config file - %w", err)
	}

	for mcpServerName, mcpServer := range mcpServerConfig.McpServers {
		// Convert ENV into "key=value" form
		env := []string{}
		for envKey, envValue := range mcpServer.Env {
			env = append(env, fmt.Sprintf("%s=%s", envKey, envValue))
		}
		log.Println("Loading MCP Client", mcpServerName, env, mcpServer.Command, mcpServer.Args)

		stdioClient, err := client.NewStdioMCPClient(mcpServer.Command, env, mcpServer.Args...)
		if err != nil {
			return nil, fmt.Errorf("failed to start transport for mcp server - %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "MCP-Go Client",
			Version: "0.1.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		_, err = stdioClient.Initialize(ctx, initRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to initialise mcp client - %w", err)
		}

		mcpClients[mcpServerName] = stdioClient
	}

	return mcpClients, nil
}
