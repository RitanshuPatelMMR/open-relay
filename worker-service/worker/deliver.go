package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DeliveryResult struct {
	StatusCode int
	Body       string
	DurationMs int
	Err        error
}

func Deliver(ctx context.Context, destinationURL string, method string, headers []byte, payload []byte, timeoutSeconds int) DeliveryResult {
	start := time.Now()

	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, method, destinationURL, bytes.NewReader(payload))
	if err != nil {
		return DeliveryResult{Err: fmt.Errorf("build request failed: %w", err)}
	}

	// restore original headers
	var originalHeaders map[string][]string
	if err := json.Unmarshal(headers, &originalHeaders); err == nil {
		for key, vals := range originalHeaders {
			for _, v := range vals {
				req.Header.Set(key, v)
			}
		}
	}

	// always set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OpenRelay-Delivery", "true")

	resp, err := client.Do(req)
	durationMs := int(time.Since(start).Milliseconds())

	if err != nil {
		return DeliveryResult{
			DurationMs: durationMs,
			Err:        fmt.Errorf("http request failed: %w", err),
		}
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

	return DeliveryResult{
		StatusCode: resp.StatusCode,
		Body:       string(respBody),
		DurationMs: durationMs,
	}
}