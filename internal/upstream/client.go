// Package upstream 封装对上游模型服务的 HTTP 调用（含流式）。
package upstream

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Chiang-Fan/llm-gateway/internal/config"
)

// ErrUpstream 表示上游返回了非 2xx。
type ErrUpstream struct {
	Status int
	Body   string
}

func (e *ErrUpstream) Error() string {
	return fmt.Sprintf("upstream returned %d: %s", e.Status, e.Body)
}

// Client 是上游客户端。
type Client struct {
	hc *http.Client
}

// New 创建一个上游客户端。timeout 为 0 表示不设总超时（流式必须如此）。
func New(timeout time.Duration) *Client {
	return &Client{hc: &http.Client{Timeout: timeout}}
}

// Do 发起一次非流式请求。
func (c *Client) Do(ctx context.Context, up config.Upstream, path string, body []byte) ([]byte, error) {
	req, err := c.build(ctx, up, path, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &ErrUpstream{Status: resp.StatusCode, Body: string(b)}
	}
	return b, nil
}

// Stream 发起一次流式请求，逐行回调。
//
// ⚠️ 这里**不缓冲**：读到一行就交给 onLine，由上层 flush 给客户端。
// 缓冲会让首包时延变差，也会破坏 SSE 的实时性。
func (c *Client) Stream(ctx context.Context, up config.Upstream, path string, body []byte, onLine func([]byte) error) error {
	req, err := c.build(ctx, up, path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return &ErrUpstream{Status: resp.StatusCode, Body: string(b)}
	}

	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		if err := onLine(bytes.TrimRight(sc.Bytes(), "\r\n")); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (c *Client) build(ctx context.Context, up config.Upstream, path string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, up.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if up.APIKey != "" {
		switch up.Protocol {
		case "anthropic":
			req.Header.Set("x-api-key", up.APIKey)
			req.Header.Set("anthropic-version", "2023-06-01")
		default:
			req.Header.Set("Authorization", "Bearer "+up.APIKey)
		}
	}
	return req, nil
}
