package http

import (
	"net/http"

	"github.com/michaelecom/practicum-metrics/internal/adapter"
	"github.com/michaelecom/practicum-metrics/internal/usecase"
)

// HTTPServer - реализация HTTP сервера с использованием стандартного net/http
type HTTPServer struct {
	handler http.Handler
}

// NewServer создаёт новый HTTP сервер
func NewServer(metricUseCase usecase.MetricUseCase) adapter.HTTPServer {
	// Создаём handler
	metricHandler := NewMetricHandler(metricUseCase)

	// Настраиваем роутер
	router := SetupRouter(metricHandler)

	return &HTTPServer{
		handler: router,
	}
}

// Start запускает HTTP сервер на указанном адресе
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}
