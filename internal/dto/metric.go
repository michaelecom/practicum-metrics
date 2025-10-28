package dto

// Metrics - DTO для передачи метрик через API (JSON)
type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	Type  string   `json:"type"`            // Тип метрики (gauge/counter)
	Delta *int64   `json:"delta,omitempty"` // Значение counter
	Value *float64 `json:"value,omitempty"` // Значение gauge
	Hash  string   `json:"hash,omitempty"`  // Хеш для проверки целостности
}
