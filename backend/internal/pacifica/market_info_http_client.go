package pacifica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type HTTPMarketInfoClient struct {
	baseURL    string
	httpClient *http.Client
	cacheTTL   time.Duration

	mu          sync.RWMutex
	cachedAt    time.Time
	cachedItems []MarketInfo
}

type marketInfoResponse struct {
	Success bool         `json:"success"`
	Data    []MarketInfo `json:"data"`
	Error   any          `json:"error"`
	Code    any          `json:"code"`
}

func NewHTTPMarketInfoClient(baseURL string, httpClient *http.Client, cacheTTL time.Duration) *HTTPMarketInfoClient {
	return &HTTPMarketInfoClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
		cacheTTL:   cacheTTL,
	}
}

func (c *HTTPMarketInfoClient) ListMarketInfo() ([]MarketInfo, error) {
	if cached, ok := c.getCached(); ok {
		return cached, nil
	}

	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/v1/info", nil)
	if err != nil {
		return nil, fmt.Errorf("build market info request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch market info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch market info: unexpected status %d", resp.StatusCode)
	}

	var payload marketInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode market info response: %w", err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("fetch market info: pacifica returned success=false")
	}

	c.setCache(payload.Data)
	return payload.Data, nil
}

func (c *HTTPMarketInfoClient) getCached() ([]MarketInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cacheTTL <= 0 || c.cachedAt.IsZero() || time.Since(c.cachedAt) > c.cacheTTL {
		return nil, false
	}

	return append([]MarketInfo(nil), c.cachedItems...), true
}

func (c *HTTPMarketInfoClient) setCache(items []MarketInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cachedAt = time.Now().UTC()
	c.cachedItems = append([]MarketInfo(nil), items...)
}
