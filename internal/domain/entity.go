package domain

// MetricType представляет тип метрики
type MetricType string

const (
	// MetricTypeGauge - тип метрики gauge (замещает предыдущее значение)
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeCounter - тип метрики counter (добавляется к предыдущему значению)
	MetricTypeCounter MetricType = "counter"
)

// Gauge - доменная сущность для gauge метрики
type Gauge struct {
	name  MetricName
	value GaugeValue
}

// NewGauge создаёт новую gauge метрику
func NewGauge(name MetricName, value GaugeValue) *Gauge {
	return &Gauge{
		name:  name,
		value: value,
	}
}

// Name возвращает имя метрики
func (g *Gauge) Name() MetricName {
	return g.name
}

// Value возвращает значение метрики
func (g *Gauge) Value() GaugeValue {
	return g.value
}

// Update обновляет значение gauge (замещает предыдущее)
func (g *Gauge) Update(newValue GaugeValue) {
	g.value = newValue
}

// Counter - доменная сущность для counter метрики
type Counter struct {
	name  MetricName
	value CounterValue
}

// NewCounter создаёт новую counter метрику
func NewCounter(name MetricName, value CounterValue) *Counter {
	return &Counter{
		name:  name,
		value: value,
	}
}

// Name возвращает имя метрики
func (c *Counter) Name() MetricName {
	return c.name
}

// Value возвращает значение метрики
func (c *Counter) Value() CounterValue {
	return c.value
}

// Increment увеличивает значение counter
func (c *Counter) Increment(delta CounterValue) {
	c.value = c.value.Add(delta)
}
