package main

import (
	"github.com/jeff-moorhead/houseProjects/cmd/app"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	a := app.NewApp(e)
	a.Run(":8080")
}
