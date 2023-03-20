package openai

type RequestToAOpenAi struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Msgs      []Message `json:"messages"`
}

type ResponseFromAOpenAi struct {
	Choices []struct {
		Msg           Message `json:"message"`
		Index         int     `json:"index"`
		FisnishReason string  `json:"finish_reason"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
