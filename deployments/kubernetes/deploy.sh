#!/bin/bash

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}开始部署Mall-Go微服务架构...${NC}"

# 创建命名空间
echo -e "${YELLOW}创建命名空间...${NC}"
kubectl apply -f namespace.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}创建命名空间失败，退出部署${NC}"
    exit 1
fi

# 部署基础设施
echo -e "${YELLOW}部署基础设施组件...${NC}"
echo -e "${YELLOW}部署MySQL...${NC}"
kubectl apply -f infrastructure/mysql.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}部署MySQL失败${NC}"
    exit 1
fi

echo -e "${YELLOW}部署Redis...${NC}"
kubectl apply -f infrastructure/redis.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}部署Redis失败${NC}"
    exit 1
fi

echo -e "${YELLOW}部署Consul...${NC}"
kubectl apply -f infrastructure/consul.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}部署Consul失败${NC}"
    exit 1
fi

# 等待基础设施就绪
echo -e "${YELLOW}等待基础设施组件就绪...${NC}"
kubectl wait --namespace=mall-go --for=condition=ready pod -l app=mysql --timeout=180s
kubectl wait --namespace=mall-go --for=condition=ready pod -l app=redis --timeout=120s
kubectl wait --namespace=mall-go --for=condition=ready pod -l app=consul --timeout=120s

# 部署微服务
echo -e "${YELLOW}部署微服务组件...${NC}"
echo -e "${YELLOW}部署用户服务(user-service)...${NC}"
kubectl apply -f services/user-service.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}部署用户服务失败${NC}"
    exit 1
fi

echo -e "${YELLOW}部署商品服务(product-service)...${NC}"
kubectl apply -f services/product-service.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}部署商品服务失败${NC}"
    exit 1
fi

# 等待核心服务就绪
echo -e "${YELLOW}等待核心服务就绪...${NC}"
kubectl wait --namespace=mall-go --for=condition=ready pod -l app=user-service --timeout=120s
kubectl wait --namespace=mall-go --for=condition=ready pod -l app=product-service --timeout=120s

echo -e "${YELLOW}部署API网关(gateway-service)...${NC}"
kubectl apply -f services/gateway-service.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}部署API网关失败${NC}"
    exit 1
fi

# 等待网关就绪
echo -e "${YELLOW}等待API网关就绪...${NC}"
kubectl wait --namespace=mall-go --for=condition=ready pod -l app=gateway-service --timeout=120s

# 部署Ingress
echo -e "${YELLOW}部署Ingress...${NC}"
kubectl apply -f ingress.yaml
if [ $? -ne 0 ]; then
    echo -e "${RED}部署Ingress失败${NC}"
    exit 1
fi

echo -e "${GREEN}mall-go微服务应用部署完成!${NC}"
echo -e "${BLUE}查看部署状态: kubectl get all -n mall-go${NC}"
echo -e "${BLUE}访问应用: http://mall-go.example.com${NC}"
echo -e "${YELLOW}注意: 请确保您已将mall-go.example.com添加到您的hosts文件中，或者已配置相应的DNS解析${NC}"

# 显示部署的资源
kubectl get all -n mall-go