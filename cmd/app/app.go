package app

// TODO: configure logging for app router

import (
	"errors"
	"github.com/jeff-moorhead/houseProjects/projects"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type App struct {
	database projects.ProjectStore
	router   *echo.Echo
}

func NewApp(e *echo.Echo, db projects.ProjectStore) *App {

	a := &App{
		database: db,
		router:   e,
	}

	// TODO: move to a JSON config
	// This variable is just to use to seed a project for testing the API
	baseProjects := []*projects.Project{
		{
			Title:        "Backyard Fence",
			Cost:         2500.00,
			DurationDays: 14,
			Description: `Fence in the back yard for a dog run. This project will involve getting a dumpster
for the garbage by the shed, moving a bunch of rocks out of the garden beds to take back some of our property,
and removing the old picket fence that only serves a decorative purpose right now.`,
		},
	}

	a.database.Init(baseProjects)
	a.setCORS()
	a.initRoutes()

	return a
}

func (a *App) setCORS() {

	// Allow cross-origin requests from the frontend server, which runs on Node at localhost:3000
	a.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:81"},
	}))
}

func (a *App) initRoutes() {

	a.router.GET("/projects", a.getAllProjects)
	a.router.GET("/projects/:title", a.getProjectByTitle)
	a.router.POST("/projects", a.createProject)
	a.router.PUT("/projects/:title", a.updateProject)
	a.router.DELETE("/projects/:title", a.deleteProject)
}

func (a *App) Run(address string) {

	a.router.Logger.Fatal(a.router.Start(address))
}

func (a *App) getAllProjects(c echo.Context) error {

	allProjects := a.database.GetAllProjects()
	return c.JSON(http.StatusOK, allProjects)
}

func (a *App) getProjectByTitle(c echo.Context) error {

	title := c.Param("title")
	if title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, projects.ErrBlankTitle.Error())
	}

	proj, err := a.database.GetProject(c.Param("title"))
	if errors.Is(err, projects.ErrProjectNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, projects.ErrProjectNotFound.Error())
	}

	return c.JSON(http.StatusOK, proj)
}

func (a *App) createProject(c echo.Context) error {

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

	c.Logger().Info(newProject)

	return c.JSON(http.StatusCreated, newProject)
}

func (a *App) updateProject(c echo.Context) error {

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
}

func (a *App) deleteProject(c echo.Context) error {

	title := c.Param("title")
	if err := a.database.DeleteProject(title); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
