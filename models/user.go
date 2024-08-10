package models

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"password"`
}

type TODO struct {
	Id        uint   `json:"id"`
	Content   string `json:"content"`
	Completed string `json:"completed"`
}
