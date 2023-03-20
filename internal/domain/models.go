package domain

type Message struct {
	Id        string `json:"id" firestore:"id"`
	HistoryId string `json:"history_id" firestore:"history_id"`
	Role      string `json:"role" firestore:"role"`
	Prompt    string `json:"prompt" firestore:"prompt"`
	Message   string `json:"message" firestore:"message"`
	Timestamp int64  `json:"timestamp" firestore:"timestamp"`
}
