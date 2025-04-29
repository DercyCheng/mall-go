FROM golang:1.18-alpine AS builder

WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./services/gateway-service/cmd/server/main.go

# 使用更小的基础镜像
FROM alpine:latest

WORKDIR /app

# 安装基本工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 复制构建好的应用
COPY --from=builder /app/gateway .
COPY --from=builder /app/services/gateway-service/configs/config.yaml ./configs/

# 启动应用
ENTRYPOINT ["./gateway"]