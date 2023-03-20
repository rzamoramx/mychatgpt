package utils

import "strings"

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
