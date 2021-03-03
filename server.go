package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"os"
)

type User struct {
	Name string `json: "name"`
	Email string `json: "email"`
}

func getIvent(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, "event " + id)
}

func show(c echo.Context) error {
	city := c.QueryParam("city")
	typeEvent := c.QueryParam("typeEvent")
	return c.JSON(http.StatusOK, "type " + typeEvent + " in city " + city)
}

func userJson(c echo.Context) error {
	u := new(User)
	err := c.Bind(u)
	if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusCreated, u.Name + u.Email)
}

func createEvent(c echo.Context) error {
	name := c.FormValue("name")
	img, err := c.FormFile("image")
	if err != nil {
		panic(err)
	}
	src, err := img.Open()
 	if err != nil {
 		return err
 	}
 	defer src.Close()
 
 	dst, err := os.Create(img.Filename)
 	if err != nil {
 		return err
 	}
 	defer dst.Close()

	return c.HTML(http.StatusOK, "<b>Thank you! " + name + "</b>")
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Main page")
	})
	e.GET("/events/:id", getIvent)
	e.GET("/show", show)
	e.GET("/users", userJson)
	e.POST("/create", createEvent)
	e.Logger.Fatal(e.Start(":1323"))
}