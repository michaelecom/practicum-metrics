package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/michaelecom/practicum-metrics/internal/usecase"
)

// MetricHandler - HTTP адаптер для обработки метрик
type MetricHandler struct {
	metricUseCase usecase.MetricUseCase
}

// NewMetricHandler создаёт новый HTTP обработчик для метрик
func NewMetricHandler(metricUseCase usecase.MetricUseCase) *MetricHandler {
	return &MetricHandler{
		metricUseCase: metricUseCase,
	}
}

// UpdateMetric обрабатывает POST /update/<type>/<name>/<value>
func (h *MetricHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверка Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "text/plain") {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	// Парсинг URL: /update/<type>/<name>/<value>
	path := strings.TrimPrefix(r.URL.Path, "/update/")
	parts := strings.Split(path, "/")

	// Валидация структуры URL
	if len(parts) < 2 || (len(parts) >= 2 && parts[1] == "") {
		http.Error(w, "Metric name is required", http.StatusNotFound)
		return
	}

	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Metric value is required", http.StatusBadRequest)
		return
	}

	metricType := parts[0]
	metricName := parts[1]
	metricValue := parts[2]

	// Обработка в зависимости от типа метрики
	var err error
	switch metricType {
	case "gauge":
		err = h.handleGaugeUpdate(metricName, metricValue)

	case "counter":
		err = h.handleCounterUpdate(metricName, metricValue)

	default:
		http.Error(w, fmt.Sprintf("Invalid metric type: %s", metricType), http.StatusBadRequest)
		return
	}

	// Обработка ошибок use case
	if err != nil {
		h.handleUseCaseError(w, err)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// handleGaugeUpdate обрабатывает обновление gauge метрики
func (h *MetricHandler) handleGaugeUpdate(name, valueStr string) error {
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("invalid gauge value: %w", err)
	}

	return h.metricUseCase.UpdateGauge(name, value)
}

// handleCounterUpdate обрабатывает обновление counter метрики
func (h *MetricHandler) handleCounterUpdate(name, valueStr string) error {
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid counter value: %w", err)
	}

	return h.metricUseCase.UpdateCounter(name, value)
}

// GetMetricValue обрабатывает GET /value/<type>/<name>
func (h *MetricHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг URL: /value/<type>/<name>
	path := strings.TrimPrefix(r.URL.Path, "/value/")
	parts := strings.Split(path, "/")

	// Валидация структуры URL
	if len(parts) < 2 || parts[1] == "" {
		http.Error(w, "Metric name is required", http.StatusNotFound)
		return
	}

	metricType := parts[0]
	metricName := parts[1]

	var value string
	var err error

	// Получение значения в зависимости от типа метрики
	switch metricType {
	case "gauge":
		var gaugeValue float64
		gaugeValue, err = h.metricUseCase.GetGauge(metricName)
		if err == nil {
			value = strconv.FormatFloat(gaugeValue, 'f', -1, 64)
		}

	case "counter":
		var counterValue int64
		counterValue, err = h.metricUseCase.GetCounter(metricName)
		if err == nil {
			value = strconv.FormatInt(counterValue, 10)
		}

	default:
		http.Error(w, fmt.Sprintf("Invalid metric type: %s", metricType), http.StatusBadRequest)
		return
	}

	// Обработка ошибок
	if err != nil {
		h.handleUseCaseError(w, err)
		return
	}

	// Успешный ответ с значением метрики
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, value)
}

// ListAllMetrics обрабатывает GET / и возвращает HTML-страницу со всеми метриками
func (h *MetricHandler) ListAllMetrics(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем только корневой путь "/"
	// Для всех остальных путей возвращаем 404 (Not Found)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Проверка метода запроса для корневого пути
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics, err := h.metricUseCase.GetAllMetrics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get metrics: %v", err), http.StatusInternalServerError)
		return
	}

	// Подготовка данных для шаблона
	data := struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	// Извлекаем данные из metrics
	if gauges, ok := metrics["gauges"].(map[string]float64); ok {
		data.Gauges = gauges
	}
	if counters, ok := metrics["counters"].(map[string]int64); ok {
		data.Counters = counters
	}

	// Рендеринг HTML через шаблон
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleUseCaseError обрабатывает ошибки из use case для net/http и возвращает соответствующий HTTP статус
func (h *MetricHandler) handleUseCaseError(w http.ResponseWriter, err error) {
	errMsg := err.Error()

	// Метрика не найдена
	if strings.Contains(errMsg, "not found") {
		http.Error(w, errMsg, http.StatusNotFound)
		return
	}

	// Ошибки валидации
	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "cannot be negative") {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// По умолчанию - внутренняя ошибка сервера
	http.Error(w, fmt.Sprintf("Failed to process metric: %v", err), http.StatusInternalServerError)
}
