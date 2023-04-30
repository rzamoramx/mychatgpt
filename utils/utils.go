package utils

import (
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func DetectLanguage(text string) string {
	// list of programming languages
	programmingLanguages := []string{
		"python",
		"javascript",
		"java",
		"c",
		"c++",
		"c#",
		"php",
		"ruby",
		"swift",
		"scala",
		"r",
		"perl",
		"objective-c",
		"matlab",
		"lua",
		"kotlin",
		"go",
		"fortran",
		"elixir",
		"dart",
		"coffeescript",
		"clojure",
		"assembly",
		"ada",
		"abap",
	}

	// detect language in text, make case insensitive
	for _, language := range programmingLanguages {
		if strings.Contains(strings.ToLower(text), strings.ToLower(language)) {
			return "language-" + language
		}
	}

	return ""
}

func GetSession(c echo.Context) (*sessions.Session, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
