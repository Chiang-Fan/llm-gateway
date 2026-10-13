package protocol

import (
	"context"
	"errors"
)

// Anthropic 实现 Anthropic Messages 协议的编解码。
//
// 与 OpenAI 的主要差异（W43 逐个实现）：
//   - 系统提示是**顶层 system 字段**，不是 role=system 的消息；
//   - 结束原因是 **stop_reason**（end_turn / max_tokens / tool_use / stop_sequence）；
//   - 思考预算是 **thinking.budget_tokens**，对应 OpenAI 的 reasoning_effort；
//   - 流式是**多种事件类型**（message_start / content_block_delta / message_delta ...），
//     需要累积状态才能压平成语义等价的一维事件流；
//   - 工具调用是 **tool_use / tool_result 内容块**，不是 tool_calls 字段。
type Anthropic struct{}

// Name 返回协议名。
func (Anthropic) Name() string { return "anthropic" }

func (Anthropic) DecodeRequest([]byte) (*Request, error) {
	return nil, errors.New("Anthropic.DecodeRequest: not implemented (W43)")
}

func (Anthropic) EncodeRequest(*Request) ([]byte, string, error) {
	// 目标路径固定为 /v1/messages。
	return nil, "", errors.New("Anthropic.EncodeRequest: not implemented (W43)")
}

func (Anthropic) DecodeResponse([]byte) (*Response, error) {
	return nil, errors.New("Anthropic.DecodeResponse: not implemented (W43)")
}

func (Anthropic) EncodeResponse(*Response) ([]byte, error) {
	return nil, errors.New("Anthropic.EncodeResponse: not implemented (W43)")
}

func (Anthropic) DecodeStreamEvent(context.Context, []byte) ([]StreamEvent, error) {
	return nil, errors.New("Anthropic.DecodeStreamEvent: not implemented (W43)")
}

func (Anthropic) EncodeStreamEvent(*StreamEvent) ([][]byte, error) {
	return nil, errors.New("Anthropic.EncodeStreamEvent: not implemented (W43)")
}

// StopReasonToIR 把 Anthropic 的 stop_reason 归一成 IR 的 FinishReason。
// 这是双向转换的字典之一，实现在 W43 与 IR 一起收口。
func StopReasonToIR(s string) FinishReason {
	switch s {
	case "end_turn", "stop_sequence":
		return FinishStop
	case "max_tokens":
		return FinishLength
	case "tool_use":
		return FinishToolUse
	default:
		return FinishStop
	}
}
