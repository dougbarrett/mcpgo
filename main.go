package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Golang Echo Scaffolder Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Tool: create_app
	createAppTool := mcp.NewTool("create_app",
		mcp.WithDescription("Instructs the LLM to scaffold a new Echo web application."),
		mcp.WithString("app_name",
			mcp.Required(),
			mcp.Description("The name of the application."),
		),
	)
	s.AddTool(createAppTool, createAppHandler)

	// Tool: create_model
	createModelTool := mcp.NewTool("create_model",
		mcp.WithDescription("Instructs the LLM to create a new GORM-compatible model and its repository files."),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model (e.g., User, Product)."),
		),
		mcp.WithString("fields",
			mcp.Required(),
			mcp.Description("A JSON array of objects, where each object has 'name' (string) and 'type' (string) for the model fields."),
		),
	)
	s.AddTool(createModelTool, createModelHandler)

	// Tool: create_model_controller
	createModelControllerTool := mcp.NewTool("create_model_controller",
		mcp.WithDescription("Instructs the LLM to create a new controller for a given model."),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model for which to create a controller (e.g., User, Product)."),
		),
	)
	s.AddTool(createModelControllerTool, createModelControllerHandler)

	// Tool: fix_app
	fixAppTool := mcp.NewTool("fix_app",
		mcp.WithDescription("Provides pointers on common issues and how to address them in an Echo web application."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application to fix."),
		),
		mcp.WithString("error_message",
			mcp.Description("The specific error message encountered."),
		),
	)
	s.AddTool(fixAppTool, fixAppHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Handler for create_app
func createAppHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	appName, err := request.RequireString("app_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting 'app_name': %v", err.Error())), nil
	}

	response := fmt.Sprintf(`
# Echo Web Application Scaffold Instructions

To scaffold the Echo web application '%[1]s', please perform the following steps:

1. Create the directory structure (or ensure it exists):
   `+"`mkdir -p %[1]s/cmd/web`"+`

2. Create or update the file at `+"`%[1]s/cmd/web/main.go`"+` with the following content:
`+"```go"+`
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
`+"```"+`

3. Initialize the Go module and fetch dependencies:
   `+"`cd %[1]s && go mod init %[1]s && go get github.com/labstack/echo/v4 && go mod tidy`"+`

4. To run the server, navigate to the application directory and execute:
   `+"`cd %[1]s && go run ./cmd/web`"+`

5. Bootstrap dependencies in `+"`%[1]s/cmd/web/main.go`"+`:
   After creating models, repositories, and controllers, you will need to create or update `+"`%[1]s/cmd/web/main.go`"+` to bootstrap these dependencies.
   This typically involves:
   - Importing `+"`gorm.io/driver/sqlite`"+` (or your chosen database driver) and `+"`gorm.io/gorm`"+`.
   - Initializing the database connection (e.g., `+"`db, err := gorm.Open(sqlite.Open(\"gorm.db\"), &gorm.Config{})`"+`).
   - Auto-migrating your models (e.g., `+"`db.AutoMigrate(&models.YourModel{})`"+`).
   - Creating instances of your repositories (e.g., `+"`userRepo := repository.NewUserRepository(db)`"+`).
   - Creating instances of your controllers, injecting their dependencies (e.g., `+"`userController := controllers.NewUserController(userRepo)`"+`).
   - Registering routes for your controllers (e.g., `+"`e.POST(\"/users\", userController.CreateUser)`"+`).

   Here's an example of how `+"`%[1]s/cmd/web/main.go`"+` might look after adding a 'User' model and controller:
   `+"```go"+`
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"%[6]s/internal/models"
	"%[6]s/internal/repository"
	"%[6]s/internal/controllers"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database initialization
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal("failed to connect database", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.User{}) // Add all your models here
	if err != nil {
		e.Logger.Fatal("failed to auto migrate models", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize controllers
	userController := controllers.NewUserController(userRepo)

	// Routes
	e.GET("/", hello)
	e.POST("/users", userController.CreateUser)
	e.GET("/users/:id", userController.GetUserByID) // Example for GetByID
	e.GET("/users", userController.ListUsers)       // Example for List
	e.PUT("/users/:id", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
`+"```"+`
`, appName, appName, appName, appName, appName, appName)

	return mcp.NewToolResultText(response), nil
}

