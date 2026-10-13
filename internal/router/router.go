// Package router 负责从多个上游中挑选一个。
//
// W42 起盘只定义接口与加权随机；自适应策略（延迟/前缀亲和）**明确不做**，
// 理由见 docs/decisions.md ADR-002。
package router

import (
	"context"
	"errors"

	"github.com/Chiang-Fan/mini-llm-gateway/internal/config"
	"github.com/Chiang-Fan/mini-llm-gateway/internal/protocol"
)

// ErrNoUpstream 表示没有可用上游。
var ErrNoUpstream = errors.New("no upstream available")

// Router 挑选一个上游。
type Router interface {
	Pick(ctx context.Context, req *protocol.Request, candidates []config.Upstream) (config.Upstream, error)
}
