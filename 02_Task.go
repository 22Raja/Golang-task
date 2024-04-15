package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name  string `json: "name"`
	Email string `json: "email"`
}

var store = make(map[string]string)

func handler_get(c echo.Context) error {
	data := c.QueryParam("data")
	return c.String(http.StatusOK, fmt.Sprintf("your %s :  %s", data, store[data]))

}

func handler_set(c echo.Context) error {
	u := new(User)
	err := c.Bind(u)
	// Access JSON data

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad data"))

	}
	store["email"] = u.Email
	store["name"] = u.Name
	return c.String(http.StatusOK, fmt.Sprintf("The data is stored ,Thanks for the data "))

}

func main() {

	e := echo.New()
	e.GET("/get", handler_get)
	e.POST("/post", handler_set)

	e.Start(":3000")
}
