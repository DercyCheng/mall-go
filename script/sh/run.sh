#!/usr/bin/env bash
# mall-go 项目启动脚本

# 项目根目录
PROJECT_ROOT="/Users/dercyc/go/src/pro/mall-go"

# 定义颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
  echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
  echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
  echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
  echo "用法: ./run.sh [命令] [服务名称]"
  echo ""
  echo "命令:"
  echo "  build       构建服务"
  echo "  run         直接运行服务（不使用Docker）"
  echo "  docker      构建Docker镜像并运行容器"
  echo "  compose     使用docker-compose启动服务"
  echo "  stop        停止服务"
  echo "  clean       清理编译文件和容器"
  echo "  all         启动所有服务 (通过docker-compose)"
  echo "  network     创建Docker网络"
  echo "  mockgen     为指定接口生成mock"
  echo "  test        运行测试"
  echo ""
  echo "服务名称:"
  echo "  gateway     网关服务"
  echo "  user        用户服务"
  echo "  all         所有服务"
  echo ""
  echo "示例:"
  echo "  ./run.sh build gateway  # 构建网关服务"
  echo "  ./run.sh run user       # 运行用户服务"
  echo "  ./run.sh docker gateway # 以Docker方式运行网关服务"
  echo "  ./run.sh compose all    # 使用docker-compose启动所有服务"
  echo "  ./run.sh all            # 启动所有服务（默认使用docker-compose）"
  echo "  ./run.sh mockgen user   # 为用户服务生成mock"
  echo "  ./run.sh test user      # 运行用户服务的测试"
}

# 构建服务
build_service() {
  cd "$PROJECT_ROOT" || exit 1
  
  local service=$1
  local service_path=""
  local output_name=""
  
  case $service in
    gateway)
      service_path="services/gateway-service"
      output_name="gateway"
      ;;
    user)
      service_path="services/user-service"
      output_name="user-service"
      ;;
    *)
      print_error "未知服务: $service"
      return 1
      ;;
  esac
  
  print_message "构建 $service 服务..."
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $output_name ./$service_path/cmd/server/main.go
  
  if [ $? -eq 0 ]; then
    print_message "$service 服务构建成功！"
  else
    print_error "$service 服务构建失败！"
    return 1
  fi
}

# 运行服务 (非Docker方式)
run_service() {
  cd "$PROJECT_ROOT" || exit 1
  
  local service=$1
  local service_path=""
  local config_path=""
  
  case $service in
    gateway)
      service_path="services/gateway-service"
      config_path="services/gateway-service/configs/config.yaml"
      ;;
    user)
      service_path="services/user-service"
      config_path="services/user-service/configs/config.yaml"
      ;;
    *)
      print_error "未知服务: $service"
      return 1
      ;;
  esac
  
  print_message "运行 $service 服务..."
  CONFIG_PATH="./$config_path" go run ./$service_path/cmd/server/main.go
}

# Docker方式运行服务
docker_service() {
  cd "$PROJECT_ROOT" || exit 1
  
  local service=$1
  local image_name="mall-go/$service"
  local container_name="mall-go-$service"
  local dockerfile_path=""
  local port=""
  
  case $service in
    gateway)
      dockerfile_path="services/gateway-service/Dockerfile"
      port="8000:8000"
      ;;
    user)
      dockerfile_path="services/user-service/Dockerfile"
      port="8080:8080"
      ;;
    *)
      print_error "未知服务: $service"
      return 1
      ;;
  esac
  
  # 停止并删除已存在的容器
  print_message "停止并删除已存在的 $service 容器..."
  docker stop $container_name >/dev/null 2>&1
  docker rm $container_name >/dev/null 2>&1
  
  # 构建Docker镜像
  print_message "构建 $service 镜像..."
  docker build -t $image_name -f $dockerfile_path .
  
  if [ $? -ne 0 ]; then
    print_error "构建 $service 镜像失败！"
    return 1
  fi
  
  # 运行容器
  print_message "启动 $service 容器..."
  docker run -d -p $port --name $container_name \
    -e TZ="Asia/Shanghai" \
    $image_name
  
  if [ $? -eq 0 ]; then
    print_message "$service 容器启动成功！"
  else
    print_error "$service 容器启动失败！"
    return 1
  fi
}

