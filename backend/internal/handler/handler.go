package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"tenders/internal/config"
	"tenders/internal/domain"
	"time"
)

func NewRouter(cfg *config.Config) {
	r := mux.NewRouter()

	r.HandleFunc("/api/ping", ping).Methods("GET")

	// Тендеры
	r.HandleFunc("/api/tenders", tendersListHandler).Methods("GET")
	r.HandleFunc("/api/tenders/new", tenderCreateHandler).Methods("POST")
	r.HandleFunc("/api/tenders/my", tendersMyHandler).Methods("GET")
	r.HandleFunc("/api/tenders/{tenderId}/edit", tenderUpdateHandler).Methods("PATCH")

	// Предложения
	r.HandleFunc("/api/bids/new", bidCreateHandler).Methods("POST")
	r.HandleFunc("/api/bids/my", bidsMyHandler).Methods("GET")
	r.HandleFunc("/api/bids/{bidId}/edit", BidsEditHandler).Methods("PATCH")

	go func() {
		for {
			err := http.ListenAndServe(cfg.ServerAddress, r)
			if err != nil {
				slog.Info("Error starting server:", err)
				// Ждем несколько секунд перед перезапуском
				time.Sleep(5 * time.Second)
				slog.Info("Error starting server:", err)
			}

			slog.Info("Error starting server:", err)
		}
	}()

	// Обработка сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("Shutting down server...")
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseOffsetParams(w http.ResponseWriter, r *http.Request) (int, int) {
	// Получаем параметры offset и limit из запроса
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	// Преобразуем параметры в целые числа
	offset := 0 // значение по умолчанию
	limit := 10 // значение по умолчанию

	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)

		}
	}

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
		}
	}

	return limit, offset
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func parseServiceTypeParams(r *http.Request) ([]string, error) {
	serviceTypes := r.URL.Query()["service_type"]
	for _, serviceType := range serviceTypes {
		switch capitalizeFirstLetter(serviceType) {
		case domain.ServiceTypeConstruction:
		case domain.ServiceTypeManufacture:
		case domain.ServiceTypeDelivery:
		default:
			return nil, errors.New("invalid service type")
		}
	}

	return serviceTypes, nil
}
