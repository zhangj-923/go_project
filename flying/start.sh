#!/bin/bash

# 斗地主游戏启动脚本

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}      斗地主游戏启动脚本        ${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

# 检查是否安装了必要的工具
check_dependencies() {
    echo -e "${YELLOW}检查依赖...${NC}"

    if ! command -v go &> /dev/null; then
        echo "错误: 未安装 Go"
        exit 1
    fi

    if ! command -v node &> /dev/null; then
        echo "错误: 未安装 Node.js"
        exit 1
    fi

    if ! command -v npm &> /dev/null; then
        echo "错误: 未安装 npm"
        exit 1
    fi

    echo -e "${GREEN}依赖检查通过${NC}"
}

# 安装前端依赖
install_frontend_deps() {
    echo -e "${YELLOW}安装前端依赖...${NC}"
    cd web
    npm install
    cd ..
    echo -e "${GREEN}前端依赖安装完成${NC}"
}

# 启动后端
start_backend() {
    echo -e "${YELLOW}启动后端服务器...${NC}"
    cd server
    go run main.go &
    BACKEND_PID=$!
    cd ..
    echo -e "${GREEN}后端服务器已启动 (PID: $BACKEND_PID)${NC}"
}

# 启动前端
start_frontend() {
    echo -e "${YELLOW}启动前端开发服务器...${NC}"
    cd web
    npm run dev &
    FRONTEND_PID=$!
    cd ..
    echo -e "${GREEN}前端开发服务器已启动 (PID: $FRONTEND_PID)${NC}"
}

# 等待用户按 Ctrl+C
wait_for_exit() {
    echo ""
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}      游戏已启动                ${NC}"
    echo -e "${GREEN}================================${NC}"
    echo ""
    echo -e "后端地址: ${YELLOW}http://localhost:8080${NC}"
    echo -e "前端地址: ${YELLOW}http://localhost:3000${NC}"
    echo ""
    echo -e "按 ${YELLOW}Ctrl+C${NC} 停止所有服务"
    echo ""

    # 捕获 Ctrl+C 信号
    trap cleanup INT

    # 等待子进程
    wait
}

# 清理函数
cleanup() {
    echo ""
    echo -e "${YELLOW}正在停止服务...${NC}"

    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null
        echo -e "${GREEN}后端服务已停止${NC}"
    fi

    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null
        echo -e "${GREEN}前端服务已停止${NC}"
    fi

    echo -e "${GREEN}所有服务已停止${NC}"
    exit 0
}

# 主流程
main() {
    check_dependencies
    install_frontend_deps
    start_backend
    sleep 2  # 等待后端启动
    start_frontend
    wait_for_exit
}

# 运行主流程
main
