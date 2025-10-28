package main

import (
	"log"
	"net/http"

	httpAdapter "github.com/michaelecom/practicum-metrics/internal/adapter/http"
	"github.com/michaelecom/practicum-metrics/internal/adapter/repository"
	"github.com/michaelecom/practicum-metrics/internal/usecase"
)

func main() {
	// Создание репозитория
	metricRepo := repository.NewMemoryRepository()

	// Создание бизнес-логики
	metricUseCase := usecase.NewMetricUseCase(metricRepo)

	// Создание HTTP обработчика
	metricHandler := httpAdapter.NewMetricHandler(metricUseCase)

	// Регистрация маршрутов
	http.HandleFunc("/update/", metricHandler.UpdateMetric)

	// Запуск сервера
	addr := "localhost:8080"
	log.Printf("Starting metrics server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
