# 构建阶段
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 优先复制依赖文件，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 静态编译，去掉调试信息以减小体积
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o server \
    ./cmd/server

# 运行阶段
FROM alpine:3.20

# 添加证书和时区数据
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/server .
COPY configs/ configs/

EXPOSE 8085

CMD ["./server", "-config", "configs/config.docker.yaml"]
