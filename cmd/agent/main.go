package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michaelecom/practicum-metrics/internal/adapter/client"
	"github.com/michaelecom/practicum-metrics/internal/adapter/collector"
	"github.com/michaelecom/practicum-metrics/internal/usecase"
)

func main() {
	var addr string
	var pollInterval int
	var reportInterval int

	// Парсинг флагов командной строки
	flag.StringVar(&addr, "a", "localhost:8080", "HTTP server endpoint address")
	flag.IntVar(&pollInterval, "p", 2, "Poll interval in seconds")
	flag.IntVar(&reportInterval, "r", 10, "Report interval in seconds")

	flag.Parse()

	baseURL := "http://" + addr

	// Создание HTTP клиента
	metricSender := client.NewHTTPClient(baseURL)

	// Создание Runtime коллектора
	metricCollector := collector.NewRuntimeCollector()

	// Создание бизнес-логики
	agentConfig := usecase.AgentConfig{
		PollInterval:   time.Duration(pollInterval) * time.Second,
		ReportInterval: time.Duration(reportInterval) * time.Second,
	}

	agent := usecase.NewAgentUseCase(metricCollector, metricSender, agentConfig)

	// Создание контекста для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan

		log.Println("Received shutdown signal, stopping agent...")
		cancel()
	}()

	// Запуск агента
	log.Println("Starting metrics agent...")
	log.Printf("Poll interval: %v", agentConfig.PollInterval)
	log.Printf("Report interval: %v", agentConfig.ReportInterval)
	if err := agent.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Agent error: %v", err)
	}

	log.Println("Agent stopped gracefully")
}
