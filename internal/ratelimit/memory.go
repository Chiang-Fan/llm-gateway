package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/Chiang-Fan/mini-llm-gateway/internal/config"
)

// window 是固定窗口计数器。
//
// ⚠️ 已知缺陷（有意保留，见 docs/decisions.md ADR-003）：
// 固定窗口在**窗口边界**会出现双倍突发 —— 窗口末尾用满配额后，
// 下一个窗口立刻又能用满，两个窗口交界处的瞬时速率可达配额的两倍。
// 换成滑动窗口可以缓解，但实现与存储成本上升；最小版先接受这个缺陷。
type window struct {
	mu     sync.Mutex
	cfg    config.RateLimitConfig
	counts map[string]*bucket
}

type bucket struct {
	used      int
	windowEnd time.Time
}

// NewInMemory 创建一个进程内的固定窗口限流器。
func NewInMemory(cfg config.RateLimitConfig) Limiter {
	return &window{cfg: cfg, counts: make(map[string]*bucket)}
}

// Check 只读判断，**不扣减** —— 对应「先放行、后付费」。
func (w *window) Check(_ context.Context, key string) (Decision, error) {
	if !w.cfg.Enabled {
		return Decision{Allowed: true, Remaining: w.cfg.Quota}, nil
	}
	now := time.Now()
	w.mu.Lock()
	defer w.mu.Unlock()

	b := w.roll(key, now)
	remaining := w.cfg.Quota - b.used
	return Decision{
		Allowed:   remaining > 0,
		Remaining: max(remaining, 0),
		ResetAt:   b.windowEnd,
	}, nil
}

// Commit 在上游返回后按真实用量扣减。
func (w *window) Commit(_ context.Context, key string, actualTokens int) error {
	if !w.cfg.Enabled || actualTokens <= 0 {
		return nil
	}
	now := time.Now()
	w.mu.Lock()
	defer w.mu.Unlock()

	w.roll(key, now).used += actualTokens
	return nil
}

// roll 返回当前窗口的桶，必要时滚动到新窗口。
func (w *window) roll(key string, now time.Time) *bucket {
	b, ok := w.counts[key]
	if !ok || !now.Before(b.windowEnd) {
		end := now.Truncate(w.cfg.Window).Add(w.cfg.Window)
		b = &bucket{windowEnd: end}
		w.counts[key] = b
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
