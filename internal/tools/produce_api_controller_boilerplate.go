package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetProduceApiControllerBoilerplateTool returns the tool definition for produce_api_controller_boilerplate
func GetProduceApiControllerBoilerplateTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("produce_api_controller_boilerplate",
		mcp.WithDescription("Instructs the LLM to output an example boilerplate for a new API controller for a given model."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application. This is used to output an example of correct import paths."),
		),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model for which to output an example a controller (e.g., User, Product)."),
		),
	)

	return tool, ProduceApiControllerBoilerplateHandler
}

// ProduceApiControllerBoilerplateHandler handles requests to generate boilerplate for an API controller
// It creates controller files with CRUD operations for a given model
func ProduceApiControllerBoilerplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
# API Controller Scaffold Instructions

To scaffold the API controller for model '%[1]s', please perform the following steps:

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
		appName,        // %[5]s - Hardcoded for now, ideally passed from generateAppBoilerplateHandler
	)

	return mcp.NewToolResultText(response), nil
}
