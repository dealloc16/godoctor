package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <tool_name> [arguments]", os.Args[0])
	}

	toolName := os.Args[1]
	args := os.Args[2:]

	client := mcp.NewClient(&mcp.Implementation{Name: "godoctor-cli", Version: "v0.0.1"}, nil)
	transport := &mcp.StreamableClientTransport{Endpoint: "http://localhost:8080"}
	session, err := client.Connect(context.Background(), transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer session.Close()

	var params *mcp.CallToolParams
	switch toolName {
	case "hello_world":
		params = &mcp.CallToolParams{
			Name: "hello_world",
		}
	case "godoc":
		if len(args) < 1 {
			log.Fatalf("Usage: %s godoc <package> [symbol]", os.Args[0])
		}
		arguments := map[string]any{"package": args[0]}
		if len(args) > 1 {
			arguments["symbol"] = args[1]
		}
		params = &mcp.CallToolParams{
			Name:      "godoc",
			Arguments: arguments,
		}
	default:
		log.Fatalf("Unknown tool: %s", toolName)
	}

	result, err := session.CallTool(context.Background(), params)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			var out map[string]any
			if err := json.Unmarshal([]byte(textContent.Text), &out); err != nil {
				fmt.Println(textContent.Text)
			} else {
				prettyJSON, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					fmt.Println(textContent.Text)
				} else {
					fmt.Println(string(prettyJSON))
				}
			}
		}
	}
}
