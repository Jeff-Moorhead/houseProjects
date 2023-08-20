package app

import (
	"bytes"
	"encoding/json"
	"github.com/jeff-moorhead/houseProjects/projects"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTest() *App {

	e := echo.New()
	db := new(mockDatabase)

	return NewApp(e, db)
}

type mockDatabase struct{}

func (mock *mockDatabase) Init([]*projects.Project) {

}

func (mock *mockDatabase) GetAllProjects() []*projects.Project {

	return []*projects.Project{
		{
			Title:        "first project",
			Cost:         1.11,
			DurationDays: 1,
			Description:  "The first project",
		},
		{
			Title:        "second project",
			Cost:         2.22,
			DurationDays: 2,
			Description:  "The second project",
		},
	}
}

func (mock *mockDatabase) GetProject(title string) (*projects.Project, error) {

	return &projects.Project{
		Title:        title,
		Cost:         1.11,
		DurationDays: 2,
		Description:  "The first project",
	}, nil
}

func (mock *mockDatabase) CreateProject(*projects.Project) error {

	return nil
}

func (mock *mockDatabase) UpdateProject(p *projects.Project) (*projects.Project, error) {

	return p, nil
}

func (mock *mockDatabase) DeleteProject(string) error {

	return nil
}
func TestApp_GetAllProjects(t *testing.T) {

	app := setupTest()

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	resp := httptest.NewRecorder()
	c := app.router.NewContext(req, resp)

	err := app.getAllProjects(c)
	if assert.NoError(t, err, "unexpected error: %v", err) {
		assert.Equal(t, http.StatusOK, resp.Code, "incorrect status code: expected %v, got %v", http.StatusOK, resp.Code)

		var got []*projects.Project
		_ = json.NewDecoder(resp.Body).Decode(&got)

		// Let's just make sure we got the right number of projects in the response body. The data is constant
		assert.Len(t, got, 2, "unexpected project count: expected %v, got %v", 2, len(got))
	}
}

func TestApp_GetProjectByTitle(t *testing.T) {

	const expectedTitle = "expected title"
	app := setupTest()

	// Path will be set later in the echo.Context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	c := app.router.NewContext(req, resp)
	c.SetPath("/projects/:title")
	c.SetParamNames("title")
	c.SetParamValues(expectedTitle)

	err := app.getProjectByTitle(c)
	if assert.NoError(t, err, "unexpected error: %v", err) {
		assert.Equal(t, http.StatusOK, resp.Code, "incorrect status code: expected %v, got %v", http.StatusOK, resp.Code)

		var got *projects.Project
		_ = json.NewDecoder(resp.Body).Decode(&got)
		assert.Equal(t, expectedTitle, got.Title, "incorrect project title: expected %v, got %v", expectedTitle, got.Title)
	}
}

func TestApp_CreateProject(t *testing.T) {

	app := setupTest()
	project := &projects.Project{
		Title:        "new project",
		Cost:         2.22,
		DurationDays: 2,
		Description:  "a new project",
	}

	projectJSON, _ := json.Marshal(project)

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(projectJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	c := app.router.NewContext(req, resp)

	err := app.createProject(c)
	if assert.NoError(t, err, "unexpected error: %v", err) {
		assert.Equal(t, http.StatusCreated, resp.Code, "incorrect status code: expected %v, got %v", http.StatusCreated, resp.Code)

		var got projects.Project
		_ = json.NewDecoder(resp.Body).Decode(&got)
		assert.Equal(t, project.Title, got.Title, "incorrect response title: expected %v, got %v", project.Title, got.Title)
	}
}

func TestApp_UpdateProject(t *testing.T) {

	app := setupTest()
	project := &projects.Project{
		Title:        "new project updated",
		Cost:         3.33,
		DurationDays: 3,
		Description:  "an updated project",
	}

	projectJSON, _ := json.Marshal(project)

	// Path will be set in the context
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(projectJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	c := app.router.NewContext(req, resp)
	c.SetPath("/projects/:title")
	c.SetParamNames("title")
	c.SetParamValues(project.Title)

	err := app.updateProject(c)
	if assert.NoError(t, err, "unexpected error: %v", err) {
		assert.Equal(t, http.StatusOK, resp.Code, "incorrect status code: expected %v, got %v", http.StatusOK, resp.Code)

		var got projects.Project
		_ = json.NewDecoder(resp.Body).Decode(&got)
		assert.Equal(t, project.Title, got.Title, "incorrect response title: expected %v, got %v", project.Title, got.Title)
	}
}

func TestApp_DeleteProject(t *testing.T) {

	app := setupTest()

	// Path will be set in the context
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	resp := httptest.NewRecorder()
	c := app.router.NewContext(req, resp)
	c.SetPath("/projects/:title")
	c.SetParamNames("title")
	c.SetParamValues("title to delete")

	err := app.deleteProject(c)
	if assert.NoError(t, err, "unexpected error: %v", err) {
		assert.Equal(t, http.StatusOK, resp.Code, "incorrect status code: expected %v, got %v", http.StatusOK, resp.Code)
	}
}
