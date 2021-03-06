package auth

import (
	"github.com/labstack/echo"
	"log"
	"net/http"

	model "github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"

)


func Register(c echo.Context) error {

	u := new(model.User)
	err := c.Bind(u)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, u.Name+u.Email)
}