package domain

type User struct {
	Id       string `json:"id" firestore:"id"`
	Username string `json:"username" firestore:"username"`
	Password string `json:"password" firestore:"pwd"`
}
