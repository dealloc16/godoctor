# godoctor

`godoctor` is a versatile Go helper tool designed to be used both as a standalone command-line application and as a component in a larger system that follows the Model Context Protocol (MCP). It provides tools for Go development, such as fetching documentation.

## For Users

This section explains how to use the `godoctor` command-line interface (CLI).

### Installation

Pre-built binaries for `godoctor-cli` are available in the `bin/` directory of this project.

### Usage

The `godoctor-cli` provides access to the tools offered by the `godoctor` server.

#### `hello_world`

A simple tool to check if the `godoctor` server is running and accessible.

**Command:**
```sh
bin/godoctor-cli hello_world
```

**Example Output:**
```json
{
  "message": "Hello, world!"
}
```

#### `godoc`

Retrieves Go documentation for a given package and an optional symbol. This tool is a wrapper around the standard `go doc` command.

**Command:**
```sh
bin/godoctor-cli godoc <package> [symbol]
```

**Arguments:**
- `<package>`: The full import path of the package (e.g., `fmt`, `github.com/modelcontextprotocol/go-sdk/mcp`).
- `[symbol]`: (Optional) The name of the function, variable, or type to get documentation for (e.g., `Println`, `NewServer`).

**Examples:**

- **Get documentation for an external package:**
  ```sh
  bin/godoctor-cli godoc fmt
  ```

- **Get documentation for a symbol in an external package:**
  ```sh
  bin/godoctor-cli godoc fmt Println
  ```

- **Get documentation for a local package:**
  ```sh
  bin/godoctor-cli godoc ./cmd/godoctor
  ```

- **Get documentation for a symbol in a local package:**
  ```sh
  bin/godoctor-cli godoc ./cmd/godoctor GodocInput
  ```

## For Developers

This section explains how to build, run, and extend the `godoctor` project.

### Building from Source

You can build the `godoctor` server and the `godoctor-cli` using the following Go commands. The resulting binaries will be placed in the `bin/` directory.

```sh
# Build the godoctor server
go build -o bin/godoctor ./cmd/godoctor

# Build the godoctor CLI
go build -o bin/godoctor-cli ./cmd/godoctor-cli
```

### Project Structure

- `cmd/godoctor/main.go`: The entry point for the `godoctor` MCP server. This is where the tools are defined and registered.
- `cmd/godoctor-cli/main.go`: The entry point for the command-line interface that interacts with the `godoctor` server.
- `bin/`: This directory contains the compiled binaries for the project.
- `go.mod`: The Go module file, which defines the project's dependencies.
- `GEMINI.md`: A file containing configuration and guidelines for the Gemini agent.

### Adding New Tools

To add a new tool to `godoctor`, you need to modify `cmd/godoctor/main.go`. Follow these steps:

1.  **Define Input and Output Structs:** Create structs for your tool's input and output, using `json` tags for serialization.

    ```go
    type MyToolInput struct {
        Param1 string `json:"param1"`
    }

    type MyToolOutput struct {
        Result string `json:"result"`
    }
    ```

2.  **Register the Tool:** Use the `mcp.AddTool` function to register your new tool with the server. Provide the tool's metadata and a handler function that implements the tool's logic.

    ```go
    mcp.AddTool(server,
        &mcp.Tool{
            Name:        "my_tool",
            Description: "A description of what my tool does.",
        },
        func(ctx context.Context, req *mcp.CallToolRequest, input MyToolInput) (*mcp.CallToolResult, MyToolOutput, error) {
            // Your tool's logic here
            result := "Hello, " + input.Param1
            return nil, MyToolOutput{Result: result}, nil
        },
    )
    ```

3.  **Build and Run:** Rebuild the `godoctor` server to include your new tool.

    ```sh
    go build -o bin/godoctor ./cmd/godoctor
    ```

You can then call your new tool using the `godoctor-cli`:

```sh
bin/godoctor-cli my_tool '{"param1": "world"}'
```
