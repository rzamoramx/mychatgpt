package domain

type History struct {
	Id     string `json:"id" firestore:"id"`
	UserId string `json:"user_id" firestore:"user_id"`
	Name   string `json:"name" firestore:"name"`
}
