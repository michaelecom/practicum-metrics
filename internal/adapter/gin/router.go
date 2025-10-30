package gin

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршруты для HTTP сервера с использованием Gin фреймворка
func SetupRouter(handler *MetricHandler) *gin.Engine {
	// Создание роутера
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Добавляем middleware для восстановления после паники
	router.Use(gin.Recovery())

	// Регистрация маршрутов
	router.POST("/update/:type/:name/:value", handler.UpdateMetric)

	router.GET("/value/:type/:name", handler.GetMetricValue)
	router.GET("/", handler.ListAllMetrics)

	return router
}
