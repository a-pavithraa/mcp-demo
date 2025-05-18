package main

import (
	"fmt"

	"github.com/apavithraa/mcp-demo/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	s := server.NewMCPServer(
		"AWS MCP server 🚀",
		"0.0.5",
		server.WithLogging(),
	)

	s.AddTool(tools.ListDynamoDbTables())
	s.AddTool(tools.GetDynamoDbTableMetadata())
	s.AddTool(tools.ListKmsKeysWithMetadata())
	s.AddTool(tools.ListS3BucketsWithMetadata())
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
