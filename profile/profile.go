package profile

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sync"
)

type UserHandler struct {
	Mu     *sync.Mutex
}



func (h *UserHandler) GetProfile(c echo.Context) (string,error) {
	defer c.Request().Body.Close()
	//authorized := false

	cookie, err := c.Cookie("SID")
	fmt.Println(cookie)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)




	return "", nil
}