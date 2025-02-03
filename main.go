package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var task string

type requestBody struct {
	Message string `json:"message"`
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	json.NewDecoder(r.Body).Decode(&body)
	task = body.Message
	fmt.Fprintln(w, "OK,", task)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello,", task)
}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()

	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	http.ListenAndServe(":8080", router)
}
