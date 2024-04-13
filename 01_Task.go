package main

import (
	"fmt"
	"net/http"
	"main"
	"github.com/labstack/echo/v4"
)

func handler(c echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, fmt.Sprintf(" Hi %s sir, I hope ypu had a great day", name))
}

func main() {
	e := echo.New()
	e.GET("/raja", handler)

	e.Start(":3000")
}
