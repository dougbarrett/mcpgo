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
		mcp.WithString("app_name",
			mcp.Description("The name of the application (e.g., mcpgo-app). This is used to generate correct import paths."),
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
	s.AddTool(createModelTool, createModelHandler)

	// Tool: create_model_controller
	createModelControllerTool := mcp.NewTool("create_model_controller",
		mcp.WithDescription("Instructs the LLM to create a new controller for a given model."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application (e.g., mcpgo-app). This is used to generate correct import paths."),
		),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model for which to create a controller (e.g., User, Product)."),
		),
	)
	s.AddTool(createModelControllerTool, createModelControllerHandler)

	// Tool: create_service
	createServiceTool := mcp.NewTool("create_service",
		mcp.WithDescription("Instructs the LLM to create a new service layer with DTOs for a given model."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application (e.g., mcpgo-app). This is used to generate correct import paths."),
		),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model for which to create a service (e.g., User, Product)."),
		),
	)
	s.AddTool(createServiceTool, createServiceHandler)

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
   After creating models, repositories, services, and controllers, you will need to create or update `+"`%[1]s/cmd/web/main.go`"+` to bootstrap these dependencies.
   This typically involves:
   - Importing `+"`gorm.io/driver/sqlite`"+` (or your chosen database driver) and `+"`gorm.io/gorm`"+`.
   - Initializing the database connection (e.g., `+"`db, err := gorm.Open(sqlite.Open(\"gorm.db\"), &gorm.Config{})`"+`).
   - Auto-migrating your models (e.g., `+"`db.AutoMigrate(&models.YourModel{})`"+`).
   - Creating instances of your repositories (e.g., `+"`userRepo := repository.NewUserRepository(db)`"+`).
   - Creating instances of your services (e.g., `+"`userService := service.NewUserService(userRepo)`"+`).
   - Creating instances of your controllers, injecting services (e.g., `+"`userController := controllers.NewUserController(userService)`"+`).
   - Registering routes for your controllers (e.g., `+"`e.POST(\"/users\", userController.CreateUser)`"+`).

   Here's an example of how `+"`%[1]s/cmd/web/main.go`"+` might look after adding a 'User' model with service layer:
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
`, appName, appName, appName, appName, appName, appName)

	return mcp.NewToolResultText(response), nil
}

// Handler for create_model
func createModelHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		appName,        // %[6]s - Hardcoded for now, ideally passed from createAppHandler
	)

	return mcp.NewToolResultText(response), nil
}

// Handler for create_model_controller
func createModelControllerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	appName := request.GetString("app_name", "") // Default app name if not provided
	if appName == "" {
		return mcp.NewToolResultError("App name is required"), nil
	}
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
	"%[5]s/internal/service"
	"%[5]s/internal/dto"
)

type %[3]sController interface {
	Create%[3]s(c echo.Context) error
	Update%[3]s(c echo.Context) error
	Delete%[3]s(c echo.Context) error
	List%[3]s(c echo.Context) error // New: List method
	Get%[3]sByID(c echo.Context) error // New: GetByID method
}

type %[3]sControllerImpl struct {
	%[4]sService service.%[3]sService
}

func New%[3]sController(%[4]sService service.%[3]sService) %[3]sController {
	return &%[3]sControllerImpl{%[4]sService: %[4]sService}
}
`+"```"+`

   b. `+"`create.go`"+` (Create method - JSON request & response):
`+"```go"+`
package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"%[5]s/internal/dto"
)

func (ctrl *%[3]sControllerImpl) Create%[3]s(c echo.Context) error {
	req := new(dto.Create%[3]sRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Add validation here if needed
	result, err := ctrl.%[4]sService.Create(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, result)
}
`+"```"+`

   c. `+"`update.go`"+` (Update method - JSON request & response):
`+"```go"+`
package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"%[5]s/internal/dto"
)

