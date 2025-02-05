package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var task string

type requestBody struct {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		http.Error(w, "Ошибка конвертации в json", http.StatusInternalServerError)
		return
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var requests []Message
	result := DB.Find(&requests)
	if result.Error != nil {
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(requests)
	if err != nil {
		http.Error(w, "Ошибка конвертации в json", http.StatusInternalServerError)
		return
	}
}

func UpdateMessages(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	result := DB.Model(&Message{}).Where("id = ?", id).Updates(Message{
		Task:   body.Message,
		IsDone: body.Status,
	})

	if result.Error != nil {
		http.Error(w, "Ошибка обновления в БД", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Запись не обновлена", http.StatusNotFound)
		return
	}

	var updatedMessage Message
	DB.First(&updatedMessage, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedMessage)
	if err != nil {
		http.Error(w, "Ошибка конвертации в JSON", http.StatusInternalServerError)
		return
	}
}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result := DB.Delete(&Message{}, id)

	if result.Error != nil {
		http.Error(w, "Ошибка удаления из БД", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()

	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", UpdateMessages).Methods("PATCH")
	router.HandleFunc("/api/messages/{id}", DeleteMessages).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
