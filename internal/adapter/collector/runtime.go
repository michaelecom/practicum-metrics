package collector

import (
	"math/rand"
	"runtime"

	"github.com/michaelecom/practicum-metrics/internal/domain"
)

// MetricCollector собирает метрики из различных источников
type MetricCollector interface {
	// CollectGauges собирает все gauge метрики
	CollectGauges() ([]*domain.Gauge, error)
	// CollectCounters собирает все counter метрики
	CollectCounters() ([]*domain.Counter, error)
	// IncrementPollCount увеличивает счётчик опросов
	IncrementPollCount()
}

// RuntimeCollector - собирает метрики из runtime пакета
type RuntimeCollector struct {
	pollCount   *domain.Counter
	randomValue *domain.Gauge
}

// NewRuntimeCollector создаёт новый коллектор runtime метрик
func NewRuntimeCollector() *RuntimeCollector {
	pollCountName, _ := domain.NewMetricName("PollCount")
	pollCountValue, _ := domain.NewCounterValue(0)
	pollCount := domain.NewCounter(pollCountName, pollCountValue)

	randomName, _ := domain.NewMetricName("RandomValue")
	randomValue, _ := domain.NewGaugeValue(0)
	random := domain.NewGauge(randomName, randomValue)

	return &RuntimeCollector{
		pollCount:   pollCount,
		randomValue: random,
	}
}

// CollectGauges собирает все gauge метрики из runtime
func (r *RuntimeCollector) CollectGauges() ([]*domain.Gauge, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gauges := make([]*domain.Gauge, 0, 28)

	// Добавляем все gauge метрики из runtime
	gauges = append(gauges, r.createGauge("Alloc", float64(memStats.Alloc)))
	gauges = append(gauges, r.createGauge("BuckHashSys", float64(memStats.BuckHashSys)))
	gauges = append(gauges, r.createGauge("Frees", float64(memStats.Frees)))
	gauges = append(gauges, r.createGauge("GCCPUFraction", memStats.GCCPUFraction))
	gauges = append(gauges, r.createGauge("GCSys", float64(memStats.GCSys)))
	gauges = append(gauges, r.createGauge("HeapAlloc", float64(memStats.HeapAlloc)))
	gauges = append(gauges, r.createGauge("HeapIdle", float64(memStats.HeapIdle)))
	gauges = append(gauges, r.createGauge("HeapInuse", float64(memStats.HeapInuse)))
	gauges = append(gauges, r.createGauge("HeapObjects", float64(memStats.HeapObjects)))
	gauges = append(gauges, r.createGauge("HeapReleased", float64(memStats.HeapReleased)))
	gauges = append(gauges, r.createGauge("HeapSys", float64(memStats.HeapSys)))
	gauges = append(gauges, r.createGauge("LastGC", float64(memStats.LastGC)))
	gauges = append(gauges, r.createGauge("Lookups", float64(memStats.Lookups)))
	gauges = append(gauges, r.createGauge("MCacheInuse", float64(memStats.MCacheInuse)))
	gauges = append(gauges, r.createGauge("MCacheSys", float64(memStats.MCacheSys)))
	gauges = append(gauges, r.createGauge("MSpanInuse", float64(memStats.MSpanInuse)))
	gauges = append(gauges, r.createGauge("MSpanSys", float64(memStats.MSpanSys)))
	gauges = append(gauges, r.createGauge("Mallocs", float64(memStats.Mallocs)))
	gauges = append(gauges, r.createGauge("NextGC", float64(memStats.NextGC)))
	gauges = append(gauges, r.createGauge("NumForcedGC", float64(memStats.NumForcedGC)))
	gauges = append(gauges, r.createGauge("NumGC", float64(memStats.NumGC)))
	gauges = append(gauges, r.createGauge("OtherSys", float64(memStats.OtherSys)))
	gauges = append(gauges, r.createGauge("PauseTotalNs", float64(memStats.PauseTotalNs)))
	gauges = append(gauges, r.createGauge("StackInuse", float64(memStats.StackInuse)))
	gauges = append(gauges, r.createGauge("StackSys", float64(memStats.StackSys)))
	gauges = append(gauges, r.createGauge("Sys", float64(memStats.Sys)))
	gauges = append(gauges, r.createGauge("TotalAlloc", float64(memStats.TotalAlloc)))

	// Обновляем RandomValue
	randomValue, _ := domain.NewGaugeValue(rand.Float64())
	r.randomValue.Update(randomValue)
	gauges = append(gauges, r.randomValue)

	return gauges, nil
}

// CollectCounters собирает все counter метрики
func (r *RuntimeCollector) CollectCounters() ([]*domain.Counter, error) {
	return []*domain.Counter{r.pollCount}, nil
}

// IncrementPollCount увеличивает счётчик опросов
func (r *RuntimeCollector) IncrementPollCount() {
	delta, _ := domain.NewCounterValue(1)
	r.pollCount.Increment(delta)
}

// createGauge вспомогательный метод для создания gauge метрики
func (r *RuntimeCollector) createGauge(name string, value float64) *domain.Gauge {
	metricName, _ := domain.NewMetricName(name)
	gaugeValue, _ := domain.NewGaugeValue(value)
	return domain.NewGauge(metricName, gaugeValue)
}
