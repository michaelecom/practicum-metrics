package repository

import (
	"fmt"
	"sync"

	"github.com/michaelecom/practicum-metrics/internal/domain"
)

// MemoryRepository - адаптер для хранения метрик в памяти
type MemoryRepository struct {
	gauges   map[string]*domain.Gauge
	counters map[string]*domain.Counter
	mu       sync.RWMutex
}

// NewMemoryRepository создаёт новый репозиторий в памяти
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		gauges:   make(map[string]*domain.Gauge),
		counters: make(map[string]*domain.Counter),
	}
}

// SaveGauge сохраняет gauge метрику
func (r *MemoryRepository) SaveGauge(metric *domain.Gauge) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := metric.Name().String()
	r.gauges[key] = metric
	return nil
}

// SaveCounter сохраняет counter метрику
func (r *MemoryRepository) SaveCounter(metric *domain.Counter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := metric.Name().String()
	r.counters[key] = metric
	return nil
}

// GetGauge возвращает gauge метрику по имени
func (r *MemoryRepository) GetGauge(name domain.MetricName) (*domain.Gauge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := name.String()
	metric, exists := r.gauges[key]
	if !exists {
		return nil, fmt.Errorf("gauge metric %s not found", key)
	}

	return metric, nil
}

// GetCounter возвращает counter метрику по имени
func (r *MemoryRepository) GetCounter(name domain.MetricName) (*domain.Counter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := name.String()
	metric, exists := r.counters[key]
	if !exists {
		return nil, fmt.Errorf("counter metric %s not found", key)
	}

	return metric, nil
}

// GetOrCreateCounter возвращает существующую counter метрику или создаёт новую с нулевым значением
func (r *MemoryRepository) GetOrCreateCounter(name domain.MetricName) (*domain.Counter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := name.String()
	if counter, exists := r.counters[key]; exists {
		return counter, nil
	}

	// Создаём новый counter с нулевым значением
	value, err := domain.NewCounterValue(0)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter value: %w", err)
	}

	newCounter := domain.NewCounter(name, value)
	r.counters[key] = newCounter

	return newCounter, nil
}
