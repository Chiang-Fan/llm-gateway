// Package server 承载 HTTP 层。
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Chiang-Fan/mini-llm-gateway/internal/config"
	"github.com/Chiang-Fan/mini-llm-gateway/internal/ratelimit"
	"github.com/Chiang-Fan/mini-llm-gateway/internal/router"
)

// Server 是网关的 HTTP 服务。
type Server struct {
	cfg        *config.Config
	httpserver *http.Server
	router     router.Router
	limiter    ratelimit.Limiter
}

// New 创建一个 Server。
func New(cfg *config.Config) *Server {
	s := &Server{
		cfg:     cfg,
		router:  router.NewWeighted(),
		limiter: ratelimit.NewInMemory(cfg.RateLimit),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/v1/chat/completions", s.handleChatCompletions)
	s.httpserver = &http.Server{
		Addr:              cfg.Listen,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		// ⚠️ 不设 WriteTimeout：流式响应可能长时间持续。
		// 超时控制交给上游请求的 context，而不是整条连接。
	}
	return s
}

// Start 启动服务（阻塞）。
func (s *Server) Start() error { return s.httpserver.ListenAndServe() }

// Shutdown 优雅关闭。
func (s *Server) Shutdown(ctx context.Context) error { return s.httpserver.Shutdown(ctx) }

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
