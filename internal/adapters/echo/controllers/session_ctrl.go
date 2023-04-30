package controllers

import (
	"fmt"
	"my_chat_gpt/internal/application"
	"my_chat_gpt/utils"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type SessionCtrl struct{}

func NewSessionCtrl() *SessionCtrl {
	return &SessionCtrl{}
}

func (class *SessionCtrl) FormLogin(c echo.Context) error {
	toRender := map[string]interface{}{
		"BaseUrl": os.Getenv("BASE_URL"),
	}
	return c.Render(http.StatusOK, "login.html", toRender)
}

func (class *SessionCtrl) Login(c echo.Context) error {
	app, ok := c.Get("MyChatgptApp").(application.MyChatGptApp)
	if !ok {
		fmt.Println("SessionCtrl -> login(): Error getting app: cannot cast to MyChatGptApp")
		return c.JSON(http.StatusInternalServerError, "Error getting app")
	}

	userId, err := app.MakeLogin(c.FormValue("user"), c.FormValue("pwd"))
	if err != nil {
		fmt.Println("Userid: ", userId)
		fmt.Println("SessionCtrl -> login(): ", err)
		return c.String(http.StatusUnauthorized, "Usuario o contraseña incorrectos o error interno, consulta al administrador")
	}

	session, err := utils.GetSession(c)
	if err != nil {
		return err
	}

	session.Values["username"] = c.FormValue("user")
	session.Options.MaxAge = 12 * 60 * 60
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/v1/home")
}
