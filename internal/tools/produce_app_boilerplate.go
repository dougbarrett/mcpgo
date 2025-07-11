package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetProduceAppBoilerplateTool returns the tool definition for produce_app_boilerplate
func GetProduceAppBoilerplateTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("start_here_produce_app_boilerplate",
		mcp.WithDescription("Instructs the LLM to output an example scaffold a new Echo web application."),
		mcp.WithString("app_name",
			mcp.Required(),
			mcp.Description("The name of the application."),
		),
	)

	return tool, ProduceAppBoilerplateHandler
}

// ProduceAppBoilerplateHandler handles requests to scaffold a new Echo web application
// It returns detailed instructions for creating the application structure and files
func ProduceAppBoilerplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

## Next Steps: Building Your Application Components

After setting up the basic application structure, you can use the following tools to create the various components of your application:

### 1. Create Models

Use the `+"`produce_model_boilerplate`"+` tool to generate model code:

`+"```"+`
produce_model_boilerplate app_name="%[1]s" model_name="User" fields="ID:uint,Name:string,Email:string,CreatedAt:time.Time,UpdatedAt:time.Time"
`+"```"+`

This will generate a model with the specified fields, along with a repository interface and implementation.

### 2. Create Services

Use the `+"`produce_service_boilerplate`"+` tool to generate service layer code:

`+"```"+`
produce_service_boilerplate app_name="%[1]s" model_name="User"
`+"```"+`

This will create a service that handles business logic for your model, connecting to the repository layer.

### 3. Create Controllers

Depending on your needs, you can create either API-based controllers or HTML-based controllers:

#### For API Controllers:

`+"```"+`
produce_api_controller_boilerplate app_name="%[1]s" model_name="User"
`+"```"+`

This will generate RESTful API endpoints for your model.

#### For HTML Controllers:

`+"```"+`
produce_html_controller_boilerplate app_name="%[1]s" model_name="User" template_engine="html/template"
`+"```"+`

This will create controllers that render HTML templates and handle form submissions.

### 4. Integrate Components

After generating these components, update your `+"`cmd/web/main.go`"+` file to:
- Import all the necessary packages
- Initialize the database connection
- Auto-migrate your models
- Create instances of repositories, services, and controllers
- Register routes for your controllers

### 5. Add Dependencies

Don't forget to add the required dependencies:

`+"```"+`
cd %[1]s && go get gorm.io/gorm gorm.io/driver/sqlite github.com/labstack/echo/v4
`+"```"+`

### 6. Run and Test

After setting up all components, run your application:

`+"```"+`
cd %[1]s && go run ./cmd/web
`+"```"+`

Test your endpoints using a tool like curl, Postman, or a web browser depending on your controller type.

`, appName, appName, appName, appName, appName, appName)

	return mcp.NewToolResultText(response), nil
}
