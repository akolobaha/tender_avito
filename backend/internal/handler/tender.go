package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tenders/db"
	"tenders/internal/domain"
)

func tendersListHandler(w http.ResponseWriter, r *http.Request) {
	db := db.GetConnection()
	defer db.Close()

	limit, offset := parseOffsetParams(r)
	serviceTypes := r.URL.Query()["service_type"]

	// Начинаем с базового запроса
	query := "SELECT id, name, description, status, service_type FROM tender"
	var args []interface{}

	if len(serviceTypes) > 0 {
		placeholders := make([]string, len(serviceTypes))
		for i := range serviceTypes {
			placeholders[i] = fmt.Sprintf("$%d", i+1) // Используем $ для плейсхолдеров
			args = append(args, serviceTypes[i])      // Добавление значений в аргументы
		}
		query += " WHERE service_type IN (" + strings.Join(placeholders, ",") + ")"
	}

	// Добавляем LIMIT и OFFSET
	query += " LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	// Выполнение запроса с подстановкой параметров
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tenders []domain.TenderResp

	// Обработка результатов
	for rows.Next() {
		var tender domain.TenderResp

		err = rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType)
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

	renderJSON(w, tenders)
}

func tenderCreateHandler(w http.ResponseWriter, r *http.Request) {
	var newTenderReq domain.Tender
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

	var newTenderResp domain.Tender

	db.QueryRow("SELECT tender.id, name, description, service_type, status, ore.organization_id, e.username FROM tender JOIN organization_responsible ore ON ore.id = tender.organization_responsible_id JOIN employee e ON e.id = ore.user_id  WHERE tender.id = $1", newTenderReq.Id).Scan(&newTenderResp.Id, &newTenderResp.Name, &newTenderResp.Description, &newTenderResp.ServiceType, &newTenderResp.Status, &newTenderResp.OrganizationId, &newTenderResp.CreatorUserName)

	renderJSON(w, newTenderResp)
}

func tendersMyHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	db := db.GetConnection()
	defer db.Close()

	var tenders []domain.Tender

	rows, err := db.Query("SELECT id, name, service_type, description, status\n\tFROM tender\n\tWHERE organization_responsible_id IN (SELECT id FROM organization_responsible WHERE user_id IN (SELECT id\n\tFROM employee\n\twhere username = $1));", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tender domain.Tender

		err = rows.Scan(&tender.Id, &tender.Name, &tender.ServiceType, &tender.Description, &tender.Status)
		if err != nil {
			log.Fatal(err)
		}

		tenders = append(tenders, tender)
	}

	renderJSON(w, tenders)
}

func tenderUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	tenderIdUuid, err := uuid.Parse(tenderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newTenderPatch domain.TenderPatch
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newTenderPatch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := db.GetConnection()
	defer db.Close()

	// Создаем слайс для хранения частей запроса и аргументов
	var setClauses []string
	var args []interface{}
	if newTenderPatch.Name != "" {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, newTenderPatch.Name)
	}
	if newTenderPatch.Description != "" {
		setClauses = append(setClauses, "description = $"+strconv.Itoa(len(args)+1))
		args = append(args, newTenderPatch.Description)
	}
	if newTenderPatch.ServiceType != "" {
		setClauses = append(setClauses, "service_type = $"+strconv.Itoa(len(args)+1))
		args = append(args, newTenderPatch.ServiceType)
	}

	if len(setClauses) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Объединяем части запроса
	query := "UPDATE tender SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, tenderIdUuid)

	_, err = db.Exec(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newTenderResp domain.Tender

	db.QueryRow("SELECT tender.id, name, description, service_type, status, ore.organization_id, e.username FROM tender JOIN organization_responsible ore ON ore.id = tender.organization_responsible_id JOIN employee e ON e.id = ore.user_id  WHERE tender.id = $1", tenderIdUuid).Scan(&newTenderResp.Id, &newTenderResp.Name, &newTenderResp.Description, &newTenderResp.ServiceType, &newTenderResp.Status, &newTenderResp.OrganizationId, &newTenderResp.CreatorUserName)

	renderJSON(w, newTenderResp)
}

func tenderRollback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["tenderId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}
