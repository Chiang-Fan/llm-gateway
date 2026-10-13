// Package protocol 定义厂商无关的中间表示（IR）与各协议编解码器。
//
// 设计原则：**先归一成 IR，再按目标协议序列化** —— 而不是写 N×N 的两两互转。
// 新增一个厂商 = 新增一个 Codec，而不是新增 N 个转换器。
package protocol

import "encoding/json"

// Role 是归一后的消息角色。
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// ContentType 是内容块类型。
type ContentType string

const (
	ContentText       ContentType = "text"
	ContentImage      ContentType = "image"
	ContentToolUse    ContentType = "tool_use"
	ContentToolResult ContentType = "tool_result"
)

// ToolUse 表示模型发起的一次工具调用。
type ToolUse struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input,omitempty"`
}

// ToolResult 表示工具调用结果。
type ToolResult struct {
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content,omitempty"`
}

// ContentBlock 是 IR 的最小内容单元。
type ContentBlock struct {
	Type       ContentType `json:"type"`
	Text       string      `json:"text,omitempty"`
	ToolUse    *ToolUse    `json:"tool_use,omitempty"`
	ToolResult *ToolResult `json:"tool_result,omitempty"`
}

// Message 是 IR 的一条消息。
type Message struct {
	Role    Role           `json:"role"`
	Content []ContentBlock `json:"content"`
}

// Tool 是 IR 的工具声明。
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

// FinishReason 是归一后的结束原因。
//
// ⚠️ 这是主要的协议差异点之一：
//   - OpenAI 用 finish_reason: stop | length | tool_calls | content_filter
//   - Anthropic 用 stop_reason:  end_turn | max_tokens | tool_use | stop_sequence
type FinishReason string

const (
	FinishStop          FinishReason = "stop"
	FinishLength        FinishReason = "length"
	FinishToolUse       FinishReason = "tool_use"
	FinishContentFilter FinishReason = "content_filter"
)

// Usage 是归一后的用量。
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	CachedTokens     int `json:"cached_tokens,omitempty"`
}

// Request 是归一后的请求。
type Request struct {
	Model    string    `json:"model"`
	System   string    `json:"system,omitempty"` // 归一：Anthropic 的顶层 system / OpenAI 的 role=system 消息都落到这里
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`

	MaxTokens   int      `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	Stream      bool     `json:"stream"`

	// ReasoningBudget 是归一后的思考预算。
	// 归一：OpenAI 的 reasoning_effort → Anthropic 的 thinking.budget_tokens。
	ReasoningBudget int `json:"reasoning_budget,omitempty"`
}

// Response 是归一后的响应。
type Response struct {
	ID           string         `json:"id"`
	Model        string         `json:"model"`
	Content      []ContentBlock `json:"content"`
	FinishReason FinishReason   `json:"finish_reason"`
	Usage        Usage          `json:"usage"`
}

// StreamEventType 是归一后的流式事件类型。
type StreamEventType string

const (
	EventStart          StreamEventType = "start"
	EventTextDelta      StreamEventType = "text_delta"
	EventReasoningDelta StreamEventType = "reasoning_delta"
	EventToolUseStart   StreamEventType = "tool_use_start"
	EventToolUseDelta   StreamEventType = "tool_use_delta"
	EventUsage          StreamEventType = "usage"
	EventDone           StreamEventType = "done"
	EventError          StreamEventType = "error"
)

// StreamEvent 是归一后的流式事件。
//
// ⚠️ 这是协议差异最大的地方：OpenAI 是统一的 data: 帧，
// Anthropic 是多事件类型（message_start / content_block_delta / message_delta ...）。
// IR 把它们压平成一维事件流，各 Codec 负责还原。
type StreamEvent struct {
	Type         StreamEventType `json:"type"`
	TextDelta    string          `json:"text_delta,omitempty"`
	ToolUseID    string          `json:"tool_use_id,omitempty"`
	ToolUseName  string          `json:"tool_use_name,omitempty"`
	InputDelta   string          `json:"input_delta,omitempty"`
	FinishReason FinishReason    `json:"finish_reason,omitempty"`
	Usage        *Usage          `json:"usage,omitempty"`
	Err          string          `json:"error,omitempty"`
}
