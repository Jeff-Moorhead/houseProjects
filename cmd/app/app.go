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

	a.router.GET("/projects", func(c echo.Context) error {
		allProjects := a.database.GetAllProjects()
		return c.JSON(http.StatusOK, allProjects)
	})

	a.router.GET("/projects/:title", func(c echo.Context) error {
		title := c.Param("title")
		if title == "" {
			return echo.NewHTTPError(http.StatusBadRequest, projects.ErrBlankTitle.Error())
		}

		proj, err := a.database.GetProject(c.Param("title"))
		if errors.Is(err, projects.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, projects.ErrProjectNotFound.Error())
		}

		return c.JSON(http.StatusOK, proj)
	})

	a.router.POST("/projects/new", func(c echo.Context) error {
		newProject := new(projects.Project)
		if err := c.Bind(newProject); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if newProject.Title == "" {
			return echo.NewHTTPError(http.StatusBadRequest, projects.ErrBlankTitle.Error())
		}

		if err := a.database.CreateProject(newProject); err != nil {
			if errors.Is(err, projects.ErrTitleExists) {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		return c.JSON(http.StatusCreated, newProject)
	})

	a.router.PUT("/projects/:title", func(c echo.Context) error {
		project := new(projects.Project)
		if err := c.Bind(project); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		p, err := a.database.UpdateProject(project)
		if err != nil {
			if errors.Is(err, projects.ErrBlankTitle) || errors.Is(err, projects.ErrProjectNotFound) {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		return c.JSON(http.StatusOK, p)
	})

	a.router.DELETE("/projects/:title", func(c echo.Context) error {
		title := c.Param("title")
		if err := a.database.DeleteProject(title); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	})
}

func (a *App) Run(address string) {

	a.router.Logger.Fatal(a.router.Start(address))
}
