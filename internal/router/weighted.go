package router

import (
	"context"
	"math/rand"

	"github.com/Chiang-Fan/llm-gateway/internal/config"
	"github.com/Chiang-Fan/llm-gateway/internal/protocol"
)

// weighted 按权重随机挑选上游。
//
// ⚠️ 为什么不直接上自适应（延迟/汤普森采样）—— 见 docs/decisions.md ADR-002：
// 样本量不足时，自适应策略的期望收益为负。
type weighted struct {
	rnd *rand.Rand
}

// NewWeighted 创建一个加权随机路由。
func NewWeighted() Router { return &weighted{rnd: rand.New(rand.NewSource(rand.Int63()))} }

func (r *weighted) Pick(_ context.Context, _ *protocol.Request, cands []config.Upstream) (config.Upstream, error) {
	if len(cands) == 0 {
		return config.Upstream{}, ErrNoUpstream
	}
	total := 0
	for _, u := range cands {
		if u.Weight > 0 {
			total += u.Weight
		}
	}
	if total == 0 {
		// 全部未配权重：退化为均匀随机。
		return cands[r.rnd.Intn(len(cands))], nil
	}
	n := r.rnd.Intn(total)
	for _, u := range cands {
		if u.Weight <= 0 {
			continue
		}
		n -= u.Weight
		if n < 0 {
			return u, nil
		}
	}
	return cands[len(cands)-1], nil
}
