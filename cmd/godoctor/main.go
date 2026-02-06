// Package main is the entry point for the godoctor MCP server application.
package main

import (
	"context"      // For managing context and cancellation
	"fmt"          // For formatting strings
	"log/slog"     // For structured logging
	"os"           // For OS-level operations like exit codes
	"os/exec"      // For executing external commands (go doc)

	"github.com/modelcontextprotocol/go-sdk/mcp" // MCP server SDK
)

// GodocInput represents the input parameters for the godoc tool.
// It specifies which Go package and optional symbol to retrieve documentation for.
type GodocInput struct {
	Package string `json:"package"`            // The Go package import path (e.g., "fmt", "github.com/...")
	Symbol  string `json:"symbol,omitempty"`   // Optional: specific symbol (function, type, etc.) to document
}

// GodocOutput represents the documentation result returned by the godoc tool.
type GodocOutput struct {
	Documentation string `json:"documentation"` // The formatted Go documentation text
}

// main is the program entry point. It initializes the MCP server and handles any startup errors.
func main() {
	// Run the server and capture any errors during initialization or execution
	if err := run(); err != nil {
		// Log the error with structured logging
		slog.Error("application error", "err", err)
		// Exit with error code 1 to indicate failure
		os.Exit(1)
	}
}

// HelloWorldInput represents the input for the hello_world tool (no parameters required).
type HelloWorldInput struct{}

// HelloWorldOutput represents the response from the hello_world tool.
type HelloWorldOutput struct {
	Message string `json:"message"` // A greeting message
}

// run initializes the MCP server with two tools (hello_world and godoc) and starts listening for requests.
func run() error {
	// Step 1: Create a new MCP server with metadata about this implementation
	// The server identifies itself as "godoctor" version "v0.0.1"
	server := mcp.NewServer(&mcp.Implementation{Name: "godoctor", Version: "v0.0.1"}, nil)
	
	// Step 2: Register the "hello_world" tool - a simple greeting tool with no inputs
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "hello_world",
			Description: "A simple tool that returns a greeting.",
		},
		// Handler function: receives context and request, returns a greeting message
		func(ctx context.Context, req *mcp.CallToolRequest, input HelloWorldInput) (*mcp.CallToolResult, HelloWorldOutput, error) {
			return nil, HelloWorldOutput{Message: "Hello, world!"}, nil
		},
	)
	
	// Step 3: Register the "godoc" tool - retrieves Go package/symbol documentation
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "godoc",
			Description: "Retrieves Go documentation for a given package and optional symbol. This tool is useful for understanding how to use a package or a specific function, variable, or type within that package. The 'package' argument is the full import path of the package (e.g., 'fmt', 'github.com/modelcontextprotocol/go-sdk/mcp'). The optional 'symbol' argument is the name of the function, variable, or type to get documentation for (e.g., 'Printf', 'NewServer').",
		},
		// Handler function: executes "go doc" command to fetch documentation
		func(ctx context.Context, req *mcp.CallToolRequest, input GodocInput) (*mcp.CallToolResult, GodocOutput, error) {
			// Step 3a: Build the go doc command arguments
			args := []string{"doc"}
			
			// Step 3b: Add package and optional symbol to the arguments
			if input.Symbol != "" {
				// If a symbol is specified, include both package and symbol
				args = append(args, input.Package, input.Symbol)
			} else {
				// Otherwise, just retrieve documentation for the entire package
				args = append(args, input.Package)
			}
			
			// Step 3c: Execute the "go doc" command with the constructed arguments
			cmd := exec.Command("go", args...)
			out, err := cmd.CombinedOutput()
			
			// Step 3d: Handle execution errors
			if err != nil {
				return nil, GodocOutput{}, fmt.Errorf("failed to execute go doc: %w, output: %s", err, string(out))
			}
			
			// Step 3e: Return the documentation output as a string
			return nil, GodocOutput{Documentation: string(out)}, nil
		},
	)
	
	// Step 4: Start the MCP server with stdio transport
	// This enables communication via standard input/output
	return server.Run(context.Background(), &mcp.StdioTransport{})
}