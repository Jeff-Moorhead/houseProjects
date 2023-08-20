package main

import (
	"github.com/jeff-moorhead/houseProjects/cmd/app"
	"github.com/jeff-moorhead/houseProjects/projects"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	db := projects.NewInMemoryProjectStore()
	a := app.NewApp(e, db)
	a.Run(":8080")
}
