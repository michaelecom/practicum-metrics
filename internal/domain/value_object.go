package domain

import (
	"errors"
	"strings"
)

// MetricName - Value Object для имени метрики
type MetricName struct {
	value string
}

// NewMetricName создаёт имя метрики
func NewMetricName(name string) (MetricName, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return MetricName{}, errors.New("metric name cannot be empty")
	}

	if len(name) > 255 {
		return MetricName{}, errors.New("metric name too long (max 255 characters)")
	}

	return MetricName{value: name}, nil
}

// String возвращает строковое значение
func (m MetricName) String() string {
	return m.value
}

// Equals проверяет равенство двух MetricName
func (m MetricName) Equals(other MetricName) bool {
	return m.value == other.value
}

// GaugeValue - Value Object для значения gauge метрики
type GaugeValue struct {
	value float64
}

// NewGaugeValue создаёт значение gauge
func NewGaugeValue(value float64) (GaugeValue, error) {
	return GaugeValue{value: value}, nil
}

// Float64 возвращает числовое значение
func (g GaugeValue) Float64() float64 {
	return g.value
}

// CounterValue - Value Object для значения counter метрики
type CounterValue struct {
	value int64
}

// NewCounterValue создаёт значение counter
func NewCounterValue(value int64) (CounterValue, error) {
	if value < 0 {
		return CounterValue{}, errors.New("counter value cannot be negative")
	}

	return CounterValue{value: value}, nil
}

// Int64 возвращает числовое значение
func (c CounterValue) Int64() int64 {
	return c.value
}

// Add складывает два значения counter
func (c CounterValue) Add(other CounterValue) CounterValue {
	return CounterValue{value: c.value + other.value}
}
