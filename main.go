package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var task string

type requestBody struct {
	Id      uint   `json:"id"`
	Message string `json:"message"`
	Status  bool   `json:"is_done"`
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	task = body.Message
	message := Message{Task: task, IsDone: false}
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

func UpdateMessages(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	id := body.Id
	task = body.Message
	done := body.Status

	result := DB.Model(&Message{}).Where("id = ?", id).Updates(Message{Task: task, IsDone: done})

	if result.Error != nil {
		http.Error(w, "Ошибка обновления в БД", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "OK,", id, task)
}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	result := DB.Delete(&Message{}, body.Id)

	if result.Error != nil {
		http.Error(w, "Ошибка удаления из БД", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "OK")
}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()

	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages", UpdateMessages).Methods("PATCH")
	router.HandleFunc("/api/messages", DeleteMessages).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
