package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetProduceHtmlControllerBoilerplateTool returns the tool definition for produce_html_controller_boilerplate
func GetProduceHtmlControllerBoilerplateTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("produce_html_controller_boilerplate",
		mcp.WithDescription("Instructs the LLM to output an example boilerplate for a new HTML controller using templUI for a given model."),
		mcp.WithString("app_name",
			mcp.Description("The name of the application. This is used to output an example of correct import paths."),
		),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("The name of the model for which to output an example HTML controller (e.g., User, Product)."),
		),
	)

	return tool, ProduceHtmlControllerBoilerplateHandler
}

// ProduceHtmlControllerBoilerplateHandler handles requests to generate boilerplate for an HTML controller using templUI
// It creates controller files with CRUD operations for a given model
func ProduceHtmlControllerBoilerplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
# HTML Controller Scaffold Instructions using templUI

To scaffold the HTML controller for model '%[1]s' using templUI, please perform the following steps:

## Prerequisites

1. Install the templUI CLI and templ:
   `+"`go install github.com/axzilla/templui/cmd/templui@latest`"+`
   `+"`go install github.com/a-h/templ/cmd/templ@latest`"+`

2. Install Tailwind CSS (on Mac):
   `+"`brew install tailwindcss`"+`

## Base Configuration

1. Create the CSS configuration file and base styles:
   `+"`mkdir -p assets/css`"+`
   Create `+"`assets/css/input.css`"+` with the following content:

`+"```css"+`
@import 'tailwindcss';

@custom-variant dark (&:where(.dark, .dark *));

