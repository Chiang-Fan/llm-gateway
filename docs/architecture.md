# 架构

## 分层

```
┌──────────────────────────────────────────────────────────┐
│ server           HTTP 层：路由注册、请求校验、优雅关闭      │
├──────────────────────────────────────────────────────────┤
│ protocol         ★ IR（中间表示）+ 各协议 Codec            │
│                  openai.go / anthropic.go                 │
├──────────────────────────────────────────────────────────┤
│ ratelimit        Check（只读判断） + Commit（事后扣减）     │
├──────────────────────────────────────────────────────────┤
│ router           从候选上游里挑一个                        │
├──────────────────────────────────────────────────────────┤
│ upstream         HTTP 客户端，含流式逐行回调               │
└──────────────────────────────────────────────────────────┘
```

## 一次请求的完整路径

```
1. 收到 POST /v1/chat/completions
2. 判定客户端协议 → Codec.DecodeRequest(body) → *protocol.Request (IR)
3. limiter.Check(key)              ← 只读判断，不扣减（ADR-001）
4. router.Pick(req, candidates)    → 选中的上游
5. Codec.EncodeRequest(req)        → 上游请求体 + 目标路径
6. upstream.Stream(...)            → 逐行回调
     └─ Codec.DecodeStreamEvent     → []StreamEvent (IR)
     └─ Codec.EncodeStreamEvent     → 客户端协议帧
     └─ flush                       ← 逐 chunk 下发，不缓冲
7. 上游结束 → 从 usage 取真实 token
8. limiter.Commit(key, actual)     ← 事后扣减
```

## 为什么要 IR

如果按「每两个协议写一个转换器」，支持 N 家厂商需要 **N×(N−1)** 个转换器。
引入 IR 后只需要 **N 个 Codec**，两两转换自动获得。

代价是 IR 必须足够表达所有协议的语义 —— 这是主要的设计难点，
也是 `internal/protocol/ir.go` 里注释最多的原因。

## 流式的难点

不在"转发字节"，而在：

1. **不缓冲** —— 缓冲会毁掉首包时延与实时性；
2. **跨 chunk 的半帧** —— TCP 分片不保证按 SSE 帧边界到达，必须自己拼行；
3. **协议形状差异** —— OpenAI 是统一 `data:` 帧，Anthropic 是多事件类型，
   IR 负责压平，各 Codec 负责还原；
4. **反向合成** —— 缓存命中时若客户端要 `stream=true`，
   需要把非流式响应**反向合成**成 SSE 帧序列（加分项）。
