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
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	task = body.Message
	message := Message{Task: task, IsDone: true}
	result := DB.Create(&message)
	if result.Error != nil {
		http.Error(w, "Ошибка сохранения в БД", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "OK,", task)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var requests []Message
	result := DB.Find(&requests)
	if result.Error != nil {
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Все сообщения: %+v", requests)

	// Отправляем текстовый ответ
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, response)
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
