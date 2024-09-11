package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"tenders/db"
	"tenders/internal/domain"
)

func tendersList(w http.ResponseWriter, r *http.Request) {
	// Настройки подключения к базе данных
	db := db.GetConnection()
	defer db.Close()

	// Выполнение запроса
	rows, err := db.Query("SELECT id, description, organization_id, employee_id, version, status FROM tender")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tenders []domain.TenderResp

	// Обработка результатов
	for rows.Next() {
		var tender domain.TenderResp

		err = rows.Scan(&tender.ID, &tender.Description, &tender.OrganizationId, &tender.EmployeeId, &tender.Version, &tender.Status)
		if err != nil {
			log.Fatal(err)
		}

		tenders = append(tenders, tender)
	}

	// Проверка на ошибки после прохода по результатам
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Установка заголовка Content-Type
	w.Header().Set("Content-Type", "application/json")
	// Сериализация данных в JSON и запись в ResponseWriter
	if err := json.NewEncoder(w).Encode(tenders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func tenderCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tenders create: %s")
}

func tendersByUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Username: %s", username)
}

func tenderUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["tenderId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}

func tenderRollback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["tenderId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}
