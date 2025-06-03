package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetProduceModelBoilerplateTool returns the tool definition for produce_model_boilerplate
func GetProduceModelBoilerplateTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("produce_model_boilerplate",
		mcp.WithDescription("Instructs the LLM to output an example boilerplate for a new GORM-compatible model and its repository files."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application. This is used to output an example of correct import paths."),
		),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model (e.g., User, Product)."),
		),
		mcp.WithString("fields",
			mcp.Required(),
			mcp.Description("A JSON array of objects, where each object has 'name' (string) and 'type' (string) for the model fields."),
		),
	)

	return tool, ProduceModelBoilerplateHandler
}

// ProduceModelBoilerplateHandler handles requests to generate boilerplate for a GORM-compatible model
// It creates the model struct and repository files with CRUD operations
func ProduceModelBoilerplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	appName := request.GetString("app_name", "") // Default app name if not provided
	if appName == "" {
		return mcp.NewToolResultError("App name is required"), nil
	}
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

Note: The model includes 'gorm.Model' which provides the following fields automatically:
- ID (uint, primary key)
- CreatedAt (time.Time)
- UpdatedAt (time.Time)
- DeletedAt (soft delete with index)

These fields don't need to be added manually to your model.

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
   After creating models, repositories, services, and controllers, you will need to create or update `+"`cmd/web/main.go`"+` to bootstrap these dependencies.
   This typically involves:
   - Importing `+"`gorm.io/driver/sqlite`"+` (or your chosen database driver) and `+"`gorm.io/gorm`"+`.
   - Initializing the database connection (e.g., `+"`db, err := gorm.Open(sqlite.Open(\"gorm.db\"), &gorm.Config{})`"+`).
   - Auto-migrating your models (e.g., `+"`db.AutoMigrate(&models.YourModel{})`"+`).
   - Creating instances of your repositories (e.g., `+"`userRepo := repository.NewUserRepository(db)`"+`).
   - Creating instances of your services, injecting repositories (e.g., `+"`userService := service.NewUserService(userRepo)`"+`).
   - Creating instances of your controllers, injecting services (e.g., `+"`userController := controllers.NewUserController(userService)`"+`).
   - Registering routes for your controllers (e.g., `+"`e.POST(\"/users\", userController.CreateUser)`"+`).

   **Important Note**: It is recommended to use a service layer between your controllers and repositories. Controllers should not communicate directly with repositories. Instead, controllers should use services, and services should use repositories. This promotes better separation of concerns and makes your code more maintainable.

   Here's an example of how `+"`cmd/web/main.go`"+` might look after adding a 'User' model with service layer:
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
	"%[6]s/internal/service"
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

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)

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
		appName,        // %[6]s - Hardcoded for now, ideally passed from generateAppBoilerplateHandler
	)

	return mcp.NewToolResultText(response), nil
}
