// Package ratelimit 实现 token 级限流。
//
// ⚠️ 核心设计决策（见 docs/decisions.md ADR-001）：
// 本包采用「**先放行、后付费**」—— 接口本身就编码了这个决策：
//   - Check 只读判断，**不扣减**；
//   - Commit 在拿到上游真实 usage 后才扣减。
//
// 因此没有 Reserve/Refund（预扣回补）方法。这不是遗漏，是选择。
package ratelimit

import (
	"context"
	"time"
)

// Decision 是一次限流判断的结果。
type Decision struct {
	Allowed    bool
	Remaining  int
	ResetAt    time.Time
	RetryAfter time.Duration
}

// Limiter 是 token 级限流器。
type Limiter interface {
	// Check 只读判断当前 key 是否还有配额。**不产生任何扣减**。
	Check(ctx context.Context, key string) (Decision, error)

	// Commit 在拿到上游真实 usage 后扣减配额（log 阶段）。
	Commit(ctx context.Context, key string, actualTokens int) error
}
