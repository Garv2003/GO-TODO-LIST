package models

type Todo struct {
	Id          string `json:"id"`
	Content     string `json:"content"`
	IsCompleted bool   `json:"isCompleted"`
	UserId      string `json:"userId"`
}
