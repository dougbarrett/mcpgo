package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetProduceServiceBoilerplateTool returns the tool definition for produce_service_boilerplate
func GetProduceServiceBoilerplateTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("produce_service_boilerplate",
		mcp.WithDescription("Instructs the LLM to output an example boilerplate for a new service layer with DTOs for a given model."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application. This is used to output an example of correct import paths."),
		),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model for which to output an example a service (e.g., User, Product)."),
		),
	)

	return tool, ProduceServiceBoilerplateHandler
}

// ProduceServiceBoilerplateHandler handles requests to generate boilerplate for a service layer
// It creates service files with DTOs (Data Transfer Objects) and business logic for a given model
func ProduceServiceBoilerplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

## Understanding DTOs (Data Transfer Objects)

**What are DTOs?**
DTOs (Data Transfer Objects) are objects that carry data between processes or layers in your application. In the context of a web API:
- They define the structure of data sent to and received from your API endpoints
- They separate your internal domain models from your external API contract
- They allow you to control exactly what data is exposed to clients

**When to use DTOs:**
- When your internal model structure differs from what you want to expose in your API
- When you need to validate or transform data before it reaches your domain model
- When you want to version your API without changing your domain models
- When you need to combine data from multiple models into a single response

**Benefits of using DTOs:**
- Decoupling: Changes to your domain models don't necessarily break your API contract
- Security: You can exclude sensitive fields from responses
- Flexibility: You can shape responses differently for different endpoints
- Validation: You can add validation rules specific to API requests

**Should you create a DTO package?**
- **Yes, create a DTO package if:**
  - Your API is public-facing or used by multiple clients
  - Your models contain sensitive fields that shouldn't be exposed
  - Your API request/response structure needs to differ from your database models
  - You need to validate API inputs separately from model validation
  - You're building a medium to large application where maintainability is important

- **You might not need a DTO package if:**
  - You're building a simple prototype or proof-of-concept
  - Your application is very small with minimal API endpoints
  - Your models map directly to your API with no sensitive fields
  - You're the only consumer of your API and don't need a strict contract

For this scaffolding, we'll create a dedicated 'dto' package to contain all your DTOs, organized by model/domain. This follows best practices for separation of concerns and maintainability in medium to large applications.

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
