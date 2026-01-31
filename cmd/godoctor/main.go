package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GodocInput struct {
	Package string `json:"package"`
	Symbol  string `json:"symbol,omitempty"`
}

type GodocOutput struct {
	Documentation string `json:"documentation"`
}

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "err", err)
		os.Exit(1)
	}
}

type HelloWorldInput struct{}

type HelloWorldOutput struct {
	Message string `json:"message"`
}

func run() error {
	server := mcp.NewServer(&mcp.Implementation{Name: "godoctor", Version: "v0.0.1"}, nil)
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "hello_world",
			Description: "A simple tool that returns a greeting.",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input HelloWorldInput) (*mcp.CallToolResult, HelloWorldOutput, error) {
			return nil, HelloWorldOutput{Message: "Hello, world!"}, nil
		},
	)
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "godoc",
			Description: "Retrieves Go documentation for a given package and optional symbol. This tool is useful for understanding how to use a package or a specific function, variable, or type within that package. The 'package' argument is the full import path of the package (e.g., 'fmt', 'github.com/modelcontextprotocol/go-sdk/mcp'). The optional 'symbol' argument is the name of the function, variable, or type to get documentation for (e.g., 'Printf', 'NewServer').",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input GodocInput) (*mcp.CallToolResult, GodocOutput, error) {
			args := []string{"doc"}
			if input.Symbol != "" {
				args = append(args, input.Package, input.Symbol)
			} else {
				args = append(args, input.Package)
			}
			cmd := exec.Command("go", args...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return nil, GodocOutput{}, fmt.Errorf("failed to execute go doc: %w, output: %s", err, string(out))
			}
			return nil, GodocOutput{Documentation: string(out)}, nil
		},
	)
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)
	return http.ListenAndServe(":8080", handler)
}