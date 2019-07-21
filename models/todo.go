package models

type Todo struct {
	Id    string `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
}
