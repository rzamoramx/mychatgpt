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

type HistoryCtrl struct{}

func NewHistoryCtrl() *HistoryCtrl {
	return &HistoryCtrl{}
}

func (class *HistoryCtrl) Delete(c echo.Context) error {
	app, ok := c.Get("MyChatgptApp").(application.MyChatGptApp)
	if !ok {
		fmt.Println("HISTORY -> DELETE: Error getting app: cannot cast to MyChatGptApp")
		return c.JSON(http.StatusInternalServerError, "Error getting app")
	}

	err := app.DeleteHistory(c.Param("id"))
	if err != nil {
		fmt.Println("HISTORY -> DELETE: Error deleting history: ", err)
		return c.JSON(http.StatusInternalServerError, "Error deleting history")
	}

	fmt.Println("HISTORY -> DELETE: History deleted")
	return c.JSON(http.StatusOK, "History deleted")
}

func (class *HistoryCtrl) Home(c echo.Context) error {
	app, ok := c.Get("MyChatgptApp").(application.MyChatGptApp)
	if !ok {
		fmt.Println("HISTORY -> HOME: Error getting app: cannot cast to MyChatGptApp")
		return c.JSON(http.StatusInternalServerError, "Error getting app")
	}

	session, err := utils.GetSession(c)
	if err != nil {
		return err
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		return c.Redirect(http.StatusSeeOther, "/v1/login")
	}

	// get all history
	histories, err := app.RetrieveAllHistories()
	if err != nil {
		fmt.Printf("HISTORY -> HOME: Error getting histories: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error getting histories")
	}

	iaModel := "0"
	if os.Getenv("OPENAI_MODEL") == "gpt-4o" {
		iaModel = "1"
	} else if os.Getenv("OPENAI_MODEL") == "gpt-4o-mini" {
		iaModel = "2"
	} else if os.Getenv("OPENAI_MODEL") == "gpt-4-turbo" {
		iaModel = "3"
	}

	toRender := map[string]interface{}{
		"Histories":   []map[string]string{},
		"Messages":    []string{},
		"BaseUrl":     os.Getenv("BASE_URL"),
		"OpenAiModel": iaModel,
		"Username":    username,
	}

	for _, history := range histories {
		toRender["Histories"] = append(toRender["Histories"].([]map[string]string), map[string]string{
			"id":   history.Id,
			"name": history.Name,
		})
	}

	fmt.Println("HISTORY -> HOME: Rendered: ", toRender)
	return c.Render(http.StatusOK, "chat.html", toRender)
}

func (class *HistoryCtrl) Get(c echo.Context) error {
	app, ok := c.Get("MyChatgptApp").(application.MyChatGptApp)
	if !ok {
		fmt.Println("HISTORY -> GET: Error getting app: cannot cast to MyChatGptApp")
		return c.JSON(http.StatusInternalServerError, "Error getting app")
	}

	historyId := c.Param("id")

	messages, err := app.GetHistoryMessages(historyId)
	if err != nil {
		fmt.Println("HISTORY -> GET: Error getting messages: ", err)
		return c.JSON(http.StatusInternalServerError, "Error getting messages: "+err.Error())
	}

	// replace break lines with html break lines for every message
	for i, message := range messages {
		messages[i]["text"] = strings.ReplaceAll(message["text"], "\n", "<br>")

		// detect if is a language name in the text that suggests what programming language is the code block
		language := utils.DetectLanguage(messages[i]["text"])

		// formatter code blocks
		var preFragment, fragment string = "", ""
		for strings.Contains(messages[i]["text"], "```") {
			preFragment = messages[i]["text"][strings.Index(messages[i]["text"], "```")+3 : strings.LastIndex(messages[i]["text"], "```")]
			fragment = "<pre><code class=\"" + language + "\">" + preFragment + "</code></pre>"
			messages[i]["text"] = strings.Replace(messages[i]["text"], "```"+preFragment+"```", fragment, 1)
		}

		messages[i] = message
	}

	response := ResponseGetHistoryId{
		Messages: messages,
	}

	fmt.Println("HISTORY -> GET: ok, responding with: ", response)
	return c.JSON(http.StatusOK, response)
}

func (class *HistoryCtrl) New(c echo.Context) error {
	app, ok := c.Get("MyChatgptApp").(application.MyChatGptApp)
	if !ok {
		fmt.Println("Error getting app: cannot cast to MyChatGptApp")
		return c.JSON(http.StatusInternalServerError, "Error getting app")
	}

	// Bind body
	var request RequestNewHistory
	if err := c.Bind(&request); err != nil {
		fmt.Println("HISTORY -> NEW: Error binding body: ", err)
		return c.JSON(http.StatusBadRequest, "Error binding body")
	}

	result, err := app.NewHistory(request.Name)
	if err != nil {
		fmt.Println("HISTORY -> NEW: Error processing message: ", err)
		return c.JSON(http.StatusInternalServerError, "Error processing message")
	}

	reponse := ResponseNewHistory{
		Id: result,
	}

	fmt.Println("HISTORY -> NEW: Ok, responding with: ", reponse)
	return c.JSON(http.StatusOK, reponse)
}
