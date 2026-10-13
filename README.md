# mini-llm-gateway

一个**最小的大模型网关**：把多家厂商的私有协议统一成 OpenAI 兼容入口，并提供 token 级限流与多上游路由。

目标是**把大模型网关的几块核心机制写清楚**，而不是做一个功能完备的产品。

---

## 它做什么

```
客户端 (OpenAI 格式)
        │
        ▼
┌───────────────────┐
│  协议适配 (IR)     │  OpenAI / Anthropic 双向转换
├───────────────────┤
│  限流              │  token 级配额，先放行、后付费
├───────────────────┤
│  路由              │  加权随机 + 失败降级
├───────────────────┤
│  上游客户端        │  SSE 流式透传，逐块 flush
└───────────────────┘
        │
        ▼
   OpenAI / Anthropic / 兼容上游
```

> 架构细节见 [`docs/architecture.md`](docs/architecture.md)；
> **关键设计取舍见 [`docs/decisions.md`](docs/decisions.md)**（先读这个）。

---

## 范围

**必做（4 项）**

1. OpenAI 兼容入口 + **Anthropic Messages 双向协议转换**
2. **SSE 流式透传** —— 不缓冲，逐 chunk flush
3. **token 级限流 + 配额** —— 固定窗口
4. **多上游路由** + 失败降级到备用模型

**加分（2 项）**

5. 用量计量与成本统计
6. 极简语义缓存

**明确不做**（理由见 `docs/decisions.md`）

- ❌ 自适应路由（延迟/前缀亲和）—— 样本量不足时期望收益为负
- ❌ 语义缓存之外的向量能力、MCP / A2A 接入、内容安全 —— 属于产品能力，不是网关内核
- ❌ 预扣回补式限流 —— 见 ADR-001

---

## 快速开始

```bash
make build
./bin/gateway -config config.yaml
curl localhost:8080/healthz
```

需要 Go 1.21+。

---

## 状态

起盘阶段。协议适配、流式透传、限流与路由的骨架已就位（接口 + 类型定义），
具体实现在 `internal/protocol/*.go` 的 `TODO` 处逐步填充。

```
internal/
├── config/     配置
├── protocol/   ★ IR + 各协议 Codec（核心）
├── ratelimit/  token 级限流
├── router/     上游选择
├── server/     HTTP 层
└── upstream/   上游客户端（含流式）
```

---

## 许可

MIT
