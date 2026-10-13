package server

import (
	"net/http"
)

// handleChatCompletions 是 /v1/chat/completions 的入口。
//
// 当前为**起盘骨架**：解析与校验尚未实现，直接返回 501。
// W43 填充：协议解码 → IR → 上游编码 → SSE 透传 → 回程转 IR → 客户端协议编码。
func (s *Server) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": map[string]any{"type": "invalid_request_error", "message": "only POST is supported"},
		})
		return
	}

	// TODO(W43): 按 Accept / 路径判定客户端协议 → Codec.DecodeRequest
	// TODO(W43): router.Pick 选上游 → Codec.EncodeRequest
	// TODO(W43): upstream.Stream → DecodeStreamEvent → EncodeStreamEvent → flush
	// TODO(W44): limiter.Check（放行）→ 结束后 limiter.Commit（按真实 usage 扣减）
	writeJSON(w, http.StatusNotImplemented, map[string]any{
		"error": map[string]any{
			"type":    "not_implemented",
			"message": "protocol adaptation is scheduled for W43; see docs/decisions.md",
		},
	})
}
