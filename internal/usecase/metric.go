package usecase

import (
	"github.com/michaelecom/practicum-metrics/internal/domain"
)

// MetricUseCase - интерфейс для работы с метриками
type MetricUseCase interface {
	// UpdateGauge обновляет gauge метрику
	UpdateGauge(name string, value float64) error

	// UpdateCounter обновляет counter метрику (добавляет к существующему значению)
	UpdateCounter(name string, delta int64) error

	// GetGauge возвращает gauge метрику
	GetGauge(name string) (float64, error)

	// GetCounter возвращает counter метрику
	GetCounter(name string) (int64, error)

	// GetAllMetrics возвращает все метрики в виде map
	GetAllMetrics() (map[string]interface{}, error)
}

// metricService - реализация use case для работы с метриками
type metricService struct {
	repo domain.MetricRepository
}

// NewMetricUseCase создаёт новый use case для метрик
func NewMetricUseCase(repo domain.MetricRepository) MetricUseCase {
	return &metricService{
		repo: repo,
	}
}

// UpdateGauge обновляет gauge метрику
func (s *metricService) UpdateGauge(name string, value float64) error {
	// Создание Value Objects с валидацией
	metricName, err := domain.NewMetricName(name)
	if err != nil {
		return err
	}

	gaugeValue, err := domain.NewGaugeValue(value)
	if err != nil {
		return err
	}

	// Создание доменной сущности
	gauge := domain.NewGauge(metricName, gaugeValue)

	// Сохранение через репозиторий
	return s.repo.SaveGauge(gauge)
}

// UpdateCounter обновляет counter метрику (добавляет к существующему значению)
func (s *metricService) UpdateCounter(name string, delta int64) error {
	// Создание Value Objects с валидацией
	metricName, err := domain.NewMetricName(name)
	if err != nil {
		return err
	}

	deltaValue, err := domain.NewCounterValue(delta)
	if err != nil {
		return err
	}

	// Получение существующей метрики или создание новой
	counter, err := s.repo.GetOrCreateCounter(metricName)
	if err != nil {
		return err
	}

	// Применение бизнес-логики: инкремент
	counter.Increment(deltaValue)

	// Сохранение обновлённой метрики
	return s.repo.SaveCounter(counter)
}

// GetGauge возвращает значение gauge метрики
func (s *metricService) GetGauge(name string) (float64, error) {
	metricName, err := domain.NewMetricName(name)
	if err != nil {
		return 0, err
	}

	gauge, err := s.repo.GetGauge(metricName)
	if err != nil {
		return 0, err
	}

	return gauge.Value().Float64(), nil
}

// GetCounter возвращает значение counter метрики
func (s *metricService) GetCounter(name string) (int64, error) {
	metricName, err := domain.NewMetricName(name)
	if err != nil {
		return 0, err
	}

	counter, err := s.repo.GetCounter(metricName)
	if err != nil {
		return 0, err
	}

	return counter.Value().Int64(), nil
}

// GetAllMetrics возвращает все метрики в виде map
func (s *metricService) GetAllMetrics() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Получаем все gauge метрики
	gauges, err := s.repo.GetAllGauges()
	if err != nil {
		return nil, err
	}

	gaugeMap := make(map[string]float64)
	for _, gauge := range gauges {
		gaugeMap[gauge.Name().String()] = gauge.Value().Float64()
	}

	result["gauges"] = gaugeMap

	// Получаем все counter метрики
	counters, err := s.repo.GetAllCounters()
	if err != nil {
		return nil, err
	}

	counterMap := make(map[string]int64)
	for _, counter := range counters {
		counterMap[counter.Name().String()] = counter.Value().Int64()
	}

	result["counters"] = counterMap

	return result, nil
}