func (ctrl *%[3]sControllerImpl) Update%[3]s(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	
	req := new(dto.Update%[3]sRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req.ID = uint(id)
	
	// Add validation here if needed
	result, err := ctrl.%[4]sService.Update(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
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
	if err := ctrl.%[4]sService.Delete(c.Request().Context(), uint(id)); err != nil {
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
	"strconv"

	"github.com/labstack/echo/v4"
)

func (ctrl *%[3]sControllerImpl) List%[3]s(c echo.Context) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	// You might want to parse query parameters for filtering here
	filters := make(map[string]interface{}) 
	// Example: filters["name"] = c.QueryParam("name")

	result, err := ctrl.%[4]sService.List(c.Request().Context(), page, limit, filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
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
	
	result, err := ctrl.%[4]sService.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
`+"```"+`
`,
		titleModelName, // %[1]s
		lowerModelName, // %[2]s
		titleModelName, // %[3]s
		lowerModelName, // %[4]s
		appName,        // %[5]s - Hardcoded for now, ideally passed from createAppHandler
	)

	return mcp.NewToolResultText(response), nil
}

// Handler for create_service
func createServiceHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	appName := request.GetString("app_name", "") // Default app name if not provided
	if appName == "" {
		return mcp.NewToolResultError("App name is required"), nil
	}
	modelName, err := request.RequireString("model_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting 'model_name': %v", err.Error())), nil
	}

	titleModelName := strings.Title(modelName)
	lowerModelName := strings.ToLower(modelName)

	response := fmt.Sprintf(`# Service Layer and DTOs Scaffold Instructions

To scaffold the service layer with DTOs for model '%[1]s', please perform the following steps:

1. Create the DTOs directory (or ensure it exists):
   mkdir -p internal/dto/%[2]s

2. Create or update the file at internal/dto/%[2]s/dto.go with the following content:

package dto

import "time"

// Create%[1]sRequest represents the request payload for creating a %[2]s
type Create%[1]sRequest struct {
	// Add your fields here based on your model
	// Example fields - replace with actual model fields:
	// Name        string `+"`json:\"name\" validate:\"required\"`"+`
	// Email       string `+"`json:\"email\" validate:\"required,email\"`"+`
	// Description string `+"`json:\"description\"`"+`
}

// Update%[1]sRequest represents the request payload for updating a %[2]s
type Update%[1]sRequest struct {
	ID uint `+"`json:\"id\" validate:\"required\"`"+`
	// Add your fields here based on your model
	// Example fields - replace with actual model fields:
	// Name        *string `+"`json:\"name,omitempty\"`"+`
	// Email       *string `+"`json:\"email,omitempty\"`"+`
	// Description *string `+"`json:\"description,omitempty\"`"+`
}

// %[1]sResponse represents the response payload for %[2]s operations
type %[1]sResponse struct {
	ID        uint      `+"`json:\"id\"`"+`
	CreatedAt time.Time `+"`json:\"created_at\"`"+`
	UpdatedAt time.Time `+"`json:\"updated_at\"`"+`
	// Add your fields here based on your model
	// Example fields - replace with actual model fields:
	// Name        string `+"`json:\"name\"`"+`
	// Email       string `+"`json:\"email\"`"+`
	// Description string `+"`json:\"description\"`"+`
}

// List%[1]sResponse represents the response payload for listing %[2]s
type List%[1]sResponse struct {
	Data  []%[1]sResponse `+"`json:\"data\"`"+`
	Total int          `+"`json:\"total\"`"+`
	Page  int          `+"`json:\"page\"`"+`
	Limit int          `+"`json:\"limit\"`"+`
}

3. Create the service directory (or ensure it exists):
   mkdir -p internal/service/%[2]s

4. Create the service files:

   a. internal/service/%[2]s/service.go (interface and constructor):

package service

import (
	"context"
	"%[3]s/internal/dto"
	"%[3]s/internal/models"
	"%[3]s/internal/repository"
)

type %[1]sService interface {
	Create(ctx context.Context, req *dto.Create%[1]sRequest) (*dto.%[1]sResponse, error)
	Update(ctx context.Context, req *dto.Update%[1]sRequest) (*dto.%[1]sResponse, error)
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*dto.%[1]sResponse, error)
	List(ctx context.Context, page, limit int, filters map[string]interface{}) (*dto.List%[1]sResponse, error)
}

type %[1]sServiceImpl struct {
	%[2]sRepo repository.%[1]sRepository
}

func New%[1]sService(%[2]sRepo repository.%[1]sRepository) %[1]sService {
	return &%[1]sServiceImpl{%[2]sRepo: %[2]sRepo}
}

// Helper function to convert model to DTO
func (s *%[1]sServiceImpl) modelToDTO(model *models.%[1]s) *dto.%[1]sResponse {
	return &dto.%[1]sResponse{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		// Map your model fields to DTO fields here
		// Example:
		// Name:        model.Name,
		// Email:       model.Email,
		// Description: model.Description,
	}
}

// Helper function to convert create DTO to model
func (s *%[1]sServiceImpl) createDTOToModel(req *dto.Create%[1]sRequest) *models.%[1]s {
	return &models.%[1]s{
		// Map your DTO fields to model fields here
		// Example:
		// Name:        req.Name,
		// Email:       req.Email,
		// Description: req.Description,
	}
}

   b. internal/service/%[2]s/create.go (Create method):

package service

import (
	"context"
	"%[3]s/internal/dto"
)

func (s *%[1]sServiceImpl) Create(ctx context.Context, req *dto.Create%[1]sRequest) (*dto.%[1]sResponse, error) {
	// Convert DTO to model
	model := s.createDTOToModel(req)
	
	// Create in repository
	if err := s.%[2]sRepo.Create(ctx, model); err != nil {
		return nil, err
	}
	
	// Convert model back to DTO and return
	return s.modelToDTO(model), nil
}

   c. internal/service/%[2]s/update.go (Update method):

package service

import (
	"context"
	"errors"
	"%[3]s/internal/dto"
)

func (s *%[1]sServiceImpl) Update(ctx context.Context, req *dto.Update%[1]sRequest) (*dto.%[1]sResponse, error) {
	// First, get the existing record
	filters := map[string]interface{}{"id": req.ID}
	existing, err := s.%[2]sRepo.Get(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(existing) == 0 {
		return nil, errors.New("%[2]s not found")
	}
	
	model := &existing[0]
	// Update only the fields that are provided (not nil)
	// Example:
	// if req.Name != nil {
	//     model.Name = *req.Name
	// }
	// if req.Email != nil {
	//     model.Email = *req.Email
	// }
	// if req.Description != nil {
	//     model.Description = *req.Description
	// }
	
	// Update in repository
	if err := s.%[2]sRepo.Update(ctx, model); err != nil {
		return nil, err
	}
	
	// Convert model back to DTO and return
	return s.modelToDTO(model), nil
}

   d. internal/service/%[2]s/delete.go (Delete method):

package service

import "context"

func (s *%[1]sServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.%[2]sRepo.Delete(ctx, id)
}

   e. internal/service/%[2]s/get_by_id.go (GetByID method):

package service

import (
	"context"
	"errors"
	"%[3]s/internal/dto"
)

func (s *%[1]sServiceImpl) GetByID(ctx context.Context, id uint) (*dto.%[1]sResponse, error) {
	filters := map[string]interface{}{"id": id}
	results, err := s.%[2]sRepo.Get(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("%[2]s not found")
	}
	
	return s.modelToDTO(&results[0]), nil
}

   f. internal/service/%[2]s/list.go (List method):

package service

import (
	"context"
	"%[3]s/internal/dto"
)

func (s *%[1]sServiceImpl) List(ctx context.Context, page, limit int, filters map[string]interface{}) (*dto.List%[1]sResponse, error) {
	// Get data from repository
	results, err := s.%[2]sRepo.Get(ctx, filters)
	if err != nil {
		return nil, err
	}
	
	// Convert models to DTOs
	dtoResults := make([]dto.%[1]sResponse, len(results))
	for i, model := range results {
		dtoResults[i] = *s.modelToDTO(&model)
	}
	
	// TODO: Implement proper pagination in repository layer
	// For now, return all results
	return &dto.List%[1]sResponse{
		Data:  dtoResults,
		Total: len(dtoResults),
		Page:  page,
		Limit: limit,
	}, nil
}

5. Update your controller to use the service layer instead of repository directly.
   The controller should now inject the service and use DTOs for request/response.

6. Bootstrap dependencies in cmd/web/main.go:
   After creating services, you will need to update cmd/web/main.go to bootstrap the service layer.
   This typically involves:
   - Creating instances of your repositories (e.g., userRepo := repository.NewUserRepository(db)).
   - Creating instances of your services, injecting repositories (e.g., userService := service.NewUserService(userRepo)).
   - Creating instances of your controllers, injecting services (e.g., userController := controllers.NewUserController(userService)).

   Here's an example of how cmd/web/main.go might look with the service layer:

package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"%[3]s/internal/models"
	"%[3]s/internal/repository"
	"%[3]s/internal/service"
	"%[3]s/internal/controllers"
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
	err = db.AutoMigrate(&models.%[1]s{}) // Add all your models here
	if err != nil {
		e.Logger.Fatal("failed to auto migrate models", err)
	}

	// Initialize repositories
	%[2]sRepo := repository.New%[1]sRepository(db)

	// Initialize services
	%[2]sService := service.New%[1]sService(%[2]sRepo)

	// Initialize controllers
	%[2]sController := controllers.New%[1]sController(%[2]sService)

	// Routes
	e.GET("/", hello)
	e.POST("/%[2]ss", %[2]sController.Create%[1]s)
	e.GET("/%[2]ss/:id", %[2]sController.Get%[1]sByID)
	e.GET("/%[2]ss", %[2]sController.List%[1]s)
	e.PUT("/%[2]ss/:id", %[2]sController.Update%[1]s)
	e.DELETE("/%[2]ss/:id", %[2]sController.Delete%[1]s)

	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
`,
		titleModelName, // %[1]s
		lowerModelName, // %[2]s
		appName,        // %[3]s
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
	responseBuilder.WriteString(fmt.Sprintf("    import (\n        \"%s/internal/models\"\n        \"%s/internal/repository\"\n        \"%s/internal/service\"\n        \"%s/internal/controllers\"\n    )\n", appName, appName, appName, appName))
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
			responseBuilder.WriteString("This error typically means Go cannot find your internal packages. Double-check your import paths to ensure they use your module name (e.g., `mcpgo-app/internal/models`) and run `go mod tidy`.\n")
		}
		// Add more specific error handling logic here if needed
	}

	return mcp.NewToolResultText(responseBuilder.String()), nil
}
