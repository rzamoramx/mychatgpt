package routes

import (
	"my_chat_gpt/internal/adapters/echo/controllers"

	"github.com/labstack/echo/v4"
)

func Boostrap(e *echo.Echo) error {
	chatCtrlV1 := controllers.NewChatCtrl()
	historyCtrlV1 := controllers.NewHistoryCtrl()

	v1 := e.Group("/v1")
	v1.PATCH("/chat/opeaimodel", chatCtrlV1.ChangeOpenAiModel)
	v1.GET("/home", historyCtrlV1.Home)
	v1.POST("/chat", chatCtrlV1.ReceiveMessage)
	v1.POST("/history", historyCtrlV1.New)
	v1.GET("/history/:id", historyCtrlV1.Get)
	v1.DELETE("/history/:id", historyCtrlV1.Delete)

	return nil
}
