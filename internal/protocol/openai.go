package protocol

import (
	"context"
	"encoding/json"
	"errors"
)

// OpenAI 实现 OpenAI Chat Completions 协议的编解码。
//
// TODO(W43): DecodeRequest / EncodeRequest / DecodeResponse / EncodeResponse
// TODO(W43): DecodeStreamEvent / EncodeStreamEvent（统一 data: 帧 + [DONE]）
type OpenAI struct{}

// Name 返回协议名。
func (OpenAI) Name() string { return "openai" }

func (OpenAI) DecodeRequest(body []byte) (*Request, error) {
	var raw struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"messages"`
		Stream bool `json:"stream"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	// TODO(W43): 把 role=system 的消息归一进 Request.System；
	// 处理 content 的 string / 数组两种形态；解析 tool_calls 与 tools。
	return nil, errors.New("OpenAI.DecodeRequest: not implemented (W43)")
}

func (OpenAI) EncodeRequest(*Request) ([]byte, string, error) {
	return nil, "", errors.New("OpenAI.EncodeRequest: not implemented (W43)")
}

func (OpenAI) DecodeResponse([]byte) (*Response, error) {
	return nil, errors.New("OpenAI.DecodeResponse: not implemented (W43)")
}

func (OpenAI) EncodeResponse(*Response) ([]byte, error) {
	return nil, errors.New("OpenAI.EncodeResponse: not implemented (W43)")
}

func (OpenAI) DecodeStreamEvent(context.Context, []byte) ([]StreamEvent, error) {
	return nil, errors.New("OpenAI.DecodeStreamEvent: not implemented (W43)")
}

func (OpenAI) EncodeStreamEvent(*StreamEvent) ([][]byte, error) {
	return nil, errors.New("OpenAI.EncodeStreamEvent: not implemented (W43)")
}
