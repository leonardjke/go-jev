package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/leonardjke/go-jev/internal/privacy"
)

type Logger interface {
	Debug(msg string, args ...any)
}

type StatusError struct {
	StatusCode int
	Body       string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("jev: api returned %d: %s", e.StatusCode, e.Body)
}

type Client struct {
	HTTP     *http.Client
	Logger   Logger
	Endpoint string
	APIKey   privacy.SensitiveString

	MaxResponseBytes int64
}

func (c *Client) Post(ctx context.Context, req, resp any) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("jev: encode request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("jev: build request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+string(c.APIKey))

	httpResp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return fmt.Errorf("jev: send request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(httpResp.Body, c.MaxResponseBytes))
	if err != nil {
		return fmt.Errorf("jev: read response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return &StatusError{StatusCode: httpResp.StatusCode, Body: string(respBody)}
	}

	c.Logger.Debug("Client", "response", string(respBody))

	if err := json.Unmarshal(respBody, resp); err != nil {
		return fmt.Errorf("jev: decode response: %w", err)
	}

	return nil
}
