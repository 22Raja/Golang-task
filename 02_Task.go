package main

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

// Data storage place
type User struct {
	key   string
	value string
}

// Handler

func send(c echo.Context) error {
	key := c.QueryParam("key")
	Data = c.Param("data")
	if Data == "json" {
		return c.JSON(http.StatusOK , map[string]string{ "Value" : data[key]}    )

}



func add(c echo.Context) error {
	data := User{}

	defer c.Request().Body.Close()
	val, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String("Error occured")
	}

	final_data := json.Unmarshal(val, &data)
	if final_data != nil {
		return c.String("Error occured")
	}
	return c.String("Json data is stored")

}


func main() {

	e := echo.New()
	// To get the JSON DATA
	e.GET("/get", send)

	e.POST("/add", add)

	e.Start(":3000")

	//fmt.Println("hi this is keerthi Raja")

}
