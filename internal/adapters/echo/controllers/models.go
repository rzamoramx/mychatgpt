package controllers

type RequestLogin struct {
	User     string `json:"user"`
	Password string `json:"pwd"`
}

type RequestNewMessage struct {
	HistoryId string `json:"history_id"`
	Message   string `json:"message"`
}

type ResponseNewMessage struct {
	Message string `json:"message"`
}

type RequestNewHistory struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

type ResponseNewHistory struct {
	Id string `json:"id"`
}

type ResponseGetHistoryId struct {
	Messages []map[string]string `json:"messages"`
}

type RequestChangeOpenAiModel struct {
	OpenAiModel string `json:"open_ai_model"`
}
