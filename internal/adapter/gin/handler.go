package gin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/michaelecom/practicum-metrics/internal/usecase"
)

// MetricHandler - Gin адаптер для обработки метрик
type MetricHandler struct {
	metricUseCase usecase.MetricUseCase
}

// NewMetricHandler создаёт новый Gin обработчик для метрик
func NewMetricHandler(metricUseCase usecase.MetricUseCase) *MetricHandler {
	return &MetricHandler{
		metricUseCase: metricUseCase,
	}
}

// UpdateMetric обрабатывает POST /update/:type/:name/:value
func (h *MetricHandler) UpdateMetric(c *gin.Context) {
	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "text/plain") {
		c.String(http.StatusBadRequest, "Invalid content type")
		return
	}

	// Получение параметров из URL
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	// Валидация имени метрики
	if metricName == "" {
		c.String(http.StatusNotFound, "Metric name is required")
		return
	}

	// Валидация значения метрики
	if metricValue == "" {
		c.String(http.StatusBadRequest, "Metric value is required")
		return
	}

	// Обработка в зависимости от типа метрики
	var err error
	switch metricType {
	case "gauge":
		err = h.handleGaugeUpdate(metricName, metricValue)

	case "counter":
		err = h.handleCounterUpdate(metricName, metricValue)

	default:
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid metric type: %s", metricType))
		return
	}

	// Обработка ошибок use case
	if err != nil {
		h.handleUseCaseError(c, err)
		return
	}

	// Успешный ответ
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Status(http.StatusOK)
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

// GetMetricValue обрабатывает GET /value/:type/:name
func (h *MetricHandler) GetMetricValue(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	// Валидация имени метрики
	if metricName == "" {
		c.String(http.StatusNotFound, "Metric name is required")
		return
	}

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
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid metric type: %s", metricType))
		return
	}

	// Обработка ошибок
	if err != nil {
		h.handleUseCaseError(c, err)
		return
	}

	// Успешный ответ с значением метрики
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, value)
}

// ListAllMetrics обрабатывает GET / и возвращает HTML-страницу со всеми метриками
func (h *MetricHandler) ListAllMetrics(c *gin.Context) {
	metrics, err := h.metricUseCase.GetAllMetrics()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get metrics: %v", err))
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
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to render template: %v", err))
		return
	}
}

// handleUseCaseError обрабатывает ошибки из use case для Gin и возвращает соответствующий HTTP статус
func (h *MetricHandler) handleUseCaseError(c *gin.Context, err error) {
	errMsg := err.Error()

	// Метрика не найдена
	if strings.Contains(errMsg, "not found") {
		c.String(http.StatusNotFound, errMsg)
		return
	}

	// Ошибки валидации
	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "cannot be negative") {
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	// По умолчанию - внутренняя ошибка сервера
	c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to process metric: %v", err))
}
