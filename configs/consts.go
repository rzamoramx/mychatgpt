package configs

import "os"

const OPENAI_BASE_URL = "https://api.openai.com/v1"
const OPENAI_CHAT_COMPLETIONS_URL = OPENAI_BASE_URL + "/chat/completions"

var OPENAI_TOKEN = os.Getenv("OPENAI_MCG_TOKEN")
var OPENAI_MODEL = "gpt-4-1106-preview" // "gpt-4-1106-preview", "gpt-4", "gpt-3.5-turbo"
