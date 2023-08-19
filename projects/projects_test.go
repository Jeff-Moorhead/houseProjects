package projects

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTestableStore() *InMemoryProjectStore {

	base := []*Project{
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

	store := NewInMemoryProjectStore()
	store.Init(base)

	return store
}

func TestInMemoryProjectStore_GetAllProjects(t *testing.T) {

	store := getTestableStore()

	expectedTitles := []string{"first project", "second project"}
	got := store.GetAllProjects()

	for _, v := range got {
		assert.Contains(t, expectedTitles, v.Title, "unexpected title found: %v", v.Title)
	}
}

func TestInMemoryProjectStore_GetProject(t *testing.T) {

	store := getTestableStore()

	_, err := store.GetProject("")
	assert.Error(t, err, "expected error for blank title")

	nonExistantTitle := "this title does not exist"
	_, err = store.GetProject(nonExistantTitle)
	assert.Error(t, err, "expected error for non-existant title")

	existingTitle := "first project"
	got, err := store.GetProject(existingTitle)
	assert.NoError(t, err, "unexpected error from store.GetProject: %v", err)

	assert.Equal(t, existingTitle, got.Title, "unexpected title from store.GetProject: expected %v, got %v", existingTitle, got.Title)
}

func TestInMemoryProjectStore_CreateProject(t *testing.T) {

	store := getTestableStore()

	newProject := &Project{
		Title:        "new project",
		Cost:         2.22,
		DurationDays: 2,
		Description:  "a new project",
	}

	err := store.CreateProject(newProject)
	assert.NoError(t, err, "unexpected error when creating project: %v", err)

	got, err := store.GetProject(newProject.Title)
	assert.NoError(t, err, "unexpected error when creating project: %v", err)

	// Double check we actually got the right project back
	assert.Equal(t, newProject.Title, got.Title, "incorrect project title: expected %v, got %v", newProject.Title, got.Title)

	existingProject := &Project{
		Title: "first project",
	}

	err = store.CreateProject(existingProject)
	assert.Error(t, err, "expected error when creating existing project")
}

func TestInMemoryProjectStore_UpdateProject(t *testing.T) {

	store := getTestableStore()

	updated := &Project{
		Title:        "first project",           // existing title, title cannot change
		Cost:         2.22,                      // updated cost
		DurationDays: 2,                         // updated duration
		Description:  "The updated description", // updated description
	}

	got, err := store.UpdateProject(updated)
	assert.NoError(t, err, "unexpected error while updating project: %v", err)
	assert.Equal(t, updated.Cost, got.Cost, "incorrect cost: expected %v, got %v", updated.Cost, got.Cost)
	assert.Equal(t, updated.DurationDays, got.DurationDays, "incorrect duration: expected %v, got %v", updated.DurationDays, got.DurationDays)
	assert.Equal(t, updated.Description, got.Description, "incorrect description: expected %v, got %v", updated.Description, got.Description)

	// Make sure blank title causes an error
	updated.Title = ""
	_, err = store.UpdateProject(updated)
	assert.Error(t, err, "expected error for blank title")

	// Make sure a non-existant title causes an error
	updated.Title = "does not exist"
	_, err = store.UpdateProject(updated)
	assert.Error(t, err, "expected error for non-existant title")
}

func TestInMemoryProjectStore_DeleteProject(t *testing.T) {

	store := getTestableStore()

	// Make sure non-existant titles do not cause errors
	err := store.DeleteProject("does not exist")
	assert.NoError(t, err, "unexpected error while deleting project with non-existant title: %v", err)

	// Make sure blank titles do not cause errors
	err = store.DeleteProject("")
	assert.NoError(t, err, "unexpected error while deleting project with blank title: %v", err)

	// Make sure we can delete an existing title
	existing := "first project"
	err = store.DeleteProject(existing)
	assert.NoError(t, err, "unexpected error while deleting project: %v", err)

	_, err = store.GetProject(existing)
	assert.Error(t, err, "expected error for already-deleted project lookup")
}
