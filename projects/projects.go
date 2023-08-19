package projects

import (
	"errors"
	"fmt"
)

var ErrBlankTitle = errors.New("project title cannot be blank")

type Project struct {
	Title        string  `json:"title"`
	DurationDays int     `json:"duration_days"`
	Cost         float64 `json:"cost"`
	Description  string  `json:"description"`
}

type ProjectStore interface {
	GetAllProjects() []*Project
	GetProject(title string) (*Project, error)
	CreateProject(p *Project) error
	UpdateProject(p *Project) (*Project, error)
	DeleteProject(title string) error
}

type InMemoryProjectStore struct {
	projects map[string]*Project // use a map for quicker lookup than a slice
}

func NewInMemoryProjectStore() *InMemoryProjectStore {

	return &InMemoryProjectStore{
		projects: make(map[string]*Project),
	}
}

func (store *InMemoryProjectStore) Init(projects []*Project) {

	for _, p := range projects {
		store.projects[p.Title] = p
	}
}

func (store *InMemoryProjectStore) GetAllProjects() []*Project {

	allProjects := make([]*Project, len(store.projects))

	idx := 0
	for _, v := range store.projects {
		allProjects[idx] = v
		idx++
	}

	return allProjects
}

func (store *InMemoryProjectStore) GetProject(title string) (*Project, error) {

	if p, ok := store.projects[title]; ok {
		return p, nil
	}

	return nil, fmt.Errorf("no project with title %v", title)
}

func (store *InMemoryProjectStore) CreateProject(p *Project) error {

	if p.Title == "" {
		return ErrBlankTitle
	}

	// Check if project title is already used
	if _, ok := store.projects[p.Title]; ok {
		return fmt.Errorf("cannot use existing title %v for new project", p.Title)
	}

	store.projects[p.Title] = p
	return nil
}

func (store *InMemoryProjectStore) UpdateProject(p *Project) (*Project, error) {

	if p.Title == "" {
		return nil, ErrBlankTitle
	}

	if existing, ok := store.projects[p.Title]; ok {
		// update existing record
		if existing.Cost != p.Cost {
			existing.Cost = p.Cost
		}

		if existing.DurationDays != p.DurationDays {
			existing.DurationDays = p.DurationDays
		}

		if existing.Description != p.Description {
			existing.Description = p.Description
		}

		return existing, nil

	}

	// If title does not exist
	return nil, fmt.Errorf("project with title %v does not exist", p.Title)
}

func (store *InMemoryProjectStore) DeleteProject(title string) error {

	delete(store.projects, title)
	return nil
}