# 使用docker-compose启动服务
compose_service() {
  cd "$PROJECT_ROOT" || exit 1
  
  local service=$1
  local compose_file=""
  
  if [ "$service" == "all" ]; then
    compose_file="script/docker/docker-compose-app.yml"
    print_message "使用docker-compose启动所有服务..."
    docker-compose -f $compose_file up -d
  else
    print_error "docker-compose仅支持启动所有服务"
    return 1
  fi
  
  if [ $? -eq 0 ]; then
    print_message "所有服务启动成功！"
  else
    print_error "启动服务失败！"
    return 1
  fi
}

# 停止服务
stop_service() {
  local service=$1
  
  if [ "$service" == "all" ]; then
    print_message "停止所有服务..."
    cd "$PROJECT_ROOT" || exit 1
    docker-compose -f script/docker/docker-compose-app.yml down
  else
    local container_name="mall-go-$service"
    print_message "停止 $service 服务..."
    docker stop $container_name
  fi
}

# 清理编译文件和容器
clean_service() {
  local service=$1
  
  if [ "$service" == "all" ]; then
    print_message "清理所有服务..."
    docker rm -f $(docker ps -a | grep mall-go | awk '{print $1}') 2>/dev/null
    docker rmi -f $(docker images | grep mall-go | awk '{print $3}') 2>/dev/null
    rm -f "$PROJECT_ROOT/gateway" "$PROJECT_ROOT/user-service"
  else
    local container_name="mall-go-$service"
    local image_name="mall-go/$service"
    local bin_name=""
    
    case $service in
      gateway)
        bin_name="gateway"
        ;;
      user)
        bin_name="user-service"
        ;;
    esac
    
    print_message "清理 $service 服务..."
    docker stop $container_name 2>/dev/null
    docker rm $container_name 2>/dev/null
    docker rmi $image_name 2>/dev/null
    rm -f "$PROJECT_ROOT/$bin_name"
  fi
}

# 启动所有依赖环境 (MySQL, Redis, Consul等)
start_dependencies() {
  cd "$PROJECT_ROOT" || exit 1
  
  print_message "启动基础环境依赖 (MySQL, Redis, Consul等)..."
  docker-compose -f script/docker/docker-compose-env.yml up -d
  
  if [ $? -eq 0 ]; then
    print_message "基础环境依赖启动成功！"
  else
    print_error "启动基础环境依赖失败！"
    return 1
  fi
}

# 创建Docker网络
create_network() {
  print_message "检查并创建Docker网络 mall-network..."
  
  # 检查网络是否已存在
  NETWORK_EXISTS=$(docker network ls | grep mall-network | wc -l)
  
  if (("$NETWORK_EXISTS" == "0")); then
    print_message "创建Docker网络 mall-network..."
    docker network create mall-network
    
    if [ $? -eq 0 ]; then
      print_message "Docker网络 mall-network 创建成功！"
    else
      print_error "创建Docker网络失败！"
      return 1
    fi
  else
    print_message "Docker网络 mall-network 已存在，无需创建。"
  fi
}

