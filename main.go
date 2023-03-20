package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"my_chat_gpt/internal/adapters/echo/routes"
	"my_chat_gpt/internal/adapters/firestore"
	"my_chat_gpt/internal/adapters/openai"
	"my_chat_gpt/internal/application"
	"my_chat_gpt/internal/ports/persistence"
	"my_chat_gpt/internal/ports/services"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Define the template registry struct
type TemplateRegistry struct {
	templates map[string]*template.Template
}

func main() {
	e := echo.New()

	// avoid this in production!!!
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
	}))

	e.Renderer = &TemplateRegistry{
		templates: registerFuncsAndTemplates(),
	}

	// Wire up the application
	wireUp(e)

	// Bootstrap routes
	err := routes.Boostrap(e)
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	// ping endpoint
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})

	e.Static("/public", "public")

	e.Logger.Fatal(e.Start(":8080"))
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]

	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

// register funcs and templates
func registerFuncsAndTemplates() map[string]*template.Template {
	// Working Directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	// funcs
	funcMap := template.FuncMap{
		"even": func(i int) bool {
			return i%2 == 0
		},
	}

	// templates
	templates := make(map[string]*template.Template)
	templates["chat.html"] = template.Must(
		template.New("").Funcs(funcMap).ParseFiles(
			wd+"/internal/adapters/echo/views/chat.html",
			wd+"/internal/adapters/echo/views/chat_jquery.html",
			wd+"/internal/adapters/echo/views/base.html"))
	//templates["about.html"] = template.Must(template.ParseFiles("view/about.html", "view/base.html"))

	return templates
}

func wireUp(e *echo.Echo) {
	// OpenAI
	oaic := openai.NewOpenaiClient()
	if oaic == nil {
		panic("OpenAI client is nil")
	}
	oaiClient, ok := oaic.(services.AiProvider)
	if !ok {
		panic("OpenAI client is not an AiProvider")
	}

	// Firestore
	fsc, err := firestore.NewFirestoreClient(os.Getenv("MSGS_COLL_NAME"), os.Getenv("HISTORY_COLL_NAME"))
	if err != nil {
		panic(err.Error())
	}
	fsClient, ok := fsc.(persistence.PersistenceProvider)
	if !ok {
		panic("Firestore client is not a PersistenceProvider")
	}

	// Application
	app, err := application.NewMyChatGptApp(oaiClient, fsClient)
	if err != nil {
		panic(err.Error())
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("MyChatgptApp", app)
			return next(c)
		}
	})
}
