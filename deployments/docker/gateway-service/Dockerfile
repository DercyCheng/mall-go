# 构建阶段
FROM golang:1.20-alpine AS builder

# 设置Go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置工作目录
WORKDIR /build

# 复制go.mod和go.sum文件
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN go build -ldflags="-s -w" -o gateway-service ./services/gateway-service/cmd/server

# 运行阶段
FROM alpine:latest

# 安装基础软件
RUN apk --no-cache add ca-certificates tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的应用
COPY --from=builder /build/gateway-service .

# 复制配置文件
COPY services/gateway-service/configs /app/configs

# 声明暴露的端口
EXPOSE 8000

# 设置健康检查
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD curl -f http://localhost:8000/health || exit 1

# 启动应用
CMD ["./gateway-service"]