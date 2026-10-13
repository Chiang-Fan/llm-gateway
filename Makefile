.PHONY: build run test fmt vet clean

BIN := bin/gateway

# 本项目为纯 Go；若本机 Go 与架构不匹配（如 amd64 Go on arm64），关掉 cgo 即可
CGO_ENABLED ?= 0
export CGO_ENABLED

# 关闭 VCS 版本戳（默认 -buildvcs=auto 会调用 git status）。
# 原因：在部分受限环境 / 无 git 的构建环境（容器、CI）中该调用会失败并中断构建。
# 代价：二进制里不再嵌入 commit hash —— 对有镜像 tag 的场景影响很小。
GOFLAGS ?= -buildvcs=false
export GOFLAGS

build:
\tgo build -o $(BIN) ./cmd/gateway

run:
\tgo run ./cmd/gateway -config config.yaml

test:
\tgo test ./...

fmt:
\tgofmt -l -w .

vet:
\tgo vet ./...

clean:
\trm -rf bin/
