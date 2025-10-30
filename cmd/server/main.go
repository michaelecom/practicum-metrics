package main

import (
	"flag"
	"log"

	adapter "github.com/michaelecom/practicum-metrics/internal/adapter/gin"

	"github.com/michaelecom/practicum-metrics/internal/adapter/repository"
	"github.com/michaelecom/practicum-metrics/internal/usecase"
)

func main() {
	var addr string

	// Парсинг флагов командной строки
	flag.StringVar(&addr, "a", "localhost:8080", "HTTP server endpoint address")

	flag.Parse()

	// Создание репозитория
	metricRepo := repository.NewMemoryRepository()

	// Создание бизнес-логики
	metricUseCase := usecase.NewMetricUseCase(metricRepo)

	// Создание HTTP обработчика
	server := adapter.NewServer(metricUseCase)

	// Запуск сервера
	log.Printf("Starting metrics server on %s", addr)
	if err := server.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
