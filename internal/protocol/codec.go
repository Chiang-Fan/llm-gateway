package protocol

import "context"

// Codec 是一个协议的编解码器。新增厂商 = 实现这个接口。
type Codec interface {
	// Name 返回协议名，如 "openai" / "anthropic"。
	Name() string

	// DecodeRequest 把入站请求体解析成 IR。
	DecodeRequest(body []byte) (*Request, error)

	// EncodeRequest 把 IR 序列化成该协议的上游请求体，并返回目标路径。
	EncodeRequest(req *Request) (body []byte, path string, err error)

	// DecodeResponse 把出站（非流式）响应体解析成 IR。
	DecodeResponse(body []byte) (*Response, error)

	// EncodeResponse 把 IR 序列化成该协议的响应体。
	EncodeResponse(resp *Response) ([]byte, error)

	// DecodeStreamEvent 解析该协议的一行流式数据，可能产出 0..N 个 IR 事件
	// （Anthropic 的多个事件类型需要累积状态）。
	DecodeStreamEvent(ctx context.Context, line []byte) ([]StreamEvent, error)

	// EncodeStreamEvent 把一个 IR 事件编码成该协议的一帧或多帧输出。
	EncodeStreamEvent(ev *StreamEvent) ([][]byte, error)
}

// registry 保存已注册的 Codec。
var registry = map[string]Codec{}

// Register 注册一个 Codec。
func Register(c Codec) { registry[c.Name()] = c }

// Get 按名字取 Codec。
func Get(name string) (Codec, bool) {
	c, ok := registry[name]
	return c, ok
}
