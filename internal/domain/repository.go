package domain

// MetricRepository - интерфейс для работы с хранилищем метрик
type MetricRepository interface {
	// SaveGauge сохраняет gauge метрику
	SaveGauge(metric *Gauge) error

	// SaveCounter сохраняет counter метрику
	SaveCounter(metric *Counter) error

	// GetGauge возвращает gauge метрику по имени
	GetGauge(name MetricName) (*Gauge, error)

	// GetCounter возвращает counter метрику по имени
	GetCounter(name MetricName) (*Counter, error)

	// GetOrCreateCounter возвращает существующую counter метрику или создаёт новую
	GetOrCreateCounter(name MetricName) (*Counter, error)

	// GetAllGauges возвращает все gauge метрики
	GetAllGauges() ([]*Gauge, error)

	// GetAllCounters возвращает все counter метрики
	GetAllCounters() ([]*Counter, error)
}
