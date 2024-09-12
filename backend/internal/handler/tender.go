package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
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
	var newTenderReq domain.TenderReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newTenderReq)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	db := db.GetConnection()
	defer db.Close()

	// Получим ответственного (есть ли такой у этой компании)
	orgResponsible := domain.OrganizationResponsible{}

	db.
		QueryRow("SELECT id FROM organization_responsible WHERE user_id = (select id from employee where username = $1) and organization_id = $2", newTenderReq.CreatorUserName, newTenderReq.OrganizationId).
		Scan(&orgResponsible.Id)

	if orgResponsible.Id == "" {
		http.Error(w, "Ответственного с указнными данными не существует", http.StatusBadRequest)
		return
	}

	err = db.QueryRow("INSERT INTO tender (name, description, service_type, status, organization_responsible_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", newTenderReq.Name, newTenderReq.Description, newTenderReq.ServiceType, strings.ToUpper(newTenderReq.Status), orgResponsible.Id).Scan(&newTenderReq.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, newTenderReq)
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
