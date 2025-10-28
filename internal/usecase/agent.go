package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/michaelecom/practicum-metrics/internal/adapter/client"
	"github.com/michaelecom/practicum-metrics/internal/adapter/collector"
)

// AgentConfig - конфигурация агента
type AgentConfig struct {
	PollInterval   time.Duration // интервал сбора метрик
	ReportInterval time.Duration // интервал отправки метрик
}

// AgentUseCase - интерфейс use case агента
type AgentUseCase interface {
	// Run запускает агент
	Run(ctx context.Context) error
}

// agentService - реализация use case агента
type agentService struct {
	collector collector.MetricCollector
	sender    client.MetricSender
	config    AgentConfig
}

// NewAgentUseCase создаёт новый use case агента
func NewAgentUseCase(
	collector collector.MetricCollector,
	sender client.MetricSender,
	config AgentConfig,
) AgentUseCase {
	return &agentService{
		collector: collector,
		sender:    sender,
		config:    config,
	}
}

// Run запускает агент с двумя goroutine: сбор и отправка метрик
func (s *agentService) Run(ctx context.Context) error {
	pollTicker := time.NewTicker(s.config.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(s.config.ReportInterval)
	defer reportTicker.Stop()

	// Запускаем первый сбор сразу
	if err := s.pollMetrics(); err != nil {
		log.Printf("Error polling metrics: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-pollTicker.C:
			if err := s.pollMetrics(); err != nil {
				log.Printf("Error polling metrics: %v", err)
			}

		case <-reportTicker.C:
			if err := s.reportMetrics(); err != nil {
				log.Printf("Error reporting metrics: %v", err)
			}
		}
	}
}

// pollMetrics собирает метрики
func (s *agentService) pollMetrics() error {
	// Увеличиваем счётчик опросов
	s.collector.IncrementPollCount()

	// Собираем метрики
	_, err := s.collector.CollectGauges()
	if err != nil {
		return fmt.Errorf("failed to collect gauges: %w", err)
	}

	_, err = s.collector.CollectCounters()
	if err != nil {
		return fmt.Errorf("failed to collect counters: %w", err)
	}

	return nil
}

// reportMetrics отправляет метрики на сервер
func (s *agentService) reportMetrics() error {
	// Собираем gauge метрики
	gauges, err := s.collector.CollectGauges()
	if err != nil {
		return fmt.Errorf("failed to collect gauges: %w", err)
	}

	// Отправляем gauge метрики
	for _, gauge := range gauges {
		if err := s.sender.SendGauge(gauge); err != nil {
			log.Printf("Failed to send gauge %s: %v", gauge.Name().String(), err)
		}
	}

	// Собираем counter метрики
	counters, err := s.collector.CollectCounters()
	if err != nil {
		return fmt.Errorf("failed to collect counters: %w", err)
	}

	// Отправляем counter метрики
	for _, counter := range counters {
		if err := s.sender.SendCounter(counter); err != nil {
			log.Printf("Failed to send counter %s: %v", counter.Name().String(), err)
		}
	}

	return nil
}
