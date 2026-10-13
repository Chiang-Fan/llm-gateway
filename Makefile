# llm-gateway

BIN := bin/gateway

# 本项目为纯 Go。若本机 Go 与 CPU 架构不匹配（如 amd64 版 Go 跑在 arm64 机器上），
# cgo 会因缺少对应架构的 libxcrun 而失败，关掉即可。
CGO_ENABLED ?= 0
export CGO_ENABLED

# 关闭 VCS 版本戳（默认 -buildvcs=auto 会调用 git status）。
# 这是稳健性选择：在容器 / CI / 受限沙箱等环境下该调用可能失败并中断构建。
# 代价很小 —— 二进制不再嵌入 commit hash。若你的环境正常，去掉这行也能构建。
GOFLAGS += -buildvcs=false
export GOFLAGS

.PHONY: build run test fmt vet tidy clean

build:
	go build -o $(BIN) ./cmd/gateway

run:
	go run ./cmd/gateway -config config.yaml

test:
	go test ./...

fmt:
	gofmt -l -w .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/
