package collector

import (
	"testing"
)

func TestRuntimeCollector_CollectGauges(t *testing.T) {
	collector := NewRuntimeCollector()

	gauges, err := collector.CollectGauges()
	if err != nil {
		t.Fatalf("CollectGauges failed: %v", err)
	}

	// Проверяем, что собрано 28 gauge метрик (27 из runtime + RandomValue)
	expectedCount := 28
	if len(gauges) != expectedCount {
		t.Errorf("Expected %d gauges, got %d", expectedCount, len(gauges))
	}

	// Проверяем, что все gauge метрики имеют имена
	for _, gauge := range gauges {
		if gauge.Name().String() == "" {
			t.Error("Gauge has empty name")
		}
	}
}

func TestRuntimeCollector_CollectCounters(t *testing.T) {
	collector := NewRuntimeCollector()

	counters, err := collector.CollectCounters()
	if err != nil {
		t.Fatalf("CollectCounters failed: %v", err)
	}

	// Проверяем, что собрана 1 counter метрика (PollCount)
	expectedCount := 1
	if len(counters) != expectedCount {
		t.Errorf("Expected %d counter, got %d", expectedCount, len(counters))
	}

	// Проверяем, что это PollCount
	if counters[0].Name().String() != "PollCount" {
		t.Errorf("Expected counter name PollCount, got %s", counters[0].Name().String())
	}

	// Проверяем начальное значение
	if counters[0].Value().Int64() != 0 {
		t.Errorf("Expected initial PollCount to be 0, got %d", counters[0].Value().Int64())
	}
}

func TestRuntimeCollector_IncrementPollCount(t *testing.T) {
	collector := NewRuntimeCollector()

	// Увеличиваем счётчик
	collector.IncrementPollCount()
	collector.IncrementPollCount()
	collector.IncrementPollCount()

	counters, err := collector.CollectCounters()
	if err != nil {
		t.Fatalf("CollectCounters failed: %v", err)
	}

	// Проверяем, что счётчик увеличился на 3
	expected := int64(3)
	if counters[0].Value().Int64() != expected {
		t.Errorf("Expected PollCount to be %d, got %d", expected, counters[0].Value().Int64())
	}
}

func TestRuntimeCollector_RandomValue(t *testing.T) {
	collector := NewRuntimeCollector()

	// Собираем метрики дважды
	gauges1, _ := collector.CollectGauges()
	gauges2, _ := collector.CollectGauges()

	// Находим RandomValue в обоих наборах
	var random1, random2 float64
	for _, g := range gauges1 {
		if g.Name().String() == "RandomValue" {
			random1 = g.Value().Float64()
		}
	}

	for _, g := range gauges2 {
		if g.Name().String() == "RandomValue" {
			random2 = g.Value().Float64()
		}
	}

	// Проверяем, что значения находятся в допустимом диапазоне [0, 1)
	if random1 < 0 || random1 >= 1 {
		t.Errorf("RandomValue out of range: %f", random1)
	}

	if random2 < 0 || random2 >= 1 {
		t.Errorf("RandomValue out of range: %f", random2)
	}
}
