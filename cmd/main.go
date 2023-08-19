package main

import (
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	// TODO: configure e with routes and initialize database access

	e.Logger.Fatal(e.Start(":8080"))
}