@theme inline {
  --color-border: var(--border);
  --color-input: var(--input);
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-destructive: var(--destructive);
  --color-destructive-foreground: var(--destructive-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(---accent-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-ring: var(--ring);

  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);

  --container-2xl: 1400px;
}

:root {
  --background: hsl(0 0% 100%);
  --foreground: hsl(240 10% 3.9%);
  --muted: hsl(240 4.8% 95.9%);
  --muted-foreground: hsl(240 3.8% 46.1%);
  --popover: hsl(0 0% 100%);
  --popover-foreground: hsl(240 10% 3.9%);
  --card: hsl(0 0% 100%);
  --card-foreground: hsl(240 10% 3.9%);
  --border: hsl(240 5.9% 90%);
  --input: hsl(240 5.9% 90%);
  --primary: hsl(240 5.9% 10%);
  --primary-foreground: hsl(0 0% 98%);
  --secondary: hsl(240 4.8% 95.9%);
  --secondary-foreground: hsl(240 5.9% 10%);
  --accent: hsl(240 4.8% 95.9%);
  --accent-foreground: hsl(240 5.9% 10%);
  --destructive: hsl(0 84.2% 60.2%);
  --destructive-foreground: hsl(0 0% 98%);
  --ring: hsl(240 5.9% 10%);
  --radius: 0.5rem;
}

.dark {
  --background: hsl(240 10% 3.9%);
  --foreground: hsl(0 0% 98%);
  --muted: hsl(240 3.7% 15.9%);
  --muted-foreground: hsl(240 5% 64.9%);
  --popover: hsl(240 10% 3.9%);
  --popover-foreground: hsl(0 0% 98%);
  --card: hsl(240 10% 3.9%);
  --card-foreground: hsl(0 0% 98%);
  --border: hsl(240 3.7% 15.9%);
  --input: hsl(240 3.7% 15.9%);
  --primary: hsl(0 0% 98%);
  --primary-foreground: hsl(240 5.9% 10%);
  --secondary: hsl(240 3.7% 15.9%);
  --secondary-foreground: hsl(0 0% 98%);
  --accent: hsl(240 3.7% 15.9%);
  --accent-foreground: hsl(0 0% 98%);
  --destructive: hsl(0 62.8% 30.6%);
  --destructive-foreground: hsl(0 0% 98%);
  --ring: hsl(240 4.9% 83.9%);
  --radius: 0.5rem;
}

@layer base {
  * {
    @apply border-border;
  }

  body {
    @apply bg-background text-foreground;
    font-feature-settings:
      "rlig" 1,
      "calt" 1;
  }
}
`+"```"+`

2. Create a Makefile for development tools:
   Create `+"`Makefile`"+` in your project root with the following content:

`+"```makefile"+`
# Run templ generation in watch mode
templ:
    templ generate --watch --proxy="http://localhost:8090" --open-browser=false

# Run air for Go hot reload
server:
    air \
    --build.cmd "go build -o tmp/bin/main ./cmd/web/main.go" \
    --build.bin "tmp/bin/main" \
    --build.delay "100" \
    --build.exclude_dir "node_modules" \
    --build.include_ext "go" \
    --build.stop_on_error "false" \
    --misc.clean_on_exit true

# Watch Tailwind CSS changes
tailwind:
    tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

# Start development server with all watchers
dev:
    make -j3 tailwind templ server
`+"```"+`

3. Initialize templUI in your project:
   `+"`templui init`"+`

4. Add required components:
   `+"`templui add button card alert checkbox input`"+`

## Create HTML Controller Structure

1. Create the directory structure:
   `+"`mkdir -p ui/layouts ui/modules ui/pages/%[2]s`"+`

2. Create the base layout:
   Create `+"`ui/layouts/base.templ`"+` with the following content:

`+"```go"+`
package layouts

import (
	"%[5]s/modules"
)

templ ThemeSwitcherScript() {
	{{ handle := templ.NewOnceHandle() }}
	@handle.Once() {
		<script nonce={ templ.GetNonce(ctx) }>
			// Initial theme setup
			document.documentElement.classList.toggle('dark', localStorage.getItem('appTheme') === 'dark');

			document.addEventListener('alpine:init', () => {
				Alpine.data('themeHandler', () => ({
					isDark: localStorage.getItem('appTheme') === 'dark',
					themeClasses() {
						return this.isDark ? 'text-white' : 'bg-white text-black'
					},
					toggleTheme() {
						this.isDark = !this.isDark;
						localStorage.setItem('appTheme', this.isDark ? 'dark' : 'light');
						document.documentElement.classList.toggle('dark', this.isDark);
					}
				}))
			})
		</script>
	}
}

templ BaseLayout() {
	<!DOCTYPE html>
	<html lang="en" class="h-full dark">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<!-- Tailwind CSS (output) -->
			<link href="/assets/css/output.css" rel="stylesheet"/>
			<!-- Alpine.js -->
			<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
			<!-- Theme switcher script -->
			@ThemeSwitcherScript()
		</head>
		<body
			x-data="themeHandler"
			x-bind:class="themeClasses"
		>
			@modules.Navbar()
			{ children... }
		</body>
	</html>
}
`+"```"+`

3. Create the navbar module:
   Create `+"`ui/modules/navbar.templ`"+` with the following content:

`+"```go"+`
package modules

templ Navbar() {
	<nav class="border-b py-3">
		<div class="container mx-auto px-4 flex justify-between items-center">
			<a href="/" class="text-xl font-bold">%[5]s</a>
			<div class="flex items-center gap-4">
				<a href="/%[2]ss" class="hover:underline">%[1]ss</a>
				@ThemeSwitcher()
			</div>
		</div>
	</nav>
}
`+"```"+`

4. Create the theme switcher module:
   Create `+"`ui/modules/theme_switcher.templ`"+` with the following content:

`+"```go"+`
package modules

import "%[5]s/components/button"
import "%[5]s/components/icon"

templ themeSwitcherHandler() {
	{{ handle := templ.NewOnceHandle() }}
	@handle.Once() {
		<script nonce={ templ.GetNonce(ctx) }>
			document.addEventListener('alpine:init', () => {
				Alpine.data('themeSwitcherHandler', () => ({
					isDarkMode() {
						return this.isDark
					},
					isLightMode() {
						return !this.isDark
					}
				}))
			}) 
		</script>
	}
}

type ThemeSwitcherProps struct {
	Class string
}

templ ThemeSwitcher(props ...ThemeSwitcherProps) {
	{{ var p ThemeSwitcherProps }}
	if len(props) > 0 {
		{{ p = props[0] }}
	}
	@themeSwitcherHandler()
	@button.Button(button.Props{
		Size:    button.SizeIcon,
		Variant: button.VariantGhost,
		Class:   p.Class,
		Attributes: templ.Attributes{
			"@click": "toggleTheme",
		},
	}) {
		@DynamicThemeIcon()
	}
}

templ DynamicThemeIcon() {
	<div x-data="themeSwitcherHandler">
		<span x-show="isDarkMode" class="block">
			@LightIcon()
		</span>
		<span x-show="isLightMode" class="block">
			@DarkIcon()
		</span>
	</div>
}

templ DarkIcon() {
	@icon.Moon()
}

templ LightIcon() {
	@icon.SunMedium()
}
`+"```"+`

5. Create the %[1]s pages:

   a. Create `+"`ui/pages/%[2]s/index.templ`"+` (List page):

`+"```go"+`
package %[2]spages

import (
	"%[5]s/layouts"
	"%[5]s/components/button"
	"%[5]s/components/alert"
	"%[5]s/components/icon"
	"%[5]s/internal/dto"
)

templ Index(items []dto.%[3]sResponse, page int, limit int, total int) {
	@layouts.BaseLayout() {
		<div class="container mx-auto px-4 py-8">
			<div class="flex justify-between items-center mb-6">
				<h1 class="text-2xl font-bold">%[1]ss</h1>
				<a href="/%[2]ss/new">
					@button.Button(button.Props{}) {
						Create %[1]s
					}
				</a>
			</div>

			<!-- Example of using Alert component -->
			<div class="mb-6">
				@alert.Alert() {
					@icon.Rocket(icon.Props{Size: 16})
					@alert.Title() {
						%[1]s Management
					}
					@alert.Description() {
						This page allows you to manage your %[1]ss. You can create, view, edit, and delete %[1]ss.
					}
				}
			</div>

			<div class="bg-card rounded-lg shadow overflow-hidden">
				<table class="min-w-full divide-y divide-border">
					<thead class="bg-muted">
						<tr>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">ID</th>
							<!-- Add your model fields here -->
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Name</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Active</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Actions</th>
						</tr>
					</thead>
					<tbody class="bg-card divide-y divide-border">
						for _, item := range items {
							<tr class="hover:bg-muted/50">
								<td class="px-6 py-4 whitespace-nowrap text-sm">{ item.ID.String() }</td>
								<!-- Add your model fields here -->
								<td class="px-6 py-4 whitespace-nowrap text-sm">{ item.Name }</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm">
									{ if item.Active { "Yes" } else { "No" } }
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm flex gap-2">
									<a href={ templ.SafeURL("/%[2]ss/" + item.ID.String()) }>
										@button.Button(button.Props{
											Variant: button.VariantOutline,
											Size: button.SizeSmall,
										}) {
											View
										}
									</a>
									<a href={ templ.SafeURL("/%[2]ss/" + item.ID.String() + "/edit") }>
										@button.Button(button.Props{
											Variant: button.VariantOutline,
											Size: button.SizeSmall,
										}) {
											Edit
										}
									</a>
									<form method="POST" action={ "/%[2]ss/" + item.ID.String() + "/delete" } onsubmit="return confirm('Are you sure you want to delete this %[2]s?')">
										@button.Button(button.Props{
											Variant: button.VariantDestructive,
											Size: button.SizeSmall,
											Type: "submit",
										}) {
											Delete
										}
									</form>
								</td>
							</tr>
						}
					</tbody>
				</table>
			</div>

			<!-- Pagination -->
			if total > 0 {
				<div class="mt-4 flex justify-between items-center">
					<div class="text-sm text-muted-foreground">
						Showing { (page-1)*limit + 1 } to { min((page)*limit, total) } of { total } entries
					</div>
					<div class="flex gap-2">
						if page > 1 {
							<a href={ templ.SafeURL(fmt.Sprintf("/%[2]ss?page=%d&limit=%d", page-1, limit)) }>
								@button.Button(button.Props{
									Variant: button.VariantOutline,
									Size: button.SizeSmall,
								}) {
									Previous
								}
							</a>
						}
						if page*limit < total {
							<a href={ templ.SafeURL(fmt.Sprintf("/%[2]ss?page=%d&limit=%d", page+1, limit)) }>
								@button.Button(button.Props{
									Variant: button.VariantOutline,
									Size: button.SizeSmall,
								}) {
									Next
								}
							</a>
						}
					</div>
				</div>
			}
		</div>
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
`+"```"+`

   b. Create `+"`ui/pages/%[2]s/show.templ`"+` (Detail page):

`+"```go"+`
package %[2]spages

import (
	"%[5]s/layouts"
	"%[5]s/components/button"
	"%[5]s/components/alert"
	"%[5]s/components/icon"
	"%[5]s/internal/dto"
)

templ Show(item dto.%[3]sResponse) {
	@layouts.BaseLayout() {
		<div class="container mx-auto px-4 py-8">
			<div class="mb-6">
				<a href="/%[2]ss">
					@button.Button(button.Props{
						Variant: button.VariantOutline,
					}) {
						← Back to %[1]ss
					}
				</a>
			</div>

			<!-- Example of using Alert component -->
			<div class="mb-6">
				@alert.Alert(alert.Props{
					Variant: alert.VariantInfo,
				}) {
					@icon.Info(icon.Props{Size: 16})
					@alert.Title() {
						%[1]s Details
					}
					@alert.Description() {
						You are viewing the details of a %[1]s. You can edit or delete this %[1]s using the buttons above.
					}
				}
			</div>

			<div class="bg-card rounded-lg shadow overflow-hidden p-6">
				<div class="flex justify-between items-center mb-6">
					<h1 class="text-2xl font-bold">%[1]s Details</h1>
					<div class="flex gap-2">
						<a href={ templ.SafeURL("/%[2]ss/" + item.ID.String() + "/edit") }>
							@button.Button(button.Props{}) {
								Edit
							}
						</a>
						<form method="POST" action={ "/%[2]ss/" + item.ID.String() + "/delete" } onsubmit="return confirm('Are you sure you want to delete this %[2]s?')">
							@button.Button(button.Props{
								Variant: button.VariantDestructive,
								Type: "submit",
							}) {
								Delete
							}
						</form>
					</div>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="space-y-2">
						<p class="text-sm font-medium text-muted-foreground">ID</p>
						<p>{ item.ID.String() }</p>
					</div>
					<!-- Add your model fields here -->
					<div class="space-y-2">
						<p class="text-sm font-medium text-muted-foreground">Name</p>
						<p>{ item.Name }</p>
					</div>
					<div class="space-y-2">
						<p class="text-sm font-medium text-muted-foreground">Active</p>
						<p>{ if item.Active { "Yes" } else { "No" } }</p>
					</div>
					<!-- Add more fields as needed -->
				</div>
			</div>
		</div>
	}
}
`+"```"+`

   c. Create `+"`ui/pages/%[2]s/form.templ`"+` (Create/Edit form):

`+"```go"+`
package %[2]spages

import (
	"%[5]s/layouts"
	"%[5]s/components/button"
	"%[5]s/components/input"
	"%[5]s/components/checkbox"
	"%[5]s/components/alert"
	"%[5]s/internal/dto"
)

type FormMode string

const (
	FormModeCreate FormMode = "create"
	FormModeEdit   FormMode = "edit"
)

templ Form(mode FormMode, item *dto.%[3]sResponse, errors map[string]string) {
	@layouts.BaseLayout() {
		<div class="container mx-auto px-4 py-8">
			<div class="mb-6">
				<a href="/%[2]ss">
					@button.Button(button.Props{
						Variant: button.VariantOutline,
					}) {
						← Back to %[1]ss
					}
				</a>
			</div>

			<div class="bg-card rounded-lg shadow overflow-hidden p-6">
				<h1 class="text-2xl font-bold mb-6">
					if mode == FormModeCreate {
						Create New %[1]s
					} else {
						Edit %[1]s
					}
				</h1>

				<form method="POST" class="space-y-6">
					<!-- Example of using Alert component for form errors -->
					if errorMsg, ok := errors["general"]; ok {
						<div class="mb-6">
							@alert.Alert(alert.Props{
								Variant: alert.VariantDestructive,
							}) {
								@icon.AlertTriangle(icon.Props{Size: 16})
								@alert.Title() {
									Error
								}
								@alert.Description() {
									{ errorMsg }
								}
							}
						</div>
					}

					<!-- Example of using Input component -->
					<div class="space-y-2">
						<label for="name" class="block text-sm font-medium">Name</label>
						@input.Input(input.Props{
							Type: input.TypeText,
							Id: "name",
							Name: "name",
							Value: item.Name,
							Placeholder: "Enter name",
							Required: true,
						})
						if errorMsg, ok := errors["name"]; ok {
							<p class="text-destructive text-sm mt-1">{ errorMsg }</p>
						}
					</div>

					<!-- Example of using Checkbox component -->
					<div class="space-y-2">
						<div class="flex items-center gap-2">
							@checkbox.Checkbox(checkbox.Props{
								Id: "active",
								Name: "active",
								Checked: item.Active,
							})
							<label for="active" class="text-sm font-medium">
								Active
							</label>
						</div>
						if errorMsg, ok := errors["active"]; ok {
							<p class="text-destructive text-sm mt-1">{ errorMsg }</p>
						}
					</div>

					<!-- Add more form fields as needed -->

					<div class="flex justify-end">
						<a href="/%[2]ss" class="mr-2">
							@button.Button(button.Props{
								Variant: button.VariantOutline,
							}) {
								Cancel
							}
						</a>
						@button.Button(button.Props{
							Type: "submit",
						}) {
							if mode == FormModeCreate {
								Create %[1]s
							} else {
								Update %[1]s
							}
						}
					</div>
				</form>
			</div>
		</div>
	}
}
`+"```"+`

6. Create the HTML controller:
   Create `+"`internal/controllers/%[2]s/html_controller.go`"+` with the following content:

`+"```go"+`
package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"%[5]s/internal/service"
	"%[5]s/internal/dto"
	"%[5]s/pages/%[2]s"
)

type %[3]sHtmlController interface {
	Index(c echo.Context) error
	Show(c echo.Context) error
	New(c echo.Context) error
	Create(c echo.Context) error
	Edit(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type %[3]sHtmlControllerImpl struct {
	%[4]sService service.%[3]sService
}

func New%[3]sHtmlController(%[4]sService service.%[3]sService) %[3]sHtmlController {
	return &%[3]sHtmlControllerImpl{%[4]sService: %[4]sService}
}

// Index renders the list page
func (ctrl *%[3]sHtmlControllerImpl) Index(c echo.Context) error {
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

	return %[2]spages.Index(result.Items, page, limit, result.Total).Render(c.Request().Context(), c.Response().Writer)
}

// Show renders the detail page
func (ctrl *%[3]sHtmlControllerImpl) Show(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	result, err := ctrl.%[4]sService.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return %[2]spages.Show(*result).Render(c.Request().Context(), c.Response().Writer)
}

// New renders the create form
func (ctrl *%[3]sHtmlControllerImpl) New(c echo.Context) error {
	// Create an empty item for the form
	item := &dto.%[3]sResponse{}
	return %[2]spages.Form(%[2]spages.FormModeCreate, item, nil).Render(c.Request().Context(), c.Response().Writer)
}

// Create handles the form submission for creating a new item
func (ctrl *%[3]sHtmlControllerImpl) Create(c echo.Context) error {
	req := new(dto.Create%[3]sRequest)
	if err := c.Bind(req); err != nil {
		// Create an empty item for the form
		item := &dto.%[3]sResponse{}
		errors := map[string]string{"general": err.Error()}
		return %[2]spages.Form(%[2]spages.FormModeCreate, item, errors).Render(c.Request().Context(), c.Response().Writer)
	}

	// Add validation here if needed
	result, err := ctrl.%[4]sService.Create(c.Request().Context(), req)
	if err != nil {
		// Return to form with errors
		item := &dto.%[3]sResponse{
			// Map request fields to response fields
			// Example: Name: req.Name,
			// Example: Active: req.Active,
		}
		errors := map[string]string{"general": err.Error()}
		return %[2]spages.Form(%[2]spages.FormModeCreate, item, errors).Render(c.Request().Context(), c.Response().Writer)
	}

	// Redirect to the detail page
	return c.Redirect(http.StatusSeeOther, "/%[2]ss/"+strconv.FormatUint(uint64(result.ID), 10))
}

// Edit renders the edit form
func (ctrl *%[3]sHtmlControllerImpl) Edit(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	result, err := ctrl.%[4]sService.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return %[2]spages.Form(%[2]spages.FormModeEdit, result, nil).Render(c.Request().Context(), c.Response().Writer)
}

// Update handles the form submission for updating an item
func (ctrl *%[3]sHtmlControllerImpl) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	req := new(dto.Update%[3]sRequest)
	if err := c.Bind(req); err != nil {
		// Get the current item for the form
		result, _ := ctrl.%[4]sService.GetByID(c.Request().Context(), uint(id))
		errors := map[string]string{"general": err.Error()}
		return %[2]spages.Form(%[2]spages.FormModeEdit, result, errors).Render(c.Request().Context(), c.Response().Writer)
	}
	req.ID = uint(id)

	// Add validation here if needed
	result, err := ctrl.%[4]sService.Update(c.Request().Context(), req)
	if err != nil {
		// Return to form with errors
		item, _ := ctrl.%[4]sService.GetByID(c.Request().Context(), uint(id))
		errors := map[string]string{"general": err.Error()}
		return %[2]spages.Form(%[2]spages.FormModeEdit, item, errors).Render(c.Request().Context(), c.Response().Writer)
	}

	// Redirect to the detail page
	return c.Redirect(http.StatusSeeOther, "/%[2]ss/"+strconv.FormatUint(uint64(result.ID), 10))
}

// Delete handles the deletion of an item
func (ctrl *%[3]sHtmlControllerImpl) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	if err := ctrl.%[4]sService.Delete(c.Request().Context(), uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect to the list page
	return c.Redirect(http.StatusSeeOther, "/%[2]ss")
}
`+"```"+`

7. Update your main.go to register the HTML routes:
   Add the following to your main.go file:

`+"```go"+`
// Initialize HTML controllers
%[4]sHtmlController := controllers.New%[3]sHtmlController(%[4]sService)

// HTML Routes
e.GET("/%[2]ss", %[4]sHtmlController.Index)
e.GET("/%[2]ss/new", %[4]sHtmlController.New)
e.POST("/%[2]ss", %[4]sHtmlController.Create)
e.GET("/%[2]ss/:id", %[4]sHtmlController.Show)
e.GET("/%[2]ss/:id/edit", %[4]sHtmlController.Edit)
e.POST("/%[2]ss/:id", %[4]sHtmlController.Update)
e.POST("/%[2]ss/:id/delete", %[4]sHtmlController.Delete)

// Serve static files
e.Static("/assets", "assets")
`+"```"+`

8. Start the development server:
   `+"`make dev`"+`

This will:
- Watch and compile templ files
- Start the Go server with hot reload
- Watch and compile Tailwind CSS changes
`,
		titleModelName, // %[1]s
		lowerModelName, // %[2]s
		titleModelName, // %[3]s
		lowerModelName, // %[4]s
		appName,        // %[5]s
	)

	return mcp.NewToolResultText(response), nil
}
