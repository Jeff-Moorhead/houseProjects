package app

// TODO: configure logging for app router
// TODO: add HTTP testing

import (
	"errors"
	"github.com/jeff-moorhead/houseProjects/projects"
	"github.com/labstack/echo/v4"
	"net/http"
)

type App struct {
	database projects.ProjectStore
	router   *echo.Echo
}

func NewApp(e *echo.Echo) *App {

	a := &App{
		database: projects.NewInMemoryProjectStore(),
		router:   e,
	}

	// TODO: move to a JSON config
	// This variable is just to use to seed a project for testing the API
	baseProjects := []*projects.Project{
		{
			Title:        "Backyard Fence",
			Cost:         2500.00,
			DurationDays: 14,
			Description:  "Fence in the back yard for a dog run",
		},
	}

	a.database.Init(baseProjects)
	a.initRoutes()

	return a
}

func (a *App) initRoutes() {

	a.router.GET("/", func(c echo.Context) error {
		allProjects := a.database.GetAllProjects()
		return c.JSON(http.StatusOK, allProjects)
	})

	a.router.GET("/:title", func(c echo.Context) error {
		title := c.Param("title")
		if title == "" {
			return c.JSON(http.StatusBadRequest, errorResponse{
				Message: projects.ErrBlankTitle.Error(),
			})
		}

		proj, err := a.database.GetProject(c.Param("title"))
		if errors.Is(err, projects.ErrProjectNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse{
				Message: projects.ErrProjectNotFound.Error(),
			})
		}

		return c.JSON(http.StatusOK, proj)
	})
}

func (a *App) Run(address string) {

	a.router.Logger.Fatal(a.router.Start(address))
}
