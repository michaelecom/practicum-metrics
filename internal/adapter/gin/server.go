package gin

import (
	"github.com/gin-gonic/gin"

	"github.com/michaelecom/practicum-metrics/internal/adapter"
	"github.com/michaelecom/practicum-metrics/internal/usecase"
)

// GinServer - реализация HTTP сервера с использованием Gin фреймворка
type GinServer struct {
	engine *gin.Engine
}

// NewServer создаёт новый Gin сервер
func NewServer(metricUseCase usecase.MetricUseCase) adapter.HTTPServer {
	// Создаём handler
	metricHandler := NewMetricHandler(metricUseCase)

	// Настраиваем роутер
	router := SetupRouter(metricHandler)

	return &GinServer{
		engine: router,
	}
}

// Start запускает Gin сервер на указанном адресе
func (s *GinServer) Start(addr string) error {
	return s.engine.Run(addr)
}
