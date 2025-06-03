# Golang Echo Scaffolder MCP Server

This project implements a **Model Context Protocol (MCP) server** in Go that exposes tools for scaffolding [Echo](https://echo.labstack.com/) web applications and related components. It is designed to be used as a **stdio MCP server**, making it compatible with a wide range of MCP clients and agentic AI workflows.

## What is MCP?

The **Model Context Protocol (MCP)** is an open protocol that standardizes how AI applications (LLMs, agents, IDEs, etc.) connect to external tools and data sources. MCP enables composable, secure, and extensible integrations between AI models and your code, files, APIs, and more.

- **Stdio MCP Server:** This project uses the [stdio transport](https://modelcontextprotocol.io/docs/concepts/transports#standard-inputoutput-stdio), meaning it communicates over standard input/output. This is the recommended approach for local integrations and is supported by most MCP clients.

For more about MCP, see [modelcontextprotocol.io](https://modelcontextprotocol.io/) or the [llms.txt](llms.txt) reference.

## Features

This server exposes the following **MCP tools** for use by LLMs and agentic clients:

- **create_app**: Scaffold a new Echo web application.
- **create_model**: Create a new GORM-compatible model and its repository files.
- **create_model_controller**: Create a new controller for a given model.
- **fix_app**: Provide pointers on common issues and how to address them in an Echo web application.

Each tool is defined with a clear input schema and returns detailed instructions or code templates for building Go web applications.

## Installation

You can install this server using Go:

```sh
go install github.com/dougbarrett/mcpgo@latest
```


## Usage

This server is intended to be launched as a **stdio MCP server**. It is typically started by an MCP client (such as Claude Desktop, Continue, Cursor, etc.) which will handle the process management and communication.

To run the server manually for testing:

```sh
go run main.go
```

Or, if installed:

```sh
mcpgo
```

> **Note:** The server will wait for MCP stdio messages on stdin and write responses to stdout. It is not intended to be run as a standalone HTTP server.

## Integrating with MCP Clients

To use this server with an MCP-compatible client:

1. Configure your client to launch the server as a stdio process (see your client's documentation for details).
2. The client will automatically discover the available tools and expose them for use in agentic workflows.

For example, in [Claude Desktop](https://claude.ai/download), you can add this server to your configuration file as follows:

```json
{
  "mcpServers": {
    "mcpgo": {
      "command": "mcpgo"
    }
  }
}
```

## Supported Tools

| Tool Name               | Description                                                        |
|-------------------------|--------------------------------------------------------------------|
| `create_app`            | Scaffold a new Echo web application.                               |
| `create_model`          | Create a new GORM-compatible model and its repository files.       |
| `create_model_controller` | Create a new controller for a given model.                       |
| `fix_app`               | Provide pointers on common issues in an Echo web application.      |

Each tool expects specific input parameters (see the code or MCP client UI for details).

## About Echo and GORM

- [Echo](https://echo.labstack.com/) is a high performance, extensible, minimalist Go web framework.
- [GORM](https://gorm.io/) is a popular ORM library for Go.

This server helps you scaffold applications using these technologies, following best practices for modular Go web development.

## Learn More

- [Model Context Protocol Documentation](https://modelcontextprotocol.io/)
- [Echo Web Framework](https://echo.labstack.com/)
- [GORM ORM](https://gorm.io/)