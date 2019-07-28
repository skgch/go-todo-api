package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
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
		rest.Get("/debug", Debug),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

func GetTodos(w rest.ResponseWriter, r *rest.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		rest.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	tokenString := r.Header.Get("Authorization")
	craims, err := parseToken(tokenString)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if craims["email"] != userID {
		rest.Error(w, "invalid user_id", http.StatusForbidden)
		return
	}

	repo := infra.NewTodoRepository()
	todos := repo.FindByUserID(userID)
	w.WriteJson(todos)
}

func GetTodo(w rest.ResponseWriter, r *rest.Request) {
	tokenString := r.Header.Get("Authorization")
	craims, err := parseToken(tokenString)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.PathParam("id")
	repo := infra.NewTodoRepository()
	todo := repo.FindById(id)
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		rest.NotFound(w, r)
		return
	} else if todo.UserID != craims["email"] {
		w.WriteHeader(http.StatusForbidden)
		rest.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	w.WriteJson(todo)
}

func PostTodo(w rest.ResponseWriter, r *rest.Request) {
	tokenString := r.Header.Get("Authorization")
	craims, err := parseToken(tokenString)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo := models.Todo{}
	err = r.DecodeJsonPayload(&todo)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todo.UserID != craims["email"] {
		w.WriteHeader(http.StatusForbidden)
		rest.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	repo := infra.NewTodoRepository()
	w.WriteJson(repo.Create(&todo))
}

func DeleteTodo(w rest.ResponseWriter, r *rest.Request) {
	tokenString := r.Header.Get("Authorization")
	craims, err := parseToken(tokenString)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.PathParam("id")
	repo := infra.NewTodoRepository()
	todo := repo.FindById(id)
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		rest.NotFound(w, r)
		return
	} else if todo.UserID != craims["email"] {
		w.WriteHeader(http.StatusForbidden)
		rest.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	repo.Delete(id)
	w.WriteHeader(http.StatusNoContent)
}

func Debug(w rest.ResponseWriter, r *rest.Request) {
	tokenString := r.Header.Get("Authorization")
	claims, err := parseToken(tokenString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(claims["email"])
}

func parseToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Cognito UserPools の署名アルゴリズムはRS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		set, err := jwk.FetchHTTP("https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_EzvKreJ2L/.well-known/jwks.json")
		if err != nil {
			return nil, err
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("expecting JWT header to have string kid")
		}

		if key := set.LookupKeyID(keyID); len(key) == 1 {
			return key[0].Materialize()
		}
		return nil, errors.New("unable to find key")
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}
