#!/bin/bash

# 设置工作目录为项目根目录
cd "$(dirname "$0")/.." || exit
ROOT_DIR=$(pwd)

echo "===== 开始编译 Protocol Buffers 文件 ====="

# 编译 user-service 的 proto 文件
echo "编译 user-service proto 文件..."
protoc --proto_path=${ROOT_DIR} \
  --go_out=${ROOT_DIR} --go_opt=paths=source_relative \
  --go-grpc_out=${ROOT_DIR} --go-grpc_opt=paths=source_relative \
  services/user-service/proto/user.proto

# 编译 gateway-service 的 proto 文件
echo "编译 gateway-service proto 文件..."
protoc --proto_path=${ROOT_DIR} \
  --go_out=${ROOT_DIR} --go_opt=paths=source_relative \
  --go-grpc_out=${ROOT_DIR} --go-grpc_opt=paths=source_relative \
  services/gateway-service/proto/gateway.proto

# 编译 order-service 的 proto 文件
echo "编译 order-service proto 文件..."
protoc --proto_path=${ROOT_DIR} \
  --go_out=${ROOT_DIR} --go_opt=paths=source_relative \
  --go-grpc_out=${ROOT_DIR} --go-grpc_opt=paths=source_relative \
  services/order-service/proto/order.proto

# 编译 product-service 的 proto 文件
echo "编译 product-service proto 文件..."
protoc --proto_path=${ROOT_DIR} \
  --go_out=${ROOT_DIR} --go_opt=paths=source_relative \
  --go-grpc_out=${ROOT_DIR} --go-grpc_opt=paths=source_relative \
  services/product-service/proto/product.proto

echo "===== Protocol Buffers 文件编译完成 ====="