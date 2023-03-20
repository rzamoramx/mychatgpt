package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"my_chat_gpt/internal/application"
	"my_chat_gpt/utils"

	"github.com/labstack/echo/v4"
)

type ChatCtrl struct{}

func NewChatCtrl() *ChatCtrl {
	return &ChatCtrl{}
}

func (class *ChatCtrl) ChangeOpenAiModel(c echo.Context) error {
	// Bind body
	var request RequestChangeOpenAiModel
	if err := c.Bind(&request); err != nil {
		fmt.Println("CHAT -> CHANGE OPENAI MODEL: Error binding body: ", err)
		return c.JSON(http.StatusBadRequest, "Error binding body")
	}

	if request.OpenAiModel == "1" {
		os.Setenv("OPENAI_MODEL", "gpt-3.5-turbo")
	} else if request.OpenAiModel == "2" {
		os.Setenv("OPENAI_MODEL", "gpt-4")
	} else {
		fmt.Println("CHAT -> CHANGE OPENAI MODEL: Invalid model: ", request.OpenAiModel)
		return c.JSON(http.StatusBadRequest, "Invalid model")
	}

	fmt.Println("CHAT -> CHANGE OPENAI MODEL: Model changed")
	return c.JSON(http.StatusOK, "Model changed")
}

func (class *ChatCtrl) ReceiveMessage(c echo.Context) error {
	app, ok := c.Get("MyChatgptApp").(application.MyChatGptApp)
	if !ok {
		fmt.Println("CHAT -> RECEIVE MESSAGE: Error getting app: cannot cast to MyChatGptApp")
		return c.JSON(http.StatusInternalServerError, "Error getting app")
	}

	// Bind body
	var request RequestNewMessage
	if err := c.Bind(&request); err != nil {
		fmt.Println("CHAT -> RECEIVE MESSAGE: Error binding body: ", err)
		return c.JSON(http.StatusBadRequest, "Error binding body")
	}

	result, err := app.ProcessMessage(request.HistoryId, request.Message)
	if err != nil {
		fmt.Println("CHAT -> RECEIVE MESSAGE: Error processing message: ", err)
		return c.JSON(http.StatusInternalServerError, "Error processing message")
	}

	// replace break lines with html break lines
	result = strings.ReplaceAll(result, "\n", "<br>")

	// detect if is a language name in the text that suggests what programming language is the code block
	language := utils.DetectLanguage(result)

	// formatter code blocks
	var preFragment, fragment string = "", ""
	for strings.Contains(result, "```") {
		preFragment = result[strings.Index(result, "```")+3 : strings.LastIndex(result, "```")]
		fragment = "<pre><code class=\"" + language + "\">" + preFragment + "</code></pre>"
		result = strings.Replace(result, "```"+preFragment+"```", fragment, 1)
	}

	reponse := ResponseNewMessage{
		Message: result,
	}

	fmt.Println("CHAT -> RECEIVE MESSAGE: Message processed")
	return c.JSON(http.StatusOK, reponse)
}
