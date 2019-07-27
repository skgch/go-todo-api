package infrastructure

import (
	"github.com/jinzhu/gorm"

	"github.com/skgch/go-todo-api/domain/repository"
	models "github.com/skgch/go-todo-api/models"
)

type TodoRepository struct {
	Db *gorm.DB
}

func NewTodoRepository() repository.TodoRepository {
	return &TodoRepository{}
}

func (r *TodoRepository) FindById(id string) *models.Todo {
	todo := models.Todo{}
	if Db.First(&todo, id).RecordNotFound() {
		return nil
	}
	return &todo
}

func (r *TodoRepository) FindAll() *[]models.Todo {
	todos := []models.Todo{}
	Db.Find(&todos)
	return &todos
}

func (r *TodoRepository) FindByUserID(userID string) *[]models.Todo {
	todos := []models.Todo{}
	Db.Where(&models.Todo{UserID: userID}).Find(&todos)
	return &todos
}

func (r *TodoRepository) Create(todo *models.Todo) *models.Todo {
	Db.Create(todo)
	return todo
}

func (r *TodoRepository) Delete(id string) {
	todo := models.Todo{}
	todo.Id = id
	Db.First(&todo)
	Db.Delete(&todo)
}
