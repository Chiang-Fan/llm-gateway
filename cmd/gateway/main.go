// Command gateway 是 llm-gateway 的入口。
//
// 最小可运行版本：只起 HTTP server 并注册 /healthz 与 /v1/chat/completions 骨架。
// 协议适配、流式透传、限流与路由在 W43/W44 逐步填充。
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chiang-Fan/llm-gateway/internal/config"
	"github.com/Chiang-Fan/llm-gateway/internal/server"
)

func main() {
	cfgPath := flag.String("config", "", "配置文件路径（可选，缺省用默认配置）")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		slog.Error("load config failed", "err", err)
		os.Exit(1)
	}

	srv := server.New(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start() }()

	slog.Info("llm-gateway started", "addr", cfg.Listen)

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			slog.Error("server exited", "err", err)
			os.Exit(1)
		}
	}

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
}
