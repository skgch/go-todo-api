package repository

import (
	models "github.com/skgch/go-todo-api/models"
)

type TodoRepository interface {
	FindById(id string) *models.Todo
	FindAll() *[]models.Todo
	Create(todo *models.Todo) *models.Todo
	Delete(id string)
}