// Handler for create_model
func createModelHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	modelName, err := request.RequireString("model_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting 'model_name': %v", err.Error())), nil
	}

	fieldsJSON, err := request.RequireString("fields") // Assuming fields are passed as a JSON string
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting 'fields': %v", err.Error())), nil
	}

	var fields []map[string]string // Use string for name and type
	err = json.Unmarshal([]byte(fieldsJSON), &fields)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid 'fields' JSON format: %v", err.Error())), nil
	}

	// Generate struct fields
	structFields := []string{}
	for _, field := range fields {
		name := field["name"]
		fieldType := field["type"]
		structFields = append(structFields, fmt.Sprintf("\t%s %s `json:\"%s\"`", strings.Title(name), fieldType, name))
	}

	modelContent := fmt.Sprintf(`package models

import "gorm.io/gorm"

type %s struct {
	gorm.Model
%s
}
`, strings.Title(modelName), strings.Join(structFields, "\n"))

	titleModelName := strings.Title(modelName)
	lowerModelName := strings.ToLower(modelName)

	response := fmt.Sprintf(`
# Model and Repository Scaffold Instructions

To scaffold the model '%[1]s' and its repository, please perform the following steps:

1. Create or update the file at `+"`internal/models/%[2]s.go`"+` with the following content:
`+"```go"+`
%[3]s
`+"```"+`

2. Create the repository directory (or ensure it exists):
   `+"`mkdir -p internal/repository/%[2]s`"+`

3. For each of the following, create or update the file in `+"`internal/repository/%[2]s/`"+` as needed:

   a. `+"`repo.go`"+` (constructor and interface for dependency injection):
`+"```go"+`
package repository

import (
	"context"
	"gorm.io/gorm"
	"%[6]s/internal/models"
)

type %[4]sRepository interface {
	Create(ctx context.Context, %[5]s *models.%[4]s) error
	Update(ctx context.Context, %[5]s *models.%[4]s) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, filters map[string]interface{}) ([]models.%[4]s, error)
}

type %[4]sRepositoryImpl struct {
	db *gorm.DB
}

func New%[4]sRepository(db *gorm.DB) %[4]sRepository {
	return &%[4]sRepositoryImpl{db: db}
}
`+"```"+`

   b. `+"`create.go`"+` (Create method):
`+"```go"+`
package repository

import (
	"context"
	"%[6]s/internal/models"
)

func (r *%[4]sRepositoryImpl) Create(ctx context.Context, %[5]s *models.%[4]s) error {
	return r.db.WithContext(ctx).Create(%[5]s).Error
}
`+"```"+`

   c. `+"`update.go`"+` (Update method):
`+"```go"+`
package repository

import (
	"context"
	"%[6]s/internal/models"
)

func (r *%[4]sRepositoryImpl) Update(ctx context.Context, %[5]s *models.%[4]s) error {
	return r.db.WithContext(ctx).Save(%[5]s).Error
}
`+"```"+`

   d. `+"`delete.go`"+` (Delete method):
`+"```go"+`
package repository

import (
	"context"
	"%[6]s/internal/models"
)

func (r *%[4]sRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.%[4]s{}, id).Error
}
`+"```"+`

   e. `+"`get.go`"+` (Get method - many-to-many with filtering):
`+"```go"+`
package repository

import (
	"context"
	"fmt"
	"%[6]s/internal/models"
)

func (r *%[4]sRepositoryImpl) Get(ctx context.Context, filters map[string]interface{}) ([]models.%[4]s, error) {
	var %[5]s []models.%[4]s
	query := r.db.WithContext(ctx)
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%%s = ?", key), value)
	}
	err := query.Find(&%[5]s).Error
	return %[5]s, err
}
`+"```"+`

4. Bootstrap dependencies in `+"`cmd/web/main.go`"+`:
   After creating models, repositories, and controllers, you will need to create or update `+"`cmd/web/main.go`"+` to bootstrap these dependencies.
   This typically involves:
   - Importing `+"`gorm.io/driver/sqlite`"+` (or your chosen database driver) and `+"`gorm.io/gorm`"+`.
   - Initializing the database connection (e.g., `+"`db, err := gorm.Open(sqlite.Open(\"gorm.db\"), &gorm.Config{})`"+`).
   - Auto-migrating your models (e.g., `+"`db.AutoMigrate(&models.YourModel{})`"+`).
   - Creating instances of your repositories (e.g., `+"`userRepo := repository.NewUserRepository(db)`"+`).
   - Creating instances of your controllers, injecting their dependencies (e.g., `+"`userController := controllers.NewUserController(userRepo)`"+`).
   - Registering routes for your controllers (e.g., `+"`e.POST(\"/users\", userController.CreateUser)`"+`).

   Here's an example of how `+"`cmd/web/main.go`"+` might look after adding a 'User' model and controller:
   `+"```go"+`
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"%[6]s/internal/models"
	"%[6]s/internal/repository"
	"%[6]s/internal/controllers"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database initialization
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal("failed to connect database", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.User{}) // Add all your models here
	if err != nil {
		e.Logger.Fatal("failed to auto migrate models", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize controllers
	userController := controllers.NewUserController(userRepo)

	// Routes
	e.GET("/", hello)
	e.POST("/users", userController.CreateUser)
	e.GET("/users/:id", userController.GetUserByID) // Example for GetByID
	e.GET("/users", userController.ListUsers)       // Example for List
	e.PUT("/users/:id", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
`+"```"+`
`,
		titleModelName, // %[1]s
		lowerModelName, // %[2]s
		modelContent,   // %[3]s
		titleModelName, // %[4]s
		lowerModelName, // %[5]s
		"mcpgo-app",    // %[6]s - Hardcoded for now, ideally passed from createAppHandler
	)

	return mcp.NewToolResultText(response), nil
}

