package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michaelecom/practicum-metrics/internal/domain"
)

func TestHTTPClient_SendGauge(t *testing.T) {
	// Создаём тестовый сервер
	var receivedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	// Создаём клиент
	client := NewHTTPClient(server.URL)

	// Создаём gauge метрику
	name, _ := domain.NewMetricName("TestGauge")
	value, _ := domain.NewGaugeValue(123.45)
	gauge := domain.NewGauge(name, value)

	// Отправляем метрику
	err := client.SendGauge(gauge)
	if err != nil {
		t.Fatalf("SendGauge failed: %v", err)
	}

	// Проверяем URL
	expectedURL := "/update/gauge/TestGauge/123.45"
	if receivedURL != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, receivedURL)
	}
}

func TestHTTPClient_SendCounter(t *testing.T) {
	// Создаём тестовый сервер
	var receivedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	// Создаём клиент
	client := NewHTTPClient(server.URL)

	// Создаём counter метрику
	name, _ := domain.NewMetricName("TestCounter")
	value, _ := domain.NewCounterValue(100)
	counter := domain.NewCounter(name, value)

	// Отправляем метрику
	err := client.SendCounter(counter)
	if err != nil {
		t.Fatalf("SendCounter failed: %v", err)
	}

	// Проверяем URL
	expectedURL := "/update/counter/TestCounter/100"
	if receivedURL != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, receivedURL)
	}
}

func TestHTTPClient_SendGauge_ServerError(t *testing.T) {
	// Создаём тестовый сервер, который возвращает ошибку
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	// Создаём клиент
	client := NewHTTPClient(server.URL)

	// Создаём gauge метрику
	name, _ := domain.NewMetricName("TestGauge")
	value, _ := domain.NewGaugeValue(123.45)
	gauge := domain.NewGauge(name, value)

	// Отправляем метрику
	err := client.SendGauge(gauge)
	if err == nil {
		t.Error("Expected error for server error response, got nil")
	}
}

func TestHTTPClient_SendCounter_ServerError(t *testing.T) {
	// Создаём тестовый сервер, который возвращает ошибку
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	defer server.Close()

	// Создаём клиент
	client := NewHTTPClient(server.URL)

	// Создаём counter метрику
	name, _ := domain.NewMetricName("TestCounter")
	value, _ := domain.NewCounterValue(100)
	counter := domain.NewCounter(name, value)

	// Отправляем метрику
	err := client.SendCounter(counter)
	if err == nil {
		t.Error("Expected error for server error response, got nil")
	}
}
