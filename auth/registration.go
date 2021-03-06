package auth

import (
	"github.com/labstack/echo"
	"log"
	"net/http"


)

func Register(c echo.Context) error {
	u := new(User)
	err := c.Bind(u)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, u.Name+u.Email)
}