// Handler for create_model_controller
func createModelControllerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	modelName, err := request.RequireString("model_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting 'model_name': %v", err.Error())), nil
	}

	titleModelName := strings.Title(modelName)
	lowerModelName := strings.ToLower(modelName)

	response := fmt.Sprintf(`
# Controller Scaffold Instructions

To scaffold the controller for model '%[1]s', please perform the following steps:

1. Create the controller directory (or ensure it exists):
   `+"`mkdir -p internal/controllers/%[2]s`"+`

2. For each of the following, create or update the file in `+"`internal/controllers/%[2]s/`"+` as needed:

   a. `+"`controller.go`"+` (interface and constructor):
`+"```go"+`
package controllers

import (
	"github.com/labstack/echo/v4"
	"%[5]s/internal/repository"
)

type %[3]sController interface {
	Create%[3]s(c echo.Context) error
	Update%[3]s(c echo.Context) error
	Delete%[3]s(c echo.Context) error
	List%[3]s(c echo.Context) error // New: List method
	Get%[3]sByID(c echo.Context) error // New: GetByID method
}

type %[3]sControllerImpl struct {
	%[4]sRepo repository.%[3]sRepository
}

func New%[3]sController(%[4]sRepo repository.%[3]sRepository) %[3]sController {
	return &%[3]sControllerImpl{%[4]sRepo: %[4]sRepo}
}
`+"```"+`

   b. `+"`create.go`"+` (Create method - JSON request & response):
`+"```go"+`
package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"%[5]s/internal/models"
)

func (ctrl *%[3]sControllerImpl) Create%[3]s(c echo.Context) error {
	%[4]s := new(models.%[3]s)
	if err := c.Bind(%[4]s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Add validation here if needed
	if err := ctrl.%[4]sRepo.Create(c.Request().Context(), %[4]s); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, %[4]s)
}
`+"```"+`

   c. `+"`update.go`"+` (Update method - JSON request & response):
`+"```go"+`
package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"%[5]s/internal/models"
)

func (ctrl *%[3]sControllerImpl) Update%[3]s(c echo.Context) error {
	%[4]s := new(models.%[3]s)
	if err := c.Bind(%[4]s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Add validation here if needed
	if err := ctrl.%[4]sRepo.Update(c.Request().Context(), %[4]s); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, %[4]s)
}
`+"```"+`

   d. `+"`delete.go`"+` (Delete method - JSON request & response):
`+"```go"+`
package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (ctrl *%[3]sControllerImpl) Delete%[3]s(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err := ctrl.%[4]sRepo.Delete(c.Request().Context(), uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
`+"```"+`

   e. `+"`list.go`"+` (List method - JSON request & response):
`+"```go"+`
package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (ctrl *%[3]sControllerImpl) List%[3]s(c echo.Context) error {
	// You might want to parse query parameters for filtering here
	filters := make(map[string]interface{}) 
	// Example: filters["name"] = c.QueryParam("name")

	%[4]s, err := ctrl.%[4]sRepo.Get(c.Request().Context(), filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, %[4]s)
}
`+"```"+`

   f. `+"`get_by_id.go`"+` (GetByID method - JSON request & response):
`+"```go"+`
package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (ctrl *%[3]sControllerImpl) Get%[3]sByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	filters := map[string]interface{}{"id": id}
	%[4]s, err := ctrl.%[4]sRepo.Get(c.Request().Context(), filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(%[4]s) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "%[1]s not found")
	}
	return c.JSON(http.StatusOK, %[4]s[0])
}
`+"```"+`
`,
		titleModelName, // %[1]s
		lowerModelName, // %[2]s
		titleModelName, // %[3]s
		lowerModelName, // %[4]s
		"mcpgo-app",    // %[5]s - Hardcoded for now, ideally passed from createAppHandler
	)

	return mcp.NewToolResultText(response), nil
}

// Handler for fix_app
func fixAppHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	responseBuilder.WriteString(fmt.Sprintf("    import (\n        \"%s/internal/models\"\n        \"%s/internal/repository\"\n        \"%s/internal/controllers\"\n    )\n", appName, appName, appName))
	responseBuilder.WriteString("    ```\n")
	responseBuilder.WriteString(fmt.Sprintf("    Make sure to replace `%s` with your actual module name (e.g., `mcpgo-app`).\n\n", appName))

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

	responseBuilder.WriteString("4.  **Repository and Controller Initialization**: Verify that you are creating instances of your repositories and controllers, and injecting the database connection (for repositories) and repositories (for controllers) correctly.\n\n")

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
			responseBuilder.WriteString("This error typically means Go cannot find your internal packages. Double-check your import paths to ensure they use your module name (e.g., `mcpgo-app/internal/models`) and run `go mod tidy`.\n")
		}
		// Add more specific error handling logic here if needed
	}

	return mcp.NewToolResultText(responseBuilder.String()), nil
}
