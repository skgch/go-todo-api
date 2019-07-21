package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	infra "github.com/skgch/go-todo-api/infrastructure"
	models "github.com/skgch/go-todo-api/models"
)

func main() {
	infra.Connect()
	defer infra.Close()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/todos", GetTodos),
		rest.Get("/todos/:id", GetTodo),
		rest.Post("/todos", PostTodo),
		rest.Delete("/todos/:id", DeleteTodo),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

func GetTodos(w rest.ResponseWriter, r *rest.Request) {
	repo := infra.NewTodoRepository()
	todos := repo.FindAll()
	w.WriteJson(todos)
}

func GetTodo(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	repo := infra.NewTodoRepository()
	todo := repo.FindById(id)
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(todo)
}

func PostTodo(w rest.ResponseWriter, r *rest.Request) {
	todo := models.Todo{}
	err := r.DecodeJsonPayload(&todo)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		repo := infra.NewTodoRepository()
		todo := repo.Create(&todo)
		w.WriteJson(todo)
	}
}

func DeleteTodo(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	repo := infra.NewTodoRepository()
	repo.Delete(id)
	w.WriteHeader(http.StatusNoContent)
}
