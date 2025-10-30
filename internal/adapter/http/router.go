package http

import (
	"net/http"
)

// SetupRouter настраивает маршруты для HTTP сервера с использованием стандартного net/http
func SetupRouter(handler *MetricHandler) http.Handler {
	// Создание роутера
	mux := http.NewServeMux()

	// Регистрация маршрутов
	mux.HandleFunc("/update/", handler.UpdateMetric)

	mux.HandleFunc("/value/", handler.GetMetricValue)
	mux.HandleFunc("/", handler.ListAllMetrics)

	return mux
}
