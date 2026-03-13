package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Test",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithLogging(),
	)

	tool := mcp.NewTool(
		"test-tool",
		mcp.WithDescription("Test Tool"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	s.AddTool(tool, testHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func testHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Test: %s", name)), nil
}