# 生成mock
generate_mocks() {
  cd "$PROJECT_ROOT" || exit 1
  
  local service=$1
  local service_path=""
  
  case $service in
    gateway)
      service_path="services/gateway-service"
      ;;
    user)
      service_path="services/user-service"
      ;;
    all)
      generate_mocks "gateway"
      generate_mocks "user"
      return 0
      ;;
    *)
      print_error "未知服务: $service"
      return 1
      ;;
  esac
  
  print_message "检查mockgen工具是否安装..."
  # 首先尝试在 PATH 中查找 mockgen
  if command -v mockgen &> /dev/null; then
    MOCKGEN_CMD="mockgen"
  # 然后尝试在 GOPATH/bin 中查找
  elif [ -f "${GOPATH}/bin/mockgen" ]; then
    MOCKGEN_CMD="${GOPATH}/bin/mockgen"
  # 最后尝试安装它
  else
    print_message "正在安装mockgen..."
    GO111MODULE=on go install github.com/golang/mock/mockgen@v1.6.0
    if [ $? -ne 0 ]; then
      print_error "安装mockgen失败！尝试直接使用go命令运行mockgen..."
      MOCKGEN_CMD="go run github.com/golang/mock/mockgen"
    else
      MOCKGEN_CMD="${GOPATH}/bin/mockgen"
    fi
  fi
  
  print_message "为 $service 服务生成mock..."
  
  if [ "$service" == "user" ]; then
    # 用户服务的接口
    $MOCKGEN_CMD -source=./$service_path/application/service/user_service.go -destination=./$service_path/mocks/user_service_mock.go -package=mocks
    $MOCKGEN_CMD -source=./$service_path/domain/repository/user_repository.go -destination=./$service_path/mocks/user_repository_mock.go -package=mocks
    
    # 检查是否存在角色仓库接口
    if [ -f "./$service_path/domain/repository/role_repository.go" ]; then
      $MOCKGEN_CMD -source=./$service_path/domain/repository/role_repository.go -destination=./$service_path/mocks/role_repository_mock.go -package=mocks
    fi
    
    # 检查是否存在缓存接口
    if [ -f "./$service_path/infrastructure/cache/user_cache.go" ]; then
      $MOCKGEN_CMD -source=./$service_path/infrastructure/cache/user_cache.go -destination=./$service_path/mocks/user_cache_mock.go -package=mocks
    fi
  elif [ "$service" == "gateway" ]; then
    # 网关服务的接口
    if [ -f "./$service_path/api/proxy/service_proxy.go" ]; then
      $MOCKGEN_CMD -source=./$service_path/api/proxy/service_proxy.go -destination=./$service_path/mocks/service_proxy_mock.go -package=mocks
    fi
  fi
  
  print_message "$service 服务的mock生成成功！"
}

# 运行测试
run_tests() {
  cd "$PROJECT_ROOT" || exit 1
  
  local service=$1
  local service_path=""
  local test_path=""
  
  case $service in
    gateway)
      service_path="services/gateway-service"
      ;;
    user)
      service_path="services/user-service"
      ;;
    all)
      run_tests "gateway"
      run_tests "user"
      return 0
      ;;
    *)
      print_error "未知服务: $service"
      return 1
      ;;
  esac
  
  print_message "运行 $service 服务的测试..."
  go test -v ./$service_path/... -coverprofile=coverage.out
  
  if [ $? -eq 0 ]; then
    print_message "$service 服务的测试运行成功！"
    go tool cover -html=coverage.out -o coverage.html
    print_message "测试覆盖率报告已生成: coverage.html"
  else
    print_error "$service 服务的测试运行失败！"
    return 1
  fi
}

# 主函数
main() {
  local command=$1
  local service=$2
  
  # 如果没有参数或者第一个参数是 -h 或 --help，显示帮助信息
  if [ $# -eq 0 ] || [ "$command" == "-h" ] || [ "$command" == "--help" ]; then
    show_help
    exit 0
  fi
  
  # 如果只有一个参数且为 "all"，默认启动所有服务
  if [ $# -eq 1 ] && [ "$command" == "all" ]; then
    create_network
    start_dependencies
    compose_service "all"
    exit 0
  fi
  
  # 如果没有指定服务，设置为默认值 "all"
  if [ -z "$service" ]; then
    service="all"
  fi
  
  # 执行相应的命令
  case $command in
    build)
      if [ "$service" == "all" ]; then
        build_service "gateway"
        build_service "user"
      else
        build_service "$service"
      fi
      ;;
    run)
      if [ "$service" == "all" ]; then
        print_error "不支持直接运行所有服务，请使用docker-compose或分别启动"
        exit 1
      else
        run_service "$service"
      fi
      ;;
    docker)
      if [ "$service" == "all" ]; then
        docker_service "gateway"
        docker_service "user"
      else
        docker_service "$service"
      fi
      ;;
    compose)
      create_network
      compose_service "$service"
      ;;
    stop)
      stop_service "$service"
      ;;
    clean)
      clean_service "$service"
      ;;
    deps)
      create_network
      start_dependencies
      ;;
    network)
      create_network
      ;;
    mockgen)
      generate_mocks "$service"
      ;;
    test)
      run_tests "$service"
      ;;
    *)
      print_error "未知命令: $command"
      show_help
      exit 1
      ;;
  esac
}

# 执行主函数
main "$@"