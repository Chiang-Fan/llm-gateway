.PHONY: build run test fmt vet clean

BIN := bin/gateway
# 本项目为纯 Go；若本机 Go 与架构不匹配（如 amd64 Go on arm64），关掉 cgo 即可
CGO_ENABLED ?= 0
export CGO_ENABLED

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
