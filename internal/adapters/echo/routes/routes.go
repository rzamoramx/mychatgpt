package routes

import (
	"my_chat_gpt/internal/adapters/echo/controllers"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Boostrap(e *echo.Echo) error {
	chatCtrlV1 := controllers.NewChatCtrl()
	historyCtrlV1 := controllers.NewHistoryCtrl()
	sessionCtrlV1 := controllers.NewSessionCtrl()

	v1 := e.Group("/v1")
	v1.PATCH("/chat/opeaimodel", chatCtrlV1.ChangeOpenAiModel)
	v1.GET("/home", historyCtrlV1.Home, authMiddleware)
	v1.POST("/chat", chatCtrlV1.ReceiveMessage)
	v1.POST("/history", historyCtrlV1.New)
	v1.GET("/history/:id", historyCtrlV1.Get)
	v1.DELETE("/history/:id", historyCtrlV1.Delete)

	v1.GET("/login", sessionCtrlV1.FormLogin)
	v1.POST("/login", sessionCtrlV1.Login)

	return nil
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}

		_, ok := sess.Values["username"].(string)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/v1/login")
		}

		return next(c)
	}
}
