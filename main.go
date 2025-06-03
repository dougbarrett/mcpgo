package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"mcpgo/internal/tools"
)

// main is the entry point for the MCP server
// This server provides tools for scaffolding Echo web applications
func main() {
	// Create a new MCP server with name, version, and capabilities
	s := server.NewMCPServer(
		"Golang Echo Scaffolder Server",   // Server name
		"1.0.0",                           // Server version
		server.WithToolCapabilities(true), // Enable tool capabilities
	)

	// Add tools from the tools package with guidance on the recommended tool sequence

	// Step 1: Produce App Boilerplate
	appBoilerplateTool, appBoilerplateHandler := tools.GetProduceAppBoilerplateTool()
	appBoilerplateTool.Description += "\n\nNext recommended step: Use 'produce_model_boilerplate' to create your data models."
	s.AddTool(appBoilerplateTool, appBoilerplateHandler)

	// Step 2: Produce Model Boilerplate
	modelBoilerplateTool, modelBoilerplateHandler := tools.GetProduceModelBoilerplateTool()
	modelBoilerplateTool.Description += "\n\nNext recommended step: Use 'produce_service_boilerplate' to create a service layer for your model."
	s.AddTool(modelBoilerplateTool, modelBoilerplateHandler)

	// Step 3: Produce Service Boilerplate
	serviceBoilerplateTool, serviceBoilerplateHandler := tools.GetProduceServiceBoilerplateTool()
	serviceBoilerplateTool.Description += "\n\nNext recommended step: Use 'produce_api_controller_boilerplate' or 'produce_html_controller_boilerplate' to create controllers for your model."
	s.AddTool(serviceBoilerplateTool, serviceBoilerplateHandler)

	// Step 4a: Produce API Controller Boilerplate
	apiControllerBoilerplateTool, apiControllerBoilerplateHandler := tools.GetProduceApiControllerBoilerplateTool()
	apiControllerBoilerplateTool.Description += "\n\nNext recommended step: If needed, use 'produce_html_controller_boilerplate' to create HTML views for your model."
	s.AddTool(apiControllerBoilerplateTool, apiControllerBoilerplateHandler)

	// Step 4b: Produce HTML Controller Boilerplate
	htmlControllerBoilerplateTool, htmlControllerBoilerplateHandler := tools.GetProduceHtmlControllerBoilerplateTool()
	htmlControllerBoilerplateTool.Description += "\n\nNext recommended step: If needed, use 'fix_app' to fix any issues with your application."
	s.AddTool(htmlControllerBoilerplateTool, htmlControllerBoilerplateHandler)

	// Utility: Fix App
	fixAppTool, fixAppHandler := tools.GetFixAppTool()
	s.AddTool(fixAppTool, fixAppHandler)

	// Serve the MCP server using stdio for communication
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
