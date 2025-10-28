package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/michaelecom/practicum-metrics/internal/domain"
)

// MetricSender отправляет метрики на сервер
type MetricSender interface {
	// SendGauge отправляет gauge метрику
	SendGauge(gauge *domain.Gauge) error
	// SendCounter отправляет counter метрику
	SendCounter(counter *domain.Counter) error
}

// HTTPClient - HTTP клиент для отправки метрик
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPClient создаёт новый HTTP клиент
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// SendGauge отправляет gauge метрику на сервер
func (c *HTTPClient) SendGauge(gauge *domain.Gauge) error {
	url := fmt.Sprintf("%s/update/gauge/%s/%v",
		c.baseURL,
		gauge.Name().String(),
		gauge.Value().Float64(),
	)

	return c.sendRequest(url)
}

// SendCounter отправляет counter метрику на сервер
func (c *HTTPClient) SendCounter(counter *domain.Counter) error {
	url := fmt.Sprintf("%s/update/counter/%s/%d",
		c.baseURL,
		counter.Name().String(),
		counter.Value().Int64(),
	)

	return c.sendRequest(url)
}

// sendRequest выполняет HTTP POST запрос
func (c *HTTPClient) sendRequest(url string) error {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
