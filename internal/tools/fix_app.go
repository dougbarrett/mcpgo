package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetFixAppTool returns the tool definition for fix_app
func GetFixAppTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("fix_app",
		mcp.WithDescription("Provides pointers on common issues and how to address them in an Echo web application."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application to fix."),
		),
		mcp.WithString("error_message",
			mcp.Description("The specific error message encountered."),
		),
	)

	return tool, FixAppHandler
}

// FixAppHandler provides guidance on common issues in Echo web applications
// It returns detailed instructions for addressing specific errors or general best practices
func FixAppHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	appName := request.GetString("app_name", "")
	errorMessage := request.GetString("error_message", "")

	var responseBuilder strings.Builder
	responseBuilder.WriteString("Here are some pointers to address common issues in your Echo web application:\n\n")

	if appName != "" {
		responseBuilder.WriteString(fmt.Sprintf("For application '%s':\n", appName))
	}

	responseBuilder.WriteString("1.  **Go Module Paths**: Ensure your Go module is correctly initialized and that internal imports use the module name.\n")
	responseBuilder.WriteString(fmt.Sprintf("    If your module is named `%s`, then imports for internal packages should look like:\n", appName))
	responseBuilder.WriteString("    ```go\n")
	responseBuilder.WriteString(fmt.Sprintf("    import (\n        \"%s/internal/models\"\n        \"%s/internal/repository\"\n        \"%s/internal/service\"\n        \"%s/internal/controllers\"\n    )\n", appName, appName, appName, appName))
	responseBuilder.WriteString("    ```\n")
	responseBuilder.WriteString(fmt.Sprintf("    Make sure to replace `%s` with your actual module name.\n\n", appName))

	responseBuilder.WriteString("2.  **Missing Dependencies**: If you see errors like \"no required module provides package...\", run `go mod tidy` in your application's root directory (")
	responseBuilder.WriteString(fmt.Sprintf("`cd %s && go mod tidy`", appName))
	responseBuilder.WriteString(") to fetch missing dependencies.\n\n")

	responseBuilder.WriteString("3.  **Database Initialization**: Ensure your `main.go` (in `cmd/web/`) correctly initializes the GORM database connection and auto-migrates all your models. For example:\n")
	responseBuilder.WriteString("    ```go\n")
	responseBuilder.WriteString(fmt.Sprintf("    import (\n        \"gorm.io/driver/sqlite\"\n        \"gorm.io/gorm\"\n        \"%s/internal/models\"\n    )\n\n", appName))
	responseBuilder.WriteString(`    func main() {
        db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
        if err != nil {
            // handle error
        }
        db.AutoMigrate(&models.Product{}, &models.Customer{}, &models.ProductCustomer{}) // Add all your models
    }
`)
	responseBuilder.WriteString("    ```\n\n")

	responseBuilder.WriteString("4.  **Repository, Service, and Controller Initialization**: Verify that you are creating instances of your repositories, services, and controllers, and injecting dependencies correctly (repositories into services, services into controllers).\n\n")

	responseBuilder.WriteString("5.  **Route Registration**: Ensure all your controller methods have corresponding routes registered in your `main.go` (in `cmd/web/`). For example:\n")
	responseBuilder.WriteString("    ```go\n")
	responseBuilder.WriteString(`    e.POST("/products", productController.CreateProduct)
    e.GET("/products/:id", productController.GetProductByID)
    // ... and so on for all CRUD operations and models
`)
	responseBuilder.WriteString("    ```\n")

	if errorMessage != "" {
		responseBuilder.WriteString(fmt.Sprintf("\n\nRegarding your specific error: \"%s\"\n", errorMessage))
		if strings.Contains(errorMessage, "is not in std") {
			responseBuilder.WriteString("This error typically means Go cannot find your internal packages. Double-check your import paths to ensure they use your module name (e.g., `[appname]/internal/models`) and run `go mod tidy`.\n")
		}
		// Add more specific error handling logic here if needed
	}

	return mcp.NewToolResultText(responseBuilder.String()), nil
}
